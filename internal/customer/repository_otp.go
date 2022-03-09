package customer

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
)

func (rc *RepositoryContext) CreateOTP(row *model.OTP) error {
	_, err := rc.stmt.OTP.Insert.ExecContext(rc.ctx, row)
	if err != nil {
		return err
	}
	return nil
}
