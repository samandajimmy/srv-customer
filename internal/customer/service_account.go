package customer

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/nbs-go/errx"
	logOption "github.com/nbs-go/nlogger/v2/option"
	"github.com/rs/xid"
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
		s.log.Error("parsing jwt token", logOption.Error(err))
		return nil, constant.InvalidJWTFormatError.Trace()
	}

	if err = jwt.Validate(t); err != nil {
		s.log.Error("error when validate", logOption.Error(err))
		return nil, constant.ExpiredJWTError.Trace()
	}

	err = jwt.Validate(t, jwt.WithIssuer(constant.JWTIssuer))
	if err != nil {
		s.log.Error("error found when validate with issuer", logOption.Error(err))
		return nil, constant.InvalidJWTIssuerError.Trace()
	}

	return t, nil
}

func (s *Service) ValidateTokenAndRetrieveUserRefID(tokenString string) (string, error) {
	// validate JWT
	token, err := s.ValidateJWT(tokenString)
	if err != nil {
		s.log.Error("error when validate JWT", logOption.Error(err))
		return "", errx.Trace(err)
	}

	accessToken, _ := token.Get("access_token")

	tokenID, _ := token.Get("id")

	// Session token
	key := fmt.Sprintf("%s:%s:%s", constant.Prefix, constant.CacheTokenJWT, tokenID)

	tokenFromCache, err := s.CacheGet(key)
	if err != nil {
		s.log.Error("error get token from cache", logOption.Error(err))
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
		s.log.Error("error when check email is exists", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	if isExists {
		s.log.Debug("Phone number has been used")
		return nil, constant.UsedPhoneNumberError.Trace()
	}

	// Find customer
	customer, err := s.repo.FindCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		s.log.Error("error when find current customer", logOption.Error(err))
		err = handleErrorRepository(err, constant.CustomerNotFoundError)
		return nil, errx.Trace(err)
	}

	// Validate payload with current data
	err = s.handleValidateChangePasswordData(payload, customer)
	if err != nil {
		return nil, err
	}

	// Check if customer has cif
	if customer.Cif != "" {
		// Format date of birth  to YYYY-MM-DD
		dob := strings.Split(payload.DateOfBirth, "-")
		payload.DateOfBirth = fmt.Sprintf("%s-%s-%s", dob[2], dob[1], dob[0])

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
		s.log.Error("error when update phone number", logOption.Error(err))
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
		s.log.Error("error when update phone number", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	return &dto.ChangePhoneNumberResult{
		PhoneNumber: customer.Phone,
	}, nil
}

func (s *Service) ValidatePin(payload *dto.ValidatePinPayload) error {
	// Validate Increasing Pin
	err := validateIncreasingPin(payload)
	if err != nil {
		s.log.Error(`new pin isn't passed validation increasing pin`, logOption.Error(err))
		return errx.Trace(err)
	}

	err = validateRepetitionPin(payload)
	if err != nil {
		s.log.Error(`new pin isn't passed validation increasing pin`, logOption.Error(err))
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
	// If not check pin
	if !payload.CheckPIN {
		return nil
	}

	// Get customer id
	customer, err := s.repo.FindCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		s.log.Error("error found when get customer repo", logOption.Error(err))
		err = handleErrorRepository(err, constant.ResourceNotFoundError)
		return errx.Trace(err)
	}

	// Get customer pin
	credential, err := s.repo.FindCredentialByCustomerID(customer.ID)
	if err != nil {
		s.log.Error("error found when get credential repo")
		err = handleErrorRepository(err, constant.AccountPINIsNotActive)
		return errx.Trace(err)
	}

	if credential.Pin == "" {
		return constant.AccountPINIsNotActive
	}

	// Check if pin is blocked
	if credential.PinBlockedStatus != constant.Unblocked {
		s.log.Error("pin is blocked")
		return constant.AccountPINIsBlocked
	}

	// Find pin with md5
	valid, err := s.repo.IsValidPin(customer.ID, nval.MD5(payload.Pin))
	if err != nil {
		s.log.Error("error found when querying valid pin", logOption.Error(err))
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
		s.log.Error("error when update credential.", logOption.Error(err))
		return errx.Trace(err)
	}

	return errMessage
}

func (s *Service) unblockPINUser(customer *model.Customer) error {
	// Get credential
	credential, err := s.repo.FindCredentialByCustomerID(customer.ID)
	if err != nil {
		s.log.Error("error found when get credential", logOption.Error(err))
		err = handleErrorRepository(err, constant.ResourceNotFoundError)
		return errx.Trace(err)
	}

	// Unblocked pin
	credential.PinCounter = 0
	credential.PinBlockedStatus = constant.Unblocked

	err = s.repo.UpdateCredential(credential)
	if err != nil {
		s.log.Error("error found when update credential", logOption.Error(err))
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
		s.log.Error("error found when handle update pin", logOption.Error(err))
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
		s.log.Error("error found when get customer repo", logOption.Error(err))
		err = handleErrorRepository(err, constant.ResourceNotFoundError)
		return errx.Trace(err)
	}

	// Get customer pin
	credential, err := s.repo.FindCredentialByCustomerID(customer.ID)
	if err != nil {
		s.log.Error("error found when get credential repo", logOption.Error(err))
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
		s.log.Error("error when update credential.", logOption.Error(err))
		return errx.Trace(err)
	}

	return nil
}

func (s *Service) CheckOTPPinCreate(payload *dto.CheckOTPPinPayload) error {
	// Get customer id
	customer, err := s.repo.FindCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		s.log.Error("error found when get customer repo", logOption.Error(err))
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
		s.log.Error("error found when get otp pin create", logOption.Error(err))
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
		s.log.Error("error found when get customer repo", logOption.Error(err))
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
		s.log.Error("error found when update customer pin", logOption.Error(err))
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
		s.log.Error("error found when get verification from repo", logOption.Error(err))
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
		s.log.Error("error found when update verification", logOption.Error(err))
		return errx.Trace(err)
	}
	return nil
}

func (s *Service) CheckOTPForgetPin(payload *dto.CheckOTPPinPayload) error {
	// Get customer id
	customer, err := s.repo.FindCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		s.log.Error("error found when get customer repo", logOption.Error(err))
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
		s.log.Error("error found when get reset pin otp validate", logOption.Error(err))
		return errx.Trace(err)
	}

	// Get from cache when switching response is not success
	if switchingResponse.ResponseCode != "00" {
		s.log.Error("error rest switching otp reset pin")
		return constant.IncorrectOTPError
	}

	return nil
}

func (s *Service) ForgetPin(payload *dto.ForgetPinPayload) error {
	// Get user by ref id
	customer, err := s.repo.FindCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		s.log.Error("error found when get customer repo", logOption.Error(err))
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
		s.log.Error("error found when check otp forget pin", logOption.Error(err))
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
		s.log.Error("error found when update customer pin", logOption.Error(err))
		return errx.Trace(err)
	}

	return nil
}

