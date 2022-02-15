package customer

import (
	"database/sql"
	"errors"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

func (rc *RepositoryContext) CreateCredential(row *model.Credential) error {
	_, err := rc.stmt.Credential.Insert.ExecContext(rc.ctx, row)
	if err != nil {
		return err
	}
	return nil
}

func (rc *RepositoryContext) FindCredentialByCustomerID(id int64) (*model.Credential, error) {
	var row model.Credential
	err := rc.stmt.Credential.FindByCustomerID.GetContext(rc.ctx, &row, id)
	return &row, err
}

func (rc *RepositoryContext) DeleteCredential(id string) error {
	result, err := rc.stmt.Credential.DeleteByID.ExecContext(rc.ctx, id)
	if err != nil {
		return ncore.TraceError("failed to delete credential", err)
	}
	if !nsql.IsUpdated(result) {
		return constant.ResourceNotFoundError
	}
	return nil
}

func (rc *RepositoryContext) UpdateCredential(row *model.Credential) error {
	result, err := rc.stmt.Credential.Update.ExecContext(rc.ctx, row)
	if err != nil {
		return ncore.TraceError("failed to update credential", err)
	}
	if !nsql.IsUpdated(result) {
		return constant.ResourceNotFoundError
	}
	return nil
}

func (rc *RepositoryContext) InsertOrUpdateCredential(row *model.Credential) error {
	// find by customer id
	credential, err := rc.FindCredentialByCustomerID(row.CustomerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			credential = nil
		}
	}
	if credential != nil {
		err = rc.UpdateCredential(credential)
		if err != nil {
			return ncore.TraceError("cannot update credential.", err)
		}
		return nil
	} else {
		err = rc.CreateCredential(row)
		if err != nil {
			return ncore.TraceError("cannot insert credential.", err)
		}
		return nil
	}
}

func (rc *RepositoryContext) IsValidPassword(customerId int64, password string) (*model.Credential, error) {
	var row model.Credential
	err := rc.stmt.Credential.FindByPasswordAndCustomerID.GetContext(rc.ctx, &row, customerId, password)
	return &row, err
}
