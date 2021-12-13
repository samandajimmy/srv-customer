package repository

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/repository/statement"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type AccessSession struct {
	db   *nsql.DB
	stmt *statement.AccessSessionStatement
}

func (a *AccessSession) HasInitialized() bool {
	return true
}

func (a *AccessSession) Init(dataSources DataSourceMap, _ contract.RepositoryMap) error {
	a.db = dataSources.DBInternal
	a.stmt = statement.NewAccessSessionStatement(a.db)
	return nil
}

func (a *AccessSession) Insert(row *model.AccessSession) error {
	_, err := a.stmt.Insert.Exec(row)
	return err
}

func (a *AccessSession) Update(row *model.AccessSession) error {
	result, err := a.stmt.Update.Exec(row)
	if err != nil {
		return err
	}
	return nsql.IsUpdated(result)
}
