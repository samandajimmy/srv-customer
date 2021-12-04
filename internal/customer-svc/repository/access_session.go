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

func (c *AccessSession) HasInitialized() bool {
	return true
}

func (a *AccessSession) Init(dataSources DataSourceMap, _ contract.RepositoryMap) error {
	a.db = dataSources.Postgres
	a.stmt = statement.NewAccessSessionStatement(a.db)
	return nil
}

func (a *AccessSession) Insert(row *model.AccessSession) (int64, error) {
	var lastInsertId int64
	err := a.stmt.Insert.QueryRow(&row).Scan(&lastInsertId)
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil
}
