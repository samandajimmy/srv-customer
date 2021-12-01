package contract

import (
	"net/http"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/dto"
)

type AuthService interface {
	ValidateClient(payload dto.ClientCredential) error
}

type CustomerService interface {
	Register(payload dto.RegisterNewCustomer) (*dto.NewRegisterResponse, error)
	RegisterStepOne(payload dto.RegisterStepOne) (*dto.RegisterStepOneResponse, error)
}

type OTPService interface {
	SendOTP(payload dto.SendOTPRequest) (*http.Response, error)
}
