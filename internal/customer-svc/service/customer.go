package service

import (
	"bytes"
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"math/rand"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/rs/xid"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/convert"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nclient"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ntime"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
)

type Customer struct {
	customerRepo        contract.CustomerRepository
	verificationOTPRepo contract.VerificationOTPRepository
	OTPRepo             contract.OTPRepository
	credentialRepo      contract.CredentialRepository
	accessSessionRepo   contract.AccessSessionRepository
	auditLoginRepo      contract.AuditLoginRepository
	verificationRepo    contract.VerificationRepository
	userExternalRepo    contract.UserExternalRepository
	otpService          contract.OTPService
	cacheService        contract.CacheService
	notificationService contract.NotificationService
	clientConfig        contract.ClientConfig
	response            *ncore.ResponseMap
	httpBaseUrl         string
	emailConfig         contract.EmailConfig
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
	c.userExternalRepo = app.Repositories.UserExternal
	c.otpService = app.Services.OTP
	c.cacheService = app.Services.Cache
	c.clientConfig = app.Config.Client
	c.notificationService = app.Services.Notification
	c.response = app.Responses
	c.httpBaseUrl = app.Config.Server.GetHttpBaseUrl()
	c.emailConfig = app.Config.Email
	return nil
}

func (c *Customer) Login(payload dto.LoginRequest) (*dto.LoginResponse, error) {

	// Check if user exists
	t := time.Now()
	customer, err := c.customerRepo.FindByEmailOrPhone(payload.Email)
	if errors.Is(err, sql.ErrNoRows) {
		// If data not found on internal database check on external database.
		_, err := c.userExternalRepo.FindByEmailOrPhone(payload.Email)
		if err != nil {
			log.Error("failed to retrieve customer not found", nlogger.Error(err))
			return nil, c.response.GetError("E_RES_1")
		}

		// TODO: MAPPING FROM DB EXTERNAL AND REGISTER TO DB INTERNAL
		return nil, ncore.TraceError(err) // TODO: REMOVE THIS
	} else if err != nil {
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
		// Set response if blocked
		message := "Akun dikunci hingga %v karena gagal login %v kali. Hubungi call center jika ini bukan kamu"
		timeBlocked := ntime.ChangeTimezone(credential.BlockedUntilAt.Time, constant.WIB).Format("02-Jan-2006 15:04:05")
		setResponse := ncore.Success.SetMessage(message, timeBlocked, credential.WrongPasswordCount)
		return nil, &setResponse
	}

	// Counter wrong password count
	passwordRequest := fmt.Sprintf("%x", md5.Sum([]byte(payload.Password)))
	if credential.Password != passwordRequest {
		err := c.HandleWrongPassword(credential)
		return nil, err
	}

	// Get token from cache
	var token string
	cacheTokenKey := fmt.Sprintf("%v:%v:%v", constant.Prefix, "token_jwt", customer.Id)
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

	// get user data

	// Unmarshall profile
	var Profile dto.CustomerProfileVO
	err = json.Unmarshal(customer.Profile, &Profile)
	if err != nil {
		log.Errorf("Error when unmarshalling profile: %v", err)
		return nil, ncore.TraceError(err)
	}

	// TODO get tabungan emas service / get tabungan emas from financial table

	// Check is force update password
	validatePassword := c.ValidatePassword(payload.Password)
	isForcePassword := false
	log.Debugf("errCode : %v", validatePassword.ErrCode)
	if validatePassword.IsValid != true {
		isForcePassword = true
	}

	return &dto.LoginResponse{
		Customer: &dto.CustomerVO{
			ID:                        nval.ParseStringFallback(customer.Id, ""),
			Cif:                       customer.Cif,
			IsKYC:                     "1",
			Nama:                      customer.FullName,
			NamaIbu:                   Profile.MaidenName,
			NoKTP:                     customer.IdentityNumber,
			Email:                     customer.Email,
			JenisKelamin:              Profile.Gender,
			TempatLahir:               Profile.PlaceOfBirth,
			TglLahir:                  Profile.DateOfBirth,
			Alamat:                    "",
			IDProvinsi:                "",
			IDKabupaten:               "",
			IDKecamatan:               "",
			IDKelurahan:               "",
			Kelurahan:                 "",
			Provinsi:                  "",
			Kabupaten:                 "",
			Kecamatan:                 "",
			KodePos:                   "",
			NoHP:                      customer.Phone,
			Avatar:                    "",
			FotoKTP:                   Profile.IdentityPhotoFile,
			IsEmailVerified:           customer.Email,
			Kewarganegaraan:           Profile.Nationality,
			JenisIdentitas:            fmt.Sprintf("%v", customer.IdentityType),
			NoIdentitas:               customer.IdentityNumber,
			TglExpiredIdentitas:       "",
			NoNPWP:                    Profile.NPWPNumber,
			NoSid:                     customer.Sid,
			FotoSid:                   Profile.SidPhotoFile,
			StatusKawin:               Profile.MarriageStatus,
			Norek:                     "",
			Saldo:                     "",
			AktifasiTransFinansial:    "",
			IsDukcapilVerified:        "",
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
			IsFirstLogin:          isFirstLogin,
			IsForceUpdatePassword: isForcePassword,
		},
		JwtToken: token,
	}, nil
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
	return tokenString, nil
}

