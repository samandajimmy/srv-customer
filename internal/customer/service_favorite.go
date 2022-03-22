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
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

func (s *Service) CreateFavorite(payload *dto.CreateFavoritePayload) (*dto.Favorite, error) {
	// Find customer
	customer, err := s.findOrFailCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		return nil, err
	}

	// Account Number exist
	favorite, err := s.repo.FindFavoriteByAccountNumberAndCustomerID(payload.AccountNumber, customer.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errx.Trace(err)
	}

	if favorite.ID != 0 {
		return composeDetailFavorite(favorite)
	}

	// Initialize data to insert
	xid, err := gonanoid.Generate(constant.AlphaNumUpperCaseRandomSet, 8)
	if err != nil {
		panic(fmt.Errorf("failed to generate xid. Error = %w", err))
	}

	favorite = &model.TransactionFavorite{
		BaseField:       model.NewBaseField(model.ToModifier(payload.Subject.ModifiedBy())),
		XID:             xid,
		CustomerID:      customer.ID,
		Type:            payload.Type,
		TypeTransaction: payload.TypeTransaction,
		AccountNumber:   payload.AccountNumber,
		AccountName:     payload.AccountName,
		BankName:        nsql.NewNullString(payload.BankName),
		BankCode:        nsql.NewNullString(payload.BankCode),
		GroupMPO:        nsql.NewNullString(payload.GroupMPO),
		ServiceCodeMPO:  nsql.NewNullString(payload.ServiceCodeMPO),
	}

	err = s.repo.CreateFavorite(favorite)
	if err != nil {
		return nil, errx.Trace(err)
	}

	return composeDetailFavorite(favorite)
}
func (s *Service) ListFavorite(userRefID string, params *dto.ListPayload) (*dto.ListFavoriteResult, error) {
	// Find customer
	customer, err := s.findOrFailCustomerByUserRefID(userRefID)
	if err != nil {
		return nil, err
	}

	// Query
	queryResult, err := s.repo.ListFavorite(customer.ID, params)
	if err != nil {
		s.log.Error("error found when get list bank account", nlogger.Error(err))
		return nil, errx.Trace(err)
	}

	// Compose response
	rowsResp := make([]dto.Favorite, len(queryResult.Rows))
	for idx, row := range queryResult.Rows {
		rowsResp[idx] = model.ToTransactionFavoriteDTO(row)
	}

	return &dto.ListFavoriteResult{
		Rows:     rowsResp,
		Metadata: dto.ToListMetadata(params, queryResult.Count),
	}, nil
}

func (s *Service) DeleteFavorite(payload *dto.GetDetailFavoritePayload) error {
	// Find customer
	customer, err := s.findOrFailCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		return err
	}

	// Ownership validate
	if customer.UserRefID.String != payload.UserRefID {
		return constant.TransactionFavoriteNotFoundError
	}

	err = s.repo.DeleteTransactionFavoriteByXID(payload.XID)
	if err != nil {
		return errx.Trace(err)
	}

	return nil
}

func composeDetailFavorite(row *model.TransactionFavorite) (*dto.Favorite, error) {
	return &dto.Favorite{
		BaseField:       model.ToBaseFieldDTO(&row.BaseField),
		XID:             row.XID,
		Type:            row.Type,
		TypeTransaction: row.TypeTransaction,
		AccountName:     row.AccountName,
		AccountNumber:   row.AccountNumber,
		BankName:        row.BankName.String,
		BankCode:        row.BankCode.String,
		GroupMPO:        row.GroupMPO.String,
		ServiceCodeMPO:  row.ServiceCodeMPO.String,
	}, nil
}
