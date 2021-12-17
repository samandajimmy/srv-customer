package service

import (
	"fmt"
	"math/rand"
	"net/http"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"

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
	httpBaseUrl  string
	emailConfig  contract.EmailConfig
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
	c.httpBaseUrl = app.Config.Server.GetHttpBaseUrl()
	c.emailConfig = app.Config.Email
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

func (c *Notification) SendEmail(payload dto.EmailPayload) (*http.Response, error) {
	var result *http.Response
	// Set payload
	reqBody := map[string]interface{}{
		"from": map[string]string{
			"name":  payload.From.Name,
			"email": payload.From.Email,
		},
		"to":         payload.To,
		"subject":    payload.Subject,
		"message":    payload.Message,
		"attachment": payload.Attachment,
		"mimeType":   payload.MimeType,
	}

	// Set header
	reqHeader := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	// Send email
	resp, err := c.client.PostData("/send-email", reqBody, reqHeader)
	if err != nil {
		log.Errorf("Error when send email. %v", err)
		return resp, err
	}

	// Set result
	result = resp

	return result, nil
}

func (c *Notification) SendNotificationRegister(data dto.NotificationRegister) error {
	// Send Email Verification
	dataEmailVerification := &dto.EmailVerification{
		FullName:        data.Customer.FullName,
		Email:           data.Customer.Email,
		VerificationUrl: fmt.Sprintf("%sauth/verify_email?t=%s", c.httpBaseUrl, data.Verification.EmailVerificationToken),
	}
	htmlMessage, err := nval.TemplateFile(dataEmailVerification, "email_verification.html")
	if err != nil {
		return err
	}

	// set payload email service
	emailPayload := dto.EmailPayload{
		Subject: fmt.Sprintf("Verifikasi Email %s", data.Customer.FullName),
		From: dto.FromEmailPayload{
			Name:  c.emailConfig.PdsEmailFromName,
			Email: c.emailConfig.PdsEmailFrom,
		},
		To:         data.Customer.Email,
		Message:    htmlMessage,
		Attachment: "",
		MimeType:   "",
	}
	_, err = c.SendEmail(emailPayload)
	if err != nil {
		log.Debugf("Error when send email verification. Payload %v", emailPayload)
	}

	// Send Notification Welcome
	id, _ := nval.ParseString(rand.Intn(100)) // TODO: insert data to notification
	var dataWelcomeMessage = map[string]string{
		"title": "Verifikasi Email",
		"body":  fmt.Sprintf(`Hai %v, Selamat datang di Pegadaian Digital Service`, data.Customer.FullName),
		"type":  constant.TypeProfile,
		"id":    id,
	}
	welcomeMessage := dto.NotificationPayload{
		Title: "Verifikasi Email",
		Body:  fmt.Sprintf(`Hai %v, Selamat datang di Pegadaian Digital Service`, data.Customer.FullName),
		Image: "",
		Token: data.Payload.FcmToken,
		Data:  dataWelcomeMessage,
	}
	_, err = c.SendNotification(welcomeMessage)
	if err != nil {
		log.Debugf("Error when send notification message: %s, phone : %s", data.RegisterOTP.RegistrationId, data.Customer.Phone)
	}

	return nil
}

func (c *Notification) SendNotificationBlock(data dto.NotificationBlock) error {

	// Send Email Block
	dataBlockEmail := &dto.EmailBlock{
		Title:        "Notifikasi Keamanan Pegadaian Digital",
		Text:         "Notifikasi Keamanan Pegadaian Digital",
		Message:      data.Message,
		LastTryLogin: data.LastTryLogin,
		BaseUrl:      c.httpBaseUrl,
	}
	htmlMessage, err := nval.TemplateFile(dataBlockEmail, "email_blocked_login.html")
	if err != nil {
		return err
	}

	// set payload email service
	emailPayload := dto.EmailPayload{
		Subject: fmt.Sprintf("Notifikasi Keamanan Pegadaian Digital %s", data.Customer.FullName),
		From: dto.FromEmailPayload{
			Name:  c.emailConfig.PdsEmailFromName,
			Email: c.emailConfig.PdsEmailFrom,
		},
		To:         data.Customer.Email,
		Message:    htmlMessage,
		Attachment: "",
		MimeType:   "",
	}
	_, err = c.SendEmail(emailPayload)
	if err != nil {
		log.Debugf("Error when send email block account. Payload %v", emailPayload)
	}

	return nil
}
