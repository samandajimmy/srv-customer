package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/constant"

	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nclient"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
)

type OTP struct {
	client       *nclient.Nclient
	pdsAPI       *contract.CorePDSConfig
	cacheService contract.CacheService
	cacheKey     string
	response     *ncore.ResponseMap
}

func (c *OTP) HasInitialized() bool {
	return true
}

func (c *OTP) Init(app *contract.PdsApp) error {
	c.pdsAPI = &app.Config.CorePDS
	c.client = nclient.NewNucleoClient(
		c.pdsAPI.CoreOauthUsername,
		c.pdsAPI.CoreClientId,
		c.pdsAPI.CoreApiUrl,
	)
	c.cacheService = app.Services.Cache
	c.cacheKey = fmt.Sprintf("%s:%s", constant.Prefix, "token_switching")
	c.response = app.Responses
	return nil
}

func (c *OTP) GetToken() (string, error) {

	// Get token if is exist on redis
	token, err := c.cacheService.Get(c.cacheKey)
	if err != nil {
		return "", err
	}

	if token != "" {
		return token, nil
	}

	// Initialise result
	var result string

	// Set payload
	reqBody := map[string]interface{}{
		"username":   c.pdsAPI.CoreOauthUsername,
		"password":   c.pdsAPI.CoreOauthPassword,
		"grant_type": c.pdsAPI.CoreOauthGrantType,
	}

	// Set header
	reqHeader := map[string]string{
		"Authorization": "Basic " + c.pdsAPI.CoreAuthorization,
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

	// Decode response body from server.
	var data dto.TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Errorf("Error when decode response data get token.")
		return result, err
	}

	// Store token to redis
	result, err = c.cacheService.SetThenGet(c.cacheKey, data.AccessToken, data.ExpiresIn)
	if err != nil {
		return "", err
	}

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
	reqBody := map[string]interface{}{
		"channelId":   "6017",
		"clientId":    c.pdsAPI.CoreClientId,
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

	// Set result
	result = resp

	return result, nil
}

func (c *OTP) VerifyOTP(payload dto.VerifyOTPRequest) (*http.Response, error) {
	var result *http.Response

	token, err := c.GetToken()
	if err != nil {
		log.Errorf("Error when trying to get Access Token. err: %s", err)
		return result, ncore.TraceError(err)
	}

	// Set payload
	reqBody := map[string]interface{}{
		"channelId":   "6017",
		"clientId":    c.pdsAPI.CoreClientId,
		"noHp":        payload.PhoneNumber,
		"requestType": payload.RequestType,
		"token":       payload.Token,
	}

	// Set header
	reqHeader := map[string]string{
		"Authorization": "Bearer " + token,
		"Accept":        "application/json",
		"Content-Type":  "application/json",
	}

	// Send OTP Rest Switching
	resp, err := c.client.PostData("/otp/check", reqBody, reqHeader)
	if err != nil {
		log.Errorf("Error when verify otp request")
		return resp, ncore.TraceError(err)
	}

	// Set result
	result = resp

	return result, nil
}
