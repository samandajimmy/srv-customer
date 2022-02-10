package service

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/nbs-go/nlogger"
	"github.com/rs/xid"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/convert"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nclient"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ntime"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
)

type Customer struct {
	customerRepo             contract.CustomerRepository
	verificationOTPRepo      contract.VerificationOTPRepository
	OTPRepo                  contract.OTPRepository
	credentialRepo           contract.CredentialRepository
	accessSessionRepo        contract.AccessSessionRepository
	auditLoginRepo           contract.AuditLoginRepository
	verificationRepo         contract.VerificationRepository
	financialDataRepo        contract.FinancialDataRepository
	addressRepo              contract.AddressRepository
	userExternalRepo         contract.UserExternalRepository
	userPinExternalRepo      contract.UserPinExternalRepository
	userRegisterExternalRepo contract.UserRegisterExternalRepository
	otpService               contract.OTPService
	cacheService             contract.CacheService
	notificationService      contract.NotificationService
	pdsAPIService            contract.PdsAPIService
	clientConfig             contract.ClientConfig
	response                 *ncore.ResponseMap
	httpBaseUrl              string
	emailConfig              contract.EmailConfig
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
	c.financialDataRepo = app.Repositories.FinancialData
	c.addressRepo = app.Repositories.Address
	c.userExternalRepo = app.Repositories.UserExternal
	c.userPinExternalRepo = app.Repositories.UserPinExternal
	c.userRegisterExternalRepo = app.Repositories.UserRegisterExternal
	c.otpService = app.Services.OTP
	c.cacheService = app.Services.Cache
	c.pdsAPIService = app.Services.PdsAPI
	c.clientConfig = app.Config.Client
	c.notificationService = app.Services.Notification
	c.response = app.Responses
	return nil
}

func (c *Customer) Login(payload dto.LoginRequest) (*dto.LoginResponse, error) {

	// Check if user exists
	t := time.Now()
	customer, err := c.customerRepo.FindByEmailOrPhone(payload.Email)
	if errors.Is(err, sql.ErrNoRows) {
		// If data not found on internal database check on external database.
		user, err := c.userExternalRepo.FindByEmailOrPhone(payload.Email)
		if err != nil {
			log.Error("failed to retrieve customer not found", nlogger.Error(err))
			return nil, c.response.GetError("E_AUTH_10")
		}

		// sync data external to internal
		customer, err = c.syncExternalToInternal(user)
		if err != nil {
			log.Error("error while sync data External to Internal", nlogger.Error(err))
			return nil, err
		}

	} else if err != nil {
		log.Errorf("failed to retrieve customer. error: %v", err)
		return nil, ncore.TraceError("error", err)
	}

	// Get credential customer
	credential, err := c.credentialRepo.FindByCustomerId(customer.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("failed to retrieve credential not found", nlogger.Error(err))
			return nil, c.response.GetError("E_AUTH_8")
		}
		log.Errorf("failed to retrieve credential. error: %v", err)
		return nil, ncore.TraceError("error", err)
	}

	// Check if account isn't blocked
	blockedUntil := ntime.ChangeTimezone(credential.BlockedUntilAt.Time, constant.WIB)
	now := ntime.ChangeTimezone(t, constant.WIB)
	if credential.BlockedUntilAt.Valid != false && blockedUntil.After(now) {
		// Set response if blocked
		message := "Akun dikunci hingga %v karena gagal login %v kali. Hubungi call center jika ini bukan kamu"
		timeBlocked := ntime.ChangeTimezone(credential.BlockedUntilAt.Time, constant.WIB).Format("02-Jan-2006 15:04:05")
		setResponse := ncore.Success.HttpStatus(401).SetMessage(message, timeBlocked, credential.WrongPasswordCount)
		return nil, &setResponse
	}

	// Counter wrong password count
	passwordRequest := fmt.Sprintf("%x", md5.Sum([]byte(payload.Password)))
	if credential.Password != passwordRequest {
		err := c.HandleWrongPassword(credential, customer)
		return nil, err
	}

	// get userRefId from external DB
	if customer.UserRefId == "" {
		registerPayload := &dto.CustomerSynchronizeRequest{
			Name:        customer.FullName,
			Email:       customer.Email,
			PhoneNumber: customer.Phone,
			Password:    credential.Password,
			FcmToken:    payload.FcmToken,
		}
		resultSync, err := c.syncInternalToExternal(registerPayload)

		if err != nil {
			log.Errorf("failed to sync to external. error: %v", err)
			return nil, ncore.TraceError("error", err)
		}
		// set userRefId
		customer.UserRefId = nval.ParseStringFallback(resultSync.UserAiid, "")
		// update customer
		err = c.customerRepo.UpdateByPhone(customer)
		if err != nil {
			log.Errorf("failed to update userRefId. error: %v", err)
			return nil, ncore.TraceError("error", err)
		}
	}

	// Get token from cache
	var token string
	cacheTokenKey := fmt.Sprintf("%v:%v:%v", constant.Prefix, "token_jwt", customer.UserRefId)
	token, err = c.SetTokenAuthentication(customer, payload.Agen, payload.Version, cacheTokenKey)

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

	// Unmarshall profile
	var profile dto.CustomerProfileVO
	err = json.Unmarshal(customer.Profile, &profile)
	if err != nil {
		log.Errorf("Error when unmarshalling profile: %v", err)
		return nil, ncore.TraceError("error", err)
	}

	// TODO get tabungan emas service / get tabungan emas from financial table

	// Check is force update password
	validatePassword := ValidatePassword(payload.Password)
	isForceUpdatePassword := false
	if validatePassword.IsValid != true {
		isForceUpdatePassword = true
	}

	// Get data address
	address, err := c.addressRepo.FindByCustomerId(customer.Id)
	if errors.Is(err, sql.ErrNoRows) {
		log.Error("failed to retrieve address not found", nlogger.Error(err))
		address = &model.Address{}
	} else if err != nil {
		log.Errorf("Error when retrieve address error: %v", err)
		return nil, c.response.GetError("E_AUTH_8")
	}

	// Get data verification
	verification, err := c.verificationRepo.FindByCustomerId(customer.Id)
	if errors.Is(err, sql.ErrNoRows) {
		log.Error("failed to retrieve verification not found", nlogger.Error(err))
		verification = &model.Verification{}
	} else if err != nil {
		log.Errorf("Error when retrieve verification error: %v", err)
		return nil, c.response.GetError("E_AUTH_8")
	}

	return c.composeLoginResponse(dto.LoginVO{
		Customer:              customer,
		Address:               address,
		Profile:               profile,
		Verification:          verification,
		IsFirstLogin:          isFirstLogin,
		IsForceUpdatePassword: isForceUpdatePassword,
		Token:                 token,
	})
}

