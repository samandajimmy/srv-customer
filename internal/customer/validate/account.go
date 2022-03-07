package validate

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/nbs-go/errx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
)

func PostLogin(p *dto.LoginPayload) error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.Email, validation.Required),
		validation.Field(&p.Password, validation.Required),
		validation.Field(&p.Agen, validation.Required),
		validation.Field(&p.Version, validation.Required),
	)

	if err != nil {
		return nhttp.BadRequestError.Trace(errx.Source(err))
	}

	return nil
}

func GetVerifyEmail(p *dto.VerificationPayload) error {
	return validation.ValidateStruct(p,
		validation.Field(&p.VerificationToken, validation.Required),
	)
}

func PostSendOTP(p *dto.SendOTPPayload) error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.Name, validation.Required, validation.Length(1, 50)),
		validation.Field(&p.Email, validation.Required, is.Email),
		validation.Field(&p.PhoneNumber, validation.Required, is.Digit),
	)

	if err != nil {
		return nhttp.BadRequestError.Trace(errx.Source(err))
	}

	return nil
}

func PostResendOTP(p *dto.RegisterResendOTPPayload) error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.PhoneNumber, validation.Required, is.Digit),
	)

	if err != nil {
		return nhttp.BadRequestError.Trace(errx.Source(err))
	}

	return nil
}

func PostVerifyOTP(p *dto.RegisterVerifyOTPPayload) error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.OTP, validation.Required, is.Digit),
		validation.Field(&p.PhoneNumber, validation.Required, is.Digit),
	)

	if err != nil {
		return nhttp.BadRequestError.Trace(errx.Source(err))
	}

	return nil
}

func PostRegister(p *dto.RegisterPayload) error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.Email, validation.Required),
		validation.Field(&p.PhoneNumber, validation.Required),
		validation.Field(&p.Password, validation.Required),
		validation.Field(&p.FcmToken, validation.Required),
		validation.Field(&p.RegistrationID, validation.Required),
	)

	if err != nil {
		return nhttp.BadRequestError.Trace(errx.Source(err))
	}

	return nil
}

func PostUpdatePasswordCheck(p *dto.UpdatePasswordCheckPayload) error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.CurrentPassword, validation.Required),
	)

	if err != nil {
		return nhttp.BadRequestError.Trace(errx.Source(err))
	}

	return nil
}

func PutUpdatePassword(p *dto.UpdatePasswordPayload) error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.CurrentPassword, validation.Required),
		validation.Field(&p.NewPassword, validation.Required),
	)

	if err != nil {
		return nhttp.BadRequestError.Trace(errx.Source(err))
	}

	return nil
}

func PostValidatePin(p *dto.ValidatePinPayload) error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.NewPin, validation.Required, validation.Length(6, 6)),
	)

	if err != nil {
		return nhttp.BadRequestError.Trace(errx.Source(err))
	}

	return nil
}

func PostCheckPin(p *dto.CheckPinPayload) error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.Pin, validation.Required, validation.Length(6, 6)),
		validation.Field(&p.UserRefID, validation.Required),
	)

	if err != nil {
		return nhttp.BadRequestError.Trace(errx.Source(err))
	}

	return nil
}

func PostUpdatePin(p *dto.UpdatePinPayload) error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.PIN, validation.Required, validation.Length(6, 6)),
		validation.Field(&p.NewPIN, validation.Required, validation.Length(6, 6)),
		validation.Field(&p.NewPINConfirmation, validation.Required, validation.Length(6, 6)),
	)
	if err != nil {
		return nhttp.BadRequestError.Trace(errx.Source(err))
	}

	// Validate pin confirmation
	err = PINConfirmation(&dto.PINConfirmation{
		NewPIN:             p.NewPIN,
		NewPINConfirmation: p.NewPINConfirmation,
	})
	if err != nil {
		return nhttp.BadRequestError.Trace(errx.Source(err))
	}

	return nil
}

func PINConfirmation(p *dto.PINConfirmation) error {
	if p.NewPIN != p.NewPINConfirmation {
		return fmt.Errorf("new_pin: Masukan PIN yang sama dengan sebelumnya")
	}

	return nil
}

func CheckPostOTPPinCreate(p *dto.CheckOTPPinPayload) error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.OTP, validation.Required),
		validation.Field(&p.UserRefID, validation.Required),
	)
	if err != nil {
		return nhttp.BadRequestError.Trace(errx.Source(err))
	}

	return nil
}

func PostCreatePin(p *dto.PostCreatePinPayload) error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.UserRefID, validation.Required),
		validation.Field(&p.NewPIN, validation.Required, validation.Length(6, 6)),
		validation.Field(&p.NewPINConfirmation, validation.Required, validation.Length(6, 6)),
		validation.Field(&p.OTP, validation.Required),
	)
	if err != nil {
		return nhttp.BadRequestError.Trace(errx.Source(err))
	}

	// Validate pin confirmation
	err = PINConfirmation(&dto.PINConfirmation{
		NewPIN:             p.NewPIN,
		NewPINConfirmation: p.NewPINConfirmation,
	})
	if err != nil {
		return nhttp.BadRequestError.Trace(errx.Source(err))
	}

	return nil
}
