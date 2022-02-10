package customer

import (
	"database/sql"
	"errors"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"

	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
)

func (rc *RepositoryContext) CreateFinancialData(row *model.FinancialData) error {
	_, err := rc.stmt.FinancialData.Insert.ExecContext(rc.ctx, row)
	if err != nil {
		return err
	}

	return nil
}

func (rc *RepositoryContext) FindFinancialDataByCustomerID(id int64) (*model.FinancialData, error) {
	var row model.FinancialData
	err := rc.stmt.FinancialData.FindByCustomerID.GetContext(rc.ctx, &row, id)
	return &row, err
}

func (rc *RepositoryContext) UpdateFinancialDataByCustomerID(row *model.FinancialData) error {
	result, err := rc.stmt.FinancialData.Update.ExecContext(rc.ctx, row)
	if err != nil {
		return ncore.TraceError("failed to update FinancialDataByCustomerId", err)
	}
	if !nsql.IsUpdated(result) {
		return constant.ResourceNotFoundError
	}
	return nil
}

func (rc *RepositoryContext) UpdateFinancialData(row *model.FinancialData) error {
	result, err := rc.stmt.FinancialData.Update.ExecContext(rc.ctx, row)
	if err != nil {
		return ncore.TraceError("failed to update financial data", err)
	}
	if !nsql.IsUpdated(result) {
		return constant.ResourceNotFoundError
	}
	return nil
}

func (rc *RepositoryContext) InsertOrUpdateFinancialData(row *model.FinancialData) error {
	// find by customer id
	financialData, err := rc.FindFinancialDataByCustomerID(row.CustomerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			financialData = nil
		}
	}
	if financialData != nil {
		err = rc.UpdateFinancialData(financialData)
		if err != nil {
			return ncore.TraceError("cannot update financial data", err)
		}
		return nil
	} else {
		err = rc.CreateFinancialData(row)
		if err != nil {
			return ncore.TraceError("cannot create financial data", err)
		}
		return nil
	}
}

func (rc *RepositoryContext) DeleteFinancialData(id string) error {
	result, err := rc.stmt.FinancialData.DeleteByID.ExecContext(rc.ctx, id)
	if err != nil {
		return ncore.TraceError("failed delete financial data", err)
	}
	if !nsql.IsUpdated(result) {
		return constant.ResourceNotFoundError
	}
	return nil
}
