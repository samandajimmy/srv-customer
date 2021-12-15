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
	FindById(id int64) (*model.Customer, error)
	FindByPhone(phone string) (*model.Customer, error)
	FindByEmail(email string) (*model.Customer, error)
	FindByEmailOrPhone(email string) (*model.Customer, error)
	BlockAccount(phone string) error
	UnBlockAccount(phone string) error
}

type AuditLoginRepository interface {
	Insert(row *model.AuditLogin) error
	CountLogin(customerId int64) (int64, error)
}

type OTPRepository interface {
	Insert(row *model.OTP) error
}

type CredentialRepository interface {
	FindByCustomerId(customerId int64) (*model.Credential, error)
	Insert(row *model.Credential) error
	InsertOrUpdate(row *model.Credential) error
	UpdateByCustomerID(row *model.Credential) error
	DeleteByID(id string) error
}

type FinancialDataRepository interface {
	FindByCustomerId(customerId int64) (*model.FinancialData, error)
	Insert(row *model.FinancialData) error
	InsertOrUpdate(row *model.FinancialData) error
	UpdateByCustomerID(row *model.FinancialData) error
	DeleteByID(id string) error
}

type AccessSessionRepository interface {
	Insert(row *model.AccessSession) error
	Update(row *model.AccessSession) error
}

type VerificationRepository interface {
	FindByCustomerId(customerId int64) (*model.Verification, error)
	FindByEmailToken(token string) (*model.Verification, error)
	Insert(row *model.Verification) error
	InsertOrUpdate(row *model.Verification) error
	UpdateByCustomerID(row *model.Verification) error
	DeleteByID(id string) error
}

type AddressRepository interface {
	Insert(row *model.Address) error
	FindPrimaryAddress(customerId int64) (*model.Address, error)
	Update(row *model.Address) error
}

type UserExternalRepository interface {
	FindByEmailOrPhone(email string) (*model.User, error)
	FindAddressByCustomerId(id string) (*model.AddressExternal, error)
}

type UserPinExternalRepository interface {
	FindByCustomerId(id int64) (*model.UserPin, error)
}
