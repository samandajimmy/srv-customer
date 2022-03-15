package customer

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/nbs-go/errx"
	"github.com/nbs-go/nlogger/v2"
	"regexp"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
	"strings"
	"time"
)

func (s *Service) ValidateJWT(token string) (jwt.Token, error) {
	// Parsing Token
	t, err := jwt.ParseString(token, jwt.WithVerify(constant.JWTSignature, []byte(s.config.JWTKey)))
	if err != nil {
		s.log.Error("parsing jwt token", nlogger.Error(err), nlogger.Context(s.ctx))
		return nil, constant.InvalidJWTFormatError.Trace()
	}

	if err = jwt.Validate(t); err != nil {
		s.log.Error("error when validate", nlogger.Error(err), nlogger.Context(s.ctx))
		return nil, constant.ExpiredJWTError.Trace()
	}

	err = jwt.Validate(t, jwt.WithIssuer(constant.JWTIssuer))
	if err != nil {
		s.log.Error("error found when validate with issuer", nlogger.Error(err), nlogger.Context(s.ctx))
		return nil, constant.InvalidJWTIssuerError.Trace()
	}

	return t, nil
}

func (s *Service) ValidateTokenAndRetrieveUserRefID(tokenString string) (string, error) {
	// Get Context
	ctx := s.ctx

	// validate JWT
	token, err := s.ValidateJWT(tokenString)
	if err != nil {
		s.log.Error("error when validate JWT", nlogger.Error(err), nlogger.Context(ctx))
		return "", errx.Trace(err)
	}

	accessToken, _ := token.Get("access_token")

	tokenID, _ := token.Get("id")

	// Session token
	key := fmt.Sprintf("%s:%s:%s", constant.Prefix, constant.CacheTokenJWT, tokenID)

	tokenFromCache, err := s.CacheGet(key)
	if err != nil {
		s.log.Error("error get token from cache", nlogger.Error(err), nlogger.Context(ctx))
		return "", errx.Trace(err)
	}

	if accessToken != tokenFromCache {
		return "", constant.InvalidTokenError
	}

	userRefID := nval.ParseStringFallback(tokenID, "")

	return userRefID, nil

}

func (s *Service) UpdatePhoneNumber(payload dto.ChangePhoneNumberPayload) (*dto.ChangePhoneNumberResult, error) {
	// Check if phone number exist // TODO: Check token if admin
	isExists, err := s.repo.PhoneNumberIsExists(payload.NewPhoneNumber)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("error when check email is exists", nlogger.Error(err))
		return nil, errx.Trace(err)
	}

	if isExists {
		s.log.Debug("Phone number has been used")
		return nil, constant.UsedPhoneNumberError.Trace()
	}

	// Find customer
	customer, err := s.repo.FindCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		s.log.Error("error when find current customer", nlogger.Error(err))
		return nil, errx.Trace(err)
	}

	// Validate payload with current data
	err = s.handleValidateChangePasswordData(payload, customer)
	if err != nil {
		return nil, err
	}

	// Check if customer has cif
	if customer.Cif != "" {
		return s.handleChangePhoneNumberCore(dto.ChangePhoneNumberRequestCore{
			CurrentPhoneNumber: payload.CurrentPhoneNumber,
			NewPhoneNumber:     payload.NewPhoneNumber,
			FullName:           payload.FullName,
			DateOfBirth:        payload.DateOfBirth,
			MaidenName:         payload.MaidenName,
			Cif:                customer.Cif,
		}, customer)
	}

	// Update phone number
	customer.Phone = payload.NewPhoneNumber
	customer.UpdatedAt = time.Now()
	customer.Version++

	// Persist
	err = s.repo.UpdateCustomerByUserRefID(customer, payload.UserRefID)
	if err != nil {
		s.log.Error("error when update phone number", nlogger.Error(err))
		return nil, errx.Trace(err)
	}

	return &dto.ChangePhoneNumberResult{
		PhoneNumber: customer.Phone,
	}, nil
}

