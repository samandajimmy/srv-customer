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
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
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

func (s *Service) UpdateAvatar(payload dto.UpdateUserFile) (*dto.UploadResponse, error) {
	// Set payload avatar
	payloadUserFile := dto.UploadUserFilePayload{
		Request:   payload.Request,
		AssetType: payload.AssetType,
	}

	// Upload Userfile
	userFile, err := s.uploadUserFile(payloadUserFile)
	if err != nil {
		s.log.Error("error found when upload user file", nlogger.Error(err))
		return nil, err
	}

	// Find customer
	customer, err := s.repo.FindCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		s.log.Error("error when find current customer", nlogger.Error(err))
		return nil, ncore.TraceError("", err)
	}

	// Remove old avatar if exist
	if photo := customer.Photos; photo != nil && photo.FileName != "" {
		_ = s.AssetRemoveFile(photo.FileName, payload.AssetType)
	}

	// Update photo entity
	customer.Photos = &model.CustomerPhoto{
		Xid:      xid.New().String(),
		FileName: userFile.FileName,
		FileSize: userFile.FileSize,
		Mimetype: userFile.MimeType,
	}
	// Update timestamp profile
	customer.Profile.ProfileUpdatedAt = time.Now().Unix()

	// Persist
	err = s.repo.UpdateCustomerByUserRefID(customer, payload.UserRefID)
	if err != nil {
		s.log.Error("error when update photo profile", nlogger.Error(err))
		return nil, ncore.TraceError("", err)
	}

	return userFile, nil
}

func (s *Service) UpdateIdentity(payload dto.UpdateUserFile) (*dto.UploadResponse, error) {
	// Set payload avatar
	payloadUserFile := dto.UploadUserFilePayload{
		Request:   payload.Request,
		AssetType: payload.AssetType,
	}

	// Upload Userfile
	userFile, err := s.uploadUserFile(payloadUserFile)
	if err != nil {
		s.log.Error("error found when upload user file", nlogger.Error(err))
		return nil, err
	}

	// Find customer
	customer, err := s.repo.FindCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		s.log.Error("error when find current customer", nlogger.Error(err))
		return nil, ncore.TraceError("", err)
	}

	// Remove old identity if exist
	if oldFile := customer.Profile.IdentityPhotoFile; oldFile != "" {
		_ = s.AssetRemoveFile(oldFile, payload.AssetType)
	}

	// Update identity photo
	customer.Profile.IdentityPhotoFile = userFile.FileName

	// Persist
	err = s.repo.UpdateCustomerByUserRefID(customer, payload.UserRefID)
	if err != nil {
		s.log.Error("error when update identity photo", nlogger.Error(err))
		return nil, ncore.TraceError("", err)
	}

	return userFile, nil
}

func (s *Service) UpdateNPWP(payload dto.UpdateNPWPRequest) (*dto.UploadResponse, error) {
	// Set payload avatar
	payloadUserFile := dto.UploadUserFilePayload{
		Request:   payload.Request,
		AssetType: constant.AssetNPWP,
	}

	// Upload Userfile
	userFile, err := s.uploadUserFile(payloadUserFile)
	if err != nil {
		s.log.Error("error found when upload user file", nlogger.Error(err))
		return nil, err
	}

	// Find customer
	customer, err := s.repo.FindCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		s.log.Error("error when find current customer", nlogger.Error(err))
		return nil, ncore.TraceError("", err)
	}

	// remove old file if exist
	if oldFile := customer.Profile.NPWPPhotoFile; oldFile != "" {
		_ = s.AssetRemoveFile(oldFile, constant.AssetNPWP)
	}

	// Update NPWP entity
	customer.Profile.NPWPPhotoFile = userFile.FileName
	customer.Profile.NPWPNumber = payload.NoNPWP
	customer.Profile.NPWPUpdatedAt = time.Now().Unix()

	// Persist
	err = s.repo.UpdateCustomerByUserRefID(customer, payload.UserRefID)
	if err != nil {
		s.log.Error("error found when update NPWP", nlogger.Error(err))
		return nil, ncore.TraceError("", err)
	}

	return userFile, nil
}

func (s *Service) UpdateSID(payload dto.UpdateSIDRequest) (*dto.UploadResponse, error) {
	// Set payload avatar
	payloadUserFile := dto.UploadUserFilePayload{
		Request:   payload.Request,
		AssetType: constant.AssetNPWP,
	}

	// Upload Userfile
	userFile, err := s.uploadUserFile(payloadUserFile)
	if err != nil {
		s.log.Error("error found when upload user file", nlogger.Error(err))
		return nil, err
	}

	// Find customer
	customer, err := s.repo.FindCustomerByUserRefID(payload.UserRefID)
	if err != nil {
		s.log.Error("error when find current customer", nlogger.Error(err))
		return nil, ncore.TraceError("", err)
	}

	// remove old file if exist
	if oldFile := customer.Profile.SidPhotoFile; oldFile != "" {
		_ = s.AssetRemoveFile(oldFile, constant.AssetNPWP)
	}

	// Update SID entity
	customer.Profile.SidPhotoFile = userFile.FileName
	customer.Sid = payload.NoSID
	customer.Profile.NPWPUpdatedAt = time.Now().Unix()

	// Persist
	err = s.repo.UpdateCustomerByUserRefID(customer, payload.UserRefID)
	if err != nil {
		s.log.Error("error found when update SID", nlogger.Error(err))
		return nil, ncore.TraceError("", err)
	}

	return userFile, nil
}

