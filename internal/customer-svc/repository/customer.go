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

func (a *Customer) HasInitialized() bool {
	return true
}

func (a *Customer) Init(dataSources DataSourceMap, _ contract.RepositoryMap) error {
	a.db = dataSources.DBInternal
	a.stmt = statement.NewCustomerStatement(a.db)
	return nil
}

func (a *Customer) Insert(row *model.Customer) (int64, error) {
	var lastInsertId int64
	err := a.stmt.Insert.QueryRow(&row).Scan(&lastInsertId)
	return lastInsertId, err
}

func (a *Customer) FindById(id int64) (*model.Customer, error) {
	var row model.Customer
	err := a.stmt.FindById.Get(&row, id)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (a *Customer) FindByPhone(phone string) (*model.Customer, error) {
	var row model.Customer
	err := a.stmt.FindByPhone.Get(&row, phone)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (a *Customer) FindByEmail(email string) (*model.Customer, error) {
	var row model.Customer
	err := a.stmt.FindByEmail.Get(&row, email)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (a *Customer) FindByEmailOrPhone(email string) (*model.Customer, error) {
	var row model.Customer
	err := a.stmt.FindByEmailOrPhone.Get(&row, email)
	return &row, err
}

func (a *Customer) UpdateByPhone(row *model.Customer) error {
	result, err := a.stmt.UpdateByPhone.Exec(row)
	if err != nil {
		return err
	}
	return nsql.IsUpdated(result)
}

func (a *Customer) BlockAccount(phone string) error {
	//TODO QUERY
	return nil
}

func (a *Customer) UnBlockAccount(phone string) error {
	//TODO QUERY
	return nil
}
