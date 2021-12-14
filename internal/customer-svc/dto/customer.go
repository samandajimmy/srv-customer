package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type ValidatePassword struct {
	IsValid bool   `json:"isValid"`
	ErrCode string `json:"errCode"`
}

type RegisterNewCustomer struct {
	Name           string `json:"nama"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"no_hp"`
	Password       string `json:"password"`
	FcmToken       string `json:"fcm_token"`
	RegistrationId string `json:"register_id"`
	Agen           string `json:"agen"`
	Version        string `json:"version"`
}

type RegisterNewCustomerResponse struct {
	User     CustomerVO  `json:"user"`
	JwtToken string      `json:"token"`
	Ekyc     *EKyc       `json:"ekyc"`
	GPoint   interface{} `json:"gpoint"`
	GCash    *GCash      `json:"gcash"`
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

type LoginResponse struct {
	Customer *CustomerVO `json:"user"`
	JwtToken string      `json:"token"`
}

type CustomerVO struct {
	ID                        string                  `json:"id"`
	Cif                       string                  `json:"cif"`
	IsKYC                     string                  `json:"isKYC"`
	Nama                      string                  `json:"nama"`
	NamaIbu                   string                  `json:"namaIbu"`
	NoKTP                     string                  `json:"noKTP"`
	Email                     string                  `json:"email"`
	JenisKelamin              string                  `json:"jenisKelamin"`
	TempatLahir               string                  `json:"tempatLahir"`
	TglLahir                  string                  `json:"tglLahir"`
	Alamat                    string                  `json:"alamat"`
	IDProvinsi                string                  `json:"idProvinsi"`
	IDKabupaten               string                  `json:"idKabupaten"`
	IDKecamatan               string                  `json:"idKecamatan"`
	IDKelurahan               string                  `json:"idKelurahan"`
	Kelurahan                 string                  `json:"kelurahan"`
	Provinsi                  string                  `json:"provinsi"`
	Kabupaten                 string                  `json:"kabupaten"`
	Kecamatan                 string                  `json:"kecamatan"`
	KodePos                   string                  `json:"kodePos"`
	NoHP                      string                  `json:"noHP"`
	Avatar                    string                  `json:"avatar"`
	FotoKTP                   string                  `json:"fotoKTP"`
	IsEmailVerified           string                  `json:"isEmailVerified"`
	Kewarganegaraan           string                  `json:"kewarganegaraan"`
	JenisIdentitas            string                  `json:"jenisIdentitas"`
	NoIdentitas               string                  `json:"noIdentitas"`
	TglExpiredIdentitas       string                  `json:"tglExpiredIdentitas"`
	NoNPWP                    string                  `json:"noNPWP"`
	FotoNPWP                  string                  `json:"fotoNPWP"`
	NoSid                     interface{}             `json:"noSid"`
	FotoSid                   interface{}             `json:"fotoSid"`
	StatusKawin               string                  `json:"statusKawin"`
	Norek                     string                  `json:"norek"`
	Saldo                     string                  `json:"saldo"`
	AktifasiTransFinansial    string                  `json:"aktifasiTransFinansial"`
	IsDukcapilVerified        string                  `json:"isDukcapilVerified"`
	IsOpenTe                  string                  `json:"isOpenTe"`
	ReferralCode              interface{}             `json:"referralCode"`
	GoldCardApplicationNumber string                  `json:"GoldCardApplicationNumber"`
	GoldCardAccountNumber     interface{}             `json:"goldCardAccountNumber"`
	KodeCabang                string                  `json:"kodeCabang"`
	TabunganEmas              *CustomerTabunganEmasVO `json:"tabunganEmas"`
	IsFirstLogin              bool                    `json:"isFirstLogin"`
	IsForceUpdatePassword     bool                    `json:"isForceUpdatePassword"`
}

type CustomerMetadata struct {
	Snapshot          string `json:"snapshot"`
	SnapshotSignature string `json:"snapshotSignature"`
}

type CustomerProfileVO struct {
	MaidenName         string `json:"maidenName"`
	Gender             string `json:"gender"`
	Nationality        string `json:"nationality"`
	DateOfBirth        string `json:"dateOfBirth"`
	PlaceOfBirth       string `json:"placeOfBirth"`
	IdentityPhotoFile  string `json:"identityPhotoFile"`
	MarriageStatus     string `json:"marriageStatus"`
	NPWPNumber         string `json:"npwpNumber"`
	NPWPPhotoFile      string `json:"npwpPhotoFile"`
	NPWPUpdatedAt      string `json:"npwpUpdatedAt"`
	ProfileUpdatedAt   string `json:"profileUpdatedAt"`
	CifLinkUpdatedAt   string `json:"cifLinkUpdatedAt"`
	CifUnlinkUpdatedAt string `json:"cifUnlinkUpdatedAt"`
	SidPhotoFile       string `json:"sidPhotoFile"`
}

type CustomerTabunganEmasVO struct {
	TotalSaldoBlokir  string      `json:"totalSaldoBlokir"`
	TotalSaldoSeluruh string      `json:"totalSaldoSeluruh"`
	TotalSaldoEfektif string      `json:"totalSaldoEfektif"`
	PrimaryRekening   interface{} `json:"primaryRekening"`
}

type LoginRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	Agen         string `json:"agen"`
	Version      string `json:"version"`
	DeviceId     string `json:"device_id"`
	IP           string `json:"ip"`
	Latitude     string `json:"latitude"`
	Longitude    string `json:"longitude"`
	Timezone     string `json:"timezone"`
	Brand        string `json:"brand"`
	OsVersion    string `json:"os_version"`
	Browser      string `json:"browser"`
	UseBiometric int64  `json:"use_biometric"`
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
		validation.Field(&d.Name, validation.Required, validation.Length(1, 50)),
		validation.Field(&d.Email, validation.Required, is.Email),
		validation.Field(&d.PhoneNumber, validation.Required, is.Digit),
	)
}

func (d LoginRequest) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Email, validation.Required),
		validation.Field(&d.Password, validation.Required),
		validation.Field(&d.Agen, validation.Required),
		validation.Field(&d.Version, validation.Required),
	)
}