func (s *Service) findOrFailCustomerByUserRefID(userRefID string) (*model.Customer, error) {
	customer, err := s.repo.FindCustomerByUserRefID(userRefID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("error found when find customer by userRefID", logOption.Error(err))
		return nil, errx.Trace(err)
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, constant.CustomerNotFoundError.Trace(nhttp.OverrideMessage("Customer not found"))
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
		s.log.Error("error found when call send OTP service", logOption.Error(err))
		return errx.Trace(err)
	}

	if resp.ResponseCode != "00" {
		s.log.Error("error rest switching otp reset pin")
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
		s.log.Error("error found when call send OTP service", logOption.Error(err))
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
		s.log.Error("error found when get customer repo", logOption.Error(err))
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
		s.log.Error("error found when get credential repo", logOption.Error(err))
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
		s.log.Error("error when update credential", logOption.Error(err))
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

	// Find customer
	customer, err := s.repo.FindCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		s.log.Error("error when find current customer", logOption.Error(err))
		err = handleErrorRepository(err, constant.ResourceNotFoundError)
		return errx.Trace(err)
	}

	// Handle if email doesn't change
	if customer.Email == payload.Email {
		return nil
	}

	// Check if email is available
	isExist, err := s.repo.EmailIsExists(payload.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("error when check email is exists", logOption.Error(err))
		return errx.Trace(err)
	}

	if isExist {
		s.log.Debug("Email has been registered")
		return constant.EmailHasBeenRegisteredError.Trace()
	}

	// Find verification
	verification, err := s.repo.FindVerificationByCustomerID(customer.ID)
	if err != nil {
		s.log.Error("error when find current customer", logOption.Error(err))
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
		s.log.Error("error when update email", logOption.Error(err))
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
		s.log.Error("error when update verification email", logOption.Error(err))
		return errx.Trace(err)
	}

	// TODO: Send Notification Email Change

	// Synchronize to PDS
	reqBody := map[string]interface{}{
		"email":                    payload.Email,
		"email_verification_token": verification.EmailVerificationToken,
		"email_verified":           verification.EmailVerifiedStatus,
	}

	// Sync customer
	resp, err := s.SynchronizeCustomer(reqBody)
	if err != nil {
		s.log.Error("sync data when ChangeEmail", logOption.Error(err))
		return errx.Trace(err)
	}

	// handle status error
	if resp.Status != constant.ResponseSuccess {
		s.log.Error("Get Error from SynchronizeCustomer")
		return nhttp.InternalError.Trace(errx.Errorf(resp.Message))
	}

	return nil
}

