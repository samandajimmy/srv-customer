package customer

import (
	"github.com/nbs-go/errx"
	logOption "github.com/nbs-go/nlogger/v2/option"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
)

func (s *Service) SendOTP(payload dto.SendOTPRequest) (*ResponseSwitchingSuccess, error) {
	// Set payload
	reqBody := map[string]interface{}{
		"noHp":        payload.PhoneNumber,
		"requestType": payload.RequestType,
	}

	sp := PostDataPayload{
		Url:  "/otp/send",
		Data: reqBody,
	}

	// send otp
	r, err := s.RestSwitchingPostData(sp)
	if err != nil {
		s.log.Error("failed when send OTP", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	// Set result
	result := r

	return result, nil
}

func (s *Service) VerifyOTP(payload dto.VerifyOTPRequest) (*ResponseSwitchingSuccess, error) {
	// Set payload
	reqBody := map[string]interface{}{
		"noHp":        payload.PhoneNumber,
		"requestType": payload.RequestType,
		"token":       payload.Token,
	}

	sp := PostDataPayload{
		Url:  "/otp/check",
		Data: reqBody,
	}

	// Send OTP Rest Switching
	r, err := s.RestSwitchingPostData(sp)
	if err != nil {
		s.log.Error("Error when verify otp request", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	// Set result
	result := r

	return result, nil
}