func (s *Service) handleValidateChangePasswordData(payload dto.ChangePhoneNumberPayload, customer *model.Customer) error {
	// Check Maiden Name
	if !strings.EqualFold(customer.Profile.MaidenName, payload.MaidenName) {
		s.log.Errorf(
			`MaidenName is not correct expected: %s actual: %s`,
			customer.Profile.MaidenName,
			payload.MaidenName,
		)
		return constant.ResourceNotFoundError
	}

	// Check Full Name
	if !strings.EqualFold(customer.FullName, payload.FullName) {
		s.log.Errorf(
			`fullName is not correct expected: %s actual: %s`,
			customer.FullName,
			payload.FullName,
		)
		return constant.ResourceNotFoundError
	}

	// Check Date Of Birth
	if customer.Profile.DateOfBirth != payload.DateOfBirth {
		s.log.Errorf(
			`dateOfBirth is not correct expected: %s actual: %s`,
			customer.Profile.DateOfBirth,
			payload.DateOfBirth,
		)
		return constant.ResourceNotFoundError
	}

	// Check Current Phone Number
	if customer.Phone != payload.CurrentPhoneNumber {
		s.log.Errorf(
			`Current phone number is not correct expected: %s actual: %s`,
			customer.Phone,
			payload.CurrentPhoneNumber,
		)
		return constant.ResourceNotFoundError
	}

	return nil
}

func (s *Service) handleChangePhoneNumberCore(payload dto.ChangePhoneNumberRequestCore, customer *model.Customer) (*dto.ChangePhoneNumberResult, error) { // Hit core change phone number
	resp, err := s.ChangePhoneNumber(payload)
	if err != nil {
		return nil, errx.Trace(err)
	}

	// Handle if hit core api not success
	if resp.ResponseCode != "00" {
		log.Debugf(`Response code: %s message: %s, Description: %s`, resp.ResponseCode, resp.Message, resp.ResponseDesc)
		return nil, constant.ChangePhoneNumberError.AddMetadata(constant.MetadataMessage, resp.ResponseDesc)
	}

	// Update phone number
	customer.Phone = payload.NewPhoneNumber
	customer.UpdatedAt = time.Now()
	customer.Version++

	// Persist update phone number
	err = s.repo.UpdateCustomerByUserRefID(customer, customer.UserRefID.String)
	if err != nil {
		s.log.Error("error when update phone number", nlogger.Error(err))
		return nil, errx.Trace(err)
	}

	return &dto.ChangePhoneNumberResult{
		PhoneNumber: customer.Phone,
	}, nil
}

func (s *Service) ValidatePin(payload *dto.ValidatePinPayload) error {
	// Get Context
	ctx := s.ctx

	// Validate Increasing Pin
	err := validateIncreasingPin(payload)
	if err != nil {
		s.log.Error(`new pin isn't passed validation increasing pin'`, nlogger.Error(err), nlogger.Context(ctx))
		return errx.Trace(err)
	}

	err = validateRepetitionPin(payload)
	if err != nil {
		s.log.Error(`new pin isn't passed validation increasing pin'`, nlogger.Error(err), nlogger.Context(ctx))
		return errx.Trace(err)
	}

	return nil
}

func validateIncreasingPin(payload *dto.ValidatePinPayload) error {
	increasingPin := regexp.MustCompile(`^[\d{6}$]1?2?3?4?5?6?7?8?9?0?$`)

	if len(increasingPin.FindStringSubmatch(payload.NewPin)) > 0 {
		return constant.InvalidIncrementalPIN
	}

	return nil
}

func validateRepetitionPin(payload *dto.ValidatePinPayload) error {
	var repetitionPin = regexp.MustCompile(`^0{6}$|^1{6}$|^2{6}$|^3{6}$|^4{6}$|^5{6}$|^6{6}$|^7{6}$|^8{6}$|^9{6}$`)
	if len(repetitionPin.FindStringSubmatch(payload.NewPin)) > 0 {
		return constant.InvalidRepeatedPIN
	}

	return nil
}

