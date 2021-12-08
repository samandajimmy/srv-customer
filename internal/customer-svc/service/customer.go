package service

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ntime"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
	"strings"
	"time"

	"github.com/rs/xid"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/convert"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nclient"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
)

type Customer struct {
	customerRepo        contract.CustomerRepository
	verificationOTPRepo contract.VerificationOTPRepository
	OTPRepo             contract.OTPRepository
	credentialRepo      contract.CredentialRepository
	accessSessionRepo   contract.AccessSessionRepository
	auditLoginRepo      contract.AuditLoginRepository
	verificationRepo    contract.VerificationRepository
	otpService          contract.OTPService
	cacheService        contract.CacheService
	clientConfig        contract.ClientConfig
	response            *ncore.ResponseMap
}

func (c *Customer) HasInitialized() bool {
	return true
}

func (c *Customer) Init(app *contract.PdsApp) error {
	c.customerRepo = app.Repositories.Customer
	c.verificationOTPRepo = app.Repositories.VerificationOTP
	c.OTPRepo = app.Repositories.OTP
	c.credentialRepo = app.Repositories.Credential
	c.accessSessionRepo = app.Repositories.AccessSession
	c.auditLoginRepo = app.Repositories.AuditLogin
	c.verificationRepo = app.Repositories.Verification
	c.otpService = app.Services.OTP
	c.cacheService = app.Services.Cache
	c.clientConfig = app.Config.Client
	c.response = app.Responses
	return nil
}

func (c *Customer) Login(payload dto.LoginRequest) (*dto.LoginResponse, error) {

	// Check if user exists
	t := time.Now()
	customer, err := c.customerRepo.FindByEmailOrPhone(payload.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("failed to retrieve customer not found", nlogger.Error(err))
			return nil, c.response.GetError("E_RES_1")
		}
		log.Errorf("failed to retrieve customer. error: %v", err)
		return nil, ncore.TraceError(err)
	}

	// Get credential customer
	credential, err := c.credentialRepo.FindByCustomerId(customer.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("failed to retrieve credential not found", nlogger.Error(err))
			return nil, c.response.GetError("E_AUTH_8")
		}
		log.Errorf("failed to retrieve credential. error: %v", err)
		return nil, ncore.TraceError(err)
	}

	// Check if account isn't blocked
	blockedUntil := ntime.ChangeTimezone(credential.BlockedUntilAt.Time, constant.WIB)
	now := ntime.ChangeTimezone(t, constant.WIB)
	if credential.BlockedUntilAt.Valid != false && blockedUntil.After(now) {
		return nil, c.response.GetError("E_AUTH_9")
	}

	// Counter wrong password count
	passwordRequest := fmt.Sprintf("%x", md5.Sum([]byte(payload.Password)))
	if credential.Password != passwordRequest {
		err := c.HandleWrongPassword(credential)
		return nil, err
	}

	// Check account is first login or not
	countAuditLog, err := c.auditLoginRepo.CountLogin(customer.Id)
	if err != nil {
		return nil, c.response.GetError("E_AUTH_1")
	}

	// Set is first login is true or false.
	var isFirstLogin = true
	if countAuditLog > 0 {
		isFirstLogin = false
	}

	// Prepare to insert audit login
	auditLogin := model.AuditLogin{
		CustomerId:   customer.Id,
		ChannelId:    GetChannelByAgen(payload.Agen),
		DeviceId:     payload.DeviceId,
		IP:           payload.IP,
		Latitude:     payload.Latitude,
		Longitude:    payload.Longitude,
		Timestamp:    t.Format(time.RFC3339),
		Timezone:     payload.Timezone,
		Brand:        payload.Brand,
		OsVersion:    payload.OsVersion,
		Browser:      payload.Browser,
		UseBiometric: payload.UseBiometric,
		Status:       1,
		Metadata:     []byte("{}"),
		ItemMetadata: model.NewItemMetadata(
			convert.ModifierDTOToModel(
				dto.Modifier{ID: "", Role: "", FullName: ""},
			),
		),
	}

	// Persist audit login
	err = c.auditLoginRepo.Insert(&auditLogin)
	if err != nil {
		log.Errorf("Error when insert audit login error: %v", err)
		return nil, c.response.GetError("E_AUTH_1")
	}

	// Get token from cache
	var token string
	cacheTokenKey := fmt.Sprintf("%v:%v:%v", constant.PREFIX, "token_jwt", customer.Id)
	token, _ = c.cacheService.Get(cacheTokenKey)
	if token == "" {
		// Generate token authentication
		token, err = c.SetTokenAuthentication(customer, payload.Agen, payload.Version, cacheTokenKey)
		if err != nil {
			log.Errorf("Failed to generate token jwt: %v", err)
			return nil, ncore.TraceError(err)
		}
	}

	// get user data

	// TODO get tabungan emas service

	// TODO check is force update password

	return &dto.LoginResponse{
		Customer: &dto.CustomerVO{
			ID:                        nval.ParseStringFallback(customer.Id, ""),
			Cif:                       customer.Cif,
			IsKYC:                     "1",
			Nama:                      customer.FullName,
			NamaIbu:                   "",
			NoKTP:                     "No ktp",
			Email:                     customer.Email,
			JenisKelamin:              "L",
			TempatLahir:               "Jakarta",
			TglLahir:                  "",
			Alamat:                    "alamat",
			IDProvinsi:                "idProvinsi",
			IDKabupaten:               "idKabupaten",
			IDKecamatan:               "idKecamatan",
			IDKelurahan:               "idKelurahan",
			Kelurahan:                 "Kelurahan",
			Provinsi:                  "provinsi",
			Kabupaten:                 "Kabupaten",
			Kecamatan:                 "Kecamantan",
			KodePos:                   "Kodepos",
			NoHP:                      customer.Phone,
			Avatar:                    "avatar",
			FotoKTP:                   "",
			IsEmailVerified:           "1",
			Kewarganegaraan:           "",
			JenisIdentitas:            "",
			NoIdentitas:               "",
			TglExpiredIdentitas:       "",
			NoNPWP:                    "npwp",
			NoSid:                     "",
			FotoSid:                   "",
			StatusKawin:               "2",
			Norek:                     "norek",
			Saldo:                     "saldo",
			AktifasiTransFinansial:    "transfinansial",
			IsDukcapilVerified:        "ok",
			IsOpenTe:                  "ok",
			ReferralCode:              "referal",
			GoldCardApplicationNumber: "",
			GoldCardAccountNumber:     "",
			KodeCabang:                "",
			TabunganEmas:              false,
			IsFirstLogin:              isFirstLogin,
			IsForceUpdatePassword:     false,
		},
		JwtToken: token,
	}, nil
}

