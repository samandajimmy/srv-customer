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
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
	"strings"
	"time"
)

func (s *Service) CustomerProfile(id string) (*dto.ProfileResponse, error) {
	// Get Context
	ctx := s.ctx

	// Get customer data
	c, err := s.repo.FindCustomerByUserRefID(id)
	if err != nil {
		s.log.Error("error found when get customer repo", nlogger.Error(err), nlogger.Context(ctx))
		return nil, ncore.TraceError("error when", err)
	}

	// Get verification data
	v, err := s.repo.FindVerificationByCustomerID(c.Id)
	if err != nil && err != sql.ErrNoRows {
		s.log.Error("error found when get verification repo", nlogger.Error(err), nlogger.Context(ctx))
		return nil, ncore.TraceError("", err)
	}

	// Get address
	a, err := s.repo.FindAddressByCustomerId(c.Id)
	if err != nil && err != sql.ErrNoRows {
		s.log.Error("error found when get address repo", nlogger.Error(err), nlogger.Context(ctx))
		return nil, ncore.TraceError("", err)
	}

	//  gold saving account
	goldSaving, err := s.getListAccountNumber(c.Cif, c.UserRefId.String)
	if err != nil {
		return nil, ncore.TraceError("error when get list gold saving account", err)
	}

	gs := &dto.GoldSavingVO{
		TotalSaldoBlokir:  goldSaving.TotalSaldoBlokir,
		TotalSaldoSeluruh: goldSaving.TotalSaldoSeluruh,
		TotalSaldoEfektif: goldSaving.TotalSaldoEfektif,
		ListTabungan:      goldSaving.ListTabungan,
		PrimaryRekening:   goldSaving.PrimaryRekening,
	}
	// Get avatar
	avatarUrl := s.AssetGetPublicUrl(constant.AssetAvatarProfile, c.Photos.FileName)
	ktpUrl := s.AssetGetPublicUrl(constant.AssetKTP, c.Profile.IdentityPhotoFile)

	// Compose response
	resp := dto.ProfileResponse{
		CustomerVO: dto.CustomerVO{
			ID:              c.UserRefId.String,
			Cif:             c.Cif,
			Nama:            c.FullName,
			NamaIbu:         c.Profile.MaidenName,
			NoKTP:           c.IdentityNumber,
			Email:           c.Email,
			JenisKelamin:    c.Profile.Gender,
			TempatLahir:     c.Profile.PlaceOfBirth,
			TglLahir:        c.Profile.DateOfBirth,
			Kewarganegaraan: c.Profile.Nationality,
			NoNPWP:          c.Profile.NPWPNumber,
			NoHP:            c.Phone,
			Alamat:          a.Line.String,
			IDProvinsi:      a.ProvinceId.Int64,
			IDKabupaten:     a.CityId.Int64,
			IDKecamatan:     a.DistrictId.Int64,
			IDKelurahan:     a.SubDistrictId.Int64,
			Provinsi:        a.ProvinceName.String,
			Kabupaten:       a.CityName.String,
			Kecamatan:       a.DistrictName.String,
			Kelurahan:       a.SubDistrictName.String,
			KodePos:         a.PostalCode.String,
			Avatar:          avatarUrl,
			FotoKTP:         ktpUrl,
			IsEmailVerified: nval.ParseStringFallback(v.EmailVerifiedStatus, ""),
			JenisIdentitas:  nval.ParseStringFallback(c.IdentityType, ""),
			TabunganEmas:    gs,
		},
	}

	return &resp, nil
}