func (s *Service) CheckPinUser(payload *dto.CheckPinPayload) error {
	// Get Context
	ctx := s.ctx

	// If not check pin
	if !payload.CheckPIN {
		return nil
	}

	// Get customer id
	customer, err := s.repo.FindCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		s.log.Error("error found when get customer repo", nlogger.Error(err), nlogger.Context(ctx))
		err = handleErrorRepository(err, constant.ResourceNotFoundError)
		return errx.Trace(err)
	}

	// Get customer pin
	credential, err := s.repo.FindCredentialByCustomerID(customer.ID)
	if err != nil {
		s.log.Error("error found when get credential repo", nlogger.Context(ctx))
		err = handleErrorRepository(err, constant.AccountPINIsNotActive)
		return errx.Trace(err)
	}

	if credential.Pin == "" {
		return constant.AccountPINIsNotActive
	}

	// Check if pin is blocked
	if credential.PinBlockedStatus != constant.Unblocked {
		s.log.Error("pin is blocked", nlogger.Context(ctx))
		return constant.AccountPINIsBlocked
	}

	// Find pin with md5
	valid, err := s.repo.IsValidPin(customer.ID, nval.MD5(payload.Pin))
	if err != nil {
		s.log.Error("error found when querying valid pin", nlogger.Error(err), nlogger.Context(ctx))
		err = handleErrorRepository(err, constant.ResourceNotFoundError)
		return errx.Trace(err)
	}

	if !valid {
		err = s.handlePinNotValid(credential)
		return errx.Trace(err)
	}

	// Unblock pin user
	err = s.unblockPINUser(customer)
	if err != nil {
		return errx.Trace(err)
	}

	return nil
}

func (s *Service) handlePinNotValid(credential *model.Credential) error {
	// Add counter
	counter := credential.PinCounter + 1
	credential.PinCounter = counter

	// Set response message
	var errMessage error

	switch counter {
	case constant.WrongPIN:
		errMessage = constant.WrongPINInput1
	case constant.WrongPIN2:
		errMessage = constant.WrongPINInput2
	case constant.MaxWrongPIN:
		// Update pin block status
		errMessage = constant.AccountPINIsBlocked
		credential.PinBlockedStatus = constant.Blocked
	default:
		errMessage = nhttp.UnauthorizedError
	}

	// Update pin last access at
	credential.PinLastAccessAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}

	err := s.repo.UpdateCredential(credential)
	if err != nil {
		s.log.Error("error when update credential.", nlogger.Error(err))
		return errx.Trace(err)
	}

	return errMessage
}

func (s *Service) unblockPINUser(customer *model.Customer) error {
	// Get credential
	credential, err := s.repo.FindCredentialByCustomerID(customer.ID)
	if err != nil {
		s.log.Error("error found when get credential", nlogger.Error(err), nlogger.Context(s.ctx))
		err = handleErrorRepository(err, constant.ResourceNotFoundError)
		return errx.Trace(err)
	}

	// Unblocked pin
	credential.PinCounter = 0
	credential.PinBlockedStatus = constant.Unblocked

	err = s.repo.UpdateCredential(credential)
	if err != nil {
		s.log.Error("error found when update credential", nlogger.Error(err), nlogger.Context(s.ctx))
		return errx.Trace(err)
	}

	return nil
}

func (s *Service) UpdatePin(payload *dto.UpdatePinPayload) (*dto.UpdatePinResult, error) {
	// If check pin true
	err := s.CheckPinUser(&dto.CheckPinPayload{
		Pin:       payload.PIN,
		UserRefID: payload.UserRefID,
		CheckPIN:  payload.CheckPIN,
	})
	if err != nil {
		return nil, err
	}

	// Validate pin user
	err = s.ValidatePin(&dto.ValidatePinPayload{
		NewPin: payload.NewPIN,
	})
	if err != nil {
		return nil, errx.Trace(err)
	}

	err = s.handleUpdatePin(payload.UserRefID, payload.NewPIN)
	if err != nil {
		s.log.Error("error found when handle update pin", nlogger.Context(s.ctx))
		return nil, errx.Trace(err)
	}

	// TODO: Audit log reset pin

	return &dto.UpdatePinResult{
		Title: "PIN Berhasil Diubah",
		Text:  "Selamat! Kamu berhasil mengubah PIN Pegadaian Digital.",
	}, nil
}

