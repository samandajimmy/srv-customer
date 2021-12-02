package repository

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/repository/statement"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type VerificationOTP struct {
	db   *nsql.DB
	stmt *statement.VerificationOTPStatement
}

func (c *VerificationOTP) HasInitialized() bool {
	return true
}

func (a *VerificationOTP) Init(dataSources DataSourceMap, _ contract.RepositoryMap) error {
	a.db = dataSources.Postgres
	a.stmt = statement.NewVerificationOTPStatement(a.db)
	return nil
}

func (a *VerificationOTP) Insert(row *model.VerificationOTP) (int64, error) {
	var lastInsertId int64
	err := a.stmt.Insert.QueryRow(&row).Scan(&lastInsertId)
	if err != nil {
		log.Fatalf("ERR INSERT %v", err)
		return 0, err
	}
	return lastInsertId, nil
}
