package service

import (
	"net/http"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nclient"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
)

type Notification struct {
	client       *nclient.Nclient
	pdsAPI       *contract.CorePDSConfig
	cacheService contract.CacheService
	cacheKey     string
	response     *ncore.ResponseMap
}

func (c *Notification) HasInitialized() bool {
	return true
}

func (c *Notification) Init(app *contract.PdsApp) error {
	c.pdsAPI = &app.Config.CorePDS
	c.client = nclient.NewNucleoClient(
		c.pdsAPI.CoreOauthUsername,
		c.pdsAPI.CoreClientId,
		app.Config.Notification.NotificationServiceUrl,
	)
	c.response = app.Responses
	return nil
}

func (c *Notification) SendNotification(payload dto.NotificationPayload) (*http.Response, error) {
	var result *http.Response

	// Set payload
	reqBody := map[string]interface{}{
		"title": payload.Title,
		"body":  payload.Body,
		"image": payload.Image,
		"token": payload.Token,
		"data":  payload.Data,
	}

	// Set header
	reqHeader := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	// Send Notification
	resp, err := c.client.PostData("/push-notification", reqBody, reqHeader)
	if err != nil {
		log.Errorf("Error when send notification")
		return resp, ncore.TraceError(err)
	}

	// Set result
	result = resp

	return result, nil
}