func (c *Customer) SetTokenAuthentication(customer *model.Customer, agen string, version string, cacheTokenKey string) (string, error) {

	// Generate access token
	accessToken := nval.Bin2Hex(nval.RandStringBytes(78))
	channelId := GetChannelByAgen(agen)
	now := time.Now()

	// Generate JWT
	token, err := jwt.NewBuilder().
		Claim("id", customer.Id).
		Claim("email", customer.Email).
		Claim("nama", customer.FullName).
		Claim("no_hp", customer.Phone).
		Claim("access_token", accessToken).
		Claim("agen", agen).
		Claim("channelId", channelId).
		Claim("version", version).
		IssuedAt(now).
		Expiration(now.Add(time.Second * time.Duration(c.clientConfig.JWTExpired))).
		Issuer("https://www.pegadaian.co.id").
		Build()
	if err != nil {
		return "", err
	}

	// Sign token
	signed, err := jwt.Sign(token, jwa.HS256, []byte(c.clientConfig.JWTKey))
	if err != nil {
		return "", err
	}
	tokenString := string(signed)

	// Set token to cache
	cacheToken, err := c.cacheService.SetThenGet(cacheTokenKey, tokenString, c.clientConfig.JWTExpired)
	if err != nil {
		return "", err
	}

	return cacheToken, nil
}

