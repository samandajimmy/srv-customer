package customer

import (
	"database/sql"
	"errors"
	"github.com/nbs-go/errx"
	logOption "github.com/nbs-go/nlogger/v2/option"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
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
		return errx.Trace(err)
	}
	if !nsql.IsUpdated(result) {
		return constant.ResourceNotFoundError
	}
	return nil
}

func (rc *RepositoryContext) UpdateFinancialData(row *model.FinancialData) error {
	result, err := rc.stmt.FinancialData.Update.ExecContext(rc.ctx, row)
	if err != nil {
		return errx.Trace(err)
	}
	if !nsql.IsUpdated(result) {
		return constant.ResourceNotFoundError
	}
	return nil
}

func (rc *RepositoryContext) InsertOrUpdateFinancialData(row *model.FinancialData) error {
	_, err := rc.FindFinancialDataByCustomerID(row.CustomerID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errx.Trace(err)
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return rc.UpdateFinancialData(row)
	}

	err = rc.CreateFinancialData(row)
	if err != nil {
		rc.log.Error("cannot create financial data", logOption.Error(err))
		return errx.Trace(err)
	}

	return nil
}

func (rc *RepositoryContext) DeleteFinancialData(id string) error {
	result, err := rc.stmt.FinancialData.DeleteByID.ExecContext(rc.ctx, id)
	if err != nil {
		return errx.Trace(err)
	}
	if !nsql.IsUpdated(result) {
		return constant.ResourceNotFoundError
	}
	return nil
}

func (rc *RepositoryContext) UpdateGoldSavingStatus(financialData *model.FinancialData) error {
	err := rc.InsertOrUpdateFinancialData(financialData)
	return err
}
