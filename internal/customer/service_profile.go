package customer

import (
	"database/sql"
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

	// TODO
	//goldSaving, err := s.getListAccountNumber(c.Cif, c.UserRefId)
	//if err != nil {
	//	return nil, ncore.TraceError("error when get list gold saving account", err)
	//}

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
			Avatar:          "", // TODO Avatar
			FotoKTP:         "", // TODO Foto KTP
			IsEmailVerified: nval.ParseStringFallback(v.EmailVerifiedStatus, ""),
			JenisIdentitas:  nval.ParseStringFallback(c.IdentityType, ""),
			TabunganEmas: &dto.GoldSavingVO{
				TotalSaldoBlokir:  "",
				TotalSaldoSeluruh: "",
				TotalSaldoEfektif: "",
				ListTabungan:      nil,
				PrimaryRekening:   nil,
			},
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
	customer.Profile.ProfileUpdatedAt = time.Now().String()

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

	s.log.Debugf("alamat %s", address.Line.String)

	// Update customer profile repo
	err = s.repo.UpdateCustomerProfile(customer, address)
	if err != nil {
		s.log.Error("error found when get customer repo", nlogger.Error(err), nlogger.Context(ctx))
		return ncore.TraceError("", err)
	}

	return nil
}
