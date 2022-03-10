package customer

import (
	"fmt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/nbs-go/errx"
	"github.com/nbs-go/nlogger/v2"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
)

func (s *Service) ListBankAccount(userRefID string, params *dto.ListPayload) (*dto.ListBankAccountResult, error) {
	// Find customer
	customer, err := s.repo.FindCustomerByUserRefID(userRefID)
	if err != nil {
		return nil, err
	}

	// Query
	queryResult, err := s.repo.ListBankAccount(customer.ID, params)
	if err != nil {
		s.log.Error("error found when get list bank account", nlogger.Error(err))
		return nil, errx.Trace(err)
	}

	// Compose response
	rowsResp := make([]dto.BankAccount, len(queryResult.Rows))
	for idx, row := range queryResult.Rows {
		rowsResp[idx] = model.ToBankAccountDTO(row)
	}

	return &dto.ListBankAccountResult{
		Rows:     rowsResp,
		Metadata: dto.ToListMetadata(params, queryResult.Count),
	}, nil
}

func (s *Service) CreateBankAccount(userRefID string, payload *dto.CreateBankAccountPayload) (*dto.GetDetailBankAccountResult, error) {
	// Find customer
	customer, err := s.repo.FindCustomerByUserRefID(userRefID)
	if err != nil {
		return nil, errx.Trace(err)
	}

	// Initialize data to insert
	xid, err := gonanoid.Generate(constant.AlphaNumUpperCaseRandomSet, 8)
	if err != nil {
		panic(fmt.Errorf("failed to generate xid. Error = %w", err))
	}

	bankAccount := model.BankAccount{
		XID:           xid,
		CustomerID:    customer.ID,
		AccountNumber: payload.AccountNumber,
		AccountName:   payload.AccountName,
		Bank:          model.ToBank(payload.Bank),
		Status:        constant.Active,
		BaseField:     model.NewBaseField(model.ToModifier(payload.Subject.ModifiedBy())),
	}

	err = s.repo.CreateBankAccount(bankAccount)
	if err != nil {
		return nil, errx.Trace(err)
	}

	return composeDetailBankAccount(&bankAccount)
}

func composeDetailBankAccount(row *model.BankAccount) (*dto.GetDetailBankAccountResult, error) {
	return &dto.GetDetailBankAccountResult{
		XID:           row.XID,
		AccountNumber: row.AccountNumber,
		AccountName:   row.AccountName,
		Bank:          model.ToBankDTO(row.Bank),
		BaseField:     model.ToBaseFieldDTO(&row.BaseField),
	}, nil
}
