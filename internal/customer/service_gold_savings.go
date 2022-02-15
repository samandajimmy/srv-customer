package customer

import (
	"database/sql"
	"encoding/json"
	"github.com/nbs-go/nlogger"
	"github.com/rs/xid"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/convert"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
)

// Tabungan Emas Service

func (s *Service) validateCIF(cif string, id string) (bool, error) {
	c, err := s.repo.FindCustomerByUserRefID(id)
	if err != nil {
		return false, err
	}

	if cif != c.Cif {
		return false, constant.DefaultError.Trace()
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
		return nil, ncore.TraceError("error when validate CIF", err)
	}

	// get customer from repo
	customer, err := s.repo.FindCustomerByPhoneOrCIF(cif)
	if err != nil {
		s.log.Error("error found when get customer repo", nlogger.Error(err), nlogger.Context(ctx))
		return nil, ncore.TraceError("failed get customer repo", err)
	}

	// Call service portfolio
	_, err = s.portfolioGoldSaving(customer.Cif)
	if err != nil {
		s.log.Error("error found when get gold saving", nlogger.Error(err), nlogger.Context(ctx))
		return nil, ncore.TraceError("failed to get portfolio gold saving", err)
	}

	// Mocking response tabungan emas
	resMock := `{
    	"responseCode": "00",
		"responseDesc": "Approved",
    	"data": "{\"totalSaldoBlokir\":\"0.0000\",\"totalSaldoSeluruh\":\"99.9672\",\"totalSaldoEfektif\":\"99.9672\",\"listTabungan\":[{\"cif\":\"1015419761\",\"kodeCabang\":\"13025\",\"namaNasabah\":\"HADIYU\",\"noBuku\":\"18051249\",\"norek\":\"1302519620000663\",\"saldoBlokir\":\"0.0000\",\"saldoEmas\":\"99.9672\",\"tglBuka\":\"2021-07-23\",\"saldoEfektif\":\"99.9672\"}]}"
	}`

	switchingResponse := ResponseSwitchingSuccess{}
	err = json.Unmarshal([]byte(resMock), &switchingResponse)
	if err != nil {
		s.log.Error("error found when unmarshal response", nlogger.Error(err), nlogger.Context(ctx))
		return nil, ncore.TraceError("error found when unmarshal response", err)
	}

	// account saving
	accountSaving := &dto.GoldSavingVO{}

	// get from cache when switching response is not success
	if switchingResponse.ResponseCode != "00" {
		accountSaving, err = s.CacheGetGoldSavings(customer.Cif)
		if err != nil {
			return nil, ncore.TraceError("error found when get cache gold saving", err)
		}
		return accountSaving, nil
	}

	// Marshal response to account saving
	err = json.Unmarshal([]byte(switchingResponse.Message), accountSaving)
	if err != nil {
		s.log.Error("error found when unmarshall portfolio", nlogger.Error(err), nlogger.Context(ctx))
		return nil, ncore.TraceError("error unmarshall data portfolio", err)
	}

	// Get financial data
	financialData, err := s.repo.FindFinancialDataByCustomerID(customer.Id)
	if err != nil && err != sql.ErrNoRows {
		s.log.Error("error found when get financial data repo", nlogger.Error(err), nlogger.Context(ctx))
		return nil, ncore.TraceError("error found get financial data", err)
	}

	if financialData == nil {
		// init model
		itemMetaData := model.NewItemMetadata(convert.ModifierDTOToModel(dto.Modifier{ID: "", Role: "", FullName: ""}))
		financialData = &model.FinancialData{
			Xid:                       xid.New().String(),
			CustomerId:                customer.Id,
			MainAccountNumber:         "",
			AccountNumber:             "",
			GoldSavingStatus:          0,
			GoldCardApplicationNumber: "",
			GoldCardAccountNumber:     "",
			Balance:                   0,
			Metadata:                  []byte("{}"),
			ItemMetadata:              itemMetaData,
		}
	}

	// Set primary account saving
	accountSaving.PrimaryRekening = s.getPrimaryAccountNumber(financialData, accountSaving.ListTabungan)

	// update is open TE
	err = s.UpdateIsOpenGoldSavings(financialData, accountSaving, customer)
	if err != nil {
		s.log.Error("error found when open gold saving.", nlogger.Error(err), nlogger.Context(ctx))
		return nil, ncore.TraceError("error found when open gold saving", err)
	}

	// Set gold saving to cache
	err = s.CacheSetGoldSavings(customer.UserRefId.String, accountSaving)
	if err != nil {
		s.log.Error("error found when set gold saving to cache", nlogger.Error(err), nlogger.Context(ctx))
		return nil, ncore.TraceError("error when set gold saving to cache", err)
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
			exists, err := s.repo.ReferralCodeExists(newReferralCode)
			if err != nil && err != sql.ErrNoRows {
				s.log.Error("error found when check referral code repo", nlogger.Error(err), nlogger.Context(s.ctx))
				return err
			}

			if exists == nil {
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
