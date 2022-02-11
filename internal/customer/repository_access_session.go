package customer

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

func (rc *RepositoryContext) CreateAccessSession(accessSession *model.AccessSession) error {
	_, err := rc.stmt.AccessSession.Insert.ExecContext(rc.ctx, accessSession)
	if err != nil {
		return err
	}
	return nil
}

func (rc *RepositoryContext) UpdateAccessSession(accessSession *model.AccessSession, currentVersion int64) error {
	result, err := rc.stmt.AccessSession.Update.ExecContext(rc.ctx, &model.UpdateAccessSession{
		AccessSession:  accessSession,
		CurrentVersion: currentVersion,
	})

	if err != nil {
		return ncore.TraceError("failed to update access session", err)
	}

	// If not updated, then return stale error
	if !nsql.IsUpdated(result) {
		return constant.StaleResourceError.Trace()
	}
	return nil
}
