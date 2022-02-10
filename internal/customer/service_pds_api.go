package customer

import (
	"net/http"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
)

func (s *Service) SynchronizeCustomer(payload dto.RegisterNewCustomer) (*http.Response, error) {
	var result *http.Response
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

	// Send email
	resp, err := s.ClientPostData("/synchronize/customer", reqBody, reqHeader)
	if err != nil {
		log.Errorf("Error when registration new submit. %v", err)
		return resp, err
	}

	// Set result
	result = resp

	return result, nil
}
