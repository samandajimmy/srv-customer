package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type RegisterNewCustomer struct {
	Name        string `json:"nama"`
	PhoneNumber string `json:"no_telepon"`
	Agen        string `json:"agen"`
	Version     int64  `json:"version"`
}

func (d RegisterNewCustomer) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Name, validation.Length(1, 50)),
		validation.Field(&d.PhoneNumber, validation.Required, is.Digit),
		validation.Field(&d.Agen, validation.Required, validation.Length(1, 10)),
	)
}

type NewRegisterResponse struct {
	Token  string `json:"token"`
	ReffId int64  `json:"reffId"`
}

type RegisterStepOne struct {
	Name        string `json:"nama"`
	Email       string `json:"email"`
	PhoneNumber string `json:"no_hp"`
}

type RegisterStepOneResponse struct {
	Action string `json:"action"`
}

type RegisterStepTwo struct {
	PhoneNumber string `json:"no_hp"`
	OTP         string `json:"otp"`
}

type RegisterStepTwoResponse struct {
	RegisterId string `json:"register_id"`
}

func (d RegisterStepTwo) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.OTP, validation.Required, is.Digit),
		validation.Field(&d.PhoneNumber, validation.Required, is.Digit),
	)
}

func (d RegisterStepOne) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Name, validation.Length(1, 50)),
		validation.Field(&d.Email, validation.Required, is.Email),
		validation.Field(&d.PhoneNumber, validation.Required, is.Digit),
	)
}
