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
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
	"strings"
	"time"
)

func (s *Service) Register(payload dto.RegisterNewCustomer) (*dto.RegisterNewCustomerResponse, error) {
	// Get context
	ctx := s.ctx

	registrationId := payload.RegistrationId
	phoneNumber := payload.PhoneNumber

	// validate exist
	var customer *model.Customer
	customer, err := s.repo.FindCustomerByEmail(payload.Email)
	if errors.Is(err, sql.ErrNoRows) {
		customer = nil
	} else if customer != nil {
		s.log.Debugf("Email already registered: %s", registrationId, nlogger.Error(err), nlogger.Context(ctx))
		return nil, s.responses.GetError("E_REG_2")
	}
	// find registerID
	registerOTP, err := s.repo.FindVerificationOTPByRegistrationIDAndPhone(registrationId, phoneNumber)
	if err != nil {
		s.log.Errorf("Registration ID not found: %s", registrationId, nlogger.Error(err), nlogger.Context(ctx))
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
			s.log.Error("error when update customer by phone", nlogger.Error(err), nlogger.Context(ctx))
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
			s.log.Error("error when marshal profile", nlogger.Error(err), nlogger.Context(ctx))
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
			s.log.Errorf("error when persist customer: %s", payload.Name, nlogger.Error(err), nlogger.Context(ctx))
			return nil, ncore.TraceError("failed to persist customer", err)
		}
		customerId = lastInsertId
		customer = insertCustomer
		customerXID = customer.CustomerXID
	}

	customer.Status = 1
	err = s.repo.UpdateCustomerByPhone(customer)
	if err != nil {
		log.Errorf("error when update customer by phone: %s", payload.Name, nlogger.Error(err), nlogger.Context(ctx))
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
		s.log.Errorf("error when persist customer verification. Customer Name: %s", payload.Name, nlogger.Error(err), nlogger.Context(ctx))
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
		s.log.Error("Error when persist customer credential", nlogger.Error(err), nlogger.Context(ctx))
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
		s.log.Error("error when persist OTP", nlogger.Error(err), nlogger.Context(ctx))
		return nil, ncore.TraceError("failed to persist OTP", err)
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
		s.log.Error("error when persist access session", nlogger.Error(err), nlogger.Context(ctx))
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
		s.log.Error("error when call login service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, ncore.TraceError("failed to login", err)
	}

	// Send Notification Register
	err = s.SendNotificationRegister(dto.NotificationRegister{
		Customer:     customer,
		Verification: verification,
		RegisterOTP:  registerOTP,
		Payload:      payload,
	})
	if err != nil {
		s.log.Error("error when send notification register", nlogger.Error(err), nlogger.Context(ctx))
	}

	// Delete OTP RegistrationId
	err = s.repo.DeleteVerificationOTP(registerOTP.RegistrationId, customer.Phone)
	if err != nil {
		s.log.Debugf("error when remove verificationOTP: %s. Phone Number : %s.",
			registerOTP.RegistrationId, customer.Phone, nlogger.Error(err), nlogger.Context(ctx),
		)
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

// Register Send OTP

func (s *Service) RegisterStepOne(payload dto.RegisterStepOne) (*dto.RegisterStepOneResponse, error) {
	// Get Context
	ctx := s.ctx

	// validate email
	emailExist, err := s.repo.FindCustomerByEmail(payload.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			emailExist = nil
		} else {
			s.log.Error("error when find customer by email", nlogger.Error(err), nlogger.Context(ctx))
			return nil, ncore.TraceError("error", err)
		}
	}
	if emailExist != nil {
		s.log.Debugf("email already registered. %s", payload.Email, nlogger.Error(err), nlogger.Context(ctx))
		return nil, s.responses.GetError("E_REG_2")
	}

	// validate phone
	phoneExist, err := s.repo.FindCustomerByPhone(payload.PhoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			phoneExist = nil
		} else {
			s.log.Error("failed when find customer by phone", nlogger.Error(err), nlogger.Context(ctx))
			return nil, ncore.TraceError("error", err)
		}
	}
	if phoneExist != nil {
		s.log.Debug("phone number already registered", nlogger.Error(err), nlogger.Context(ctx))
		return nil, s.responses.GetError("E_REG_3")
	}

	sendOtpRequest := dto.SendOTPRequest{
		PhoneNumber: payload.PhoneNumber,
		RequestType: constant.RequestTypeRegister,
	}

	// Send OTP To Phone Number
	resp, err := s.SendOTP(sendOtpRequest)
	if err != nil {
		s.log.Error("error when send OTP", nlogger.Error(err), nlogger.Context(ctx))
		return nil, ncore.TraceError("failed to send OTP", err)
	}

	if resp.ResponseCode == "15" {
		return nil, s.responses.GetError("E_OTP_3")
	}

	if resp.Message != "" {
		s.log.Debugf("Debug: RegisterStepOne OTP CODE %s", resp.Message, nlogger.Context(ctx))
	}

	return &dto.RegisterStepOneResponse{
		Action: resp.ResponseDesc,
	}, nil
}