func (s *Service) handleUpdatePin(userRefID string, newPin string) error {
	// Get customer id
	customer, err := s.repo.FindCustomerByUserRefID(userRefID)
	if err != nil {
		s.log.Error("error found when get customer repo", nlogger.Error(err), nlogger.Context(s.ctx))
		err = handleErrorRepository(err, constant.ResourceNotFoundError)
		return errx.Trace(err)
	}

	// Get customer pin
	credential, err := s.repo.FindCredentialByCustomerID(customer.ID)
	if err != nil {
		s.log.Error("error found when get credential repo", nlogger.Context(s.ctx))
		err = handleErrorRepository(err, constant.AccountPINIsNotActive)
		return errx.Trace(err)
	}

	// Prepare update pin
	credential.Pin = nval.MD5(newPin)
	credential.UpdatedAt = time.Now()
	credential.PinUpdatedAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}
	credential.Version++

	err = s.repo.UpdateCredential(credential)
	if err != nil {
		s.log.Error("error when update credential.", nlogger.Error(err))
		return errx.Trace(err)
	}

	return nil
}

func (s *Service) CheckOTPPinCreate(payload *dto.CheckOTPPinPayload) error {
	// Get customer id
	customer, err := s.repo.FindCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		s.log.Error("error found when get customer repo", nlogger.Error(err), nlogger.Context(s.ctx))
		err = handleErrorRepository(err, constant.ResourceNotFoundError)
		return errx.Trace(err)
	}

	// Set payload
	payloadSwitching := &dto.RestSwitchingOTPPinCreate{
		Cif:  customer.Cif,
		OTP:  payload.OTP,
		NoHp: customer.Phone,
	}

	// Rest switching customer
	switchingResponse, err := s.customerActivation(payloadSwitching)
	if err != nil {
		s.log.Error("error found when get otp pin create", nlogger.Error(err), nlogger.Context(s.ctx))
		return errx.Trace(err)
	}

	// Get from cache when switching response is not success
	if switchingResponse.ResponseCode != "00" {
		return constant.IncorrectOTPError
	}

	return nil
}

func (s *Service) CreatePinUser(payload *dto.PostCreatePinPayload) error {
	// Get user by ref id
	customer, err := s.repo.FindCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		s.log.Error("error found when get customer repo", nlogger.Error(err), nlogger.Context(s.ctx))
		err = handleErrorRepository(err, constant.ResourceNotFoundError)
		return errx.Trace(err)
	}

	// Check OTP
	err = s.CheckOTPPinCreate(&dto.CheckOTPPinPayload{
		OTP:       payload.OTP,
		UserRefID: payload.UserRefID,
	})
	if err != nil {
		return errx.Trace(err)
	}

	// Unblock pin user
	err = s.unblockPINUser(customer)
	if err != nil {
		return errx.Trace(err)
	}

	// Update pin
	_, err = s.UpdatePin(&dto.UpdatePinPayload{
		UserRefID: customer.UserRefID.String,
		NewPIN:    payload.NewPIN,
		CheckPIN:  false,
	})
	if err != nil {
		s.log.Error("error found when update customer pin", nlogger.Error(err), nlogger.Context(s.ctx))
		return errx.Trace(err)
	}

	// Aktivasi finansial
	err = s.activateFinancialStatus(customer)
	if err != nil {
		return errx.Trace(err)
	}

	// TODO: Audit log aktifasi finansial

	return nil
}

