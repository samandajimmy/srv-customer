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
		return "", err
	}

	err = validateRepetitionPin(payload)
	if err != nil {
		s.log.Error(`new pin isn't passed validation increasing pin'`, nlogger.Error(err), nlogger.Context(ctx))
		return "", err
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
		return "", err
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
		return "", err
	}

	// Unblock pin user
	err = s.unblockPINUser(customer)
	if err != nil {
		return "", err
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
			return errors.New("credential not found")
		}
		return err
	}

	// Unblocked pin
	credential.PinCounter = 0
	credential.PinBlockedStatus = constant.Unblocked

	err = s.repo.UpdateCredential(credential)
	if err != nil {
		s.log.Error("error found when update credential", nlogger.Error(err), nlogger.Context(s.ctx))
		return err
	}

	return nil
}
