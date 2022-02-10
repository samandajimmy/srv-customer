package contract

import (
	"net/http"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
)

type Service interface {
	ValidateClient(payload dto.ClientCredential) error

	Login(payload dto.LoginRequest) (*dto.LoginResponse, error)
	Register(payload dto.RegisterNewCustomer) (*dto.RegisterNewCustomerResponse, error)
	RegisterStepOne(payload dto.RegisterStepOne) (*dto.RegisterStepOneResponse, error)
	RegisterResendOTP(payload dto.RegisterResendOTP) (*dto.RegisterResendOTPResponse, error)
	RegisterStepTwo(payload dto.RegisterStepTwo) (*dto.RegisterStepTwoResponse, error)

	SendOTP(payload dto.SendOTPRequest) (*http.Response, error)
	VerifyOTP(payload dto.VerifyOTPRequest) (*http.Response, error)

	SynchronizeCustomer(payload dto.RegisterNewCustomer) (*http.Response, error)

	CacheGet(key string) (string, error)
	CacheSetThenGet(key string, value string, expire int64) (string, error)

	SendNotification(payload dto.NotificationPayload) (*http.Response, error)
	SendEmail(payload dto.EmailPayload) (*http.Response, error)
	SendNotificationRegister(data dto.NotificationRegister) error
	SendNotificationBlock(data dto.NotificationBlock) error

	VerifyEmailCustomer(payload dto.VerificationPayload) (string, error)
}

type AuthService interface {
	ValidateClient(payload dto.ClientCredential) error
}

type CustomerService interface {
	Login(payload dto.LoginRequest) (*dto.LoginResponse, error)
	Register(payload dto.RegisterNewCustomer) (*dto.RegisterNewCustomerResponse, error)
	RegisterStepOne(payload dto.RegisterStepOne) (*dto.RegisterStepOneResponse, error)
	RegisterResendOTP(payload dto.RegisterResendOTP) (*dto.RegisterResendOTPResponse, error)
	RegisterStepTwo(payload dto.RegisterStepTwo) (*dto.RegisterStepTwoResponse, error)
}

type OTPService interface {
	SendOTP(payload dto.SendOTPRequest) (*http.Response, error)
	VerifyOTP(payload dto.VerifyOTPRequest) (*http.Response, error)
}

type CacheService interface {
	Get(key string) (string, error)
	SetThenGet(key string, value string, expire int64) (string, error)
}

type NotificationService interface {
	SendNotification(payload dto.NotificationPayload) (*http.Response, error)
	SendEmail(payload dto.EmailPayload) (*http.Response, error)
	SendNotificationRegister(data dto.NotificationRegister) error
	SendNotificationBlock(data dto.NotificationBlock) error
}

type VerificationService interface {
	VerifyEmailCustomer(payload dto.VerificationPayload) (string, error)
}

type PdsAPIService interface {
	SynchronizeCustomer(payload dto.RegisterNewCustomer) (*http.Response, error)
}
