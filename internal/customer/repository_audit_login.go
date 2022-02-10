package customer

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
)

func (rc *RepositoryContext) CreateAuditLogin(row *model.AuditLogin) error {
	_, err := rc.stmt.AuditLogin.Insert.ExecContext(rc.ctx, row)
	if err != nil {
		return err
	}
	return nil
}

func (rc *RepositoryContext) CountAuditLogin(customerId int64) (int64, error) {
	var count int64
	err := rc.stmt.AuditLogin.CountLogin.QueryRowContext(rc.ctx, customerId).Scan(&count)
	return count, err
}
