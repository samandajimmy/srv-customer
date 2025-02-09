package customer

import (
	"fmt"
	"github.com/nbs-go/errx"
	logOption "github.com/nbs-go/nlogger/v2/option"
	"math/rand"
	"net/http"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
)

func (s *Service) SendNotification(payload dto.NotificationPayload) (*http.Response, error) {
	var result *http.Response

	// Set payload
	reqBody := map[string]interface{}{
		"userId": "",
		"options": map[string]interface{}{
			"fcm": map[string]interface{}{
				"title":    payload.Title,
				"body":     payload.Body,
				"imageUrl": payload.Image,
				"token":    payload.Token,
				"data":     payload.Data,
			},
		},
	}

	// Set header
	reqHeader := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	sp := &NotificationPostDataPayload{
		URL:    "/notifications",
		Data:   reqBody,
		Header: &reqHeader,
	}

	// Create Notification
	resp, err := s.CreateNotificationPostData(sp)
	if err != nil {
		log.Error("Error when send notification")
		return resp, errx.Trace(err)
	}

	// Set result
	result = resp

	return result, nil
}

func (s *Service) SendEmail(payload dto.EmailPayload) (*http.Response, error) {
	var result *http.Response
	// Set payload
	reqBody := map[string]interface{}{
		"userId": payload.UserID,
		"options": map[string]interface{}{
			"smtp": map[string]interface{}{
				"from": map[string]string{
					"name":  payload.From.Name,
					"email": payload.From.Email,
				},
				"to":         payload.To,
				"subject":    payload.Subject,
				"message":    payload.Message,
				"attachment": payload.Attachment,
				"mimeType":   payload.MimeType,
			},
		},
	}

	// Set header
	reqHeader := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	sp := &NotificationPostDataPayload{
		URL:    "/notifications",
		Data:   reqBody,
		Header: &reqHeader,
	}

	// Send email
	resp, err := s.CreateNotificationPostData(sp)
	if err != nil {
		s.log.Error("Error when send email", logOption.Error(err))
		return resp, err
	}

	// Set result
	result = resp

	return result, nil
}

func (s *Service) SendEmailAndNotification(data dto.EmailAndNotificationPayload) (*http.Response, error) {
	var result *http.Response
	// Set payload
	reqBody := map[string]interface{}{
		"userId": data.UserID,
		"options": map[string]interface{}{
			"fcm": map[string]interface{}{
				"title":    data.Title,
				"body":     data.Body,
				"imageUrl": data.Image,
				"token":    data.Token,
				"data":     data.Data,
			},
			"smtp": map[string]interface{}{
				"from": map[string]string{
					"name":  data.From.Name,
					"email": data.From.Email,
				},
				"to":         data.To,
				"subject":    data.Subject,
				"message":    data.Message,
				"attachment": data.Attachment,
				"mimeType":   data.MimeType,
			},
		},
	}

	// Set header
	reqHeader := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	sp := &NotificationPostDataPayload{
		URL:    "/notifications",
		Data:   reqBody,
		Header: &reqHeader,
	}

	// Send email
	resp, err := s.CreateNotificationPostData(sp)
	if err != nil {
		s.log.Error("Error when send email", logOption.Error(err))
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
		VerificationURL: fmt.Sprintf("%s/accounts/verify-email?t=%s", s.config.GetHTTPBaseURL(), verification.EmailVerificationToken),
	}
	htmlMessage, err := nval.TemplateFile(dataEmailVerification, "email_verification.html")
	if err != nil {
		return errx.Trace(err)
	}

	// set payload email service
	emailPayload := dto.EmailPayload{
		UserID:  customer.Phone,
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

	// Send Notification Welcome
	id, _ := nval.ParseString(rand.Intn(100)) //nolint:gosec
	var dataWelcomeMessage = map[string]string{
		"title": "Verifikasi Email",
		"body":  fmt.Sprintf(`Hai %v, Selamat datang di Pegadaian Digital Service`, customer.FullName),
		"type":  constant.TypeProfile,
		"id":    id,
	}
	notificationPayload := dto.NotificationPayload{
		Title: "Verifikasi Email",
		Body:  fmt.Sprintf(`Hai %v, Selamat datang di Pegadaian Digital Service`, customer.FullName),
		Image: "",
		Token: data.Payload.FcmToken,
		Data:  dataWelcomeMessage,
	}

	// Send Email Welcome
	titleWelcome := "Pegadaian Digital Service"
	dataWelcome := dto.EmailWelcome{
		Title:    titleWelcome,
		Email:    customer.Email,
		FullName: customer.FullName,
		BaseURL:  s.config.GetHTTPBaseURL(),
	}
	htmlWelcomeMessage, err := nval.TemplateFile(dataWelcome, "email_welcome.html")
	if err != nil {
		return errx.Trace(err)
	}

	emailWelcomePayload := dto.EmailPayload{
		UserID:  customer.Phone,
		Subject: titleWelcome,
		From: dto.FromEmailPayload{
			Name:  s.config.EmailConfig.PdsEmailFromName,
			Email: s.config.EmailConfig.PdsEmailFrom,
		},
		To:         customer.Email,
		Message:    htmlWelcomeMessage,
		Attachment: "",
		MimeType:   "",
	}
	respSendWelcome, err := s.SendEmail(emailWelcomePayload)
	if err != nil {
		s.log.Debug("error found when send email message", logOption.Error(err))
	}
	defer handleClose(respSendWelcome.Body)

	respSend, err := s.SendEmailAndNotification(dto.EmailAndNotificationPayload{
		EmailPayload:        emailPayload,
		NotificationPayload: notificationPayload,
	})
	if err != nil {
		s.log.Debug("error found when send notification message", logOption.Error(err))
	}
	defer handleClose(respSend.Body)

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
		UserID:  customer.Phone,
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
		s.log.Debug("Error when send email block account. Payload %v", logOption.Format(emailPayload), logOption.Error(err))
	}
	defer handleClose(sendEmail.Body)

	return nil
}
