package customer

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

func (rc *RepositoryContext) CreateAccessSession(row *model.AccessSession) error {
	_, err := rc.stmt.AccessSession.Insert.ExecContext(rc.ctx, row)
	if err != nil {
		return err
	}
	return nil
}

func (rc *RepositoryContext) UpdateAccessSession(row *model.AccessSession) error {
	result, err := rc.stmt.AccessSession.Update.ExecContext(rc.ctx, row)
	if err != nil {
		return ncore.TraceError("failed to update access session", err)
	}
	if !nsql.IsUpdated(result) {
		return constant.ResourceNotFoundError
	}
	return nil
}
