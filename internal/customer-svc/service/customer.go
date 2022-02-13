package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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
