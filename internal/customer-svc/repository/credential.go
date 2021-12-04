package repository

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/repository/statement"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type Credential struct {
	db   *nsql.DB
	stmt *statement.CredentialStatement
}

func (c *Credential) HasInitialized() bool {
	return true
}

func (a *Credential) Init(dataSources DataSourceMap, _ contract.RepositoryMap) error {
	a.db = dataSources.Postgres
	a.stmt = statement.NewCredentialStatement(a.db)
	return nil
}

func (a *Credential) Insert(row *model.Credential) (int64, error) {
	var lastInsertId int64
	err := a.stmt.Insert.QueryRow(&row).Scan(&lastInsertId)
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil
}
