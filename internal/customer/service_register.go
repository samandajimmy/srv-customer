package customer

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/nbs-go/errx"
	"github.com/nbs-go/nlogger"
	"github.com/rs/xid"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
	"strings"
	"time"
)

func (s *Service) Register(payload dto.RegisterNewCustomer) (*dto.RegisterNewCustomerResponse, error) {
	// Get context
	ctx := s.ctx

	// Set Registration ID and Phone Number
	registrationID := payload.RegistrationID
	phoneNumber := payload.PhoneNumber

	// Validate exist
	var customer *model.Customer
	customer, err := s.repo.FindCustomerByEmail(payload.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("error found when get customer by email", nlogger.Error(err), nlogger.Context(ctx))
		return nil, constant.RegistrationFailedError.Trace()
	}

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		customer = nil
	}

	if customer != nil {
		s.log.Debugf("Email already registered: %s", registrationID)
		return nil, constant.EmailHasBeenRegisteredError.Trace()
	}

	// Find registerID
	registerOTP, err := s.repo.FindVerificationOTPByRegistrationIDAndPhone(registrationID, phoneNumber)
	if err != nil {
		s.log.Errorf("Registration ID not found: %s", registrationID)
		return nil, constant.RegistrationFailedError.Trace()
	}

	// Get data user
	customer, err = s.repo.FindCustomerByPhone(phoneNumber)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("error found when get customer", nlogger.Error(err))
		return nil, constant.RegistrationFailedError.Trace()
	}

	tx, err := s.repo.conn.BeginTxx(ctx, nil)
	if err != nil {
		return nil, errx.Trace(err)
	}
	defer s.repo.ReleaseTx(tx, &err)

	var customerXID string
	var customerID int64

	// Update name if customer exists
	if customer != nil {
		customer.FullName = payload.Name
		customer.Email = payload.Email
		err = s.repo.UpdateCustomerByPhone(customer)
		if err != nil {
			s.log.Error("error when update customer by phone", nlogger.Error(err), nlogger.Context(ctx))
			return nil, constant.UsedPhoneNumberError.Trace()
		}

		// Update customerId and CustomerXID
		customerXID = customer.CustomerXID
		customerID = customer.ID
	} else {
		// Create new profile
		customerXID = strings.ToUpper(xid.New().String())
		profile := &dto.CustomerProfileVO{
			MaidenName:         "",
			Gender:             "",
			Nationality:        "",
			DateOfBirth:        "",
			PlaceOfBirth:       "",
			MarriageStatus:     "",
			NPWPNumber:         "",
			IdentityPhotoFile:  "",
			NPWPPhotoFile:      "",
			SidPhotoFile:       "",
			NPWPUpdatedAt:      0,
			ProfileUpdatedAt:   0,
			CifLinkUpdatedAt:   0,
			CifUnlinkUpdatedAt: 0,
		}

		insertCustomer := &model.Customer{
			CustomerXID:    customerXID,
			FullName:       payload.Name,
			Phone:          payload.PhoneNumber,
			Email:          payload.Email,
			Status:         0,
			IdentityType:   0,
			IdentityNumber: "",
			UserRefID:      sql.NullString{},
			Cif:            "",
			Sid:            "",
			ReferralCode:   "",
			Profile:        model.ToCustomerProfile(profile),
			Photos:         nil,
			BaseField:      model.EmptyBaseField,
		}

		// Persist Customer
		lastInsertId, errInsert := s.repo.CreateCustomer(insertCustomer)
		if errInsert != nil {
			s.log.Errorf("error when persist customer: %s", payload.Name, nlogger.Error(errInsert), nlogger.Context(ctx))
			return nil, errx.Trace(errInsert)
		}
		customerID = lastInsertId
		customer = insertCustomer
		customerXID = customer.CustomerXID
	}

	customer.Status = 1
	err = s.repo.UpdateCustomerByPhone(customer)
	if err != nil {
		log.Errorf("error when update customer by phone: %s", payload.Name, nlogger.Error(err), nlogger.Context(ctx))
		return nil, constant.RegistrationFailedError.Trace()
	}

	// Prepare verification model
	verification := &model.Verification{
		Xid:                             strings.ToUpper(xid.New().String()),
		CustomerID:                      customerID,
		KycVerifiedStatus:               0,
		KycVerifiedAt:                   sql.NullTime{},
		EmailVerificationToken:          nval.RandomString(78),
		EmailVerifiedStatus:             0,
		EmailVerifiedAt:                 sql.NullTime{},
		DukcapilVerifiedStatus:          0,
		DukcapilVerifiedAt:              sql.NullTime{},
		FinancialTransactionStatus:      0,
		FinancialTransactionActivatedAt: sql.NullTime{},
		BaseField:                       model.EmptyBaseField,
	}

	// Update verification
	err = s.repo.InsertOrUpdateVerification(verification)
	if err != nil {
		s.log.Errorf("error when persist customer verification. Customer Name: %s", payload.Name, nlogger.Error(err), nlogger.Context(ctx))
		return nil, constant.RegistrationFailedError.Trace()
	}

	// Set metadata credential
	var Format dto.MetadataCredential
	Format.TryLoginAt = ""
	Format.PinCreatedAt = ""
	Format.PinBlockedAt = ""

	metadata, err := json.Marshal(&Format)
	if err != nil {
		s.log.Error("error when marshal metadata", nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	credentialBaseField := model.EmptyBaseField
	credentialBaseField.Metadata = metadata

	// Create credential
	credentialXID := strings.ToUpper(xid.New().String())
	credentialInsert := &model.Credential{
		Xid:                 credentialXID,
		CustomerID:          customerID,
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
		BiometricDeviceID:   "",
		BaseField:           credentialBaseField,
	}

	// Insert or update credential
	err = s.repo.InsertOrUpdateCredential(credentialInsert)
	if err != nil {
		s.log.Error("Error when persist customer credential", nlogger.Error(err), nlogger.Context(ctx))
		return nil, constant.RegistrationFailedError.Trace()
	}

	// Insert OTP
	insertOTP := &model.OTP{
		CustomerID: customerID,
		Content:    "",
		Type:       "registrasi_user",
		Data:       "",
		Status:     "",
		UpdatedAt:  time.Now(),
	}

	// Create OTP repo
	err = s.repo.CreateOTP(insertOTP)
	if err != nil {
		s.log.Error("error when persist OTP", nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	// Init model access session
	insertAccessSession := &model.AccessSession{
		Xid:                  customerXID,
		CustomerID:           customerID,
		ExpiredAt:            time.Now().Add(time.Hour),
		NotificationToken:    payload.FcmToken,
		NotificationProvider: 1,
		BaseField:            model.EmptyBaseField,
	}

	// Create access session
	err = s.repo.CreateAccessSession(insertAccessSession)
	if err != nil {
		s.log.Error("error when persist access session", nlogger.Error(err), nlogger.Context(ctx))
		return nil, constant.RegistrationFailedError.Trace()
	}

	// Init payload Login
	payloadLogin := dto.LoginRequest{
		Email:    payload.Email,
		Password: payload.Password,
		Agen:     payload.Agen,
		Version:  payload.Version,
		FcmToken: payload.FcmToken,
	}
	// Call login service
	loginResponse, err := s.Login(payloadLogin)
	if err != nil {
		s.log.Error("error found when call login service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	// Delete OTP RegistrationId
	err = s.repo.DeleteVerificationOTP(registerOTP.RegistrationID, customer.Phone)
	if err != nil {
		s.log.Errorf("error when remove verificationOTP: %s. Phone Number : %s.",
			registerOTP.RegistrationID, customer.Phone, nlogger.Error(err), nlogger.Context(ctx),
		)
		return nil, constant.RegistrationFailedError.Trace()
	}

	// TODO: Fix payload endpoint to unified create notification service
	// Send Notification Register
	err = s.SendNotificationRegister(dto.NotificationRegister{
		Customer:     customer,
		Verification: verification,
		RegisterOTP:  registerOTP,
		Payload:      payload,
	})
	if err != nil {
		err = nil
		s.log.Error("error when send notification register", nlogger.Context(ctx))
	}

	// Compose response
	resp := s.composeRegisterResponse(loginResponse)

	return resp, nil
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
			return nil, errx.Trace(err)
		}
	}
	if emailExist != nil {
		s.log.Debugf("email already registered. %s", payload.Email)
		return nil, constant.EmailHasBeenRegisteredError.Trace()
	}

	// validate phone
	phoneExist, err := s.repo.FindCustomerByPhone(payload.PhoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			phoneExist = nil
		} else {
			s.log.Error("failed when find customer by phone", nlogger.Error(err), nlogger.Context(ctx))
			return nil, errx.Trace(err)
		}
	}
	if phoneExist != nil {
		s.log.Debug("phone number already registered", nlogger.Error(err))
		return nil, constant.InvalidPasswordError.Trace()
	}

	sendOtpRequest := dto.SendOTPRequest{
		PhoneNumber: payload.PhoneNumber,
		RequestType: constant.RequestTypeRegister,
	}

	// Send OTP To Phone Number
	resp, err := s.SendOTP(sendOtpRequest)
	if err != nil {
		s.log.Error("error found when call send OTP service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	if resp.ResponseCode == "15" {
		return nil, constant.OTPReachResendLimitError.Trace()
	}

	if resp.Message != "" {
		s.log.Debugf("Debug: RegisterStepOne OTP CODE %s", resp.Message)
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
		return nil, constant.InvalidPasswordError.Trace()
	}

	// Send OTP To Phone Number
	resp, err := s.SendOTP(request)
	if err != nil {
		s.log.Error("error when send otp", nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	if resp.ResponseCode == "15" {
		return nil, constant.OTPReachResendLimitError.Trace()
	}

	if resp.Message != "" {
		s.log.Debugf("Debug: RegisterResendOTP OTP CODE: %s", resp.Message)
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
		RequestType: constant.RequestTypeRegister,
	}
	// validate phone
	phoneExist, err := s.repo.FindCustomerByPhone(payload.PhoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			phoneExist = nil
		} else {
			s.log.Error("error when find customer by phone", nlogger.Error(err), nlogger.Context(ctx))
			return nil, errx.Trace(err)
		}
	}
	if phoneExist != nil {
		s.log.Debug("phone already registered", nlogger.Error(err), nlogger.Context(ctx))
		return nil, constant.InvalidPasswordError.Trace()
	}

	// Verify OTP To Phone Number
	resp, err := s.VerifyOTP(request)
	if err != nil {
		return nil, errx.Trace(err)
	}

	// handle Expired OTP
	if resp.ResponseCode == "12" {
		s.log.Errorf("Expired OTP. Phone Number : %s", payload.PhoneNumber, nlogger.Error(err), nlogger.Context(ctx))
		return nil, constant.ExpiredOTPError.Trace()
	}
	// handle Wrong OTP
	if resp.ResponseCode == "14" {
		s.log.Errorf("Wrong OTP. Phone Number : %s", payload.PhoneNumber, nlogger.Error(err), nlogger.Context(ctx))
		return nil, constant.IncorrectOTPError.Trace()
	}

	if resp.ResponseCode != "00" {
		s.log.Errorf("Wrong OTP. Phone Number : %s", payload.PhoneNumber, nlogger.Error(err), nlogger.Context(ctx))
		return nil, constant.IncorrectOTPError.Trace()
	}

	// insert verification otp
	vOTP := &model.VerificationOTP{
		CreatedAt:      time.Now(),
		Phone:          payload.PhoneNumber,
		RegistrationID: xid.New().String(),
	}
	_, err = s.repo.CreateVerificationOTP(vOTP)
	if err != nil {
		s.log.Errorf("error when create verification otp. Phone %s", payload.PhoneNumber, nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	return &dto.RegisterStepTwoResponse{
		RegisterID: vOTP.RegistrationID,
	}, nil
}

func (s *Service) composeRegisterResponse(loginResponse *dto.LoginResponse) *dto.RegisterNewCustomerResponse {
	return &dto.RegisterNewCustomerResponse{
		LoginResponse: loginResponse,
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
	}
}
