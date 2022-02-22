package customer

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/nbs-go/nlogger"
	"github.com/rs/xid"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
	"time"
)

func (s *Service) CustomerProfile(id string) (*dto.ProfileResponse, error) {
	// Get customer data
	customer, err := s.repo.FindCustomerByUserRefID(id)
	if err != nil {
		s.log.Error("error found when get customer repo", nlogger.Error(err))
		return nil, ncore.TraceError("", err)
	}

	// Get verification data
	verification, err := s.repo.FindVerificationByCustomerID(customer.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("error found when get verification repo", nlogger.Error(err))
		return nil, ncore.TraceError("", err)
	}

	// Get financial data
	financial, err := s.repo.FindFinancialDataByCustomerID(customer.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("error found when get financial repo", nlogger.Error(err))
		return nil, ncore.TraceError("", err)
	}

	// Get address
	address, err := s.repo.FindAddressByCustomerId(customer.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("error found when get address repo", nlogger.Error(err))
		return nil, ncore.TraceError("", err)
	}

	// Get Gold saving account
	goldSaving, err := s.getListAccountNumber(customer.Cif, customer.UserRefID.String)
	if err != nil {
		s.log.Error("error found when get list gold saving service", nlogger.Error(err))
		return nil, ncore.TraceError("", err)
	}

	var gs interface{}
	if goldSaving == nil {
		gs = false
	} else {
		gs = &dto.GoldSavingVO{
			TotalSaldoBlokir:  goldSaving.TotalSaldoBlokir,
			TotalSaldoSeluruh: goldSaving.TotalSaldoSeluruh,
			TotalSaldoEfektif: goldSaving.TotalSaldoEfektif,
			ListTabungan:      goldSaving.ListTabungan,
			PrimaryRekening:   goldSaving.PrimaryRekening,
		}
	}

	// Compose response
	resp := s.composeProfileResponse(customer, address, financial, verification, gs)

	return &resp, nil
}

func (s *Service) UpdateCustomerProfile(id string, payload dto.UpdateProfileRequest) error {
	// Get current customer data
	customer, err := s.repo.FindCustomerByUserRefID(id)
	if err != nil {
		s.log.Error("error when get customer data", nlogger.Error(err))
		return ncore.TraceError("", err)
	}

	// Update customer profile repo
	err = s.repo.UpdateCustomerProfile(customer, payload)
	if err != nil {
		s.log.Error("error found when get customer repo", nlogger.Error(err))
		return ncore.TraceError("", err)
	}

	return nil
}

func (s *Service) isValidPassword(userRefID string, password string) (bool, error) {
	// Find customer
	c, err := s.repo.FindCustomerByUserRefID(userRefID)
	if err != nil {
		return false, err
	}

	pw := stringToMD5(password)

	// Check is valid
	_, err = s.repo.IsValidPassword(c.ID, pw)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("error when check password match", nlogger.Error(err))
		return false, ncore.TraceError("error when validate password match", err)
	}

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	return true, nil
}

func (s *Service) UpdatePassword(userRefID string, payload dto.UpdatePasswordRequest) error {
	// Find customer
	customer, err := s.repo.FindCustomerByUserRefID(userRefID)
	if err != nil {
		return err
	}

	// Set current password to md5
	currentPassword := stringToMD5(payload.CurrentPassword)

	// Validate current password
	_, err = s.repo.IsValidPassword(customer.ID, currentPassword)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("error found when check password match", nlogger.Error(err))
		return ncore.TraceError("error when validate password match", err)
	}

	// If password doesn't match
	if errors.Is(err, sql.ErrNoRows) {
		return s.responses.GetError("E_USR_1")
	}

	// Set password to md5
	password := stringToMD5(payload.NewPassword)

	// Update password
	err = s.repo.UpdatePassword(customer.ID, password)
	if err != nil {
		s.log.Error("error found when check password match", nlogger.Error(err))
		return ncore.TraceError("failed to update password", err)
	}

	return nil
}

func (s *Service) UpdateAvatar(userRefID string, uploaded *dto.UploadResponse) error {
	// Find customer
	customer, err := s.repo.FindCustomerByUserRefID(userRefID)
	if err != nil {
		s.log.Error("error when find current customer", nlogger.Error(err))
		return ncore.TraceError("", err)
	}

	// Remove old avatar if exist
	if photo := customer.Photos; photo != nil && photo.FileName != "" {
		_ = s.AssetRemoveFile(photo.FileName, constant.AssetAvatarProfile)
	}

	// Create new photo
	newPhoto := &model.CustomerPhoto{
		Xid:      xid.New().String(),
		FileName: uploaded.FileName,
		FileSize: uploaded.FileSize,
		Mimetype: uploaded.MimeType,
	}

	customer.Photos = newPhoto
	customer.Profile.ProfileUpdatedAt = time.Now().Unix()

	// Persist
	err = s.repo.UpdateCustomerByUserRefID(customer, userRefID)
	if err != nil {
		s.log.Error("error when update photo profile", nlogger.Error(err))
		return ncore.TraceError("", err)
	}

	return nil
}

func (s *Service) UpdateIdentity(userRefID string, uploaded *dto.UploadResponse) error {
	// Find customer
	customer, err := s.repo.FindCustomerByUserRefID(userRefID)
	if err != nil {
		s.log.Error("error when find current customer", nlogger.Error(err))
		return ncore.TraceError("", err)
	}

	// Remove old avatar if exist
	if oldFile := customer.Profile.IdentityPhotoFile; oldFile != "" {
		_ = s.AssetRemoveFile(oldFile, constant.AssetKTP)
	}

	// Update identity photo
	customer.Profile.IdentityPhotoFile = uploaded.FileName

	// Persist
	err = s.repo.UpdateCustomerByUserRefID(customer, userRefID)
	if err != nil {
		s.log.Error("error when update identity photo", nlogger.Error(err))
		return ncore.TraceError("", err)
	}

	return nil
}

