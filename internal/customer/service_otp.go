package customer

import (
	"github.com/nbs-go/errx"
	"github.com/nbs-go/nlogger"
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
		s.log.Error("failed when send OTP", nlogger.Error(err), nlogger.Context(s.ctx))
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
		s.log.Error("Error when verify otp request", nlogger.Error(err), nlogger.Context(s.ctx))
		return nil, errx.Trace(err)
	}

	// Set result
	result := r

	return result, nil
}