func (s *Service) activateFinancialStatus(customer *model.Customer) error {
	// Get verification
	verification, err := s.repo.FindVerificationByCustomerID(customer.ID)
	if err != nil {
		s.log.Error("error found when get verification from repo", nlogger.Error(err), nlogger.Context(s.ctx))
		err = handleErrorRepository(err, constant.ResourceNotFoundError)
		return errx.Trace(err)
	}

	// Update financial transaction status
	verification.FinancialTransactionStatus = constant.Enabled
	verification.FinancialTransactionActivatedAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}
	verification.UpdatedAt = time.Now()
	verification.Version++

	err = s.repo.UpdateVerification(verification)
	if err != nil {
		s.log.Error("error found when update verification", nlogger.Error(err), nlogger.Context(s.ctx))
		return errx.Trace(err)
	}
	return nil
}

func (s *Service) CheckOTPForgetPin(payload *dto.CheckOTPPinPayload) error {
	// Get customer id
	customer, err := s.repo.FindCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		s.log.Error("error found when get customer repo", nlogger.Error(err), nlogger.Context(s.ctx))
		err = handleErrorRepository(err, constant.ResourceNotFoundError)
		return errx.Trace(err)
	}

	// Set payload
	payloadSwitching := &dto.RestSwitchingOTPForgetPin{
		Cif:         customer.Cif,
		Flag:        "K",
		NoHp:        customer.Phone,
		NoRek:       customer.Phone,
		RequestType: constant.RequestTypePinReset,
		OTP:         payload.OTP,
	}

	// Rest switching customer
	switchingResponse, err := s.otpValidate(payloadSwitching)
	if err != nil {
		s.log.Error("error found when get reset pin otp validate", nlogger.Error(err), nlogger.Context(s.ctx))
		return errx.Trace(err)
	}

	// Get from cache when switching response is not success
	if switchingResponse.ResponseCode != "00" {
		s.log.Error("error rest switching otp reset pin", nlogger.Context(s.ctx))
		return constant.IncorrectOTPError
	}

	return nil
}

func (s *Service) ForgetPin(payload *dto.ForgetPinPayload) error {
	// Get user by ref id
	customer, err := s.repo.FindCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		s.log.Error("error found when get customer repo", nlogger.Error(err), nlogger.Context(s.ctx))
		err = handleErrorRepository(err, constant.ResourceNotFoundError)
		return errx.Trace(err)
	}

	// Check OTP Reset PIN
	payloadOTP := &dto.CheckOTPPinPayload{
		OTP:       payload.OTP,
		UserRefID: payload.UserRefID,
	}

	err = s.CheckOTPForgetPin(payloadOTP)
	if err != nil {
		s.log.Error("error found when check otp forget pin", nlogger.Error(err), nlogger.Context(s.ctx))
		return errx.Trace(err)
	}

	// Unblock pin user
	err = s.unblockPINUser(customer)
	if err != nil {
		return errx.Trace(err)
	}

	// Set payload update pin
	userRefID := ""
	if customer.UserRefID.Valid {
		userRefID = customer.UserRefID.String
	}

	payloadUpdatePIN := &dto.UpdatePinPayload{
		UserRefID: userRefID,
		NewPIN:    payload.NewPIN,
		CheckPIN:  false,
	}

	// Update pin
	_, err = s.UpdatePin(payloadUpdatePIN)
	if err != nil {
		s.log.Error("error found when update customer pin", nlogger.Error(err), nlogger.Context(s.ctx))
		return errx.Trace(err)
	}

	return nil
}

func (s *Service) findOrFailCustomerByUserRefID(userRefID string) (*model.Customer, error) {
	customer, err := s.repo.FindCustomerByUserRefID(userRefID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("error found when find customer by userRefID", nlogger.Error(err))
		return nil, errx.Trace(err)
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, constant.CustomerNotFoundError
	}
	return customer, nil
}

func handleErrorRepository(errRepo error, errMsg error) error {
	if errors.Is(errRepo, sql.ErrNoRows) {
		return errMsg
	}
	return errRepo
}

func (s *Service) SendOTPResetPassword(payload dto.OTPResetPasswordPayload) error {
	// Send OTP To Phone Number
	resp, err := s.SendOTP(dto.SendOTPRequest{
		PhoneNumber: payload.Email,
		RequestType: constant.RequestResetPassword,
	})
	if err != nil {
		s.log.Error("error found when call send OTP service", nlogger.Error(err), nlogger.Context(s.ctx))
		return errx.Trace(err)
	}

	s.log.Debugf("Debug: reset password message %s", resp.Message)

	if resp.ResponseCode != "00" {
		s.log.Error("error rest switching otp reset pin", nlogger.Context(s.ctx))
		return constant.FailedResendOTP.Trace()
	}

	return nil
}

