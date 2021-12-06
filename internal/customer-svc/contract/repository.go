package contract

import "repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/model"

type VerificationOTPRepository interface {
	Insert(row *model.VerificationOTP) (int64, error)
	FindByRegistrationIdAndPhone(id string, phone string) (*model.VerificationOTP, error)
	Delete(id string, phone string) error
}

type CustomerRepository interface {
	Insert(row *model.Customer) (int64, error)
	UpdateByPhone(row *model.Customer) error
	FindByPhone(phone string) (*model.Customer, error)
	FindByEmail(email string) (*model.Customer, error)
	FindByEmailOrPhone(phone string) *model.CustomerAuthentication
	BlockAccount(phone string) error
	UnBlockAccount(phone string) error
}

type OTPRepository interface {
	Insert(row *model.OTP) error
}

type CredentialRepository interface {
	FindByCustomerId(customerId int64) (*model.Credential, error)
	Insert(row *model.Credential) error
	InsertOrUpdate(row *model.Credential) error
	DeleteByID(id string) error
}

type AccessSessionRepository interface {
	Insert(row *model.AccessSession) error
	Update(row *model.AccessSession) error
}
