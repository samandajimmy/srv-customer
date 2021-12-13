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

type Credential struct {
	db   *nsql.DB
	stmt *statement.CredentialStatement
}

func (a *Credential) HasInitialized() bool {
	return true
}

func (a *Credential) Init(dataSources DataSourceMap, _ contract.RepositoryMap) error {
	a.db = dataSources.DBInternal
	a.stmt = statement.NewCredentialStatement(a.db)
	return nil
}

func (a *Credential) Insert(row *model.Credential) error {
	_, err := a.stmt.Insert.Exec(row)
	return err
}

func (a *Credential) FindByCustomerId(id int64) (*model.Credential, error) {
	var row model.Credential
	err := a.stmt.FindByCustomerID.Get(&row, id)
	return &row, err
}

func (a *Credential) DeleteByID(id string) error {
	_, err := a.stmt.DeleteByID.Exec(id)
	return ncore.TraceError(err)
}

func (a *Credential) UpdateByCustomerID(row *model.Credential) error {
	_, err := a.stmt.Update.Exec(row)
	return ncore.TraceError(err)
}

func (a *Credential) InsertOrUpdate(row *model.Credential) error {
	// find by customer id
	credential, err := a.FindByCustomerId(row.CustomerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			credential = nil
			log.Errorf("Credential.FindByCustomerId: %v no rows.", row.CustomerId)
		} else {
			log.Errorf("Error FindByCustomerId.", row.CustomerId)
		}
	}
	if credential != nil {
		result, err := a.stmt.Update.Exec(row)
		if err != nil {
			log.Errorf("Update credential by customerId error.")
			return err
		}
		return nsql.IsUpdated(result)

	} else {
		err = a.Insert(row)
		if err != nil {
			log.Errorf("Insert credential by customerId error. %v", err)
			return err
		}
		return nil
	}
}
