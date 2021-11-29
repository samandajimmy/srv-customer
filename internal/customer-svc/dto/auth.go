package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/constant"
)

type Subject struct {
	SubjectID    string
	SubjectRefID int64
	SubjectRole  string
	SubjectType  constant.SubjectType
	ModifiedBy   Modifier
	Metadata     map[string]string
	SessionID    int64
}

type ClientCredential struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

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
