package validate

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/nbs-go/errx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
)

func PostCreateBankAccount(p *dto.CreateBankAccountPayload) error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.AccountNumber, validation.Required),
		validation.Field(&p.AccountName, validation.Required),
		validation.Field(&p.Bank, validation.Required),
	)
	if err != nil {
		return nhttp.BadRequestError.Trace(errx.Source(err))
	}

	// Validate Bank
	if err = Bank(p.Bank); err != nil {
		return nhttp.BadRequestError.Trace(errx.Source(err))
	}

	return nil
}

func PutUpdateBankAccount(p *dto.UpdateBankAccountPayload) error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.AccountNumber, validation.Required),
		validation.Field(&p.AccountName, validation.Required),
		validation.Field(&p.Bank, validation.Required),
		validation.Field(&p.Status, validation.Required),
		validation.Field(&p.Version, validation.Required, validation.Min(1)),
	)
	if err != nil {
		return nhttp.BadRequestError.Trace(errx.Source(err))
	}

	// Validate Bank
	if err = Bank(p.Bank); err != nil {
		return nhttp.BadRequestError.Trace(errx.Source(err))
	}

	return nil
}

func Bank(p *dto.Bank) error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.ID, validation.Required),
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.Title, validation.Required),
		validation.Field(&p.Code, validation.Required),
		validation.Field(&p.Thumbnail, validation.Required),
	)
	if err != nil {
		return nhttp.BadRequestError.Trace(errx.Source(err))
	}

	return nil
}
