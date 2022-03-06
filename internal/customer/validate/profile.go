package validate

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/nbs-go/errx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
)

func PutUpdateProfile(p *dto.UpdateProfilePayload) error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.Nama, validation.Required),
		validation.Field(&p.Alamat, validation.Required),
		validation.Field(&p.NamaIbu, validation.Required),
		validation.Field(&p.TempatLahir, validation.Required),
		validation.Field(&p.TglLahir, validation.Required),
		validation.Field(&p.NoKtp, validation.Required),
		validation.Field(&p.IDKelurahan, validation.Required),
		validation.Field(&p.JenisKelamin, validation.Required),
		validation.Field(&p.JenisIdentitas, validation.Required),
		validation.Field(&p.Kewarganegaraan, validation.Required),
		validation.Field(&p.Agama, validation.Required),
	)
	if err != nil {
		return nhttp.BadRequestError.Trace(errx.Source(err))
	}
	return nil
}

func PostUpdateSID(p *dto.UpdateSIDPayload) error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.NoSID, validation.Required, validation.Length(15, 15)),
	)
	if err != nil {
		return nhttp.BadRequestError.Trace(errx.Source(err))
	}
	return nil
}

func PostUpdateNPWP(p *dto.UpdateNPWPPayload) error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.NoNPWP, validation.Required, validation.Length(15, 15)),
	)
	if err != nil {
		return nhttp.BadRequestError.Trace(errx.Source(err))
	}
	return nil
}
