package contract

import (
	"net/http"

	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/dto"
)

type AuthService interface {
	ValidateClient(payload dto.ClientCredential) error
}

type CustomerService interface {
	Register(payload dto.RegisterNewCustomer) (*dto.RegisterNewCustomerResponse, error)
	RegisterStepOne(payload dto.RegisterStepOne) (*dto.RegisterStepOneResponse, error)
	RegisterResendOTP(payload dto.RegisterResendOTP) (*dto.RegisterResendOTPResponse, error)
	RegisterStepTwo(payload dto.RegisterStepTwo) (*dto.RegisterStepTwoResponse, error)
}

type OTPService interface {
	SendOTP(payload dto.SendOTPRequest) (*http.Response, error)
	VerifyOTP(payload dto.VerifyOTPRequest) (*http.Response, error)
}
