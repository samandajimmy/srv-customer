package customer

import (
	"github.com/nbs-go/errx"
	"github.com/nbs-go/nlogger/v2"
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
