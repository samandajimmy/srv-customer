package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type VerificationPayload struct {
	VerificationToken string
}

func (d VerificationPayload) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.VerificationToken, validation.Required),
	)
}
