package customer

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
)

func (rc *RepositoryContext) CreateUserRegister(row *model.UserRegister) error {
	_, err := rc.stmt.UserRegister.Insert.ExecContext(rc.ctx, row)
	if err != nil {
		return err
	}
	return nil
}
