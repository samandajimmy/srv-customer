package customer

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/imdario/mergo"
	"github.com/nbs-go/errx"
	"github.com/nbs-go/nlogger"
	"net/http"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nclient"
)

var cacheKeyRestSwitching = fmt.Sprintf("%s:%s", constant.Prefix, constant.CacheTokenSwitching)

// Rest Switching Section Start

type ResponseSwitchingError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type ResponseSwitchingSuccess struct {
	ResponseCode string `json:"responseCode"`
	ResponseDesc string `json:"responseDesc"`
	Message      string `json:"data"`
}

type PostDataPayload struct {
	Url     string
	Data    map[string]interface{}
	Header  *map[string]string
	Counter int64
}

func (s *Service) clientRestSwitching() *nclient.Nclient {
	s.client.ClientId = s.config.CoreClientID
	s.client.BaseUrl = s.config.CoreAPIURL

	return s.client
}

func (s *Service) RestSwitchingPostData(payload PostDataPayload) (*ResponseSwitchingSuccess, error) {
	// Get context
	ctx := s.ctx
	// Default counter for retrying hit to url
	payload.Counter = 0

	// Get access token
	token, err := s.restSwitchingGetToken()
	if err != nil {
		return nil, errx.Trace(err)
	}

	// Set request body
	reqBody := map[string]interface{}{
		"channelId": constant.ChannelMobile,
		"clientId":  s.config.CoreOauthUsername,
	}

	// Merge Request Body
	err = mergo.Merge(&reqBody, payload.Data)
	if err != nil {
		s.log.Error("error when merge request body", nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	// Set request header
	reqHeader := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
		"Accept":        "application/json",
		"Content-Type":  "application/json",
	}

	if payload.Header != nil {
		// Merge Request Header
		err = mergo.Merge(&reqHeader, payload.Header)
		if err != nil {
			return nil, errx.Trace(err)
		}
	}

	// Sent Request Rest Switching
	restResponse, err := s.clientRestSwitching().PostData(payload.Url, reqBody, reqHeader)
	if err != nil {
		return nil, errx.Trace(err)
	}
	defer handleClose(restResponse.Body)

	var invalidResponse string
	if restResponse.StatusCode == 401 {
		responseSwitching, err := s.restSwitchingErrorResponse(restResponse)
		if err != nil {
			return nil, err
		}
		invalidResponse = responseSwitching.Error
	}

	responseSwitching, err := s.restSwitchingSuccessResponse(restResponse)
	if err != nil {
		return nil, err
	}

	payload.Counter++
	responseAfterRefreshToken, isRefreshSuccess := s.restSwitchingRefreshToken(invalidResponse, payload)
	if isRefreshSuccess {
		return responseAfterRefreshToken, nil
	}

	return responseSwitching, nil
}

func (s *Service) restSwitchingSuccessResponse(restResponse *http.Response) (*ResponseSwitchingSuccess, error) {
	var responseSwitching *ResponseSwitchingSuccess
	err := json.NewDecoder(restResponse.Body).Decode(&responseSwitching)
	if err != nil {
		s.log.Error("error when get response rest switching.", nlogger.Error(err), nlogger.Context(s.ctx))
		return nil, errx.Trace(err)
	}

	return responseSwitching, nil
}

func (s *Service) restSwitchingErrorResponse(restResponse *http.Response) (*ResponseSwitchingError, error) {
	var responseSwitching *ResponseSwitchingError
	err := json.NewDecoder(restResponse.Body).Decode(&responseSwitching)
	if err != nil {
		s.log.Error("error when get response rest switching.", nlogger.Error(err), nlogger.Context(s.ctx))
		return nil, errx.Trace(err)
	}

	return responseSwitching, nil
}

// restSwitchingGetToken Get Access Token from rest switching service (PDS)
func (s *Service) restSwitchingGetToken() (string, error) {
	token, err := s.CacheGet(cacheKeyRestSwitching)
	if err != nil {
		return "", errx.Trace(err)
	}

	if token != "" {
		return token, nil
	}

	newToken, err := s.restSwitchingNewToken()
	if err != nil {
		return "", errx.Trace(err)
	}

	return newToken, nil
}

// restSwitchingNewToken Generate New access token from rest switching service (PDS)
func (s *Service) restSwitchingNewToken() (string, error) {
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
	resp, err := s.clientRestSwitching().PostData("/oauth/token", reqBody, reqHeader)
	if err != nil {
		s.log.Error("Error when request oauth token")
		return result, errx.Trace(err)
	}

	// Decode response body from server.
	var data dto.TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		s.log.Error("Error when decode response data get token.")
		return result, err
	}

	// Store token to redis
	result, err = s.CacheSetThenGet(cacheKeyRestSwitching, data.AccessToken, data.ExpiresIn)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (s *Service) restSwitchingRefreshToken(responseError string, payload PostDataPayload) (*ResponseSwitchingSuccess, bool) {
	if responseError == constant.RestSwitchingInvalidToken {
		_, _ = s.restSwitchingNewToken()
		resp, _ := s.RestSwitchingPostData(payload)
		return resp, true
	}

	return nil, false
}

// Rest Switching Section End

// PDS API Section Start

type ResponsePdsAPI struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func (s *Service) clientPDS() *nclient.Nclient {
	return s.pdsClient
}

func (s *Service) PdsPostData(payload PostDataPayload) (*ResponsePdsAPI, error) {
	// Get context
	ctx := s.ctx

	// Set request body
	reqBody := map[string]interface{}{}

	// Merge Request Body
	err := mergo.Merge(&reqBody, payload.Data)
	if err != nil {
		s.log.Error("error when merge request body", nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	// Set request header
	reqHeader := map[string]string{}

	if payload.Header != nil {
		// Merge Request Header
		err = mergo.Merge(&reqHeader, payload.Header)
		if err != nil {
			return nil, errx.Trace(err)
		}
	}

	// Sent Request Rest Switching
	restResponse, err := s.clientPDS().PostData(payload.Url, reqBody, reqHeader)
	if err != nil {
		return nil, errx.Trace(err)
	}
	defer handleClose(restResponse.Body)

	response, err := s.pdsAPISuccessResponse(restResponse)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Service) pdsAPISuccessResponse(restResponse *http.Response) (*ResponsePdsAPI, error) {
	var responseSwitching *ResponsePdsAPI
	err := json.NewDecoder(restResponse.Body).Decode(&responseSwitching)
	if err != nil {
		s.log.Error("error when get response rest switching.", nlogger.Error(err), nlogger.Context(s.ctx))
		return nil, errx.Trace(err)
	}

	return responseSwitching, nil
}

// PDS API Section End

type NotificationPostDataPayload struct {
	URL    string
	Data   map[string]interface{}
	Header *map[string]string
}

// Notification Service Section Start
func (s *Service) clientNotification() *nclient.Nclient {
	s.client.BaseUrl = s.config.NotificationServiceURL

	return s.client
}

func (s *Service) CreateNotificationPostData(payload *NotificationPostDataPayload) (*http.Response, error) {
	rawBasicAuth := fmt.Sprintf("%s:%s", s.config.NotificationServiceAppXid, s.config.NotificationServiceAppAPIKey)
	basicAuth := base64.StdEncoding.EncodeToString([]byte(rawBasicAuth))

	// Set request header
	reqHeader := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", basicAuth),
		"Accept":        "application/json",
		"Content-Type":  "application/json",
	}

	if payload.Header != nil {
		// Merge Request Header
		err := mergo.Merge(&reqHeader, payload.Header)
		if err != nil {
			return nil, errx.Trace(err)
		}
	}

	_, err := s.clientNotification().PostData(payload.URL, payload.Data, reqHeader)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Notification Service Section End
