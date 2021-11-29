package contract

import "repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/model"

type CustomerRepository interface {
	Insert(row *model.Customer) (int64, error)
	FindByPhone(phone string) *model.Customer
}
