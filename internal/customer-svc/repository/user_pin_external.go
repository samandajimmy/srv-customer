package repository

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/repository/statement"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type UserPinExternal struct {
	db   *nsql.DB
	stmt *statement.UserPinStatement
}

func (a *UserPinExternal) HasInitialized() bool {
	return true
}

func (a *UserPinExternal) Init(dataSources DataSourceMap, _ contract.RepositoryMap) error {
	a.db = dataSources.DBExternal
	a.stmt = statement.NewUserPinStatement(a.db)
	return nil
}

func (a *UserPinExternal) FindByCustomerId(id int64) (*model.UserPin, error) {
	var row model.UserPin
	err := a.stmt.FindByCustomerID.Get(&row, id)
	return &row, err
}
