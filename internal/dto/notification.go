package dto

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
)

type NotificationPayload struct {
	Title string            `json:"title"`
	Body  string            `json:"body"`
	Image string            `json:"image"`
	Token string            `json:"token"`
	Data  map[string]string `json:"data"`
}

type NotificationRegister struct {
	Customer     *model.Customer
	Verification *model.Verification
	RegisterOTP  *model.VerificationOTP
	Payload      RegisterNewCustomer
}

type NotificationBlock struct {
	Customer     *model.Customer
	Message      string
	LastTryLogin string
}
