package customer

import (
	"github.com/nbs-go/errx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

func (rc *RepositoryContext) CreateAccessSession(accessSession *model.AccessSession) error {
	_, err := rc.stmt.AccessSession.Insert.ExecContext(rc.ctx, accessSession)
	if err != nil {
		return err
	}
	return nil
}

func (rc *RepositoryContext) UpdateAccessSession(accessSession *model.AccessSession) error {
	result, err := rc.stmt.AccessSession.Update.ExecContext(rc.ctx, accessSession)

	if err != nil {
		return errx.Trace(err)
	}

	// If not updated, then return stale error
	if !nsql.IsUpdated(result) {
		return constant.StaleResourceError.Trace()
	}
	return nil
}
