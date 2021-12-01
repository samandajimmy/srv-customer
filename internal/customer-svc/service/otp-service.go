package service

import (
	"encoding/json"
	"net/http"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nclient"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
)

type OTP struct {
	client   *nclient.Nclient
	pdsAPI   *contract.CorePDSConfig
	response *ncore.ResponseMap
}

func (c *OTP) HasInitialized() bool {
	return true
}

func (c *OTP) Init(app *contract.PdsApp) error {
	c.pdsAPI = &app.Config.CorePDS
	c.client = nclient.NewNucleoClient(
		c.pdsAPI.CORE_OAUTH_USERNAME,
		c.pdsAPI.CORE_CLIENT_ID,
		c.pdsAPI.CORE_API_URL,
	)
	c.response = app.Responses
	return nil
}

func (c *OTP) GetToken() (string, error) {

	// TODO: Get token is exist on redis

	// Initialise result
	var result string

	// Set payload
	reqBody := map[string]string{
		"username":   c.pdsAPI.CORE_OAUTH_USERNAME,
		"password":   c.pdsAPI.CORE_OAUTH_PASSWORD,
		"grant_type": c.pdsAPI.CORE_OAUTH_GRANT_TYPE,
	}

	// Set header
	reqHeader := map[string]string{
		"Authorization": "Basic " + c.pdsAPI.CORE_AUTHORIZATION,
		"Accept":        "application/json",
		"Content-Type":  "application/x-www-form-urlencoded",
	}

	// Send OTP Rest Switching
	resp, err := c.client.PostData("/oauth/token", reqBody, reqHeader)
	if err != nil {
		log.Errorf("Error when request oauth token")
		return result, ncore.TraceError(err)
	}
	defer resp.Body.Close()

	// TODO: Store token to redis

	// Decode response body from server.
	var data dto.TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Errorf("Error when decode response data get token.")
		return result, err
	}

	result = data.AccessToken

	return result, nil
}

func (c *OTP) SendOTP(payload dto.SendOTPRequest) (*http.Response, error) {
	var result *http.Response

	token, err := c.GetToken()
	if err != nil {
		log.Errorf("Error when trying to get Access Token. err: %s", err)
		return result, ncore.TraceError(err)
	}

	// Set payload
	reqBody := map[string]string{
		"channelId":   "6017",
		"clientId":    c.pdsAPI.CORE_OAUTH_USERNAME,
		"noHp":        payload.PhoneNumber,
		"requestType": payload.RequestType,
	}

	// Set header
	reqHeader := map[string]string{
		"Authorization": "Bearer " + token,
		"Accept":        "application/json",
		"Content-Type":  "application/json",
	}

	// Send OTP Rest Switching
	resp, err := c.client.PostData("/otp/send", reqBody, reqHeader)
	if err != nil {
		log.Errorf("Error when send otp to phone number")
		return resp, ncore.TraceError(err)
	}
	defer resp.Body.Close()

	// Set result
	result = resp

	return result, nil
}
