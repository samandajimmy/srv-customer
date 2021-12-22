package repository

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/repository/statement"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type UserRegisterExternal struct {
	db   *nsql.DB
	stmt *statement.UserRegisterStatement
}

func (a *UserRegisterExternal) HasInitialized() bool {
	return true
}

func (a *UserRegisterExternal) Init(dataSources DataSourceMap, _ contract.RepositoryMap) error {
	a.db = dataSources.DBExternal
	a.stmt = statement.NewUserRegisterStatement(a.db)
	return nil
}

func (a *UserRegisterExternal) Insert(row *model.UserRegister) error {
	_, err := a.stmt.Insert.Exec(row)
	return err
}
