package customer

import (
	"database/sql"
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
)

func (s *Service) CustomerProfile(id int64) (*dto.ProfileResponse, error) {
	// Get Context
	ctx := s.ctx

	// Get customer data
	c, err := s.repo.FindCustomerByRefID(id)
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

	// TODO Tabungan Emas

	// Compose response
	resp := dto.ProfileResponse{
		CustomerVO: dto.CustomerVO{
			ID:              nval.ParseStringFallback(c.UserRefId, ""),
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
			IDProvinsi:      a.ProvinceId.String,
			IDKabupaten:     a.CityId.String,
			IDKecamatan:     a.DistrictId.String,
			IDKelurahan:     a.SubDistrictId.String,
			Provinsi:        a.ProvinceName.String,
			Kabupaten:       a.CityName.String,
			Kecamatan:       a.DistrictName.String,
			Kelurahan:       a.SubDistrictName.String,
			KodePos:         a.PostalCode.String,
			Avatar:          "", // TODO Avatar
			FotoKTP:         "", // TODO Foto KTP
			IsEmailVerified: nval.ParseStringFallback(v.EmailVerifiedStatus, ""),
			JenisIdentitas:  nval.ParseStringFallback(c.IdentityType, ""),
			TabunganEmas: &dto.CustomerTabunganEmasVO{
				TotalSaldoBlokir:  "",
				TotalSaldoSeluruh: "",
				TotalSaldoEfektif: "",
				PrimaryRekening:   "",
			},
		},
	}

	return &resp, nil
}