func (c *Customer) HandleWrongPassword(credential *model.Credential) error {
	var code string
	t := time.Now()

	// Unmarshalling metadata credential to get tryLoginAt
	var Metadata dto.MetadataCredential
	err := json.Unmarshal(credential.Metadata, &Metadata)
	if err != nil {
		log.Errorf("Cannot unmarshaling metadata credential. err: %v", err)
		return ncore.TraceError(err)
	}
	// Parse time from metadata string to time
	tryLoginAt, _ := time.Parse(time.RFC3339, Metadata.TryLoginAt)
	now := ntime.ChangeTimezone(t, constant.WIB)

	// If user is not trying to login after 1 day, set wrongPassword to 0
	if now.After(tryLoginAt.Add(time.Hour * time.Duration(24))) {
		credential.WrongPasswordCount = 0
	}

	// If now is after than blockedUntilAt set wrong password to 0 and unblock account
	blockedUntil := ntime.ChangeTimezone(credential.BlockedUntilAt.Time, constant.WIB)
	if credential.WrongPasswordCount == constant.MAX_WRONG_PASSWORD && now.After(blockedUntil) {
		credential.BlockedAt = sql.NullTime{}
		credential.BlockedUntilAt = sql.NullTime{}
		credential.WrongPasswordCount = 0
	}

	wrongCount := credential.WrongPasswordCount + 1

	switch wrongCount {
	case constant.WARN_2X_WRONG_PASSWORD:
		code = "E_AUTH_6"
		credential.WrongPasswordCount = wrongCount
	case constant.WARN_4X_WRONG_PASSWORD:
		code = "E_AUTH_7"
		credential.WrongPasswordCount = wrongCount
	case constant.MIN_WRONG_PASSWORD:
		code = "E_AUTH_9"
		// Set block account
		hour := 1 // Block for 1 hours
		duration := time.Hour * time.Duration(hour)
		credential.BlockedAt = sql.NullTime{
			Time:  t,
			Valid: true,
		}
		credential.BlockedUntilAt = sql.NullTime{
			Time:  t.Add(duration),
			Valid: true,
		}
		credential.WrongPasswordCount = wrongCount
		// TODO sendNotificationBlockedLoginOneHour
		break
	case constant.MAX_WRONG_PASSWORD:
		code = "E_AUTH_9"
		// Set block account
		hour := 24 // Block for 24 hours
		duration := time.Hour * time.Duration(hour)
		credential.BlockedAt = sql.NullTime{
			Time:  t,
			Valid: true,
		}
		credential.BlockedUntilAt = sql.NullTime{
			Time:  t.Add(duration),
			Valid: true,
		}
		credential.WrongPasswordCount = wrongCount
		// TODO sendNotificationBlockedLoginOneDay
		break
	default:
		code = "E_AUTH_8"
		credential.WrongPasswordCount = wrongCount
		break
	}

	// Set trying login at to metadata
	var Format dto.MetadataCredential
	Format.TryLoginAt = t.Format(time.RFC3339)
	Format.PinCreatedAt = Metadata.PinCreatedAt
	Format.PinBlockedAt = Metadata.PinBlockedAt
	MetadataCredential, _ := json.Marshal(&Format)
	credential.Metadata = MetadataCredential

	err = c.credentialRepo.UpdateByCustomerID(credential)
	if err != nil {
		log.Errorf("Error when update credential when password is invalid.")
		return ncore.TraceError(err)
	}

	return c.response.GetError(code)
}

func GetChannelByAgen(agen string) string {

	// Generalize agen
	agen = strings.ToLower(agen)

	if agen == constant.AGEN_ANDROID {
		return constant.CHANNEL_ANDROID
	}

	if agen == constant.AGEN_MOBILE {
		return constant.CHANNEL_MOBILE
	}

	if agen == constant.AGEN_WEB {
		return constant.CHANNEL_WEB
	}

	return ""
}