func (s *Service) CheckStatus(userRefID string) (*dto.CheckStatusResponse, error) {
	// Find customer
	customer, verification, credential, err := s.repo.FindCombineCustomerDataByUserRefID(userRefID)
	if err != nil {
		s.log.Error("error when find current customer", nlogger.Error(err))
		return nil, ncore.TraceError("", err)
	}

	// Check pin available
	var pinAvailable = false
	if credential.Pin != "" {
		pinAvailable = true
	}
	var status = &dto.CheckStatusResponse{
		Cif:                    customer.Cif,
		EmailVerified:          nval.ParseBooleanFallback(verification.EmailVerifiedStatus, false),
		KycVerified:            nval.ParseBooleanFallback(verification.KycVerifiedStatus, false),
		PinAvailable:           pinAvailable,
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

func (s *Service) uploadUserFile(payload dto.UploadUserFilePayload) (*dto.UploadResponse, error) {
	// Parse multipart file
	file, err := nhttp.GetFile(payload.Request, constant.KeyUserFile, nhttp.MaxFileSizeImage, nhttp.MimeTypesImage)
	if err != nil {
		return nil, ncore.TraceError("", err)
	}

	// Upload file payload
	filePayload := dto.UploadRequest{
		AssetType: payload.AssetType,
		File:      file,
	}

	// Upload a file
	uploaded, err := s.AssetUploadFile(filePayload)
	if err != nil {
		s.log.Error("error found when call service", nlogger.Error(err))
		return nil, constant.UploadFileError.Wrap(err)
	}

	return uploaded, nil
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

func (s *Service) composeProfileResponse(customer *model.Customer, address *model.Address, financial *model.FinancialData,
	verification *model.Verification, gs interface{}) dto.ProfileResponse {
	avatarURL := s.AssetGetPublicURL(constant.AssetAvatarProfile, customer.Photos.FileName)
	ktpURL := s.AssetGetPublicURL(constant.AssetKTP, customer.Profile.IdentityPhotoFile)
	npwpURL := s.AssetGetPublicURL(constant.AssetNPWP, customer.Profile.NPWPPhotoFile)
	sidURL := s.AssetGetPublicURL(constant.AssetNPWP, customer.Profile.SidPhotoFile)

	return dto.ProfileResponse{
		CustomerVO: dto.CustomerVO{
			ID:                        customer.UserRefID.String,
			Cif:                       customer.Cif,
			Nama:                      customer.FullName,
			NamaIbu:                   customer.Profile.MaidenName,
			NoKTP:                     customer.IdentityNumber,
			Email:                     customer.Email,
			JenisKelamin:              customer.Profile.Gender,
			TempatLahir:               customer.Profile.PlaceOfBirth,
			TglLahir:                  customer.Profile.DateOfBirth,
			ReferralCode:              customer.ReferralCode,
			NoHP:                      customer.Phone,
			Kewarganegaraan:           customer.Profile.Nationality,
			NoIdentitas:               customer.IdentityNumber,
			TglExpiredIdentitas:       customer.Profile.IdentityExpiredAt,
			StatusKawin:               customer.Profile.MarriageStatus,
			NoNPWP:                    customer.Profile.NPWPNumber,
			NoSid:                     customer.Sid,
			IsKYC:                     nval.ParseStringFallback(verification.KycVerifiedStatus, ""),
			JenisIdentitas:            nval.ParseStringFallback(customer.IdentityType, ""),
			FotoNPWP:                  npwpURL,
			FotoSid:                   sidURL,
			Avatar:                    avatarURL,
			FotoKTP:                   ktpURL,
			Alamat:                    address.Line.String,
			IDProvinsi:                address.ProvinceID.Int64,
			IDKabupaten:               address.CityID.Int64,
			IDKecamatan:               address.DistrictID.Int64,
			IDKelurahan:               address.SubDistrictID.Int64,
			Kelurahan:                 address.SubDistrictName.String,
			Provinsi:                  address.ProvinceName.String,
			Kabupaten:                 address.CityName.String,
			Kecamatan:                 address.DistrictName.String,
			KodePos:                   address.PostalCode.String,
			Norek:                     financial.AccountNumber,
			GoldCardApplicationNumber: financial.GoldCardApplicationNumber,
			GoldCardAccountNumber:     financial.GoldCardAccountNumber,
			Saldo:                     nval.ParseStringFallback(financial.Balance, "0"),
			IsOpenTe:                  nval.ParseStringFallback(financial.GoldSavingStatus, "0"),
			IsEmailVerified:           nval.ParseStringFallback(verification.EmailVerifiedStatus, "0"),
			IsDukcapilVerified:        nval.ParseStringFallback(verification.DukcapilVerifiedStatus, "0"),
			AktifasiTransFinansial:    nval.ParseStringFallback(verification.FinancialTransactionStatus, "0"),
			KodeCabang:                "", // TODO Branch Code
			TabunganEmas:              gs,
		},
	}
}
