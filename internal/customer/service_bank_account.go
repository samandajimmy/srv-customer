package customer

import (
	"database/sql"
	"errors"
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
	customer, err := s.findOrFailCustomerByUserRefID(userRefID)
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

func (s *Service) CreateBankAccount(payload *dto.CreateBankAccountPayload) (*dto.GetDetailBankAccountResult, error) {
	// Find customer
	customer, err := s.findOrFailCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		return nil, err
	}

	// Account Number exist
	bankAccount, err := s.repo.FindBankAccountByAccountNumberAndCustomerID(payload.AccountNumber, customer.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errx.Trace(err)
	}

	if bankAccount.ID != 0 {
		return composeDetailBankAccount(bankAccount)
	}

	// Initialize data to insert
	xid, err := gonanoid.Generate(constant.AlphaNumUpperCaseRandomSet, 8)
	if err != nil {
		panic(fmt.Errorf("failed to generate xid. Error = %w", err))
	}

	bankAccount = &model.BankAccount{
		XID:           xid,
		CustomerID:    customer.ID,
		AccountNumber: payload.AccountNumber,
		AccountName:   payload.AccountName,
		Bank:          model.ToBank(payload.Bank),
		Status:        constant.Enabled,
		BaseField:     model.NewBaseField(model.ToModifier(payload.Subject.ModifiedBy())),
	}

	err = s.repo.CreateBankAccount(bankAccount)
	if err != nil {
		return nil, errx.Trace(err)
	}

	return composeDetailBankAccount(bankAccount)
}

func (s *Service) GetDetailBankAccount(payload *dto.GetDetailBankAccountPayload) (*dto.GetDetailBankAccountResult, error) {
	// Find customer
	customer, err := s.findOrFailCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		return nil, err
	}

	// Ownership validate
	if customer.UserRefID.String != payload.UserRefID {
		return nil, constant.BankAccountNotFoundError
	}

	// Get bank account by xid
	bankAccount, err := s.repo.FindBankAccountByXID(payload.XID)
	if err != nil {
		return nil, constant.BankAccountNotFoundError
	}

	return composeDetailBankAccount(bankAccount)
}

func (s *Service) UpdateBankAccount(payload *dto.UpdateBankAccountPayload) (*dto.GetDetailBankAccountResult, error) {
	// Find customer
	customer, err := s.findOrFailCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		return nil, err
	}

	// Get from repo
	bankAccount, err := s.repo.FindBankAccountByXID(payload.XID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errx.Trace(err)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, constant.ResourceNotFoundError
	}

	// Ownership validate
	if customer.UserRefID.String != payload.UserRefID {
		return nil, constant.ResourceNotFoundError.AddMetadata("message", "Bank account not found")
	}

	// Version validate
	if bankAccount.Version != payload.Version {
		return nil, constant.StaleResourceError.Trace()
	}

	// Update model
	bankAccount.AccountName = payload.AccountName
	bankAccount.AccountNumber = payload.AccountNumber
	bankAccount.Status = payload.Status
	bankAccount.Bank = model.ToBank(payload.Bank)
	bankAccount.BaseField = bankAccount.Upgrade(model.ToModifier(payload.Subject.ModifiedBy()))

	// Persist
	err = s.repo.UpdateBankAccount(bankAccount)
	if err != nil {
		return nil, errx.Trace(err)
	}

	return composeDetailBankAccount(bankAccount)
}

func (s *Service) DeleteBankAccount(payload *dto.GetDetailBankAccountPayload) error {
	// Find customer
	customer, err := s.findOrFailCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		return err
	}

	// Ownership validate
	if customer.UserRefID.String != payload.UserRefID {
		return constant.BankAccountNotFoundError
	}

	err = s.repo.DeleteBankAccountByXID(payload.XID)
	if err != nil {
		return errx.Trace(err)
	}

	return nil
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