func (c *Customer) ValidatePassword(password string) *dto.ValidatePassword {
	var validation dto.ValidatePassword

	lowerCase, _ := regexp.Compile(`[a-z]+`)
	upperCase, _ := regexp.Compile(`[A-Z]+`)
	allNumber, _ := regexp.Compile(`[0-9]+`)

	if len(lowerCase.FindStringSubmatch(password)) < 1 {
		validation.IsValid = false
		validation.ErrCode = "isLower"
		return &validation
	}

	if len(upperCase.FindStringSubmatch(password)) < 1 {
		validation.IsValid = false
		validation.ErrCode = "isUpper"
		return &validation
	}

	if len(allNumber.FindStringSubmatch(password)) < 1 {
		validation.IsValid = false
		validation.ErrCode = "isNumber"
		return &validation
	}

	if len(password) < 8 {
		validation.IsValid = false
		validation.ErrCode = "length"
		return &validation
	}

	if strings.Contains(password, "gadai") {
		validation.IsValid = false
		validation.ErrCode = "containsGadai"
		return &validation
	}

	validation.IsValid = true
	validation.ErrCode = ""

	return &validation
}

func (c *Customer) HandleWrongPassword(credential *model.Credential) error {
	var resp error
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
		message := "Akun dikunci hingga %v karena gagal login %v kali. Hubungi call center jika ini bukan kamu"
		timeBlocked := ntime.ChangeTimezone(credential.BlockedUntilAt.Time, constant.WIB).Format("02-Jan-2006 15:04:05")
		setResponse := ncore.Success.SetMessage(message, timeBlocked, credential.WrongPasswordCount)
		resp = &setResponse

		// TODO sendNotificationBlockedLoginOneHour
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
		message := "Akun dikunci hingga %v karena gagal login %v kali. Hubungi call center jika ini bukan kamu"
		timeBlocked := ntime.ChangeTimezone(credential.BlockedUntilAt.Time, constant.WIB).Format("02-Jan-2006 15:04:05")
		setResponse := ncore.Success.SetMessage(message, timeBlocked, credential.WrongPasswordCount)
		resp = &setResponse

		// TODO sendNotificationBlockedLoginOneDay
		break
	default:
		resp = c.response.GetError("E_AUTH_8")
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

	return resp
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
			return nil, ncore.TraceError(err)
		}

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
		PinCif:              "",
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

	// call login service
	payloadLogin := dto.LoginRequest{
		Email:    payload.Email,
		Password: payload.Password,
		Agen:     payload.Agen,
		Version:  payload.Version,
	}
	res, err := c.Login(payloadLogin)
	if err != nil {
		return nil, ncore.TraceError(err)
	}
	user := res.Customer
	// Send Email Verification
	dataEmailVerification := &dto.EmailVerification{
		FullName:        customer.FullName,
		Email:           customer.Email,
		VerificationUrl: fmt.Sprintf("%sauth/verify_email?t=%s", c.httpBaseUrl, verification.EmailVerificationToken),
	}
	htmlMessage, err := templateFile(dataEmailVerification, "email_verification.html")
	if err != nil {
		return nil, err
	}

	// set payload email service
	emailPayload := dto.EmailPayload{
		Subject: fmt.Sprintf("Verifikasi Email %s", customer.FullName),
		From: dto.FromEmailPayload{
			Name:  c.emailConfig.PdsEmailFromName,
			Email: c.emailConfig.PdsEmailFrom,
		},
		To:         customer.Email,
		Message:    htmlMessage,
		Attachment: "",
		MimeType:   "",
	}
	_, err = c.notificationService.SendEmail(emailPayload)
	if err != nil {
		log.Debugf("Error when send email verification. Payload %v", emailPayload)
	}

	// Send Notification Welcome
	id, _ := nval.ParseString(rand.Intn(100)) // TODO: insert data to notification
	var dataWelcomeMessage = map[string]string{
		"title": "Verifikasi Email",
		"body":  fmt.Sprintf(`Hai %v, Selamat datang di Pegadaian Digital Service`, user.Nama),
		"type":  constant.TypeProfile,
		"id":    id,
	}
	welcomeMessage := dto.NotificationPayload{
		Title: "Verifikasi Email",
		Body:  fmt.Sprintf(`Hai %v, Selamat datang di Pegadaian Digital Service`, user.Nama),
		Image: "",
		Token: payload.FcmToken,
		Data:  dataWelcomeMessage,
	}
	_, err = c.notificationService.SendNotification(welcomeMessage)
	if err != nil {
		log.Debugf("Error when send notification message: %s, phone : %s", registerOTP.RegistrationId, customer.Phone)
	}

	// Delete OTP RegistrationId
	err = c.verificationOTPRepo.Delete(registerOTP.RegistrationId, customer.Phone)
	if err != nil {
		log.Debugf("Error when remove by registration id : %s, phone : %s", registerOTP.RegistrationId, customer.Phone)
		return nil, c.response.GetError("E_REG_1")
	}

	return &dto.RegisterNewCustomerResponse{
		User: dto.CustomerVO{
			ID:                        user.ID,
			Cif:                       user.Cif,
			IsKYC:                     user.IsKYC,
			Nama:                      user.Nama,
			NamaIbu:                   user.NamaIbu,
			NoKTP:                     user.NoKTP,
			Email:                     user.Email,
			JenisKelamin:              user.JenisKelamin,
			TempatLahir:               user.TempatLahir,
			TglLahir:                  user.TglLahir,
			Alamat:                    user.Alamat,
			IDProvinsi:                user.IDProvinsi,
			IDKabupaten:               user.IDKabupaten,
			IDKecamatan:               user.IDKecamatan,
			IDKelurahan:               user.IDKelurahan,
			Kelurahan:                 user.Kelurahan,
			Provinsi:                  user.Provinsi,
			Kabupaten:                 user.Kabupaten,
			Kecamatan:                 user.Kecamatan,
			KodePos:                   user.KodePos,
			NoHP:                      user.NoHP,
			Avatar:                    user.Avatar,
			FotoKTP:                   user.FotoKTP,
			IsEmailVerified:           user.IsEmailVerified,
			Kewarganegaraan:           user.Kewarganegaraan,
			JenisIdentitas:            user.JenisIdentitas,
			NoIdentitas:               user.NoIdentitas,
			TglExpiredIdentitas:       user.TglExpiredIdentitas,
			NoNPWP:                    user.NoNPWP,
			FotoNPWP:                  user.FotoNPWP,
			NoSid:                     user.NoSid,
			FotoSid:                   user.FotoSid,
			StatusKawin:               user.StatusKawin,
			Norek:                     user.Norek,
			Saldo:                     user.Saldo,
			AktifasiTransFinansial:    user.AktifasiTransFinansial,
			IsDukcapilVerified:        user.IsDukcapilVerified,
			IsOpenTe:                  user.IsOpenTe,
			ReferralCode:              user.ReferralCode,
			GoldCardApplicationNumber: user.GoldCardApplicationNumber,
			GoldCardAccountNumber:     user.GoldCardAccountNumber,
			KodeCabang:                user.KodeCabang,
			IsFirstLogin:              user.IsFirstLogin,
			IsForceUpdatePassword:     user.IsForceUpdatePassword,
			// TODO Load Tabungan
			TabunganEmas: &dto.CustomerTabunganEmasVO{
				TotalSaldoBlokir:  user.TabunganEmas.TotalSaldoBlokir,
				TotalSaldoSeluruh: user.TabunganEmas.TotalSaldoSeluruh,
				TotalSaldoEfektif: user.TabunganEmas.TotalSaldoEfektif,
				PrimaryRekening:   user.TabunganEmas.PrimaryRekening,
			},
		},
		JwtToken: res.JwtToken,
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

func templateFile(data interface{}, htmlFile string) (string, error) {
	var filepath = path.Join("web/templates", htmlFile)
	var tmpl, err = template.ParseFiles(filepath)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)

	if err := tmpl.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
