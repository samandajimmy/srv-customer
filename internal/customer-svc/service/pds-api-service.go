package service

import (
	"net/http"

	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nclient"
)

type PdsAPI struct {
	client *nclient.Nclient
	pdsAPI *contract.CorePDSConfig
}

func (s *PdsAPI) HasInitialized() bool {
	return true
}

func (s *PdsAPI) Init(app *contract.PdsApp) error {
	s.pdsAPI = &app.Config.CorePDS
	s.client = nclient.NewNucleoClient(
		s.pdsAPI.CoreOauthUsername,
		s.pdsAPI.CoreClientId,
		app.Config.ClientEndpoint.PdsApiServiceUrl,
	)
	return nil
}

func (s *PdsAPI) StepOneRegistration(payload dto.RegisterStepOne) (*http.Response, error) {
	var result *http.Response
	// Set payload
	reqBody := map[string]interface{}{
		"nama":  payload.Name,
		"email": payload.Email,
		"no_hp": payload.PhoneNumber,
	}

	// Set header
	reqHeader := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	// Send email
	resp, err := s.client.PostData("/auth/register/", reqBody, reqHeader)
	if err != nil {
		log.Errorf("Error when registration step one. %v", err)
		return resp, err
	}

	// Set result
	result = resp

	return result, nil
}

func (s *PdsAPI) StepTwoRegistration(payload dto.RegisterStepTwo) (*http.Response, error) {
	var result *http.Response
	// Set payload
	reqBody := map[string]interface{}{
		"no_hp": payload.PhoneNumber,
		"otp":   payload.OTP,
	}

	// Set header
	reqHeader := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	// Send email
	resp, err := s.client.PostData("/auth/register_new/check", reqBody, reqHeader)
	if err != nil {
		log.Errorf("Error when registration step two. %v", err)
		return resp, err
	}

	// Set result
	result = resp

	return result, nil
}

func (s *PdsAPI) Register(payload dto.RegisterNewCustomer) (*http.Response, error) {
	var result *http.Response
	// Set payload
	reqBody := map[string]interface{}{
		"nama":        payload.Name,
		"email":       payload.Email,
		"no_hp":       payload.PhoneNumber,
		"password":    payload.Password,
		"fcm_token":   payload.FcmToken,
		"register_id": payload.RegistrationId,
	}

	// Set header
	reqHeader := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	// Send email
	resp, err := s.client.PostData("/auth/register_new", reqBody, reqHeader)
	if err != nil {
		log.Errorf("Error when registration new submit. %v", err)
		return resp, err
	}

	// Set result
	result = resp

	return result, nil
}
