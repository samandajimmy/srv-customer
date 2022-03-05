package customer

import (
	"database/sql"
	"github.com/nbs-go/errx"
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

func (rc *RepositoryContext) FindAddressByCustomerId(id int64) (*model.Address, error) {
	var row model.Address
	err := rc.stmt.Address.FindByCustomerID.GetContext(rc.ctx, &row, id)
	return &row, err
}

func (rc *RepositoryContext) CreateAddress(row *model.Address) error {
	_, err := rc.stmt.Address.Insert.ExecContext(rc.ctx, &row)
	return err
}

func (rc *RepositoryContext) UpdateAddress(row *model.Address) error {
	result, err := rc.stmt.Address.Update.ExecContext(rc.ctx, row)
	if err != nil {
		return errx.Trace(err)
	}
	if !nsql.IsUpdated(result) {
		return constant.ResourceNotFoundError
	}
	return nil
}

func (rc *RepositoryContext) FindAddressPrimary(customerId int64) (*model.Address, error) {
	var row model.Address
	err := rc.stmt.Address.FindPrimaryByCustomerID.GetContext(rc.ctx, &row, customerId)
	if err != nil {
		return nil, errx.Trace(err)
	}
	return &row, nil
}

func (rc *RepositoryContext) InsertOrUpdateAddress(row *model.Address) error {
	address, err := rc.FindAddressByCustomerId(row.CustomerID)
	if err != nil && err != sql.ErrNoRows {
		return errx.Trace(err)
	}

	if err != sql.ErrNoRows {
		err = rc.UpdateAddress(address)
		if err != nil {
			rc.log.Error("cannot update address.", nlogger.Error(err), nlogger.Context(rc.ctx))
			return errx.Trace(err)
		}
		return nil

	} else {
		err = rc.CreateAddress(row)
		if err != nil {
			rc.log.Error("cannot create address.", nlogger.Error(err), nlogger.Context(rc.ctx))
			return errx.Trace(err)
		}
		return nil
	}
}
