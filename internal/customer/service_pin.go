package customer

import (
	"database/sql"
	"errors"
	"github.com/nbs-go/errx"
	"github.com/nbs-go/nlogger/v2"
	"regexp"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
	"time"
)

func (s *Service) ValidatePin(payload *dto.ValidatePinPayload) (string, error) {
	// Get Context
	ctx := s.ctx

	// Validate Increasing Pin
	err := validateIncreasingPin(payload)
	if err != nil {
		s.log.Error(`new pin isn't passed validation increasing pin'`, nlogger.Error(err), nlogger.Context(ctx))
		return "", errx.Trace(err)
	}

	err = validateRepetitionPin(payload)
	if err != nil {
		s.log.Error(`new pin isn't passed validation increasing pin'`, nlogger.Error(err), nlogger.Context(ctx))
		return "", errx.Trace(err)
	}

	return "PIN sudah valid", nil
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

func (s *Service) CheckPinUser(payload *dto.CheckPinPayload) (string, error) {
	// Get Context
	ctx := s.ctx

	// Get customer id
	customer, err := s.repo.FindCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		s.log.Error("error found when get customer repo", nlogger.Error(err), nlogger.Context(ctx))
		if errors.Is(err, sql.ErrNoRows) {
			return "", constant.ResourceNotFoundError
		}
		return "", errx.Trace(err)
	}

	// Get customer pin
	credential, err := s.repo.FindCredentialByCustomerID(customer.ID)
	if err != nil || credential.Pin == "" {
		s.log.Error("error found when get credential repo", nlogger.Context(ctx))
		if errors.Is(err, sql.ErrNoRows) || credential.Pin == "" {
			return "", constant.AccountPINIsNotActive
		}
		return "", errx.Trace(err)
	}

	// Check if pin is blocked
	if credential.PinBlockedStatus != constant.Unblocked {
		s.log.Error("pin is blocked", nlogger.Context(ctx))
		return "", constant.AccountPINIsBlocked
	}

	// Find pin with md5
	valid, err := s.repo.IsValidPin(customer.ID, nval.MD5(payload.Pin))
	if err != nil {
		s.log.Error("error found when querying valid pin", nlogger.Error(err), nlogger.Context(ctx))
		return "", errx.Trace(err)
	}

	if !valid {
		err = s.handlePinNotValid(credential)
		return "", errx.Trace(err)
	}

	// Unblock pin user
	err = s.unblockPINUser(customer)
	if err != nil {
		return "", errx.Trace(err)
	}

	return "PIN Tersedia dan sudah di aktifasi finansial", nil
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
		if errors.Is(err, sql.ErrNoRows) {
			return constant.ResourceNotFoundError
		}
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

func (s *Service) UpdatePin(payload *dto.UpdatePinPayload) (*dto.UpdatePinResult, string, error) {
	// If check pin is true
	if payload.CheckPIN { // Prepare payload check pin
		payloadCheckPIN := &dto.CheckPinPayload{
			Pin: payload.PIN,
		}
		// Check PIN user
		_, err := s.CheckPinUser(payloadCheckPIN)
		if err != nil {
			return nil, "", errx.Trace(err)
		}
	}

	// Validate pin user
	_, err := s.ValidatePin(&dto.ValidatePinPayload{
		NewPin: payload.NewPIN,
	})
	if err != nil {
		return nil, "", errx.Trace(err)
	}

	err = s.handleUpdatePin(payload.UserRefID, payload.NewPIN)
	if err != nil {
		s.log.Error("error found when handle update pin", nlogger.Context(s.ctx))
		return nil, "", errx.Trace(err)
	}

	// TODO: Audit log reset pin

	return &dto.UpdatePinResult{
		Title: "PIN Berhasil Diubah",
		Text:  "Selamat! Kamu berhasil mengubah PIN Pegadaian Digital.",
	}, "PIN Berhasil Diubah", nil
}

func (s *Service) handleUpdatePin(userRefID string, newPin string) error {
	// Get customer id
	customer, err := s.repo.FindCustomerByUserRefID(userRefID)
	if err != nil {
		s.log.Error("error found when get customer repo", nlogger.Error(err), nlogger.Context(s.ctx))
		if errors.Is(err, sql.ErrNoRows) {
			return constant.ResourceNotFoundError
		}
		return errx.Trace(err)
	}

	// Get customer pin
	credential, err := s.repo.FindCredentialByCustomerID(customer.ID)
	if err != nil || credential.Pin == "" {
		s.log.Error("error found when get credential repo", nlogger.Context(s.ctx))
		if errors.Is(err, sql.ErrNoRows) || credential.Pin == "" {
			return constant.AccountPINIsNotActive
		}
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

func (s *Service) CheckOTPPinCreate(payload *dto.CheckOTPPinPayload) (string, error) {
	// Get customer id
	customer, err := s.repo.FindCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		s.log.Error("error found when get customer repo", nlogger.Error(err), nlogger.Context(s.ctx))
		if errors.Is(err, sql.ErrNoRows) {
			return "", constant.ResourceNotFoundError
		}
		return "", errx.Trace(err)
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
		return "", errx.Trace(err)
	}

	// Get from cache when switching response is not success
	if switchingResponse.ResponseCode != "00" {
		return "", constant.IncorrectOTPError
	}

	return "OTP yang dimasukan valid!", nil
}

func (s *Service) CreatePinUser(payload *dto.PostCreatePinPayload) (string, error) {
	// Get user by ref id
	customer, err := s.repo.FindCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		s.log.Error("error found when get customer repo", nlogger.Error(err), nlogger.Context(s.ctx))
		if errors.Is(err, sql.ErrNoRows) {
			return "", constant.ResourceNotFoundError
		}
		return "", errx.Trace(err)
	}

	// Check OTP
	msg, err := s.CheckOTPPinCreate(&dto.CheckOTPPinPayload{
		OTP:       payload.OTP,
		UserRefID: payload.UserRefID,
	})
	if err != nil {
		return msg, errx.Trace(err)
	}

	// Unblock pin user
	err = s.unblockPINUser(customer)
	if err != nil {
		return "", errx.Trace(err)
	}

	// Update pin
	_, _, err = s.UpdatePin(&dto.UpdatePinPayload{
		UserRefID: customer.UserRefID.String,
		NewPIN:    payload.NewPIN,
		CheckPIN:  false,
	})
	if err != nil {
		s.log.Error("error found when update customer pin", nlogger.Error(err), nlogger.Context(s.ctx))
		return "", errx.Trace(err)
	}

	// Aktivasi finansial
	err = s.activateFinancialStatus(customer)
	if err != nil {
		return "", errx.Trace(err)
	}

	// TODO: Audit log aktifasi finansial

	return "PIN sudah valid", nil
}

func (s *Service) activateFinancialStatus(customer *model.Customer) error {
	// Get verification
	verification, err := s.repo.FindVerificationByCustomerID(customer.ID)
	if err != nil {
		s.log.Error("error found when get verification from repo", nlogger.Error(err), nlogger.Context(s.ctx))
		if errors.Is(err, sql.ErrNoRows) {
			return constant.ResourceNotFoundError
		}
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

func (s *Service) CheckOTPForgetPin(payload *dto.CheckOTPPinPayload) (string, error) { // Get customer id
	customer, err := s.repo.FindCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		s.log.Error("error found when get customer repo", nlogger.Error(err), nlogger.Context(s.ctx))
		if errors.Is(err, sql.ErrNoRows) {
			return "", constant.ResourceNotFoundError
		}
		return "", errx.Trace(err)
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
		return "", errx.Trace(err)
	}

	// Get from cache when switching response is not success
	if switchingResponse.ResponseCode != "00" {
		s.log.Error("error rest switching otp reset pin", nlogger.Context(s.ctx))
		return "", constant.IncorrectOTPError
	}

	return "OTP yang dimasukan valid!", nil
}

func (s *Service) ForgetPin(payload *dto.ForgetPinPayload) (string, error) {
	// Get user by ref id
	customer, err := s.repo.FindCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		s.log.Error("error found when get customer repo", nlogger.Error(err), nlogger.Context(s.ctx))
		if errors.Is(err, sql.ErrNoRows) {
			return "", constant.ResourceNotFoundError
		}
		return "", errx.Trace(err)
	}

	// Check OTP Reset PIN
	payloadOTP := &dto.CheckOTPPinPayload{
		OTP:       payload.OTP,
		UserRefID: payload.UserRefID,
	}

	msg, err := s.CheckOTPForgetPin(payloadOTP)
	if err != nil {
		s.log.Error("error found when check otp forget pin", nlogger.Error(err), nlogger.Context(s.ctx))
		return msg, errx.Trace(err)
	}

	// Unblock pin user
	err = s.unblockPINUser(customer)
	if err != nil {
		return "", errx.Trace(err)
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
	_, _, err = s.UpdatePin(payloadUpdatePIN)
	if err != nil {
		s.log.Error("error found when update customer pin", nlogger.Error(err), nlogger.Context(s.ctx))
		return "", errx.Trace(err)
	}

	return "PIN berhasil diubah", nil
}
