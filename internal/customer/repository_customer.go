package customer

import (
	"github.com/nbs-go/nlogger"
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

func (rc *RepositoryContext) FindCustomerByUserRefID(id string) (*model.Customer, error) {
	var row model.Customer
	err := rc.stmt.Customer.FindByRefId.GetContext(rc.ctx, &row, id)
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

func (rc *RepositoryContext) FindCustomerByPhoneOrCIF(cif string) (*model.Customer, error) {
	var row model.Customer
	err := rc.stmt.Customer.FindByPhoneOrCIF.GetContext(rc.ctx, &row, cif, cif)
	return &row, err
}

func (rc *RepositoryContext) ReferralCodeExists(referralCode string) (*model.Customer, error) {
	var row model.Customer
	err := rc.stmt.Customer.ReferralCodeExist.GetContext(rc.ctx, &row, referralCode)
	return &row, err
}

func (rc *RepositoryContext) UpdateCustomerByCIF(customer *model.Customer, cif string) error {
	result, err := rc.stmt.Customer.UpdateByCIF.ExecContext(rc.ctx, &model.UpdateByCIF{
		Customer: customer,
		Cif:      cif,
	})
	if err != nil {
		return ncore.TraceError("cannot update customer by cif", err)
	}
	if !nsql.IsUpdated(result) {
		return constant.ResourceNotFoundError
	}
	return nil
}

func (rc *RepositoryContext) UpdateCustomerProfile(customer *model.Customer, address *model.Address) error {
	tx, err := rc.conn.BeginTxx(rc.ctx, nil)
	if err != nil {
		return ncore.TraceError("", err)
	}
	defer rc.ReleaseTx(tx, &err)

	// Update the data to repositories
	customerUpdate, err := rc.stmt.Customer.UpdateByUserRefID.ExecContext(rc.ctx, &model.UpdateByID{
		Customer: customer,
		ID:       customer.Id,
	})
	if err != nil {
		return ncore.TraceError("", err)
	}
	if !nsql.IsUpdated(customerUpdate) {
		return constant.ResourceNotFoundError
	}

	err = rc.InsertOrUpdateAddress(address)
	if err != nil {
		return ncore.TraceError("", err)
	}

	return nil
}

func (rc *RepositoryContext) UpdateCustomerByPhone(customer *model.Customer) error {
	result, err := rc.stmt.Customer.UpdateByPhone.ExecContext(rc.ctx, customer)
	if err != nil {
		return ncore.TraceError("cannot update customer by phone", err)
	}
	if !nsql.IsUpdated(result) {
		return constant.ResourceNotFoundError
	}
	return nil
}

func (rc *RepositoryContext) UpdateCustomerByUserRefID(customer *model.Customer, userRefID string) error {
	result, err := rc.stmt.Customer.UpdateByPhone.ExecContext(rc.ctx, &model.UpdateCustomerByUserRefID{
		Customer:  customer,
		UserRefId: userRefID,
	})
	if err != nil {
		rc.log.Error("error found when update customer by UserRefID", nlogger.Error(err))
		return ncore.TraceError("cannot update customer by phone", err)
	}
	if !nsql.IsUpdated(result) {
		return constant.ResourceNotFoundError
	}

	return nil
}
