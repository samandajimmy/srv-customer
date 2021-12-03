package repository

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/repository/statement"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type Customer struct {
	db   *nsql.DB
	stmt *statement.CustomerStatement
}

func (c *Customer) HasInitialized() bool {
	return true
}

func (a *Customer) Init(dataSources DataSourceMap, _ contract.RepositoryMap) error {
	a.db = dataSources.Postgres
	a.stmt = statement.NewCustomerStatement(a.db)
	return nil
}

func (a *Customer) Insert(row *model.Customer) (int64, error) {
	var lastInsertId int64
	err := a.stmt.Insert.QueryRow(&row).Scan(&lastInsertId)
	return lastInsertId, err
}

func (a *Customer) FindByPhone(phone string) *model.Customer {
	var row model.Customer
	_ = a.stmt.FindByPhone.Get(&row, phone)
	return &row
}

func (a *Customer) FindByEmailOrPhone(phone string) *model.CustomerAuthentication {
	var row model.CustomerAuthentication
	_ = a.stmt.FindByEmailOrPhone.Get(&row, phone)
	return &row
}

func (a *Customer) BlockAccount(phone string) error {
	return nil
}

func (a *Customer) UnBlockAccount(phone string) error {
	return nil
}
