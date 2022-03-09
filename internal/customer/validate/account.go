package validate

import (
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
