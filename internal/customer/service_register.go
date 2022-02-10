package customer

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/nbs-go/nlogger"
	"github.com/rs/xid"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/convert"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nclient"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
	"strings"
	"time"
)

func (s *Service) Register(payload dto.RegisterNewCustomer) (*dto.RegisterNewCustomerResponse, error) {

	registrationId := payload.RegistrationId
	phoneNumber := payload.PhoneNumber

	// validate exist
	var customer *model.Customer
	customer, err := s.repo.FindCustomerByEmail(payload.Email)
	if errors.Is(err, sql.ErrNoRows) {
		customer = nil
	} else if customer != nil {
		log.Debugf("Email already registered: %s", registrationId)
		return nil, s.responses.GetError("E_REG_2")
	}
	// find registerID
	registerOTP, err := s.repo.FindVerificationOTPByRegistrationIDAndPhone(registrationId, phoneNumber)
	if err != nil {
		log.Errorf("Registration ID not found: %s", registrationId)
		return nil, s.responses.GetError("E_REG_1")
	}

	// Get data user
	customer, _ = s.repo.FindCustomerByPhone(phoneNumber)

	var customerXID string
	var customerId int64
	// update name if customer exists
	if customer != nil {
		customer.FullName = payload.Name
		customer.Email = payload.Email
		err := s.repo.UpdateCustomerByPhone(customer)
		if err != nil {
			return nil, s.responses.GetError("E_AUTH_5")
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
		lastInsertId, err := s.repo.CreateCustomer(insertCustomer)
		if err != nil {
			log.Errorf("Error when persist customer : %s", payload.Name)
			return nil, ncore.TraceError("error", err)
		}
		customerId = lastInsertId
		customer = insertCustomer
		customerXID = customer.CustomerXID
	}

	customer.Status = 1
	err = s.repo.UpdateCustomerByPhone(customer)
	if err != nil {
		log.Errorf("Error when update customer : %s", payload.Name)
		return nil, s.responses.GetError("E_REG_1")
	}

	// prepare verification model
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
	// update verification
	err = s.repo.InsertOrUpdateVerification(verification)
	if err != nil {
		log.Errorf("Error when persist customer verification : %s . Err : %v", payload.Name, err)
		return nil, s.responses.GetError("E_REG_1")
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
	err = s.repo.InsertOrUpdateCredential(credentialInsert)
	if err != nil {
		log.Errorf("Error when persist customer credential err: %v", err)
		return nil, s.responses.GetError("E_REG_1")
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
	err = s.repo.CreateOTP(insertOTP)
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
	err = s.repo.CreateAccessSession(insertAccessSession)
	if err != nil {
		log.Errorf("Error when persist access session: %s. Err: %s", payload.Name, err)
		return nil, s.responses.GetError("E_REG_1")
	}

	// call login service
	payloadLogin := dto.LoginRequest{
		Email:    payload.Email,
		Password: payload.Password,
		Agen:     payload.Agen,
		Version:  payload.Version,
		FcmToken: payload.FcmToken,
	}
	res, err := s.Login(payloadLogin)
	if err != nil {
		return nil, ncore.TraceError("error", err)
	}

	// Send Notification Register
	err = s.SendNotificationRegister(dto.NotificationRegister{
		Customer:     customer,
		Verification: verification,
		RegisterOTP:  registerOTP,
		Payload:      payload,
	})
	if err != nil {
		s.log.Debugf("Error when send notification: %v", err)
	}

	// Delete OTP RegistrationId
	err = s.repo.DeleteVerificationOTP(registerOTP.RegistrationId, customer.Phone)
	if err != nil {
		log.Debugf("Error when remove by registration id : %s, phone : %s", registerOTP.RegistrationId, customer.Phone)
		return nil, s.responses.GetError("E_REG_1")
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

func (s *Service) RegisterStepOne(payload dto.RegisterStepOne) (*dto.RegisterStepOneResponse, error) {

	// validate email
	emailExist, err := s.repo.FindCustomerByEmail(payload.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			emailExist = nil
		} else {
			s.log.Error("failed when query check email.", nlogger.Error(err))
			return nil, ncore.TraceError("error", err)
		}
	}
	if emailExist != nil {
		s.log.Debugf("Email already registered")
		return nil, s.responses.GetError("E_REG_2")
	}

	// validate phone
	phoneExist, err := s.repo.FindCustomerByPhone(payload.PhoneNumber)
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
		return nil, s.responses.GetError("E_REG_3")
	}

	// Set request
	request := dto.SendOTPRequest{
		PhoneNumber: payload.PhoneNumber,
		RequestType: constant.RequestTypeRegister,
	}

	// Send OTP To Phone Number
	resp, err := s.SendOTP(request)
	if err != nil {
		return nil, ncore.TraceError("error", err)
	}

	// Extract response from server
	data, err := nclient.GetResponseData(resp)

	if data.ResponseCode == "15" {
		return nil, s.responses.GetError("E_OTP_3")
	}

	if data.Message != "" {
		log.Errorf("Debug: RegisterStepOne OTP CODE %s", data.Message)
	}

	return &dto.RegisterStepOneResponse{
		Action: data.ResponseDesc,
	}, nil
}

func (s *Service) RegisterStepTwo(payload dto.RegisterStepTwo) (*dto.RegisterStepTwoResponse, error) {
	// Set request
	request := dto.VerifyOTPRequest{
		PhoneNumber: payload.PhoneNumber,
		Token:       payload.OTP,
		RequestType: "register",
	}
	// validate phone
	phoneExist, err := s.repo.FindCustomerByPhone(payload.PhoneNumber)
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
		return nil, s.responses.GetError("E_REG_3")
	}

	// Verify OTP To Phone Number
	resp, err := s.VerifyOTP(request)
	if err != nil {
		return nil, ncore.TraceError("error", err)
	}

	// Extract response from server
	data, err := nclient.GetResponseData(resp)

	// handle Expired OTP
	if data.ResponseCode == "12" {
		log.Errorf("Expired OTP. Phone Number : %s", payload.PhoneNumber)
		return nil, s.responses.GetError("E_OTP_4")
	}
	// handle Wrong OTP
	if data.ResponseCode == "14" {
		log.Errorf("Wrong OTP. Phone Number : %s", payload.PhoneNumber)
		return nil, s.responses.GetError("E_OTP_1")
	}

	if data.ResponseCode != "00" {
		log.Errorf("Wrong OTP. Phone Number : %s", payload.PhoneNumber)
		return nil, s.responses.GetError("E_OTP_1")
	}

	// insert verification otp
	registrationId := xid.New().String()
	insert := &model.VerificationOTP{
		CreatedAt:      time.Now(),
		Phone:          payload.PhoneNumber,
		RegistrationId: registrationId,
	}
	_, err = s.repo.CreateVerificationOTP(insert)
	if err != nil {
		log.Errorf("Error when persist verificationOTP. Phone Number: %s", payload.PhoneNumber)
		return nil, ncore.TraceError("error", err)
	}

	return &dto.RegisterStepTwoResponse{
		RegisterId: registrationId,
	}, nil
}

func (s *Service) RegisterResendOTP(payload dto.RegisterResendOTP) (*dto.RegisterResendOTPResponse, error) {
	// Set request
	request := dto.SendOTPRequest{
		PhoneNumber: payload.PhoneNumber,
		RequestType: constant.RequestTypeRegister,
	}

	// validate phone
	phoneExist, err := s.repo.FindCustomerByPhone(payload.PhoneNumber)
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
		return nil, s.responses.GetError("E_REG_3")
	}

	// Send OTP To Phone Number
	resp, err := s.SendOTP(request)
	if err != nil {
		return nil, ncore.TraceError("error", err)
	}

	// Extract response from server
	data, err := nclient.GetResponseData(resp)

	if data.ResponseCode == "15" {
		return nil, s.responses.GetError("E_OTP_3")
	}

	if data.Message != "" {
		log.Errorf("Debug: RegisterResendOTP OTP CODE %s", data.Message)
	}

	return &dto.RegisterResendOTPResponse{
		Action: data.ResponseDesc,
	}, nil
}

func TokenIsExists(headers map[interface{}]interface{}) bool {
	if headers == nil {
		return false
	}

	return nval.KeyExists("Authorization", headers)
}
