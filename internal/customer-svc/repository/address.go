package repository

import (
	"database/sql"
	"errors"

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

func (a *Address) FindByCustomerId(id int64) (*model.Address, error) {
	var row model.Address
	err := a.stmt.FindByCustomerId.Get(&row, id)
	return &row, err
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

func (a *Address) InsertOrUpdate(row *model.Address) error {
	// find by customer id
	address, err := a.FindByCustomerId(row.CustomerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			address = nil
			log.Errorf("Address.FindByCustomerId: %v no rows.", row.CustomerId)
		} else {
			log.Errorf("Error FindByCustomerId.", row.CustomerId)
		}
	}

	if address != nil {
		result, err := a.stmt.Update.Exec(row)
		if err != nil {
			log.Errorf("Update address by customerId error.")
			return err
		}
		return nsql.IsUpdated(result)

	} else {
		err = a.Insert(row)
		if err != nil {
			log.Errorf("Insert address by customerId error. %v", err)
			return err
		}
		return nil
	}
}