func (c *Customer) Register(payload dto.RegisterNewCustomer) (*dto.RegisterNewCustomerResponse, error) {
	// validate exist
	customer, err := c.customerRepo.FindByEmail(payload.Email)
	if err != nil {
		log.Errorf("error while retrieve by email: %s", payload.Email)
		return nil, c.response.GetError("E_REG_1")
	}
	if customer != nil {
		log.Debugf("Email already registered: %s", payload.RegistrationId)
		return nil, c.response.GetError("E_REG_2")
	}
	// find registerID
	registerOTP, err := c.verificationOTPRepo.FindByRegistrationIdAndPhone(payload.RegistrationId, payload.PhoneNumber)
	if err != nil {
		log.Errorf("Registration ID not found: %s", payload.RegistrationId)
		return nil, c.response.GetError("E_REG_1")
	}

	// Get data user
	customer, _ = c.customerRepo.FindByPhone(payload.PhoneNumber)

	var customerXID string
	var customerId int64
	// update name if customer exists
	if customer != nil {
		customer.FullName = payload.Name
		customer.Email = payload.Email
		err := c.customerRepo.UpdateByPhone(customer)
		if err != nil {
			return nil, c.response.GetError("E_AUTH_5")
		}
		customerXID = customer.CustomerXID
		customerId = customer.Id
	} else {
		// create new one
		customerXID = strings.ToUpper(xid.New().String())
		metaData := model.NewItemMetadata(
			convert.ModifierDTOToModel(
				dto.Modifier{ID: "", Role: "", FullName: ""},
			),
		)
		insertCustomer := &model.Customer{
			CustomerXID:    customerXID,
			FullName:       payload.Name,
			Phone:          payload.PhoneNumber,
			Email:          payload.Email,
			Status:         0,
			IdentityType:   0,
			IdentityNumber: "",
			UserRefId:      0,
			Photos:         []byte("{}"),
			Profile:        []byte("{}"),
			Cif:            "",
			Sid:            "",
			ReferralCode:   "",
			Metadata:       []byte("{}"),
			ItemMetadata:   metaData,
		}
		lastInsertId, err := c.customerRepo.Insert(insertCustomer)
		if err != nil {
			log.Errorf("Error when persist customer : %s", payload.Name)
			return nil, ncore.TraceError(err)
		}
		customerId = lastInsertId
		customer = insertCustomer
		customerXID = customer.CustomerXID
	}

	customer.Status = 1
	err = c.customerRepo.UpdateByPhone(customer)
	if err != nil {
		log.Errorf("Error when update customer : %s", payload.Name)
		return nil, c.response.GetError("E_REG_1")
	}

	// set metadata credential
	var Format dto.MetadataCredential
	Format.TryLoginAt = ""
	Format.PinCreatedAt = ""
	Format.PinBlockedAt = ""
	Metadata, _ := json.Marshal(&Format)

	// create credential
	credentialXID := strings.ToUpper(xid.New().String())
	credentialInsert := &model.Credential{
		Xid:                 credentialXID,
		CustomerId:          customerId,
		Password:            fmt.Sprintf("%x", md5.Sum([]byte(payload.Password))),
		NextPasswordResetAt: nil,
		Pin:                 "",
		PinCif:              "",
		PinUpdatedAt:        nil,
		PinLastAccessAt:     nil,
		PinCounter:          0,
		PinBlockedStatus:    0,
		IsLocked:            0,
		LoginFailCount:      0,
		WrongPasswordCount:  0,
		BlockedAt:           sql.NullTime{},
		BlockedUntilAt:      sql.NullTime{},
		BiometricLogin:      0,
		BiometricDeviceId:   "",
		Metadata:            Metadata,
		ItemMetadata:        model.NewItemMetadata(convert.ModifierDTOToModel(dto.Modifier{ID: "", Role: "", FullName: ""})),
	}
	err = c.credentialRepo.InsertOrUpdate(credentialInsert)
	if err != nil {
		log.Errorf("Error when persist customer credential : %s", payload.Name)
		return nil, c.response.GetError("E_REG_1")
	}

	// insert OTP
	insertOTP := &model.OTP{
		CustomerId: customerId,
		Content:    "",
		Type:       "registrasi_user",
		Data:       "",
		Status:     "",
		UpdatedAt:  time.Now(),
	}
	err = c.OTPRepo.Insert(insertOTP)
	if err != nil {
		log.Errorf("Error when persist OTP. Err: %s", err)
		return nil, ncore.TraceError(err)
	}

	// insert session
	insertAccessSession := &model.AccessSession{
		Xid:                  customerXID,
		CustomerId:           customerId,
		ExpiredAt:            time.Now().Add(1 * time.Hour),
		NotificationToken:    payload.FcmToken,
		NotificationProvider: 1,
		Metadata:             []byte("{}"),
		ItemMetadata:         model.NewItemMetadata(convert.ModifierDTOToModel(dto.Modifier{ID: "", Role: "", FullName: ""})),
	}
	err = c.accessSessionRepo.Insert(insertAccessSession)
	if err != nil {
		log.Errorf("Error when persist access session: %s. Err: %s", payload.Name, err)
		return nil, c.response.GetError("E_REG_1")
	}

	// TODO email_verification_token

	// TODO call login service

	// Delete OTP RegistrationId
	err = c.verificationOTPRepo.Delete(registerOTP.RegistrationId, customer.Phone)
	if err != nil {
		log.Debugf("Error when remove by registration id : %s, phone : %s", registerOTP.RegistrationId, customer.Phone)
		return nil, c.response.GetError("E_REG_1")
	}

	return &dto.RegisterNewCustomerResponse{
		Name:        customer.FullName,
		Email:       customer.Email,
		PhoneNumber: customer.Phone,
	}, nil
}

