package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type RegisterNewCustomer struct {
	Name           string `json:"nama"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"no_hp"`
	Password       string `json:"password"`
	FcmToken       string `json:"fcm_token"`
	RegistrationId string `json:"register_id"`
}

func (d RegisterNewCustomer) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Name, validation.Required),
		validation.Field(&d.Email, validation.Required),
		validation.Field(&d.PhoneNumber, validation.Required),
		validation.Field(&d.Password, validation.Required),
		validation.Field(&d.FcmToken, validation.Required),
		validation.Field(&d.RegistrationId, validation.Required),
	)
}

type RegisterNewCustomerResponse struct {
	User  CustomerVO `json:"user"`
	Token string     `json:"token"`
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

type RegisterResendOTP struct {
	PhoneNumber string `json:"no_hp"`
}

type RegisterResendOTPResponse struct {
	Action string `json:"action"`
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

func (d RegisterResendOTP) Validate() error {
	return validation.ValidateStruct(&d,
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

type CustomerVO struct {
	ID                        string      `json:"id"`
	Cif                       string      `json:"cif"`
	IsKYC                     string      `json:"isKYC"`
	Nama                      string      `json:"nama"`
	NamaIbu                   interface{} `json:"namaIbu"`
	NoKTP                     string      `json:"noKTP"`
	Email                     string      `json:"email"`
	JenisKelamin              string      `json:"jenisKelamin"`
	TempatLahir               interface{} `json:"tempatLahir"`
	TglLahir                  interface{} `json:"tglLahir"`
	Alamat                    interface{} `json:"alamat"`
	IDProvinsi                interface{} `json:"idProvinsi"`
	IDKabupaten               interface{} `json:"idKabupaten"`
	IDKecamatan               interface{} `json:"idKecamatan"`
	IDKelurahan               interface{} `json:"idKelurahan"`
	Kelurahan                 interface{} `json:"kelurahan"`
	Provinsi                  interface{} `json:"provinsi"`
	Kabupaten                 interface{} `json:"kabupaten"`
	Kecamatan                 interface{} `json:"kecamatan"`
	KodePos                   interface{} `json:"kodePos"`
	NoHP                      string      `json:"noHP"`
	Avatar                    interface{} `json:"avatar"`
	FotoKTP                   string      `json:"fotoKTP"`
	IsEmailVerified           string      `json:"isEmailVerified"`
	Kewarganegaraan           string      `json:"kewarganegaraan"`
	JenisIdentitas            string      `json:"jenisIdentitas"`
	NoIdentitas               string      `json:"noIdentitas"`
	TglExpiredIdentitas       string      `json:"tglExpiredIdentitas"`
	NoNPWP                    string      `json:"noNPWP"`
	FotoNPWP                  string      `json:"fotoNPWP"`
	NoSid                     interface{} `json:"noSid"`
	FotoSid                   interface{} `json:"fotoSid"`
	StatusKawin               string      `json:"statusKawin"`
	Norek                     string      `json:"norek"`
	Saldo                     string      `json:"saldo"`
	AktifasiTransFinansial    string      `json:"aktifasiTransFinansial"`
	IsDukcapilVerified        interface{} `json:"isDukcapilVerified"`
	IsOpenTe                  string      `json:"isOpenTe"`
	ReferralCode              interface{} `json:"referralCode"`
	GoldCardApplicationNumber interface{} `json:"GoldCardApplicationNumber"`
	GoldCardAccountNumber     interface{} `json:"goldCardAccountNumber"`
	KodeCabang                string      `json:"kodeCabang"`
	TabunganEmas              bool        `json:"tabunganEmas"`
	IsFirstLogin              bool        `json:"isFirstLogin"`
	IsForceUpdatePassword     bool        `json:"isForceUpdatePassword"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Agen     string `json:"agen"`
	Version  string `json:"version"`
}

func (d LoginRequest) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Email, validation.Required),
		validation.Field(&d.Password, validation.Required),
		validation.Field(&d.Agen, validation.Required),
		validation.Field(&d.Version, validation.Required),
	)
}
