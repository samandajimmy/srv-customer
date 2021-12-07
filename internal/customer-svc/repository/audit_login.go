package repository

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/repository/statement"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type AuditLogin struct {
	db   *nsql.DB
	stmt *statement.AuditStatement
}

func (a *AuditLogin) HasInitialized() bool {
	return true
}

func (a *AuditLogin) Init(dataSources DataSourceMap, _ contract.RepositoryMap) error {
	a.db = dataSources.Postgres
	a.stmt = statement.NewAuditLoginStatement(a.db)
	return nil
}

func (a *AuditLogin) Insert(row *model.AuditLogin) error {
	_, err := a.stmt.Insert.Exec(row)
	return err
}

func (a *AuditLogin) CountLogin(customerId int64) (int64, error) {
	var count int64
	err := a.stmt.CountLogin.QueryRow(customerId).Scan(&count)
	return count, err
}
