package validate

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/nbs-go/errx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
)

func PostCreateFavorite(p *dto.CreateFavoritePayload) error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.AccountName, validation.Required, validation.Length(0, 32)),
		validation.Field(&p.AccountNumber, validation.Required, validation.Length(0, 16)),
		validation.Field(&p.Type, validation.Required, validation.Length(0, 32)),
		validation.Field(&p.TypeTransaction, validation.Required, validation.Length(0, 16)),
		validation.Field(&p.BankCode, validation.Length(0, 5)),
	)
	if err != nil {
		return nhttp.BadRequestError.Trace(errx.Source(err))
	}
	return nil
}