func (s *Service) RegisterResendOTP(payload dto.RegisterResendOTP) (*dto.RegisterResendOTPResponse, error) {
	// Get Context
	ctx := s.ctx

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
			s.log.Error("failed when query check phone.", nlogger.Error(err), nlogger.Context(ctx))
			return nil, err
		}
	}
	if phoneExist != nil {
		s.log.Debug("Phone already registered", nlogger.Error(err), nlogger.Context(ctx))
		return nil, s.responses.GetError("E_REG_3")
	}

	// Send OTP To Phone Number
	resp, err := s.SendOTP(request)
	if err != nil {
		s.log.Error("error when send otp", nlogger.Error(err), nlogger.Context(ctx))
		return nil, ncore.TraceError("failed to send otp", err)
	}

	if resp.ResponseCode == "15" {
		return nil, s.responses.GetError("E_OTP_3")
	}

	if resp.Message != "" {
		s.log.Debugf("Debug: RegisterResendOTP OTP CODE: %s", resp.Message, nlogger.Error(err), nlogger.Context(ctx))
	}

	return &dto.RegisterResendOTPResponse{
		Action: resp.ResponseDesc,
	}, nil
}

func (s *Service) RegisterStepTwo(payload dto.RegisterStepTwo) (*dto.RegisterStepTwoResponse, error) {
	// Get Context
	ctx := s.ctx

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
			s.log.Error("error when find customer by phone", nlogger.Error(err), nlogger.Context(ctx))
			return nil, ncore.TraceError("failed to find customer", err)
		}
	}
	if phoneExist != nil {
		s.log.Debug("phone already registered", nlogger.Error(err), nlogger.Context(ctx))
		return nil, s.responses.GetError("E_REG_3")
	}

	// Verify OTP To Phone Number
	resp, err := s.VerifyOTP(request)
	if err != nil {
		return nil, ncore.TraceError("error", err)
	}

	// handle Expired OTP
	if resp.ResponseCode == "12" {
		s.log.Errorf("Expired OTP. Phone Number : %s", payload.PhoneNumber, nlogger.Error(err), nlogger.Context(ctx))
		return nil, s.responses.GetError("E_OTP_4")
	}
	// handle Wrong OTP
	if resp.ResponseCode == "14" {
		s.log.Errorf("Wrong OTP. Phone Number : %s", payload.PhoneNumber, nlogger.Error(err), nlogger.Context(ctx))
		return nil, s.responses.GetError("E_OTP_1")
	}

	if resp.ResponseCode != "00" {
		s.log.Errorf("Wrong OTP. Phone Number : %s", payload.PhoneNumber, nlogger.Error(err), nlogger.Context(ctx))
		return nil, s.responses.GetError("E_OTP_1")
	}

	// insert verification otp
	vOTP := &model.VerificationOTP{
		CreatedAt:      time.Now(),
		Phone:          payload.PhoneNumber,
		RegistrationId: xid.New().String(),
	}
	_, err = s.repo.CreateVerificationOTP(vOTP)
	if err != nil {
		s.log.Errorf("error when create verification otp. Phone %s", payload.PhoneNumber, nlogger.Error(err), nlogger.Context(ctx))
		return nil, ncore.TraceError("failed to persist verification otp", err)
	}

	return &dto.RegisterStepTwoResponse{
		RegisterId: vOTP.RegistrationId,
	}, nil
}
