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

type Verification struct {
	db   *nsql.DB
	stmt *statement.VerificationStatement
}

func (a *Verification) HasInitialized() bool {
	return true
}

func (a *Verification) Init(dataSources DataSourceMap, _ contract.RepositoryMap) error {
	a.db = dataSources.DBInternal
	a.stmt = statement.NewVerificationStatement(a.db)
	return nil
}

func (a *Verification) Insert(row *model.Verification) error {
	_, err := a.stmt.Insert.Exec(row)
	return err
}

func (a *Verification) FindByCustomerId(id int64) (*model.Verification, error) {
	var row model.Verification
	err := a.stmt.FindByCustomerID.Get(&row, id)
	return &row, err
}

func (a *Verification) DeleteByID(id string) error {
	_, err := a.stmt.DeleteByID.Exec(id)
	return ncore.TraceError(err)
}

func (a *Verification) UpdateByCustomerID(row *model.Verification) error {
	_, err := a.stmt.Update.Exec(row)
	return ncore.TraceError(err)
}

func (a *Verification) InsertOrUpdate(row *model.Verification) error {
	// find by customer id
	verification, err := a.FindByCustomerId(row.CustomerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			verification = nil
			log.Errorf("Verification.FindByCustomerId: %v no rows.", row.CustomerId)
		} else {
			log.Errorf("Error FindByCustomerId.", row.CustomerId)
		}
	}

	if verification != nil {
		result, err := a.stmt.Update.Exec(row)
		if err != nil {
			log.Errorf("Update verification by customerId error.")
			return err
		}
		return nsql.IsUpdated(result)

	} else {
		err = a.Insert(row)
		if err != nil {
			log.Errorf("Insert verification by customerId error. %v", err)
			return err
		}
		return nil
	}
}