func (c *Customer) RegisterStepOne(payload dto.RegisterStepOne) (*dto.RegisterStepOneResponse, error) {

	// validate email
	emailExist, err := c.customerRepo.FindByEmail(payload.Email)
	if err != nil {
		log.Errorf("error while retrieve by email: %s", payload.Email)
		return nil, c.response.GetError("E_REG_1")
	}
	if emailExist != nil {
		log.Debugf("Email already registered")
		return nil, c.response.GetError("E_REG_2")
	}

	// validate phone
	phoneExist, err := c.customerRepo.FindByPhone(payload.PhoneNumber)
	if err != nil {
		log.Errorf("error while retrieve by phone: %s", payload.PhoneNumber)
		return nil, c.response.GetError("E_REG_1")
	}
	if phoneExist != nil {
		log.Debugf("Phone already registered")
		return nil, c.response.GetError("E_REG_3")
	}

	// Set request
	request := dto.SendOTPRequest{
		PhoneNumber: payload.PhoneNumber,
		RequestType: "register",
	}

	// Send OTP To Phone Number
	resp, err := c.otpService.SendOTP(request)
	if err != nil {
		return nil, ncore.TraceError(err)
	}

	// Extract response from server
	data, err := nclient.GetResponseData(resp)

	if data.ResponseCode == "15" {
		return nil, c.response.GetError("E_OTP_3")
	}

	if data.Message != "" {
		log.Errorf("Debug: RegisterStepOne OTP CODE %s", data.Message)
	}

	return &dto.RegisterStepOneResponse{
		Action: data.ResponseDesc,
	}, nil
}

func (c *Customer) RegisterStepTwo(payload dto.RegisterStepTwo) (*dto.RegisterStepTwoResponse, error) {
	// Set request
	request := dto.VerifyOTPRequest{
		PhoneNumber: payload.PhoneNumber,
		Token:       payload.OTP,
		RequestType: "register",
	}
	// validate phone
	phoneExist, err := c.customerRepo.FindByPhone(payload.PhoneNumber)
	if err != nil {
		log.Errorf("error while retrieve by phone: %s", payload.PhoneNumber)
		return nil, c.response.GetError("E_REG_1")
	}
	if phoneExist != nil {
		log.Debugf("Phone already registered")
		return nil, c.response.GetError("E_REG_3")
	}

	// Verify OTP To Phone Number
	resp, err := c.otpService.VerifyOTP(request)
	if err != nil {
		return nil, ncore.TraceError(err)
	}

	// Extract response from server
	data, err := nclient.GetResponseData(resp)

	// wrong otp handle
	if data.ResponseCode != "00" {
		log.Errorf("Wrong OTP. Phone Number : %s", payload.PhoneNumber)
		return nil, c.response.GetError("E_OTP_1")
	}

	// insert verification otp
	registrationId := xid.New().String()
	insert := &model.VerificationOTP{
		CreatedAt:      time.Now(),
		Phone:          payload.PhoneNumber,
		RegistrationId: registrationId,
	}
	_, err = c.verificationOTPRepo.Insert(insert)
	if err != nil {
		log.Errorf("Error when persist verificationOTP. Phone Number: %s", payload.PhoneNumber)
		return nil, ncore.TraceError(err)
	}

	return &dto.RegisterStepTwoResponse{
		RegisterId: registrationId,
	}, nil
}

func (c *Customer) RegisterResendOTP(payload dto.RegisterResendOTP) (*dto.RegisterResendOTPResponse, error) {
	// Set request
	request := dto.SendOTPRequest{
		PhoneNumber: payload.PhoneNumber,
		RequestType: "register",
	}

	// Send OTP To Phone Number
	resp, err := c.otpService.SendOTP(request)
	if err != nil {
		return nil, ncore.TraceError(err)
	}

	// Extract response from server
	data, err := nclient.GetResponseData(resp)

	if data.ResponseCode == "15" {
		return nil, c.response.GetError("E_OTP_3")
	}

	if data.Message != "" {
		log.Errorf("Debug: RegisterResendOTP OTP CODE %s", data.Message)
	}

	return &dto.RegisterResendOTPResponse{
		Action: data.ResponseDesc,
	}, nil
}
