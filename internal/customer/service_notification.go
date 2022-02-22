package customer

import (
	"fmt"
	"github.com/nbs-go/nlogger"
	"math/rand"
	"net/http"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
)

func (s *Service) SendNotification(payload dto.NotificationPayload) (*http.Response, error) {
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

	// TODO: Refactor endpoint & payload data for [Notification Service]
	// Send Notification
	resp, err := s.ClientPostData("/push-notification", reqBody, reqHeader)
	if err != nil {
		log.Error("Error when send notification")
		return resp, ncore.TraceError("error", err)
	}

	// Set result
	result = resp

	return result, nil
}

func (s *Service) SendEmail(payload dto.EmailPayload) (*http.Response, error) {
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
	resp, err := s.ClientPostData("/send-email", reqBody, reqHeader)
	if err != nil {
		s.log.Errorf("Error when send email. %v", err)
		return resp, err
	}

	// Set result
	result = resp

	return result, nil
}

func (s *Service) SendNotificationRegister(data dto.NotificationRegister) error {
	customer := data.Customer.(*model.Customer)
	verification := data.Verification.(*model.Verification)

	// Send Email Verification
	dataEmailVerification := &dto.EmailVerification{
		FullName:        customer.FullName,
		Email:           customer.Email,
		VerificationURL: fmt.Sprintf("%sauth/verify_email?t=%s", s.config.GetHTTPBaseURL(), verification.EmailVerificationToken),
	}
	htmlMessage, err := nval.TemplateFile(dataEmailVerification, "email_verification.html")
	if err != nil {
		return err
	}

	// set payload email service
	emailPayload := dto.EmailPayload{
		Subject: fmt.Sprintf("Verifikasi Email %s", customer.FullName),
		From: dto.FromEmailPayload{
			Name:  s.config.EmailConfig.PdsEmailFromName,
			Email: s.config.EmailConfig.PdsEmailFrom,
		},
		To:         customer.Email,
		Message:    htmlMessage,
		Attachment: "",
		MimeType:   "",
	}
	respEmail, err := s.SendEmail(emailPayload)
	if err != nil {
		log.Debugf("Error when send email verification. Payload %v", emailPayload)
	}
	defer respEmail.Body.Close()

	// Send Notification Welcome
	id, _ := nval.ParseString(rand.Intn(100)) //nolint:gosec
	var dataWelcomeMessage = map[string]string{
		"title": "Verifikasi Email",
		"body":  fmt.Sprintf(`Hai %v, Selamat datang di Pegadaian Digital Service`, customer.FullName),
		"type":  constant.TypeProfile,
		"id":    id,
	}
	welcomeMessage := dto.NotificationPayload{
		Title: "Verifikasi Email",
		Body:  fmt.Sprintf(`Hai %v, Selamat datang di Pegadaian Digital Service`, customer.FullName),
		Image: "",
		Token: data.Payload.FcmToken,
		Data:  dataWelcomeMessage,
	}
	respSend, err := s.SendNotification(welcomeMessage)
	if err != nil {
		s.log.Debugf("Error when send notification message.", nlogger.Error(err), nlogger.Context(s.ctx))
	}
	defer respSend.Body.Close()

	return nil
}

func (s *Service) SendNotificationBlock(data dto.NotificationBlock) error {
	customer := data.Customer.(*model.Customer)

	baseURL := s.config.GetHTTPBaseURL()
	// Send Email Block
	dataBlockEmail := &dto.EmailBlock{
		Title:        "Notifikasi Keamanan Pegadaian Digital",
		Text:         "Notifikasi Keamanan Pegadaian Digital",
		Message:      data.Message,
		LastTryLogin: data.LastTryLogin,
		BaseURL:      baseURL,
	}
	htmlMessage, err := nval.TemplateFile(dataBlockEmail, "email_blocked_login.html")
	if err != nil {
		return err
	}

	// set payload email service
	emailPayload := dto.EmailPayload{
		Subject: fmt.Sprintf("Notifikasi Keamanan Pegadaian Digital %s", customer.FullName),
		From: dto.FromEmailPayload{
			Name:  s.config.EmailConfig.PdsEmailFromName,
			Email: s.config.EmailConfig.PdsEmailFrom,
		},
		To:         customer.Email,
		Message:    htmlMessage,
		Attachment: "",
		MimeType:   "",
	}
	sendEmail, err := s.SendEmail(emailPayload)
	if err != nil {
		s.log.Debugf("Error when send email block account. Payload %v", emailPayload)
	}
	defer sendEmail.Body.Close()

	return nil
}
