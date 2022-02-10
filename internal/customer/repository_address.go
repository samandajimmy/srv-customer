package customer

import (
	"database/sql"
	"errors"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

func (rc *RepositoryContext) FindAddressByCustomerId(id int64) (*model.Address, error) {
	var row model.Address
	err := rc.stmt.Address.FindByCustomerId.GetContext(rc.ctx, &row, id)
	return &row, err
}

func (rc *RepositoryContext) CreateAddress(row *model.Address) error {
	_, err := rc.stmt.Address.Insert.ExecContext(rc.ctx, row)
	if err != nil {
		return err
	}
	return nil
}

func (rc *RepositoryContext) UpdateAddress(row *model.Address) error {
	result, err := rc.stmt.Address.Update.ExecContext(rc.ctx, row)
	if err != nil {
		return ncore.TraceError("failed to update address", err)
	}
	if !nsql.IsUpdated(result) {
		return constant.ResourceNotFoundError
	}
	return nil
}

func (rc *RepositoryContext) FindAddressPrimary(customerId int64) (*model.Address, error) {
	var row model.Address
	err := rc.stmt.Address.FindPrimaryByCustomerId.GetContext(rc.ctx, &row, customerId)
	if err != nil {
		return nil, ncore.TraceError("failed to find primary address", err)
	}
	return &row, nil
}

func (rc *RepositoryContext) InsertOrUpdateAddress(row *model.Address) error {
	address, err := rc.FindAddressByCustomerId(row.CustomerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			address = nil
		}
	}

	if address != nil {
		err = rc.UpdateAddress(address)
		if err != nil {
			return ncore.TraceError("cannot update address.", err)
		}
		return nil

	} else {
		err = rc.CreateAddress(row)
		if err != nil {
			return ncore.TraceError("cannot create address.", err)
		}
		return nil
	}
}
