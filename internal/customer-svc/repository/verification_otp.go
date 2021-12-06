package repository

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/repository/statement"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type VerificationOTP struct {
	db   *nsql.DB
	stmt *statement.VerificationOTPStatement
}

func (a *VerificationOTP) HasInitialized() bool {
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
		return 0, err
	}
	return lastInsertId, nil
}

func (a *VerificationOTP) FindByRegistrationIdAndPhone(id string, phone string) (*model.VerificationOTP, error) {
	var row model.VerificationOTP
	err := a.stmt.FindByRegistrationId.Get(&row, id, phone)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (a *VerificationOTP) Delete(id string, phone string) error {
	_, err := a.stmt.Delete.Exec(id, phone)
	return ncore.TraceError(err)
}
