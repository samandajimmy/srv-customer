package customer

import (
	"database/sql"
	"errors"
	"github.com/nbs-go/errx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

func (rc *RepositoryContext) InsertVerification(row *model.Verification) error {
	_, err := rc.stmt.Verification.Insert.ExecContext(rc.ctx, row)
	if err != nil {
		return err
	}
	return nil
}

func (rc *RepositoryContext) UpdateVerification(row *model.Verification) error {
	result, err := rc.stmt.Verification.Update.ExecContext(rc.ctx, row)
	if err != nil {
		return errx.Trace(err)
	}
	if !nsql.IsUpdated(result) {
		return constant.ResourceNotFoundError
	}
	return nil
}

func (rc *RepositoryContext) FindVerificationByCustomerID(id int64) (*model.Verification, error) {
	var row model.Verification
	err := rc.stmt.Verification.FindByCustomerID.GetContext(rc.ctx, &row, id)
	return &row, err
}

func (rc *RepositoryContext) FindVerificationByEmailToken(token string) (*model.Verification, error) {
	var row model.Verification
	err := rc.stmt.Verification.FindByEmailToken.GetContext(rc.ctx, &row, token)
	return &row, err
}

func (rc *RepositoryContext) DeleteVerification(id string) error {
	result, err := rc.stmt.Verification.DeleteByID.ExecContext(rc.ctx, id)
	if err != nil {
		return errx.Trace(err)
	}
	if !nsql.IsUpdated(result) {
		return constant.ResourceNotFoundError
	}
	return nil
}

func (rc *RepositoryContext) UpdateVerificationByCustomerID(row *model.Verification) error {
	result, err := rc.stmt.Verification.Update.ExecContext(rc.ctx, row)
	if err != nil {
		return errx.Trace(err)
	}
	if !nsql.IsUpdated(result) {
		return constant.ResourceNotFoundError
	}
	return nil
}

func (rc *RepositoryContext) InsertOrUpdateVerification(row *model.Verification) error {
	// find by customer id
	verification, err := rc.FindVerificationByCustomerID(row.CustomerID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errx.Trace(err)
	}

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		err = rc.InsertVerification(row)
		if err != nil {
			return errx.Trace(err)
		}

		return nil
	}

	err = rc.UpdateVerification(verification)
	if err != nil {
		return errx.Trace(err)
	}

	return nil
}