func (s *Service) UpdateCustomerProfile(id string, payload dto.UpdateProfileRequest) error {
	ctx := s.ctx

	// Get current customer data
	customer, err := s.repo.FindCustomerByUserRefID(id)
	if err != nil {
		s.log.Error("error when get customer data", nlogger.Error(err), nlogger.Context(ctx))
		return ncore.TraceError("", err)
	}

	// update customer model
	customer.FullName = payload.Nama
	customer.Profile.MaidenName = payload.NamaIbu
	customer.Profile.PlaceOfBirth = payload.TempatLahir
	customer.Profile.DateOfBirth = payload.TglLahir
	customer.IdentityType = nval.ParseInt64Fallback(payload.JenisIdentitas, 10)
	customer.IdentityNumber = payload.NoKtp
	customer.Profile.MarriageStatus = payload.StatusKawin
	customer.Profile.Gender = payload.JenisKelamin
	customer.Profile.Nationality = payload.Kewarganegaraan
	customer.Profile.IdentityExpiredAt = payload.TanggalExpiredIdentitas
	customer.Profile.Religion = payload.Agama
	customer.Profile.ProfileUpdatedAt = time.Now().Unix()

	// Get current address data
	address, errAddress := s.repo.FindAddressByCustomerId(customer.Id)
	if errAddress != nil && errAddress != sql.ErrNoRows {
		s.log.Error("error when get customer data", nlogger.Error(err), nlogger.Context(ctx))
		return ncore.TraceError("", err)
	}

	// Update address model
	address.Line = nsql.NewNullString(payload.Alamat)
	address.ProvinceId = sql.NullInt64{Int64: nval.ParseInt64Fallback(payload.IdProvinsi, 0), Valid: true}
	address.CityId = sql.NullInt64{Int64: nval.ParseInt64Fallback(payload.IdKabupaten, 0), Valid: true}
	address.DistrictId = sql.NullInt64{Int64: nval.ParseInt64Fallback(payload.IdKecamatan, 0), Valid: true}
	address.SubDistrictId = sql.NullInt64{Int64: nval.ParseInt64Fallback(payload.IdKelurahan, 0), Valid: true}
	address.ProvinceName = sql.NullString{String: payload.NamaProvinsi, Valid: true}
	address.CityName = sql.NullString{String: payload.NamaKabupaten, Valid: true}
	address.DistrictName = sql.NullString{String: payload.NamaKecamatan, Valid: true}
	address.SubDistrictName = sql.NullString{String: payload.NamaKelurahan, Valid: true}
	address.PostalCode = sql.NullString{String: payload.KodePos, Valid: true}

	// if empty create new address
	if errAddress == sql.ErrNoRows {
		address.CustomerId = customer.Id
		address.Xid = strings.ToUpper(xid.New().String())
		address.Metadata = []byte("{}")
		address.Purpose = constant.IdentityCard
		address.IsPrimary = sql.NullBool{Bool: true, Valid: true}
		address.ItemMetadata = model.NewItemMetadata(convert.ModifierDTOToModel(dto.Modifier{ID: "", Role: "", FullName: ""}))
	}

	// Update customer profile repo
	err = s.repo.UpdateCustomerProfile(customer, address)
	if err != nil {
		s.log.Error("error found when get customer repo", nlogger.Error(err), nlogger.Context(ctx))
		return ncore.TraceError("", err)
	}

	return nil
}

func (s *Service) isValidPassword(tokenString string, password string) (bool, error) {
	// Get Context
	ctx := s.ctx

	// Get UserRefID
	userRefId, err := s.validateTokenAndRetrieveUserRefID(tokenString)
	if err != nil {
		s.log.Error("error when validate token", nlogger.Error(err), nlogger.Context(ctx))
		return false, ncore.TraceError("", err)
	}

	// Find customer
	c, err := s.repo.FindCustomerByUserRefID(userRefId)
	if err != nil {
		return false, err
	}

	pw := stringToMD5(password)

	// Check is valid
	_, err = s.repo.IsValidPassword(c.Id, pw)
	if err != nil && err != sql.ErrNoRows {
		s.log.Error("error when check password match", nlogger.Error(err), nlogger.Context(ctx))
		return false, ncore.TraceError("error when validate password match", err)
	}

	if err == sql.ErrNoRows {
		return false, nil
	}

	return true, nil
}

func (s *Service) UpdatePassword(userRefId string, payload dto.UpdatePasswordRequest) error {
	// Get Context
	ctx := s.ctx

	// Find customer
	customer, err := s.repo.FindCustomerByUserRefID(userRefId)
	if err != nil {
		return err
	}

	// Set current password to md5
	currentPassword := stringToMD5(payload.CurrentPassword)

	// Validate current password
	_, err = s.repo.IsValidPassword(customer.Id, currentPassword)
	if err != nil && err != sql.ErrNoRows {
		s.log.Error("error found when check password match", nlogger.Error(err), nlogger.Context(ctx))
		return ncore.TraceError("error when validate password match", err)
	}

	// If password doesn't match
	if err == sql.ErrNoRows {
		err = nil
		return s.responses.GetError("E_USR_1")
	}

	// Set password to md5
	password := stringToMD5(payload.NewPassword)

	// Update password
	err = s.repo.UpdatePassword(customer.Id, password)
	if err != nil {
		s.log.Error("error found when check password match", nlogger.Error(err), nlogger.Context(ctx))
		return ncore.TraceError("failed to update password", err)
	}

	return nil
}