func (s *Service) VerifyOTPResetPassword(payload dto.VerifyOTPResetPasswordPayload) error {
	// Send OTP To Phone Number
	resp, err := s.VerifyOTP(dto.VerifyOTPRequest{
		PhoneNumber: payload.Email,
		RequestType: constant.RequestResetPassword,
		Token:       payload.OTP,
	})
	if err != nil {
		s.log.Error("error found when call send OTP service", nlogger.Error(err), nlogger.Context(s.ctx))
		return errx.Trace(err)
	}

	s.log.Debugf("Debug: verify reset password otp message %s", resp.Message)

	if resp.ResponseCode != "00" {
		s.log.Error("error rest switching otp reset pin")
		return constant.InvalidOTPError.Trace()
	}

	return nil
}

func (s *Service) ResetPasswordByOTP(payload dto.ResetPasswordByOTPPayload) error {
	// Get customer
	customer, err := s.repo.FindCustomerByEmailOrPhone(payload.Email)
	if err != nil {
		s.log.Error("error found when get customer repo", nlogger.Error(err), nlogger.Context(s.ctx))
		err = handleErrorRepository(err, constant.ResourceNotFoundError)
		return errx.Trace(err)
	}

	if customer.Status != constant.Enabled {
		return constant.AccountIsNotActiveError
	}

	// Check OTP Reset Password
	err = s.VerifyOTPResetPassword(dto.VerifyOTPResetPasswordPayload{
		Email: payload.Email,
		OTP:   payload.OTP,
	})
	if err != nil {
		s.log.Error("Error when check otp reset password")
		return errx.Trace(err)
	}

	// Get credential customer
	credential, err := s.repo.FindCredentialByCustomerID(customer.ID)
	if err != nil {
		s.log.Error("error found when get credential repo", nlogger.Error(err), nlogger.Context(s.ctx))
		err = handleErrorRepository(err, constant.ResourceNotFoundError)
		return errx.Trace(err)
	}

	// Update credential
	password := nval.MD5(payload.Password)
	credential.Password = password
	credential.BiometricLogin = constant.Disabled
	credential.BiometricDeviceID = ""
	credential.Version++

	err = s.repo.UpdateCredential(credential)
	if err != nil {
		s.log.Error("error when update credential.", nlogger.Error(err), nlogger.Context(s.ctx))
		return errx.Trace(err)
	}

	// Synchronize Password PDS
	err = s.HandleSynchronizePassword(customer, password)
	if err != nil {
		return errx.Trace(err)
	}

	return nil
}

func (s *Service) ChangeEmail(payload dto.EmailChangePayload) error {

	// TODO: Validation check if token is admin

	// Check if email is available
	isExist, err := s.repo.EmailIsExists(payload.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("error when check email is exists", nlogger.Error(err))
		return errx.Trace(err)
	}

	if isExist {
		s.log.Debug("Email has been registered")
		return constant.EmailHasBeenRegisteredError.Trace()
	}

	// Find customer
	customer, err := s.repo.FindCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		s.log.Error("error when find current customer", nlogger.Error(err))
		err = handleErrorRepository(err, constant.ResourceNotFoundError)
		return errx.Trace(err)
	}

	// Find verification
	verification, err := s.repo.FindVerificationByCustomerID(customer.ID)
	if err != nil {
		s.log.Error("error when find current customer", nlogger.Error(err))
		err = handleErrorRepository(err, constant.ResourceNotFoundError)
		return errx.Trace(err)
	}

	// Begin transaction
	tx, err := s.repo.conn.BeginTxx(s.ctx, nil)
	if err != nil {
		return errx.Trace(err)
	}
	defer s.repo.ReleaseTx(tx, &err)

	// Update email
	customer.Email = payload.Email
	customer.UpdatedAt = time.Now()
	customer.Version++

	// Persist update email
	err = s.repo.UpdateCustomerByUserRefID(customer, payload.UserRefID)
	if err != nil {
		s.log.Error("error when update email", nlogger.Error(err))
		return errx.Trace(err)
	}

	// Update verification
	verification.EmailVerificationToken = nval.RandomString(74)
	verification.EmailVerifiedStatus = 0
	verification.EmailVerifiedAt = sql.NullTime{}
	verification.UpdatedAt = time.Now()
	verification.Version++

	err = s.repo.UpdateVerificationByCustomerID(verification)
	if err != nil {
		s.log.Error("error when update verification email", nlogger.Error(err))
		return errx.Trace(err)
	}

	// TODO: Send Notification Email Change

	return nil
}

