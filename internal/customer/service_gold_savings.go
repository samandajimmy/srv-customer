package customer

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/nbs-go/errx"
	logOption "github.com/nbs-go/nlogger/v2/option"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
)

// Tabungan Emas Service

func (s *Service) validateCIF(cif string, id string) (bool, error) {
	c, err := s.repo.FindCustomerByUserRefID(id)
	if err != nil {
		return false, errx.Trace(err)
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
	_, err := s.validateCIF(cif, userRefId)
	if err != nil {
		return nil, errx.Trace(err)
	}

	// get customer from repo
	customer, err := s.repo.FindCustomerByPhoneOrCIF(cif)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("error found when get customer repo", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	if customer.Cif == "" {
		return nil, nil
	}

	// Call service portfolio
	switchingResponse, err := s.portfolioGoldSaving(customer.Cif)
	if err != nil {
		s.log.Error("error found when get gold saving", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	// Gold saving account
	accountSaving := &dto.GoldSavingVO{}

	// Get from cache when switching response is not success
	if switchingResponse.ResponseCode != "00" {
		return s.CacheGetGoldSavings(customer.Cif)
	}

	// Marshal response to account saving
	err = json.Unmarshal([]byte(switchingResponse.Message), accountSaving)
	if err != nil {
		s.log.Error("error found when unmarshall portfolio", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	// Get financial data
	financialData, err := s.repo.FindFinancialDataByCustomerID(customer.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("error found when get financial data repo", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		// init model
		financialData = model.EmptyFinancialData
		financialData.CustomerID = customer.ID
	}
	// Set primary account saving
	accountSaving.PrimaryRekening = s.getPrimaryAccountNumber(financialData, accountSaving.ListTabungan)

	// Update is open TE
	err = s.UpdateIsOpenGoldSavings(financialData, accountSaving, customer)
	if err != nil {
		s.log.Error("error found when open gold saving", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	// Set gold saving to cache
	err = s.CacheSetGoldSavings(customer.UserRefID.String, accountSaving)
	if err != nil {
		s.log.Error("error found when set gold saving to cache", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	return accountSaving, nil
}

func (s *Service) UpdateIsOpenGoldSavings(fd *model.FinancialData, as *dto.GoldSavingVO, c *model.Customer) error {
	if as == nil {
		return nil
	}
	// Set default status
	fd.GoldSavingStatus = 0

	if len(as.ListTabungan) > 0 {
		fd.GoldSavingStatus = 1
	}

	// Update financial data to repo
	err := s.repo.UpdateGoldSavingStatus(fd)
	if err != nil {
		s.log.Error("error found when update gold saving status", logOption.Error(err))
	}

	// Update referral code
	err = s.updateReferralCode(c)
	if err != nil {
		return errx.Trace(err)
	}

	return nil
}

func (s *Service) updateReferralCode(customer *model.Customer) error {
	if customer.ReferralCode.String != "" {
		return nil
	}

	// Set prefix referral code
	prefixReferralCode := "PDS"
	// Init referral code
	var newReferralCode string
	for {
		newReferralCode = generateReferralCode(prefixReferralCode)
		_, err := s.repo.ReferralCodeExists(newReferralCode)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			s.log.Error("error found when check referral code repo", logOption.Error(err))
			return errx.Trace(err)
		}

		if err != nil && errors.Is(err, sql.ErrNoRows) {
			break
		}
	}

	// Update referral code customer repo
	customer.ReferralCode = sql.NullString{String: newReferralCode, Valid: true}
	err := s.repo.UpdateCustomerByCIF(customer, customer.Cif)
	if err != nil {
		s.log.Error("error found when update customer repo", logOption.Error(err))
		return errx.Trace(err)
	}

	return nil
}
