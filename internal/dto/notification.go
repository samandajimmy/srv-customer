package dto

import "encoding/json"

type NotificationPayload struct {
	Title string            `json:"title"`
	Body  string            `json:"body"`
	Image string            `json:"imageUrl"`
	Token string            `json:"token"`
	Data  map[string]string `json:"data"`
}

type NotificationRegister struct {
	Customer     interface{} // *model.Customer
	Verification interface{} // *model.Verification
	RegisterOTP  interface{} // *model.VerificationOTP
	Payload      RegisterPayload
}

type NotificationBlock struct {
	Customer     interface{} // *model.Customer
	Message      string
	LastTryLogin string
}

type NotificationOptionVO struct {
	FCM  *FCMOption  `json:"fcm"`
	SMTP *SMTPOption `json:"smtp"`
}

type FCMOption struct {
	UserID   string            `json:"userId"`
	Title    string            `json:"title"`
	Body     string            `json:"body"`
	ImageURL string            `json:"imageUrl"`
	Token    string            `json:"token"`
	Metadata json.RawMessage   `json:"metadata"`
	Data     map[string]string `json:"data"`
}

type SMTPOption struct {
	UserID     string   `json:"userId"`
	Subject    string   `json:"subject"`
	Message    string   `json:"message"`
	From       MailFrom `json:"from"`
	To         string   `json:"to"`
	Attachment string   `json:"attachment"`
	MimeType   string   `json:"mimeType"`
}

type MailFrom struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateNotificationRequest struct {
	UserID  string               `json:"userId"`
	Options NotificationOptionVO `json:"options"`
}

type EmailAndNotificationPayload struct {
	EmailPayload
	NotificationPayload
}