func (s *Service) UpdateAvatar(userRefId string, uploaded *dto.UploadResponse) error {
	// Get context
	ctx := s.ctx

	// Find customer
	customer, err := s.repo.FindCustomerByUserRefID(userRefId)
	if err != nil {
		s.log.Error("error when find current customer", nlogger.Error(err), nlogger.Context(ctx))
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
	err = s.repo.UpdateCustomerByUserRefID(customer, userRefId)
	if err != nil {
		s.log.Error("error when update photo profile", nlogger.Error(err), nlogger.Context(ctx))
		return ncore.TraceError("", err)
	}

	return nil
}

func (s *Service) UpdateIdentity(userRefId string, uploaded *dto.UploadResponse) error {
	// Get context
	ctx := s.ctx

	// Find customer
	customer, err := s.repo.FindCustomerByUserRefID(userRefId)
	if err != nil {
		s.log.Error("error when find current customer", nlogger.Error(err), nlogger.Context(ctx))
		return ncore.TraceError("", err)
	}

	// Remove old avatar if exist
	if oldFile := customer.Profile.IdentityPhotoFile; oldFile != "" {
		_ = s.AssetRemoveFile(oldFile, constant.AssetKTP)
	}

	// Update identity photo
	customer.Profile.IdentityPhotoFile = uploaded.FileName

	// Persist
	err = s.repo.UpdateCustomerByUserRefID(customer, userRefId)
	if err != nil {
		s.log.Error("error when update identity photo", nlogger.Error(err), nlogger.Context(ctx))
		return ncore.TraceError("", err)
	}

	return nil
}

func (s *Service) UpdateNPWP(userRefId string, npwpNumber string, uploaded *dto.UploadResponse) error {
	// Get context
	ctx := s.ctx

	// Find customer
	customer, err := s.repo.FindCustomerByUserRefID(userRefId)
	if err != nil {
		s.log.Error("error when find current customer", nlogger.Error(err), nlogger.Context(ctx))
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
	err = s.repo.UpdateCustomerByUserRefID(customer, userRefId)
	if err != nil {
		s.log.Error("error found when update NPWP", nlogger.Error(err), nlogger.Context(ctx))
		return ncore.TraceError("", err)
	}

	return nil
}

func (s *Service) UpdateSID(userRefId string, sidNumber string, uploaded *dto.UploadResponse) error {
	// Get context
	ctx := s.ctx

	// Find customer
	customer, err := s.repo.FindCustomerByUserRefID(userRefId)
	if err != nil {
		s.log.Error("error when find current customer", nlogger.Error(err), nlogger.Context(ctx))
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
	err = s.repo.UpdateCustomerByUserRefID(customer, userRefId)
	if err != nil {
		s.log.Error("error found when update NPWP", nlogger.Error(err), nlogger.Context(ctx))
		return ncore.TraceError("", err)
	}

	return nil
}

func (s *Service) CheckStatus(userRefId string) (*dto.CheckStatusResponse, error) {
	// Get context
	ctx := s.ctx

	// Find customer
	customer, err := s.repo.FindCustomerByUserRefID(userRefId)
	if err != nil {
		s.log.Error("error when find current customer", nlogger.Error(err), nlogger.Context(ctx))
		return nil, ncore.TraceError("", err)
	}

	// Find Verification
	verification, err := s.repo.FindVerificationByCustomerID(customer.Id)

	// Find Credential
	credential, err := s.repo.FindCredentialByCustomerID(customer.Id)

	var status = &dto.CheckStatusResponse{}

	status.Cif = customer.Cif
	status.EmailVerified = nval.ParseBooleanFallback(verification.EmailVerifiedStatus, false)
	status.KycVerified = false
	status.PinAvailable = nval.ParseBooleanFallback(credential.Pin, false)
	status.AktifasiTransFinansial = nval.ParseStringFallback(verification.FinancialTransactionStatus, "0")

	// If cif is empty
	if customer.Cif == "" {
		return status, nil
	}

	// Check CIF
	checkCifResponse, err := s.CheckCIF(customer.Cif)
	if err != nil {
		s.log.Error("error found when check CIF", nlogger.Error(err), nlogger.Context(ctx))
		return nil, ncore.TraceError("", err)
	}

	if checkCifResponse.ResponseCode != "00" {
		s.log.Error("error response code from check CIF", nlogger.Error(err), nlogger.Context(ctx))
		return status, nil
	}

	if checkCifResponse.Message == "" {
		s.log.Error("error response message from check CIF", nlogger.Error(err), nlogger.Context(ctx))
		return status, nil
	}

	customerInquiry := &dto.CustomerInquiryVO{}
	err = json.Unmarshal([]byte(checkCifResponse.Message), customerInquiry)
	if err != nil {
		s.log.Error("error marshall customer inquiry", nlogger.Error(err), nlogger.Context(ctx))
		return nil, constant.DefaultError.Trace()
	}

	// Update KYC

	// -- Update time if kyc was recently activated
	if customerInquiry.StatusKyc == "1" && verification.KycVerifiedStatus == 0 {
		verification.KycVerifiedAt = sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		}
	}

	verification.KycVerifiedStatus = nval.ParseInt64Fallback(customerInquiry.StatusKyc, 0)
	err = s.repo.UpdateVerification(verification)
	if err != nil {
		s.log.Error("error found when update verification repo", nlogger.Error(err), nlogger.Context(ctx))
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
		s.log.Error("error found when get gold savings", nlogger.Error(err), nlogger.Context(s.ctx))
		return nil, ncore.TraceError("error found when get gold savings", err)
	}

	return data, nil
}