func (s *Service) UpdateNPWP(userRefID string, npwpNumber string, uploaded *dto.UploadResponse) error {
	// Find customer
	customer, err := s.repo.FindCustomerByUserRefID(userRefID)
	if err != nil {
		s.log.Error("error when find current customer", nlogger.Error(err))
		return ncore.TraceError("", err)
	}

	// remove old file if exist
	if oldFile := customer.Profile.NPWPPhotoFile; oldFile != "" {
		_ = s.AssetRemoveFile(oldFile, constant.AssetNPWP)
	}

	// Update NPWP entity
	customer.Profile.NPWPPhotoFile = uploaded.FileName
	customer.Profile.NPWPNumber = npwpNumber
	customer.Profile.NPWPUpdatedAt = time.Now().Unix()

	// Persist
	err = s.repo.UpdateCustomerByUserRefID(customer, userRefID)
	if err != nil {
		s.log.Error("error found when update NPWP", nlogger.Error(err))
		return ncore.TraceError("", err)
	}

	return nil
}

func (s *Service) UpdateSID(userRefID string, sidNumber string, uploaded *dto.UploadResponse) error {
	// Find customer
	customer, err := s.repo.FindCustomerByUserRefID(userRefID)
	if err != nil {
		s.log.Error("error when find current customer", nlogger.Error(err))
		return ncore.TraceError("", err)
	}

	// remove old file if exist
	if oldFile := customer.Profile.SidPhotoFile; oldFile != "" {
		_ = s.AssetRemoveFile(oldFile, constant.AssetNPWP)
	}

	// Update NPWP entity
	customer.Profile.SidPhotoFile = uploaded.FileName
	customer.Sid = sidNumber
	customer.Profile.NPWPUpdatedAt = time.Now().Unix()

	// Persist
	err = s.repo.UpdateCustomerByUserRefID(customer, userRefID)
	if err != nil {
		s.log.Error("error found when update NPWP", nlogger.Error(err))
		return ncore.TraceError("", err)
	}

	return nil
}

func (s *Service) CheckStatus(userRefID string) (*dto.CheckStatusResponse, error) {
	// Find customer
	customer, verification, credential, err := s.repo.FindCombineCustomerDataByUserRefID(userRefID)
	if err != nil {
		s.log.Error("error when find current customer", nlogger.Error(err))
		return nil, ncore.TraceError("", err)
	}

	var status = &dto.CheckStatusResponse{
		Cif:                    customer.Cif,
		EmailVerified:          nval.ParseBooleanFallback(verification.EmailVerifiedStatus, false),
		KycVerified:            false,
		PinAvailable:           nval.ParseBooleanFallback(credential.Pin, false),
		AktifasiTransFinansial: nval.ParseStringFallback(verification.FinancialTransactionStatus, "0"),
	}

	// If cif is empty
	if customer.Cif == "" {
		return status, nil
	}

	// Check CIF
	checkCifResponse, err := s.CheckCIF(customer.Cif)
	if err != nil {
		s.log.Error("error found when check CIF", nlogger.Error(err))
		return nil, ncore.TraceError("", err)
	}

	if checkCifResponse.ResponseCode != "00" {
		s.log.Error("error response code from check CIF", nlogger.Error(err))
		return status, nil
	}

	if checkCifResponse.Message == "" {
		s.log.Error("error response message from check CIF", nlogger.Error(err))
		return status, nil
	}

	customerInquiry := &dto.CustomerInquiryVO{}
	err = json.Unmarshal([]byte(checkCifResponse.Message), customerInquiry)
	if err != nil {
		s.log.Error("error marshall customer inquiry", nlogger.Error(err))
		return nil, constant.DefaultError.Trace()
	}

	// Update KYC
	status, err = s.profileUpdateKyc(customerInquiry, verification, status)
	if err != nil {
		s.log.Error("error when update kyc", nlogger.Error(err))
		return nil, err
	}

	return status, nil
}

func (s *Service) profileUpdateKyc(customerInquiry *dto.CustomerInquiryVO, verification *model.Verification, status *dto.CheckStatusResponse) (*dto.CheckStatusResponse, error) {
	if customerInquiry.StatusKyc == "1" && verification.KycVerifiedStatus == 0 {
		verification.KycVerifiedAt = sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		}
	}

	verification.KycVerifiedStatus = nval.ParseInt64Fallback(customerInquiry.StatusKyc, 0)
	err := s.repo.UpdateVerification(verification)
	if err != nil {
		s.log.Error("error found when update verification repo", nlogger.Error(err))
		return nil, ncore.TraceError("", err)
	}

	status.KycVerified = nval.ParseBooleanFallback(verification.KycVerifiedStatus, false)

	return status, nil
}

// Endpoint POST /customer/inquiry

func (s *Service) CheckCIF(cif string) (*ResponseSwitchingSuccess, error) {
	// Check CIF
	reqBody := map[string]interface{}{
		"cif": cif,
	}

	sp := PostDataPayload{
		Url:  "/customer/inquiry",
		Data: reqBody,
	}

	data, err := s.RestSwitchingPostData(sp)
	if err != nil {
		s.log.Error("error found when get gold savings", nlogger.Error(err))
		return nil, ncore.TraceError("error found when get gold savings", err)
	}

	return data, nil
}
