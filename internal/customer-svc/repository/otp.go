package repository

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/repository/statement"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type OTP struct {
	db   *nsql.DB
	stmt *statement.OTPStatement
}

func (a *OTP) HasInitialized() bool {
	return true
}

func (a *OTP) Init(dataSources DataSourceMap, _ contract.RepositoryMap) error {
	a.db = dataSources.Postgres
	a.stmt = statement.NewOTPStatement(a.db)
	return nil
}

func (a *OTP) Insert(row *model.OTP) error {
	_, err := a.stmt.Insert.Exec(row)
	return err
}