func (c *Customer) SetTokenAuthentication(customer *model.Customer, agen string, version string, cacheTokenKey string) (string, error) {

	var accessToken string
	accessToken, _ = c.cacheService.Get(cacheTokenKey)
	if accessToken == "" {
		// Generate access token
		accessToken = nval.Bin2Hex(nval.RandStringBytes(78))
		// Set token to cache
		cacheToken, err := c.cacheService.SetThenGet(cacheTokenKey, accessToken, c.clientConfig.JWTExpired)
		if err != nil {
			return "", err
		}
		// Set access token
		accessToken = cacheToken
	}

	channelId := GetChannelByAgen(agen)
	now := time.Now()

	// Generate JWT
	token, err := jwt.NewBuilder().
		Claim("id", customer.UserRefId).
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
	return tokenString, nil
}

func (c *Customer) ValidateToken(token string) error {
	// Parsing Token
	t, err := jwt.ParseString(token, jwt.WithVerify(jwa.ES256, c.clientConfig.JWTKey))
	if err != nil {
		return err
	}

	if err = jwt.Validate(t); err != nil {
		return err
	}

	err = jwt.Validate(t, jwt.WithIssuer("https://www.pegadaian.co.id"))
	if err != nil {
		return err
	}

	return nil
}

func ValidatePassword(password string) *dto.ValidatePassword {
	var validation dto.ValidatePassword

	lowerCase, _ := regexp.Compile(`[a-z]+`)
	upperCase, _ := regexp.Compile(`[A-Z]+`)
	allNumber, _ := regexp.Compile(`[0-9]+`)

	if len(lowerCase.FindStringSubmatch(password)) < 1 {
		validation.IsValid = false
		validation.ErrCode = "isLower"
		validation.Message = "Password harus terdapat satu huruf kecil."
		return &validation
	}

	if len(upperCase.FindStringSubmatch(password)) < 1 {
		validation.IsValid = false
		validation.ErrCode = "isUpper"
		validation.Message = "Password harus terdapat satu huruf kapital.."
		return &validation
	}

	if len(allNumber.FindStringSubmatch(password)) < 1 {
		validation.IsValid = false
		validation.ErrCode = "isNumber"
		validation.Message = "Password harus terdapat angka."
		return &validation
	}

	if len(password) < 8 {
		validation.IsValid = false
		validation.ErrCode = "length"
		validation.Message = "Pasword harus terdapat minimal 8 karakter."
		return &validation
	}

	if strings.Contains(strings.ToLower(password), "gadai") {
		validation.IsValid = false
		validation.ErrCode = "containsGadai"
		validation.Message = "Hindari menggunakan kata gadai."
		return &validation
	}

	validation.IsValid = true
	validation.ErrCode = ""

	return &validation
}

