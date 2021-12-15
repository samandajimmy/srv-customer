package repository

import (
	"database/sql"
	"errors"

	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/repository/statement"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type FinancialData struct {
	db   *nsql.DB
	stmt *statement.FinancialDataStatement
}

func (a *FinancialData) HasInitialized() bool {
	return true
}

func (a *FinancialData) Init(dataSources DataSourceMap, _ contract.RepositoryMap) error {
	a.db = dataSources.DBInternal
	a.stmt = statement.NewFinancialDataStatement(a.db)
	return nil
}

func (a *FinancialData) Insert(row *model.FinancialData) error {
	_, err := a.stmt.Insert.Exec(row)
	return err
}

func (a *FinancialData) FindByCustomerId(id int64) (*model.FinancialData, error) {
	var row model.FinancialData
	err := a.stmt.FindByCustomerID.Get(&row, id)
	return &row, err
}

func (a *FinancialData) DeleteByID(id string) error {
	_, err := a.stmt.DeleteByID.Exec(id)
	return ncore.TraceError(err)
}

func (a *FinancialData) UpdateByCustomerID(row *model.FinancialData) error {
	_, err := a.stmt.Update.Exec(row)
	return ncore.TraceError(err)
}

func (a *FinancialData) InsertOrUpdate(row *model.FinancialData) error {
	// find by customer id
	financialData, err := a.FindByCustomerId(row.CustomerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			financialData = nil
			log.Errorf("FinancialData.FindByCustomerId: %v no rows.", row.CustomerId)
		} else {
			log.Errorf("Error FindByCustomerId.", row.CustomerId)
		}
	}
	if financialData != nil {
		result, err := a.stmt.Update.Exec(row)
		if err != nil {
			log.Errorf("Update financialData by customerId error.")
			return err
		}
		return nsql.IsUpdated(result)

	} else {
		err = a.Insert(row)
		if err != nil {
			log.Errorf("Insert financialData by customerId error. %v", err)
			return err
		}
		return nil
	}
}