func (s *Service) PostUpdateSmartAccess(payload dto.UpdateSmartAccessPayload) error {
	// Get customer
	customer, err := s.repo.FindCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		s.log.Error("error when find current customer", logOption.Error(err))
		err = handleErrorRepository(err, constant.ResourceNotFoundError)
		return errx.Trace(err)
	}

	// Get credential
	credential, err := s.repo.FindCredentialByCustomerID(customer.ID)
	if err != nil {
		s.log.Error("error when find current credential", logOption.Error(err))
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
		s.log.Error("error when activate biometric", logOption.Error(err))
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
		s.log.Error("error when unset biometric", logOption.Error(err))
		return errx.Trace(err)
	}

	return nil
}

func (s *Service) GetSmartAccessStatus(payload dto.GetSmartAccessStatusPayload) (*dto.GetSmartAccessStatusResult, error) {
	// Get customer
	customer, err := s.repo.FindCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		s.log.Error("error when find current customer", logOption.Error(err))
		err = handleErrorRepository(err, constant.ResourceNotFoundError)
		return nil, errx.Trace(err)
	}

	// Get credential
	credential, err := s.repo.FindCredentialByCustomerID(customer.ID)
	if err != nil {
		s.log.Error("error when find current credential", logOption.Error(err))
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
		s.log.Error("error found when sync data customer via API PDS", logOption.Error(err))
		return errx.Trace(err)
	}

	// handle status error
	if resp.Status != constant.ResponseSuccess {
		s.log.Error("Get Error from SynchronizePassword")
		return nhttp.InternalError.Trace(errx.Errorf(resp.Message))
	}

	return nil
}

func (s *Service) HandleSynchronizeVerification(customer *model.Customer, verification *model.Verification) error {
	requestBodySyncCustomer := map[string]interface{}{
		"no_hp":customer.Phone,
		"email_verified": verification.EmailVerifiedStatus,
	}

	resp, err := s.SynchronizeCustomer(requestBodySyncCustomer)
	if err != nil {
		s.log.Error("error found when sync data customer via API PDS", logOption.Error(err))
		return errx.Trace(err)
	}

	// handle status error
	if resp.Status != constant.ResponseSuccess {
		s.log.Error("Get Error from SynchronizePassword")
		return nhttp.InternalError.Trace(errx.Errorf(resp.Message))
	}

	return nil
}

