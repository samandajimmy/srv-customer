package customer

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/nbs-go/errx"
	"github.com/nbs-go/nlogger/v2"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
)

// Tabungan Emas Service

func (s *Service) validateCIF(cif string, id string) (bool, error) {
	c, err := s.repo.FindCustomerByUserRefID(id)
	if err != nil {
		return false, err
	}

	if cif != c.Cif {
		return false, nhttp.InternalError.Trace(errx.Errorf("unexpected cif value. Actual = %s", cif))
	}

	return true, nil
}

func (s *Service) getPrimaryAccountNumber(financial *model.FinancialData, list []dto.AccountSavingVO) *dto.AccountSavingVO {
	if len(list) == 0 {
		return nil
	}

	if financial == nil || financial.MainAccountNumber == "" {
		return &list[0]
	}

	for _, saving := range list {
		if saving.Norek == financial.MainAccountNumber {
			return &saving
		}
	}

	return nil
}

func (s *Service) getListAccountNumber(cif string, userRefId string) (*dto.GoldSavingVO, error) {
	// Get context
	ctx := s.ctx

	_, err := s.validateCIF(cif, userRefId)
	if err != nil {
		return nil, errx.Trace(err)
	}

	// get customer from repo
	customer, err := s.repo.FindCustomerByPhoneOrCIF(cif)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("error found when get customer repo", nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	if customer.Cif == "" {
		return nil, nil
	}

	// Call service portfolio
	switchingResponse, err := s.portfolioGoldSaving(customer.Cif)
	if err != nil {
		s.log.Error("error found when get gold saving", nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	// Gold saving account
	accountSaving := &dto.GoldSavingVO{}

	// Get from cache when switching response is not success
	if switchingResponse.ResponseCode != "00" {
		accountSaving, err = s.CacheGetGoldSavings(customer.Cif)
		if err != nil {
			return nil, errx.Trace(err)
		}
		return accountSaving, nil
	}

	// Marshal response to account saving
	err = json.Unmarshal([]byte(switchingResponse.Message), accountSaving)
	if err != nil {
		s.log.Error("error found when unmarshall portfolio", nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	// Get financial data
	financialData, err := s.repo.FindFinancialDataByCustomerID(customer.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("error found when get financial data repo", nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	if errors.Is(err, sql.ErrNoRows) {
		// init model
		financialData = model.EmptyFinancialData
		financialData.CustomerID = customer.ID
	}
	// Set primary account saving
	accountSaving.PrimaryRekening = s.getPrimaryAccountNumber(financialData, accountSaving.ListTabungan)

	// update is open TE
	err = s.UpdateIsOpenGoldSavings(financialData, accountSaving, customer)
	if err != nil {
		s.log.Error("error found when open gold saving.", nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	// Set gold saving to cache
	err = s.CacheSetGoldSavings(customer.UserRefID.String, accountSaving)
	if err != nil {
		s.log.Error("error found when set gold saving to cache", nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	return accountSaving, nil
}

func (s *Service) UpdateIsOpenGoldSavings(fd *model.FinancialData, as *dto.GoldSavingVO, c *model.Customer) error {
	// Update is open te
	if as != nil {
		if len(as.ListTabungan) > 0 {
			fd.GoldSavingStatus = 1
		} else {
			fd.GoldSavingStatus = 0
		}

		// Update financial data to repo
		err := s.repo.UpdateGoldSavingStatus(fd)
		if err != nil {
			s.log.Error("error found when update gold saving status", nlogger.Error(err), nlogger.Context(s.ctx))
		}

		// Update referral code
		err = s.updateReferralCode(c, fd)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) updateReferralCode(customer *model.Customer, financialData *model.FinancialData) error {
	// Set Prefix Referral code
	prefixReferralCode := "PDS"

	if customer.ReferralCode == "" && financialData.GoldSavingStatus == 1 {
		var newReferralCode string

		for {
			newReferralCode = generateReferralCode(prefixReferralCode)
			_, err := s.repo.ReferralCodeExists(newReferralCode)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				s.log.Error("error found when check referral code repo", nlogger.Error(err), nlogger.Context(s.ctx))
				return errx.Trace(err)
			}

			if errors.Is(err, sql.ErrNoRows) {
				break
			}
		}

		// Update referral code customer repo
		customer.ReferralCode = newReferralCode
		err := s.repo.UpdateCustomerByCIF(customer, customer.Cif)
		if err != nil {
			s.log.Error("error found when update customer repo", nlogger.Error(err), nlogger.Context(s.ctx))
			return err
		}
	}

	return nil
}
