package contract

import "repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/model"

type VerificationOTPRepository interface {
	Insert(row *model.VerificationOTP) (int64, error)
}

type CustomerRepository interface {
	Insert(row *model.Customer) (int64, error)
	FindByPhone(phone string) *model.Customer
	FindByEmailOrPhone(phone string) *model.CustomerAuthentication
	BlockAccount(phone string) error
	UnBlockAccount(phone string) error
}

type OTPRepository interface {
	Insert(row *model.OTP) (int64, error)
}

type CredentialRepository interface {
	Insert(row *model.Credential) (int64, error)
}
