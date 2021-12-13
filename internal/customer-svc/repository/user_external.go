package repository

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/repository/statement"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type UserExternal struct {
	db   *nsql.DB
	stmt *statement.UserStatement
}

func (a *UserExternal) HasInitialized() bool {
	return true
}

func (a *UserExternal) Init(dataSources DataSourceMap, _ contract.RepositoryMap) error {
	a.db = dataSources.DBExternal
	a.stmt = statement.NewUserStatement(a.db)
	return nil
}

func (a *UserExternal) FindByEmailOrPhone(email string) (*model.User, error) {
	var row model.User
	err := a.stmt.FindByEmailOrPhone.Get(&row, email, email)
	return &row, err
}
