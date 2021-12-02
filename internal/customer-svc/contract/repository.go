package contract

import "repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/model"

type VerificationOTPRepository interface {
	Insert(row *model.VerificationOTP) (int64, error)
}

type CustomerRepository interface {
	Insert(row *model.Customer) (int64, error)
	FindByPhone(phone string) *model.Customer
}
