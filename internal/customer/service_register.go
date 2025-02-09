package customer

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/nbs-go/errx"
	logOption "github.com/nbs-go/nlogger/v2/option"
	"github.com/rs/xid"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
	"strings"
	"time"
)

func (s *Service) Register(payload dto.RegisterPayload) (*dto.RegisterResult, error) {
	// Set Registration ID and Phone Number
	registrationID := payload.RegistrationID
	phoneNumber := payload.PhoneNumber

	// Validate exist
	var customer *model.Customer
	customer, err := s.repo.FindCustomerByEmail(payload.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("error found when get customer by email", logOption.Error(err))
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
		s.log.Error("Registration ID not found: %s", logOption.Format(registrationID), logOption.Error(err))
		return nil, constant.RegistrationFailedError.Trace()
	}

	// Get data user
	customer, err = s.repo.FindCustomerByPhone(phoneNumber)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("error found when get customer", logOption.Error(err))
		return nil, constant.RegistrationFailedError.Trace()
	}

	tx, err := s.repo.conn.BeginTxx(s.ctx, nil)
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
			s.log.Error("error when update customer by phone", logOption.Error(err))
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
			ReferralCode:   sql.NullString{},
			Profile:        model.ToCustomerProfile(profile),
			Photos:         nil,
			BaseField:      model.EmptyBaseField,
		}

		// Persist Customer
		lastInsertID, errInsert := s.repo.CreateCustomer(insertCustomer)
		if errInsert != nil {
			s.log.Error("error when persist customer: %s", logOption.Format(payload.Name), logOption.Error(errInsert))
			return nil, errx.Trace(errInsert)
		}
		customerID = lastInsertID
		customer = insertCustomer
		customerXID = customer.CustomerXID
	}

	customer.Status = 1
	err = s.repo.UpdateCustomerByPhone(customer)
	if err != nil {
		s.log.Error("error when update customer by phone: %s", logOption.Format(payload.Name), logOption.Error(err))
		return nil, constant.RegistrationFailedError.Trace()
	}

	// Prepare verification model
	verification := &model.Verification{
		Xid:                             strings.ToUpper(xid.New().String()),
		CustomerID:                      customerID,
		KycVerifiedStatus:               0,
		KycVerifiedAt:                   sql.NullTime{},
		EmailVerificationToken:          nval.Bin2Hex(nval.RandomString(78)),
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
		s.log.Error("error when persist customer verification. Customer Name: %s", logOption.Format(payload.Name), logOption.Error(err))
		return nil, constant.RegistrationFailedError.Trace()
	}

	// Set metadata credential
	var Format dto.MetadataCredential
	Format.TryLoginAt = ""
	Format.PinCreatedAt = ""
	Format.PinBlockedAt = ""

	metadata, err := json.Marshal(&Format)
	if err != nil {
		s.log.Error("error when marshal metadata", logOption.Error(err))
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
		s.log.Error("Error when persist customer credential", logOption.Error(err))
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
		s.log.Error("error when persist OTP", logOption.Error(err))
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
		s.log.Error("error when persist access session", logOption.Error(err))
		return nil, constant.RegistrationFailedError.Trace()
	}

	// Init payload Login
	payloadLogin := dto.LoginPayload{
		Email:    payload.Email,
		Password: payload.Password,
		Agen:     payload.Agen,
		Version:  payload.Version,
		FcmToken: payload.FcmToken,
	}
	// Call login service
	loginResponse, err := s.Login(payloadLogin)
	if err != nil {
		s.log.Error("error found when call login service", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	// Delete OTP RegistrationId
	err = s.repo.DeleteVerificationOTP(registerOTP.RegistrationID, customer.Phone)
	if err != nil {
		s.log.Errorf("error when remove verificationOTP: %s. Phone Number : %s.",
			registerOTP.RegistrationID, customer.Phone, logOption.Error(err),
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
		s.log.Error("error when send notification register", logOption.Error(err))
	}

	// Compose response
	resp := s.composeRegisterResult(loginResponse)

	return resp, nil
}

// Register Send OTP

func (s *Service) RegisterStepOne(payload dto.SendOTPPayload) (*dto.SendOTPResult, error) {
	// validate email
	emailExist, err := s.repo.FindCustomerByEmail(payload.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("error when find customer by email", logOption.Error(err))
		return nil, errx.Trace(err)
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		emailExist = nil
	}
	if emailExist != nil {
		s.log.Debugf("email already registered")
		return nil, constant.EmailHasBeenRegisteredError
	}

	// validate phone
	phoneExist, err := s.repo.FindCustomerByPhone(payload.PhoneNumber)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("failed when find customer by phone", logOption.Error(err))
		return nil, errx.Trace(err)
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		phoneExist = nil
	}
	if phoneExist != nil {
		s.log.Debug("phone number already registered")
		return nil, constant.PhoneHasBeenRegisteredError
	}

	sendOtpRequest := dto.SendOTPRequest{
		PhoneNumber: payload.PhoneNumber,
		RequestType: constant.RequestTypeRegister,
	}

	// Send OTP To Phone Number
	resp, err := s.SendOTP(sendOtpRequest)
	if err != nil {
		s.log.Error("error found when call send OTP service", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	if resp.ResponseCode == "15" {
		return nil, constant.OTPReachResendLimitError.Trace()
	}

	return &dto.SendOTPResult{
		Action: resp.ResponseDesc,
	}, nil
}

func (s *Service) RegisterResendOTP(payload dto.RegisterResendOTPPayload) (*dto.RegisterResendOTPResult, error) {
	// validate phone
	phoneExist, err := s.repo.FindCustomerByPhone(payload.PhoneNumber)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("failed when find customer by phone", logOption.Error(err))
		return nil, errx.Trace(err)
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		phoneExist = nil
	}
	if phoneExist != nil {
		s.log.Debug("phone number already registered")
		return nil, constant.PhoneHasBeenRegisteredError
	}

	// Set request
	request := dto.SendOTPRequest{
		PhoneNumber: payload.PhoneNumber,
		RequestType: constant.RequestTypeRegister,
	}

	// Send OTP To Phone Number
	resp, err := s.SendOTP(request)
	if err != nil {
		s.log.Error("error when send otp", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	if resp.ResponseCode == "15" {
		return nil, constant.OTPReachResendLimitError.Trace()
	}

	if resp.Message != "" {
		s.log.Debugf("Debug: RegisterResendOTP OTP CODE: %s", resp.Message)
	}

	return &dto.RegisterResendOTPResult{
		Action: resp.ResponseDesc,
	}, nil
}

func (s *Service) RegisterStepTwo(payload dto.RegisterVerifyOTPPayload) (*dto.RegisterVerifyOTPResult, error) {
	// validate phone
	phoneExist, err := s.repo.FindCustomerByPhone(payload.PhoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			phoneExist = nil
		} else {
			s.log.Error("error when find customer by phone", logOption.Error(err))
			return nil, errx.Trace(err)
		}
	}
	if phoneExist != nil {
		s.log.Debug("phone already registered")
		return nil, constant.InvalidPasswordError.Trace()
	}

	// Set request
	request := dto.VerifyOTPRequest{
		PhoneNumber: payload.PhoneNumber,
		Token:       payload.OTP,
		RequestType: constant.RequestTypeRegister,
	}

	// Verify OTP To Phone Number
	resp, err := s.VerifyOTP(request)
	if err != nil {
		return nil, errx.Trace(err)
	}

	// handle Expired OTP
	if resp.ResponseCode == "12" {
		s.log.Errorf("Expired OTP. Phone Number : %s", payload.PhoneNumber)
		return nil, constant.ExpiredOTPError.Trace()
	}
	// handle Wrong OTP
	if resp.ResponseCode == "14" {
		s.log.Errorf("Wrong OTP. Phone Number : %s", payload.PhoneNumber)
		return nil, constant.IncorrectOTPError.Trace()
	}

	if resp.ResponseCode != "00" {
		s.log.Errorf("Wrong OTP. Phone Number : %s", payload.PhoneNumber)
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
		s.log.Errorf("error when create verification otp. Phone %s", payload.PhoneNumber, logOption.Error(err))
		return nil, errx.Trace(err)
	}

	return &dto.RegisterVerifyOTPResult{
		RegisterID: vOTP.RegistrationID,
	}, nil
}

func (s *Service) composeRegisterResult(loginResponse *dto.LoginResult) *dto.RegisterResult {
	return &dto.RegisterResult{
		LoginResult: loginResponse,
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