func (s *Service) PutSynchronizeCustomer(payload dto.PutSynchronizeCustomerPayload) (*dto.PutSynchronizeCustomerResult, error) {

	// Get customer
	customer, err := s.repo.FindCustomerByPhone(payload.Customer.Phone)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("error when find current customer", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	// Check on external database
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		customer, err = s.handleCheckDataOnExternal(payload.Customer.Phone)
	}
	// Handle error external database
	if err != nil {
		s.log.Error("error when find customer", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	// Get credential
	credential, err := s.repo.FindCredentialByCustomerID(customer.ID)
	if err != nil {
		s.log.Error("error when find credential", logOption.Error(err))
		err = handleErrorRepository(err, constant.CredentialNotFoundError)
		return nil, errx.Trace(err)
	}

	// Get financial data
	financial, err := s.repo.FindFinancialDataByCustomerID(customer.ID)
	if err != nil {
		s.log.Error("error when find financial data", logOption.Error(err))
		err = handleErrorRepository(err, constant.FinancialNotFoundError)
		return nil, errx.Trace(err)
	}

	// Get verification data
	verification, err := s.repo.FindVerificationByCustomerID(customer.ID)
	if err != nil {
		s.log.Error("error when find verification", logOption.Error(err))
		err = handleErrorRepository(err, constant.VerificationNotFoundError)
		return nil, errx.Trace(err)
	}

	// Get address data
	address, err := s.repo.FindAddressByCustomerId(customer.ID)
	if err != nil {
		s.log.Error("error when find address", logOption.Error(err))
		err = handleErrorRepository(err, constant.AddressNotFoundError)
		return nil, errx.Trace(err)
	}

	// Prepare model customer
	customer = prepareModelCustomerSync(customer, payload.Customer)

	// Persist customer
	err = s.repo.InsertOrUpdateCustomer(customer)
	if err != nil {
		s.log.Error("failed persist to customer", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	if payload.Credential != nil {
		credential, err = s.handleCredentialUpdate(credential, payload.Credential)
	}
	// Handle error persist credential
	if err != nil {
		s.log.Error("failed persist credential", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	if payload.Financial != nil {
		financial, err = s.handleFinancialUpdate(financial, payload.Financial)
	}
	// Handle error persist financial
	if err != nil {
		s.log.Error("failed persist financial", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	if payload.Verification != nil {
		verification, err = s.handleVerificationUpdate(verification, payload.Verification)
	}
	// Handle error persist verification
	if err != nil {
		s.log.Error("failed persist verification", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	if payload.Address != nil {
		address, err = s.handleAddressUpdate(address, payload.Address)
	}
	// Handle error persist address
	if err != nil {
		s.log.Error("failed persist address", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	return s.composeCustomerSynchronizeResult(&model.PostSynchronizeCustomerModel{
		Customer:     customer,
		Credential:   credential,
		Financial:    financial,
		Verification: verification,
		Address:      address,
	})
}

func prepareModelCustomerSync(customer *model.Customer, customerUpdate *dto.CustomerSyncVO) *model.Customer {
	// Update customer model
	customer.FullName = customerUpdate.FullName
	customer.Phone = customerUpdate.Phone
	customer.Email = customerUpdate.Email
	customer.IdentityType = customerUpdate.IdentityType
	customer.IdentityNumber = customerUpdate.IdentityNumber
	customer.Cif = customerUpdate.Cif
	customer.Sid = customerUpdate.Sid
	customer.ReferralCode = sql.NullString{
		String: customerUpdate.ReferralCode,
		Valid:  true,
	}
	customer.Photos = &model.CustomerPhoto{
		Xid:      xid.New().String(),
		FileName: customerUpdate.Photos.FileName,
		FileSize: customerUpdate.Photos.FileSize,
		Mimetype: customerUpdate.Photos.MimeType,
	}
	customer.Status = customerUpdate.Status
	customer.UpdatedAt = time.Now()
	customer.Version++
	// customer.Photos = customerUpdate.Photos TODO: update customer photo

	// Customer profile
	profilePayload := customerUpdate.Profile
	profileUpdate := &model.CustomerProfile{
		MaidenName:         profilePayload.MaidenName,
		Gender:             profilePayload.Gender,
		Nationality:        profilePayload.Nationality,
		DateOfBirth:        profilePayload.DateOfBirth,
		PlaceOfBirth:       profilePayload.PlaceOfBirth,
		IdentityPhotoFile:  profilePayload.IdentityPhotoFile,
		MarriageStatus:     profilePayload.MarriageStatus,
		NPWPNumber:         profilePayload.NPWPNumber,
		NPWPPhotoFile:      profilePayload.NPWPPhotoFile,
		NPWPUpdatedAt:      profilePayload.NPWPUpdatedAt,
		ProfileUpdatedAt:   profilePayload.ProfileUpdatedAt,
		CifLinkUpdatedAt:   profilePayload.CifLinkUpdatedAt,
		CifUnlinkUpdatedAt: profilePayload.CifUnlinkUpdatedAt,
		SidPhotoFile:       profilePayload.SidPhotoFile,
		Religion:           profilePayload.Religion,
	}
	customer.Profile = profileUpdate

	return customer
}

func (s *Service) prepareModelCredentialSync(credential *model.Credential, credentialUpdate *dto.CredentialSyncVO) (*model.Credential, error) {

	credential.Pin = credentialUpdate.Pin
	credential.Password = credentialUpdate.Password
	credential.NextPasswordResetAt = sql.NullTime{
		Time:  time.Unix(credentialUpdate.NextPasswordResetAt, 0),
		Valid: true,
	}
	credential.Pin = credentialUpdate.Pin
	credential.PinUpdatedAt = sql.NullTime{
		Time:  time.Unix(credentialUpdate.PinUpdatedAt, 0),
		Valid: true,
	}
	credential.PinLastAccessAt = sql.NullTime{
		Time:  time.Unix(credentialUpdate.PinLastAccessAt, 0),
		Valid: true,
	}
	credential.PinCounter = credentialUpdate.PinCounter
	credential.PinBlockedStatus = credentialUpdate.PinBlockedStatus
	credential.IsLocked = credentialUpdate.IsLocked
	credential.LoginFailCount = credentialUpdate.LoginFailCount
	credential.WrongPasswordCount = credentialUpdate.WrongPasswordCount
	credential.BlockedAt = sql.NullTime{
		Time:  time.Unix(credentialUpdate.BlockedAt, 0),
		Valid: true,
	}
	credential.BlockedUntilAt = sql.NullTime{
		Time:  time.Unix(credentialUpdate.BlockedUntilAt, 0),
		Valid: true,
	}
	credential.BiometricLogin = constant.ControlStatus(credentialUpdate.BiometricLogin)
	credential.BiometricDeviceID = credentialUpdate.BiometricDeviceID

	// Credential metadata
	payloadCredentialMetadata := credentialUpdate.Metadata
	credentialMetadata, err := json.Marshal(dto.MetadataCredentialVO{
		TryLoginAt:   payloadCredentialMetadata.TryLoginAt,
		PinCreatedAt: payloadCredentialMetadata.PinCreatedAt,
		PinBlockedAt: payloadCredentialMetadata.PinBlockedAt,
	})
	if err != nil {
		s.log.Error("error when find current customer", logOption.Error(err))
		return nil, errx.Trace(err)
	}
	credential.Metadata = credentialMetadata
	credential.UpdatedAt = time.Now()
	credential.Version++

	return credential, err
}

func prepareModelFinancialSync(financial *model.FinancialData, financialUpdate *dto.FinancialSyncVO) *model.FinancialData {
	financial.MainAccountNumber = financialUpdate.MainAccountNumber
	financial.AccountNumber = financialUpdate.AccountNumber
	financial.GoldSavingStatus = financialUpdate.GoldSavingStatus
	financial.GoldCardApplicationNumber = financialUpdate.GoldCardApplicationNumber
	financial.GoldCardAccountNumber = financialUpdate.GoldCardAccountNumber
	financial.Balance = financialUpdate.Balance
	financial.UpdatedAt = time.Now()
	financial.Version++
	return financial
}

func prepareModelVerificationSync(verification *model.Verification, verificationUpdate *dto.VerificationSyncVO) *model.Verification {
	verification.KycVerifiedStatus = verificationUpdate.KycVerifiedStatus
	verification.EmailVerificationToken = verificationUpdate.EmailVerificationToken
	verification.EmailVerifiedStatus = verificationUpdate.EmailVerifiedStatus
	verification.DukcapilVerifiedStatus = verificationUpdate.DukcapilVerifiedStatus
	verification.FinancialTransactionStatus = constant.ControlStatus(verificationUpdate.FinancialTransactionStatus)
	verification.FinancialTransactionActivatedAt = sql.NullTime{
		Time:  time.Unix(verificationUpdate.FinancialTransactionActivatedAt, 0),
		Valid: true,
	}
	verification.UpdatedAt = time.Now()
	verification.Version++
	return verification
}

func prepareAddressModelSync(address *model.Address, addressUpdate *dto.AddressSyncVO) *model.Address {
	address.Purpose = addressUpdate.Purpose
	address.ProvinceID = sql.NullInt64{
		Int64: addressUpdate.ProvinceID,
		Valid: true,
	}
	address.ProvinceName = sql.NullString{
		String: addressUpdate.ProvinceName,
		Valid:  true,
	}
	address.CityID = sql.NullInt64{
		Int64: addressUpdate.CityID,
		Valid: true,
	}
	address.CityName = sql.NullString{
		String: addressUpdate.CityName,
		Valid:  true,
	}
	address.DistrictID = sql.NullInt64{
		Int64: addressUpdate.DistrictID,
		Valid: true,
	}
	address.DistrictName = sql.NullString{
		String: addressUpdate.DistrictName,
		Valid:  true,
	}
	address.SubDistrictID = sql.NullInt64{
		Int64: addressUpdate.SubDistrictID,
		Valid: true,
	}
	address.SubDistrictName = sql.NullString{
		String: addressUpdate.SubDistrictName,
		Valid:  true,
	}
	address.Line = sql.NullString{
		String: addressUpdate.Line,
		Valid:  true,
	}
	address.PostalCode = sql.NullString{
		String: addressUpdate.PostalCode,
		Valid:  true,
	}
	address.IsPrimary = sql.NullBool{
		Bool:  addressUpdate.IsPrimary,
		Valid: true,
	}
	return address
}

func (s *Service) composeCustomerSynchronizeResult(mSynchronize *model.PostSynchronizeCustomerModel) (*dto.PutSynchronizeCustomerResult, error) {

	credential, err := model.ToCredentialSyncVO(mSynchronize.Credential)
	if err != nil {
		s.log.Error("error when find current customer", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	return &dto.PutSynchronizeCustomerResult{
		Customer:     model.ToCustomerSyncVO(mSynchronize.Customer),
		Financial:    model.ToFinancialSyncVO(mSynchronize.Financial),
		Credential:   credential,
		Verification: model.ToVerificationSyncVO(mSynchronize.Verification),
		Address:      model.ToAddressSyncVO(mSynchronize.Address),
	}, nil
}

func (s *Service) handleCredentialUpdate(credential *model.Credential, payload *dto.CredentialSyncVO) (*model.Credential, error) {
	// Prepare model credential
	credential, err := s.prepareModelCredentialSync(credential, payload)
	if err != nil {
		s.log.Error("error when prepare credential model", logOption.Error(err))
		return nil, errx.Trace(err)
	}
	// Persist credential
	err = s.repo.InsertOrUpdateCredential(credential)
	if err != nil {
		s.log.Error("failed persist to credential", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	return credential, nil
}

func (s *Service) handleFinancialUpdate(financial *model.FinancialData, financialUpdate *dto.FinancialSyncVO) (*model.FinancialData, error) {
	// Prepare model financial
	financial = prepareModelFinancialSync(financial, financialUpdate)
	// Persist financial data
	err := s.repo.InsertOrUpdateFinancialData(financial)
	if err != nil {
		s.log.Error("failed persist to financial data", logOption.Error(err))
		return nil, errx.Trace(err)
	}
	return financial, nil
}

func (s *Service) handleVerificationUpdate(verification *model.Verification, verificationUpdate *dto.VerificationSyncVO) (*model.Verification, error) {
	// Prepare model verification
	verification = prepareModelVerificationSync(verification, verificationUpdate)
	// persist verification
	err := s.repo.InsertOrUpdateVerification(verification)
	if err != nil {
		s.log.Error("failed persist verification.", logOption.Error(err))
		return nil, errx.Trace(err)
	}
	return verification, nil
}

func (s *Service) handleAddressUpdate(address *model.Address, addressUpdate *dto.AddressSyncVO) (*model.Address, error) {
	// Prepare model address
	address = prepareAddressModelSync(address, addressUpdate)
	// Persist address
	err := s.repo.InsertOrUpdateAddress(address)
	if err != nil {
		s.log.Error("failed persist to address", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	return address, nil
}

func (s *Service) handleCheckDataOnExternal(phone string) (*model.Customer, error) {
	// If data not found on internal database check on external database.
	user, err := s.repoExternal.FindUserExternalByEmailOrPhone(phone)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("error found when query find by email or phone", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	if errors.Is(err, sql.ErrNoRows) {
		s.log.Debug("Phone or email is not registered")
		return nil, constant.NoPhoneEmailError.Trace()
	}

	// sync data external to internal
	customer, err := s.syncExternalToInternal(user)
	if err != nil {
		s.log.Error("error while sync data External to Internal", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	return customer, nil
}
