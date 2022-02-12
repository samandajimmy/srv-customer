package customer

import (
	"encoding/json"
	"fmt"
	"github.com/nbs-go/nlogger"
	"net/http"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
)

func (s *Service) GetToken() (string, error) {
	key := fmt.Sprintf("%s:%s", constant.Prefix, "token_switching")
	// Get token if is exists on redis
	token, err := s.CacheGet(key)
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
		"username":   s.config.CoreOauthUsername,
		"password":   s.config.CoreOauthPassword,
		"grant_type": s.config.CoreOauthGrantType,
	}

	// Set header
	reqHeader := map[string]string{
		"Authorization": "Basic " + s.config.CoreAuthorization,
		"Accept":        "application/json",
		"Content-Type":  "application/x-www-form-urlencoded",
	}

	// Send OTP Rest Switching
	resp, err := s.ClientPostData("/oauth/token", reqBody, reqHeader)
	if err != nil {
		s.log.Error("Error when request oauth token", nlogger.Error(err))
		return result, ncore.TraceError("error", err)
	}
	defer resp.Body.Close()

	// Decode response body from server.
	var data dto.TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		s.log.Error("Error when decode response data get token.", nlogger.Error(err))
		return result, ncore.TraceError("failed to decode response data get token", err)
	}

	cacheKey := fmt.Sprintf("%s:%s", constant.Prefix, "token_switching")
	// Store token to redis
	result, err = s.CacheSetThenGet(cacheKey, data.AccessToken, data.ExpiresIn)
	if err != nil {
		s.log.Errorf("Error when store access token to redis", nlogger.Error(err))
		return "", ncore.TraceError("failed to store access token to redis", err)
	}

	return result, nil
}

func (s *Service) SendOTP(payload dto.SendOTPRequest) (*http.Response, error) {
	var result *http.Response

	token, err := s.GetToken()
	if err != nil {
		s.log.Errorf("Error when trying to get Access Token. err: %s", nlogger.Error(err))
		return result, ncore.TraceError("failed to get access token", err)
	}

	// Set payload
	reqBody := map[string]interface{}{
		"channelId":   "6017",
		"clientId":    s.config.CoreClientID,
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
	resp, err := s.ClientPostData("/otp/send", reqBody, reqHeader)
	if err != nil {
		s.log.Error("Error when send otp to phone number", nlogger.Error(err))
		return resp, ncore.TraceError("failed to send otp", err)
	}

	// Set result
	result = resp

	return result, nil
}

func (s *Service) VerifyOTP(payload dto.VerifyOTPRequest) (*http.Response, error) {
	var result *http.Response

	token, err := s.GetToken()
	if err != nil {
		s.log.Errorf("Error when trying to get Access Token. err: %s", err)
		return result, ncore.TraceError("failed to get access token", err)
	}

	// Set payload
	reqBody := map[string]interface{}{
		"channelId":   "6017",
		"clientId":    s.config.CoreClientID,
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
	resp, err := s.ClientPostData("/otp/check", reqBody, reqHeader)
	if err != nil {
		s.log.Error("Error when verify otp request", nlogger.Error(err))
		return resp, ncore.TraceError("failed to send verify otp", err)
	}

	// Set result
	result = resp

	return result, nil
}
