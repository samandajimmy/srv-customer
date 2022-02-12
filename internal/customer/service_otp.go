package customer

import (
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
)

func (s *Service) SendOTP(payload dto.SendOTPRequest) (*ResponseSwitchingSuccess, error) {
	// Set payload
	reqBody := map[string]interface{}{
		"noHp":        payload.PhoneNumber,
		"requestType": payload.RequestType,
	}

	sp := SwitchingPostDataPayload{
		Url:  "/otp/send",
		Data: reqBody,
	}

	// send otp
	r, err := s.RestSwitchingPostData(sp)
	if err != nil {
		s.log.Error("failed when send OTP", nlogger.Error(err), nlogger.Context(s.ctx))
		return nil, ncore.TraceError("failed to send otp", err)
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

	sp := SwitchingPostDataPayload{
		Url:  "/otp/check",
		Data: reqBody,
	}

	// Send OTP Rest Switching
	r, err := s.RestSwitchingPostData(sp)
	if err != nil {
		s.log.Error("Error when verify otp request", nlogger.Error(err), nlogger.Context(s.ctx))
		return nil, ncore.TraceError("failed to send verify otp", err)
	}

	// Set result
	result := r

	return result, nil
}
