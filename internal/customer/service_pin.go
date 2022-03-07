package customer

import (
	"github.com/nbs-go/nlogger/v2"
	"regexp"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
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
