package customer

import (
	"github.com/nbs-go/nlogger/v2"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
)

func (s *Service) SynchronizeCustomer(payload dto.RegisterPayload) (*ResponsePdsAPI, error) {
	// Set payload
	reqBody := map[string]interface{}{
		"nama":      payload.Name,
		"email":     payload.Email,
		"no_hp":     payload.PhoneNumber,
		"password":  payload.Password,
		"fcm_token": payload.FcmToken,
	}

	// Set header
	reqHeader := map[string]string{
		"Accept":       "*/*",
		"Content-Type": "application/x-www-form-urlencoded",
	}

	// Set payload
	postDataPayload := PostDataPayload{
		Url:    "/synchronize/customer",
		Data:   reqBody,
		Header: &reqHeader,
	}

	resp, err := s.PdsPostData(postDataPayload)
	if err != nil {
		s.log.Error("error found when sync customer to PDS API", nlogger.Error(err), nlogger.Context(s.ctx))
		return nil, err
	}

	// Set result
	result := resp

	return result, nil
}

func (s *Service) SynchronizePassword(payload dto.RegisterPayload) (*ResponsePdsAPI, error) {
	// Set payload
	reqBody := map[string]interface{}{
		"no_hp":    payload.PhoneNumber,
		"password": payload.Password,
	}

	// Set header
	reqHeader := map[string]string{
		"Accept":       "*/*",
		"Content-Type": "application/x-www-form-urlencoded",
	}

	// Set payload
	postDataPayload := PostDataPayload{
		Url:    "/synchronize/customer",
		Data:   reqBody,
		Header: &reqHeader,
	}

	resp, err := s.PdsPostData(postDataPayload)
	if err != nil {
		s.log.Error("error found when sync customer password to PDS API", nlogger.Error(err), nlogger.Context(s.ctx))
		return nil, err
	}

	// Set result
	result := resp

	return result, nil
}
