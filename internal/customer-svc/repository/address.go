package repository

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/repository/statement"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type Address struct {
	db   *nsql.DB
	stmt *statement.AddressStatement
}

func (a *Address) HasInitialized() bool {
	return true
}

func (a *Address) Init(dataSources DataSourceMap, _ contract.RepositoryMap) error {
	a.db = dataSources.DBInternal
	a.stmt = statement.NewAddressStatement(a.db)
	return nil
}

func (a *Address) Insert(row *model.Address) error {
	_, err := a.stmt.Insert.Exec(row)
	return err
}

func (a *Address) Update(row *model.Address) error {
	result, err := a.stmt.Update.Exec(row)
	if err != nil {
		return err
	}
	return nsql.IsUpdated(result)
}

func (a *Address) FindPrimaryAddress(customerId int64) (*model.Address, error) {
	var row model.Address
	err := a.stmt.FindPrimaryByCustomerId.Get(&row, customerId)
	if err != nil {
		return nil, err
	}
	return &row, nil
}
