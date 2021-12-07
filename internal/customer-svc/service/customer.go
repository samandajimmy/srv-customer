package service

import (
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nlogger"
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
	CredentialRepo      contract.CredentialRepository
	AccessSessionRepo   contract.AccessSessionRepository
	auditLoginRepo      contract.AuditLoginRepository
	otpService          contract.OTPService
	response            *ncore.ResponseMap
}

func (c *Customer) HasInitialized() bool {
	return true
}

func (c *Customer) Init(app *contract.PdsApp) error {
	c.customerRepo = app.Repositories.Customer
	c.verificationOTPRepo = app.Repositories.VerificationOTP
	c.OTPRepo = app.Repositories.OTP
	c.CredentialRepo = app.Repositories.Credential
	c.AccessSessionRepo = app.Repositories.AccessSession
	c.otpService = app.Services.OTP
	c.auditLoginRepo = app.Repositories.AuditLogin
	c.response = app.Responses
	return nil
}

func (c *Customer) Login(payload dto.LoginRequest) (*dto.CustomerVO, error) {

	// Check if user exists
	customer, err := c.customerRepo.FindByEmailOrPhone(payload.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("failed to retrieve customer not found", nlogger.Error(err))
			return nil, c.response.GetError("E_RES_1")
		}
		log.Errorf("failed to retrieve customer. error: %v", err)
		return nil, ncore.TraceError(err)
	}

	// Get Auth

	// counter wrong password count
	//customer.WrongPasswordCount += 1
	//
	//if customer.WrongPasswordCount == 2 {
	//	return nil, c.response.GetError("E_AUTH_6")
	//} else if customer.WrongPasswordCount == 4 {
	//	return nil, c.response.GetError("E_AUTH_7")
	//}

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
	t := time.Now()
	auditLog := model.AuditLogin{
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

	// Persist audit loging
	err = c.auditLoginRepo.Insert(&auditLog)
	if err != nil {
		log.Errorf("Error when insert audit login error: %v", err)
		return nil, c.response.GetError("E_AUTH_1")
	}

	// TODO: update user_model -> try_login_date = now()

	// check user account is blocked or not
	// if password doesn't match
	// 	cek setBlockedUser function
	//    if blocked_to_date > now()
	//    	return err_account_locked message

	//if customer.Password != payload.Password {
	//	//
	//}

	// if password is matched
	//    update user_model
	//   		set blocked_date = null
	//        set blocked_to_date = null
	//        wrong_password_count = 0

	// set token authentication

	// get user data

	// get tabungan emas service

	// check is force update password

	// return response user and token

	return &dto.CustomerVO{
		ID:                        "1",
		Cif:                       "Cif",
		IsKYC:                     "1",
		Nama:                      customer.FullName,
		NamaIbu:                   "Nama ibu",
		NoKTP:                     "No ktp",
		Email:                     "Email",
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
	}, nil
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

	// create credential
	credentialInsert := &model.Credential{
		Xid:                 customerXID,
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
		BlockedAt:           nil,
		BlockedUntilAt:      nil,
		BiometricLogin:      0,
		BiometricDeviceId:   "",
		Metadata:            []byte("{}"),
		ItemMetadata:        model.NewItemMetadata(convert.ModifierDTOToModel(dto.Modifier{ID: "", Role: "", FullName: ""})),
	}
	err = c.CredentialRepo.InsertOrUpdate(credentialInsert)
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
	err = c.AccessSessionRepo.Insert(insertAccessSession)
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