func (c *Customer) HandleWrongPassword(credential *model.Credential, customer *model.Customer) error {
	var resp error
	t := time.Now()

	// Unmarshalling metadata credential to get tryLoginAt
	var Metadata dto.MetadataCredential
	err := json.Unmarshal(credential.Metadata, &Metadata)
	if err != nil {
		log.Errorf("Cannot unmarshaling metadata credential. err: %v", err)
		return ncore.TraceError("error", err)
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
	if credential.WrongPasswordCount == constant.MaxWrongPassword && now.After(blockedUntil) {
		credential.BlockedAt = sql.NullTime{}
		credential.BlockedUntilAt = sql.NullTime{}
		credential.WrongPasswordCount = 0
	}

	wrongCount := credential.WrongPasswordCount + 1

	switch wrongCount {
	case constant.Warn2XWrongPassword:
		resp = c.response.GetError("E_AUTH_6")
		credential.WrongPasswordCount = wrongCount
	case constant.Warn4XWrongPassword:
		resp = c.response.GetError("E_AUTH_7")
		credential.WrongPasswordCount = wrongCount
	case constant.MinWrongPassword:

		// Set block account
		hour := 1 // Block for 1 hours
		duration := time.Hour * time.Duration(hour)
		addDuration := t.Add(duration)
		credential.BlockedAt = sql.NullTime{
			Time:  t,
			Valid: true,
		}
		credential.BlockedUntilAt = sql.NullTime{
			Time:  addDuration,
			Valid: true,
		}
		credential.WrongPasswordCount = wrongCount

		// Set response if blocked for 1 hour
		message := "Akun dikunci hingga %v WIB karena gagal login %v kali. Hubungi call center jika ini bukan kamu"
		timeBlocked := ntime.ChangeTimezone(credential.BlockedUntilAt.Time, constant.WIB).Format("02-Jan-2006 15:04:05")
		setResponse := ncore.Success.HttpStatus(401).SetMessage(message, timeBlocked, credential.WrongPasswordCount)
		resp = &setResponse

		// Send OTP To Phone Number
		request := dto.SendOTPRequest{
			PhoneNumber: customer.Phone,
			RequestType: constant.RequestTypeBlockOneHour,
		}
		_, err := c.otpService.SendOTP(request)
		if err != nil {
			log.Debugf("Error when sending otp block one hour. err: %v", err)
		}

		// Send Notification Blocked Login One Hour
		err = c.notificationService.SendNotificationBlock(dto.NotificationBlock{
			Customer:     customer,
			Message:      fmt.Sprintf(message, timeBlocked, credential.WrongPasswordCount),
			LastTryLogin: ntime.NewTimeWIB(tryLoginAt).Format("02-Jan-2006 15:04:05"),
		})

		break
	case constant.MaxWrongPassword:
		// Set block account
		hour := 24 // Block for 24 hours
		duration := time.Hour * time.Duration(hour)
		addDuration := t.Add(duration)
		credential.BlockedAt = sql.NullTime{
			Time:  t,
			Valid: true,
		}
		credential.BlockedUntilAt = sql.NullTime{
			Time:  addDuration,
			Valid: true,
		}
		credential.WrongPasswordCount = wrongCount

		// Set response if blocked for 24 hour
		message := "Akun dikunci hingga %v WIB karena gagal login %v kali. Hubungi call center jika ini bukan kamu"
		timeBlocked := ntime.ChangeTimezone(credential.BlockedUntilAt.Time, constant.WIB).Format("02-Jan-2006 15:04:05")
		setResponse := ncore.Success.HttpStatus(401).SetMessage(message, timeBlocked, credential.WrongPasswordCount)
		resp = &setResponse

		// Send OTP To Phone Number
		request := dto.SendOTPRequest{
			PhoneNumber: customer.Phone,
			RequestType: constant.RequestTypeBlockOneDay,
		}
		_, err := c.otpService.SendOTP(request)
		if err != nil {
			log.Debugf("Error when sending otp block one hour. err: %v", err)
		}

		// Send Notification Blocked Login One Day
		err = c.notificationService.SendNotificationBlock(dto.NotificationBlock{
			Customer:     customer,
			Message:      fmt.Sprintf(message, timeBlocked, credential.WrongPasswordCount),
			LastTryLogin: ntime.NewTimeWIB(tryLoginAt).Format("02-Jan-2006 15:04:05"),
		})
		break
	default:
		resp = c.response.GetError("E_AUTH_8")
		credential.WrongPasswordCount = wrongCount
		break
	}

	// Handle notification error
	if err != nil {
		log.Debugf("Error when sending notification block: %v", err)
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
		return ncore.TraceError("error", err)
	}

	return resp
}

func (c *Customer) Register(payload dto.RegisterNewCustomer) (*dto.RegisterNewCustomerResponse, error) {
	// validate exist
	var customer *model.Customer
	customer, err := c.customerRepo.FindByEmail(payload.Email)
	if errors.Is(err, sql.ErrNoRows) {
		customer = nil
	} else if customer != nil {
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

		customerProfile := dto.CustomerProfileVO{
			MaidenName:         "",
			Gender:             "",
			Nationality:        "",
			DateOfBirth:        "",
			PlaceOfBirth:       "",
			IdentityPhotoFile:  "",
			MarriageStatus:     "",
			NPWPNumber:         "",
			NPWPPhotoFile:      "",
			NPWPUpdatedAt:      "",
			ProfileUpdatedAt:   "",
			CifLinkUpdatedAt:   "",
			CifUnlinkUpdatedAt: "",
			SidPhotoFile:       "",
		}
		profile, err := json.Marshal(customerProfile)
		if err != nil {
			return nil, ncore.TraceError("error", err)
		}

		insertCustomer := &model.Customer{
			CustomerXID:    customerXID,
			FullName:       payload.Name,
			Phone:          payload.PhoneNumber,
			Email:          payload.Email,
			Status:         0,
			IdentityType:   0,
			IdentityNumber: "",
			UserRefId:      "",
			Photos:         []byte("{}"),
			Profile:        profile,
			Cif:            "",
			Sid:            "",
			ReferralCode:   "",
			Metadata:       []byte("{}"),
			ItemMetadata:   metaData,
		}
		lastInsertId, err := c.customerRepo.Insert(insertCustomer)
		if err != nil {
			log.Errorf("Error when persist customer : %s", payload.Name)
			return nil, ncore.TraceError("error", err)
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

	// create verification
	verification := &model.Verification{
		Xid:                             strings.ToUpper(xid.New().String()),
		CustomerId:                      customerId,
		KycVerifiedStatus:               0,
		KycVerifiedAt:                   sql.NullTime{},
		EmailVerificationToken:          nval.RandomString(78),
		EmailVerifiedStatus:             0,
		EmailVerifiedAt:                 sql.NullTime{},
		DukcapilVerifiedStatus:          0,
		DukcapilVerifiedAt:              sql.NullTime{},
		FinancialTransactionStatus:      0,
		FinancialTransactionActivatedAt: sql.NullTime{},
		Metadata:                        []byte("{}"),
		ItemMetadata:                    model.NewItemMetadata(convert.ModifierDTOToModel(dto.Modifier{ID: "", Role: "", FullName: ""})),
	}
	err = c.verificationRepo.InsertOrUpdate(verification)
	if err != nil {
		log.Errorf("Error when persist customer verification : %s . Err : %v", payload.Name, err)
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
		Password:            nval.MD5(payload.Password),
		NextPasswordResetAt: sql.NullTime{},
		Pin:                 "",
		PinCif:              sql.NullString{},
		PinUpdatedAt:        sql.NullTime{},
		PinLastAccessAt:     sql.NullTime{},
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
		log.Errorf("Error when persist customer credential err: %v", err)
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
		return nil, ncore.TraceError("error", err)
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

	// call login service
	payloadLogin := dto.LoginRequest{
		Email:    payload.Email,
		Password: payload.Password,
		Agen:     payload.Agen,
		Version:  payload.Version,
		FcmToken: payload.FcmToken,
	}
	res, err := c.Login(payloadLogin)
	if err != nil {
		return nil, ncore.TraceError("error", err)
	}

	// Send Notification Register
	err = c.notificationService.SendNotificationRegister(dto.NotificationRegister{
		Customer:     customer,
		Verification: verification,
		RegisterOTP:  registerOTP,
		Payload:      payload,
	})
	if err != nil {
		log.Debugf("Error when send notification: %v", err)
	}

	// Delete OTP RegistrationId
	err = c.verificationOTPRepo.Delete(registerOTP.RegistrationId, customer.Phone)
	if err != nil {
		log.Debugf("Error when remove by registration id : %s, phone : %s", registerOTP.RegistrationId, customer.Phone)
		return nil, c.response.GetError("E_REG_1")
	}

	return &dto.RegisterNewCustomerResponse{
		LoginResponse: res,
		// TODO EKYC
		Ekyc: &dto.EKyc{
			AccountType: "",
			Status:      "",
			Screen:      "",
		},
		// TODO GPoint null response by default
		GPoint: nil,
		// TODO GCash
		GCash: &dto.GCash{
			TotalSaldo: 0,
			Va: []dto.GCashVa{{
				ID:             "",
				UserAIID:       "",
				Amount:         "",
				KodeBank:       "",
				TrxID:          "",
				TglExpired:     "",
				VirtualAccount: "",
				VaNumber:       "",
				CreatedAt:      "",
				LastUpdate:     "",
				NamaBank:       "",
				Thumbnail:      "",
			}},
			VaAvailable: []string{},
		},
	}, nil
}

func (c *Customer) RegisterStepOne(payload dto.RegisterStepOne) (*dto.RegisterStepOneResponse, error) {

	// validate email
	emailExist, err := c.customerRepo.FindByEmail(payload.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			emailExist = nil
		} else {
			log.Error("failed when query check email.", nlogger.Error(err))
			return nil, ncore.TraceError("error", err)
		}
	}
	if emailExist != nil {
		log.Debugf("Email already registered")
		return nil, c.response.GetError("E_REG_2")
	}

	// validate phone
	phoneExist, err := c.customerRepo.FindByPhone(payload.PhoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			phoneExist = nil
		} else {
			log.Error("failed when query check phone.", nlogger.Error(err))
			return nil, ncore.TraceError("error", err)
		}
	}
	if phoneExist != nil {
		log.Debugf("Phone already registered")
		return nil, c.response.GetError("E_REG_3")
	}

	// Set request
	request := dto.SendOTPRequest{
		PhoneNumber: payload.PhoneNumber,
		RequestType: constant.RequestTypeRegister,
	}

	// Send OTP To Phone Number
	resp, err := c.otpService.SendOTP(request)
	if err != nil {
		return nil, ncore.TraceError("error", err)
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
		if errors.Is(err, sql.ErrNoRows) {
			phoneExist = nil
		} else {
			log.Error("failed when query check phone.", nlogger.Error(err))
			return nil, ncore.TraceError("", err)
		}
	}
	if phoneExist != nil {
		log.Debugf("Phone already registered")
		return nil, c.response.GetError("E_REG_3")
	}

	// Verify OTP To Phone Number
	resp, err := c.otpService.VerifyOTP(request)
	if err != nil {
		return nil, ncore.TraceError("error", err)
	}

	// Extract response from server
	data, err := nclient.GetResponseData(resp)

	// handle Expired OTP
	if data.ResponseCode == "12" {
		log.Errorf("Expired OTP. Phone Number : %s", payload.PhoneNumber)
		return nil, c.response.GetError("E_OTP_4")
	}
	// handle Wrong OTP
	if data.ResponseCode == "14" {
		log.Errorf("Wrong OTP. Phone Number : %s", payload.PhoneNumber)
		return nil, c.response.GetError("E_OTP_1")
	}

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
		return nil, ncore.TraceError("error", err)
	}

	return &dto.RegisterStepTwoResponse{
		RegisterId: registrationId,
	}, nil
}

func (c *Customer) RegisterResendOTP(payload dto.RegisterResendOTP) (*dto.RegisterResendOTPResponse, error) {
	// Set request
	request := dto.SendOTPRequest{
		PhoneNumber: payload.PhoneNumber,
		RequestType: constant.RequestTypeRegister,
	}

	// validate phone
	phoneExist, err := c.customerRepo.FindByPhone(payload.PhoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			phoneExist = nil
		} else {
			log.Error("failed when query check phone.", nlogger.Error(err))
			return nil, ncore.TraceError("error", err)
		}
	}
	if phoneExist != nil {
		log.Debugf("Phone already registered")
		return nil, c.response.GetError("E_REG_3")
	}

	// Send OTP To Phone Number
	resp, err := c.otpService.SendOTP(request)
	if err != nil {
		return nil, ncore.TraceError("error", err)
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

func GetChannelByAgen(agen string) string {

	// Generalize agen
	agen = strings.ToLower(agen)

	if agen == constant.AgenAndroid {
		return constant.ChannelAndroid
	}

	if agen == constant.AgenMobile {
		return constant.ChannelMobile
	}

	if agen == constant.AgenWeb {
		return constant.ChannelWeb
	}

	return ""
}

func (c *Customer) composeLoginResponse(data dto.LoginVO) (*dto.LoginResponse, error) {

	return &dto.LoginResponse{
		Customer: &dto.CustomerVO{
			ID:                        nval.ParseStringFallback(data.Customer.Id, ""),
			Cif:                       data.Customer.Cif,
			IsKYC:                     nval.ParseStringFallback(data.Verification.KycVerifiedStatus, "0"),
			Nama:                      data.Customer.FullName,
			NamaIbu:                   data.Profile.MaidenName,
			NoKTP:                     data.Customer.IdentityNumber,
			Email:                     data.Customer.Email,
			JenisKelamin:              data.Profile.Gender,
			TempatLahir:               data.Profile.PlaceOfBirth,
			TglLahir:                  data.Profile.DateOfBirth,
			Alamat:                    data.Address.Line.String,
			IDProvinsi:                data.Address.ProvinceId.String,
			IDKabupaten:               data.Address.CityId.String,
			IDKecamatan:               data.Address.DistrictId.String,
			IDKelurahan:               data.Address.SubDistrictId.String,
			Kelurahan:                 data.Address.DistrictName.String,
			Provinsi:                  data.Address.ProvinceName.String,
			Kabupaten:                 data.Address.CityName.String,
			Kecamatan:                 data.Address.DistrictName.String,
			KodePos:                   data.Address.PostalCode.String,
			NoHP:                      data.Customer.Phone,
			Avatar:                    "",
			FotoKTP:                   data.Profile.IdentityPhotoFile,
			IsEmailVerified:           nval.ParseStringFallback(data.Verification.EmailVerifiedStatus, "0"),
			Kewarganegaraan:           data.Profile.Nationality,
			JenisIdentitas:            fmt.Sprintf("%v", data.Customer.IdentityType),
			NoIdentitas:               data.Customer.IdentityNumber,
			TglExpiredIdentitas:       "",
			NoNPWP:                    data.Profile.NPWPNumber,
			FotoNPWP:                  data.Profile.NPWPPhotoFile,
			NoSid:                     data.Customer.Sid,
			FotoSid:                   data.Profile.SidPhotoFile,
			StatusKawin:               data.Profile.MarriageStatus,
			Norek:                     "",
			Saldo:                     "",
			AktifasiTransFinansial:    nval.ParseStringFallback(data.Verification.FinancialTransactionStatus, ""),
			IsDukcapilVerified:        nval.ParseStringFallback(data.Verification.DukcapilVerifiedStatus, "0"),
			IsOpenTe:                  "",
			ReferralCode:              "",
			GoldCardApplicationNumber: "",
			GoldCardAccountNumber:     "",
			KodeCabang:                "",
			TabunganEmas: &dto.CustomerTabunganEmasVO{
				TotalSaldoBlokir:  "",
				TotalSaldoSeluruh: "",
				TotalSaldoEfektif: "",
				PrimaryRekening:   "",
			},
			IsFirstLogin:          data.IsFirstLogin,
			IsForceUpdatePassword: data.IsForceUpdatePassword,
		},
		JwtToken: data.Token,
	}, nil
}

func (c *Customer) syncExternalToInternal(user *model.User) (*model.Customer, error) {
	// prepare customer
	customer, err := convert.ModelUserToCustomer(user)
	if err != nil {
		log.Error("failed to convert to model customer", nlogger.Error(err))
		return nil, err
	}

	// Check has userPin or not
	userPin, err := c.userPinExternalRepo.FindByCustomerId(customer.Id)
	if errors.Is(err, sql.ErrNoRows) {
		userPin = &model.UserPin{}
	} else if err != nil {
		log.Error("failed retrieve user pin from external database", nlogger.Error(err))
		return nil, err
	}

	// Check user address on external database
	addressExternal, err := c.userExternalRepo.FindAddressByCustomerId(user.UserAiid)
	if errors.Is(err, sql.ErrNoRows) {
		addressExternal = &model.AddressExternal{}
	} else if err != nil {
		log.Error("failed retrieve address from external database", nlogger.Error(err))
		return nil, err
	}

	// Prepare credential
	credential, err := convert.ModelUserToCredential(user, userPin)
	if err != nil {
		log.Error("failed convert to credential model", nlogger.Error(err))
		return nil, err
	}

	// Prepare financial data
	financialData, err := convert.ModelUserToFinancialData(user)
	if err != nil {
		log.Error("failed convert to financial data", nlogger.Error(err))
		return nil, err
	}

	// Prepare verification
	verification, err := convert.ModelUserToVerification(user)
	if err != nil {
		log.Error("failed convert to verification", nlogger.Error(err))
		return nil, err
	}

	// Prepare address
	address, err := convert.ModelUserToAddress(user, addressExternal)
	if err != nil {
		log.Error("failed convert address", nlogger.Error(err))
		return nil, err
	}

	// Persist customer data
	customerId, err := c.customerRepo.Insert(customer)
	if err != nil {
		return nil, err
	}

	// Set credential customer id to last inserted
	credential.CustomerId = customerId
	financialData.CustomerId = customerId
	verification.CustomerId = customerId
	address.CustomerId = customerId

	// persist credential
	err = c.credentialRepo.InsertOrUpdate(credential)
	if err != nil {
		log.Error("failed persist to credential", nlogger.Error(err))
		return nil, err
	}

	// persist financial data
	err = c.financialDataRepo.InsertOrUpdate(financialData)
	if err != nil {
		log.Error("failed persist to financial data", nlogger.Error(err))
		return nil, err
	}

	// persist verification
	err = c.verificationRepo.InsertOrUpdate(verification)
	if err != nil {
		log.Errorf("failed persist verification err: %v", nlogger.Error(err))
		return nil, err
	}

	// persist address
	err = c.addressRepo.InsertOrUpdate(address)
	if err != nil {
		log.Error("failed persist address", nlogger.Error(err))
		return nil, err
	}

	customer, err = c.customerRepo.FindById(customerId)
	if err != nil {
		log.Error("failed to retrieve customer not found", nlogger.Error(err))
		return nil, c.response.GetError("E_RES_1")
	}

	return customer, nil
}

func (c *Customer) syncInternalToExternal(payload *dto.CustomerSynchronizeRequest) (*dto.UserVO, error) {

	// call register pds api
	registerCustomer := dto.RegisterNewCustomer{
		Name:        payload.Name,
		Email:       payload.Email,
		PhoneNumber: payload.PhoneNumber,
		Password:    payload.Password,
		FcmToken:    payload.FcmToken,
	}
	// sync
	sync, err := c.pdsAPIService.SynchronizeCustomer(registerCustomer)
	if err != nil {
		log.Error("Error when SynchronizeCustomer.", nlogger.Error(err))
		return nil, ncore.TraceError("error", err)
	}

	// set response data
	resp, err := nclient.GetResponseDataPdsAPI(sync)
	if err != nil {
		log.Error("Cannot parsing from SynchronizeCustomer response.", nlogger.Error(err))
		return nil, ncore.TraceError("error", err)
	}

	// handle status error
	if resp.Status != "success" {
		log.Error("Get Error from SynchronizeCustomer.", nlogger.Error(err))
		return nil, ncore.NewError(resp.Message)
	}

	// parsing response
	var user dto.UserVO
	err = json.Unmarshal(resp.Data, &user)
	if err != nil {
		log.Errorf("Cannot unmarshall data login pds. err: %v", err)
		return nil, ncore.TraceError("error", err)
	}

	// set result
	result := &user

	return result, nil
}