func (s *Service) PostUpdateSmartAccess(payload dto.UpdateSmartAccessPayload) error {
	// Get customer
	customer, err := s.repo.FindCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		s.log.Error("error when find current customer", nlogger.Error(err))
		err = handleErrorRepository(err, constant.ResourceNotFoundError)
		return errx.Trace(err)
	}

	// Get credential
	credential, err := s.repo.FindCredentialByCustomerID(customer.ID)
	if err != nil {
		s.log.Error("error when find current credential", nlogger.Error(err))
		err = handleErrorRepository(err, constant.ResourceNotFoundError)
		return errx.Trace(err)
	}

	// Unset biometric
	err = s.unsetBiometric(credential)
	if err != nil {
		return err
	}

	// Activate biometric
	credential.BiometricLogin = payload.UseBiometric
	credential.BiometricDeviceID = payload.DeviceID
	err = s.repo.UpdateCredential(credential)
	if err != nil {
		s.log.Error("error when activate biometric", nlogger.Error(err))
		return errx.Trace(err)
	}

	return nil
}

func (s *Service) unsetBiometric(credential *model.Credential) error {
	// Unset biometric
	credential.BiometricDeviceID = ""
	credential.BiometricLogin = constant.Disabled
	err := s.repo.UpdateCredential(credential)
	if err != nil {
		s.log.Error("error when unset biometric", nlogger.Error(err))
		return errx.Trace(err)
	}

	return nil
}

func (s *Service) GetSmartAccessStatus(payload dto.GetSmartAccessStatusPayload) (*dto.GetSmartAccessStatusResult, error) {
	// Get customer
	customer, err := s.repo.FindCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		s.log.Error("error when find current customer", nlogger.Error(err))
		err = handleErrorRepository(err, constant.ResourceNotFoundError)
		return nil, errx.Trace(err)
	}

	// Get credential
	credential, err := s.repo.FindCredentialByCustomerID(customer.ID)
	if err != nil {
		s.log.Error("error when find current credential", nlogger.Error(err))
		err = handleErrorRepository(err, constant.ResourceNotFoundError)
		return nil, errx.Trace(err)
	}

	// validate device id
	isMatchDevice := true
	if credential.BiometricDeviceID != payload.DeviceID {
		isMatchDevice = false
	}

	return &dto.GetSmartAccessStatusResult{
		UserRefID:            customer.UserRefID.String,
		DeviceID:             credential.BiometricDeviceID,
		IsSetBiometric:       isMatchDevice,
		IsSetBiometricDevice: credential.BiometricLogin,
	}, nil
}

func (s *Service) HandleSynchronizePassword(customer *model.Customer, password string) error {
	resp, err := s.SynchronizePassword(dto.RegisterPayload{
		PhoneNumber: customer.Phone,
		Password:    password,
	})
	if err != nil {
		s.log.Error("error found when sync data customer via API PDS", nlogger.Error(err), nlogger.Context(s.ctx))
		return errx.Trace(err)
	}

	// handle status error
	if resp.Status != "success" {
		s.log.Error("Get Error from SynchronizePassword.", nlogger.Error(err))
		return nhttp.InternalError.Trace(errx.Errorf(resp.Message))
	}

	return nil
}
