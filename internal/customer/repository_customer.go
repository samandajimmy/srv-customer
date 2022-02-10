package customer

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

func (rc *RepositoryContext) CreateCustomer(row *model.Customer) (int64, error) {
	var lastInsertId int64
	err := rc.stmt.Customer.Insert.QueryRowContext(rc.ctx, &row).Scan(&lastInsertId)
	return lastInsertId, err
}

func (rc *RepositoryContext) FindCustomerByID(id int64) (*model.Customer, error) {
	var row model.Customer
	err := rc.stmt.Customer.FindById.GetContext(rc.ctx, &row, id)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (rc *RepositoryContext) FindCustomerByPhone(phone string) (*model.Customer, error) {
	var row model.Customer
	err := rc.stmt.Customer.FindByPhone.GetContext(rc.ctx, &row, phone)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (rc *RepositoryContext) FindCustomerByEmail(email string) (*model.Customer, error) {
	var row model.Customer
	err := rc.stmt.Customer.FindByEmail.GetContext(rc.ctx, &row, email)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (rc *RepositoryContext) FindCustomerByEmailOrPhone(email string) (*model.Customer, error) {
	var row model.Customer
	err := rc.stmt.Customer.FindByEmailOrPhone.GetContext(rc.ctx, &row, email)
	return &row, err
}

func (rc *RepositoryContext) UpdateCustomerByPhone(row *model.Customer) error {
	result, err := rc.stmt.Customer.UpdateByPhone.ExecContext(rc.ctx, row)
	if err != nil {
		return ncore.TraceError("cannot update customer by phone", err)
	}
	if !nsql.IsUpdated(result) {
		return constant.ResourceNotFoundError
	}
	return nil
}
