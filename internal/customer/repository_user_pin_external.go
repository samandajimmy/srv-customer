package customer

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
)

func (rc *RepositoryContext) FindUserPINByCustomerID(id int64) (*model.UserPin, error) {
	var row model.UserPin
	err := rc.stmt.UserPin.FindByCustomerID.GetContext(rc.ctx, &row, id)
	return &row, err
}
