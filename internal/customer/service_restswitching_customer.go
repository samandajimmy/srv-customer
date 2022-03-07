package customer

import (
	"github.com/nbs-go/errx"
	"github.com/nbs-go/nlogger/v2"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
)

// Endpoint POST /customer/activation
func (s *Service) customerActivation(payload *dto.RestSwitchingOTPPinCreate) (*ResponseSwitchingSuccess, error) {
	// Set payload
	reqBody := map[string]interface{}{
		"cif":      payload.Cif,
		"token":    payload.OTP,
		"username": payload.NoHp,
	}

	sp := PostDataPayload{
		Url:  "/customer/activation",
		Data: reqBody,
	}

	data, err := s.RestSwitchingPostData(sp)
	if err != nil {
		s.log.Error("error found when get otp pin activation", nlogger.Error(err), nlogger.Context(s.ctx))
		return nil, errx.Trace(err)
	}

	return data, nil
}
