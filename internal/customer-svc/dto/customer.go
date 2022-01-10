package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/model"
)

type ValidatePassword struct {
	IsValid bool   `json:"isValid"`
	ErrCode string `json:"errCode"`
	Message string `json:"message"`
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
	*LoginResponse
	Ekyc   *EKyc       `json:"ekyc"`
	GPoint interface{} `json:"gpoint"`
	GCash  *GCash      `json:"gcash"`
}

type NewRegisterResponse struct {
	Token  string `json:"token"`
	ReffId int64  `json:"reffId"`
}

type RegisterStepOne struct {
	Name        string `json:"nama"`
	Email       string `json:"email"`
	PhoneNumber string `json:"no_hp"`
	Agen        string `json:"agen"`
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

type CustomerSynchronizeRequest struct {
	Name        string `json:"nama"`
	Email       string `json:"email"`
	PhoneNumber string `json:"no_hp"`
	Password    string `json:"password"`
	FcmToken    string `json:"fcm_token"`
}

type CustomerSynchronizeResponse struct {
	Customer *UserVO `json:"user"`
}

type LoginVO struct {
	Customer              *model.Customer
	Address               *model.Address
	Profile               CustomerProfileVO
	Verification          *model.Verification
	IsFirstLogin          bool
	IsForceUpdatePassword bool
	Token                 string
}

type EmailPayload struct {
	Subject    string           `json:"subject"`
	From       FromEmailPayload `json:"from"`
	To         string           `json:"to"`
	Message    string           `json:"message"`
	Attachment string           `json:"attachment"`
	MimeType   string           `json:"mimeType"`
}

type FromEmailPayload struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CustomerVO struct {
	ID                        string                  `json:"id"`
	Cif                       string                  `json:"cif"`
	IsKYC                     string                  `json:"is_kyc"`
	Nama                      string                  `json:"nama"`
	NamaIbu                   string                  `json:"nama_ibu"`
	NoKTP                     string                  `json:"no_ktp"`
	Email                     string                  `json:"email"`
	JenisKelamin              string                  `json:"jenis_kelamin"`
	TempatLahir               string                  `json:"tempat_lahir"`
	TglLahir                  string                  `json:"tgl_lahir"`
	Alamat                    string                  `json:"alamat"`
	IDProvinsi                string                  `json:"id_provinsi"`
	IDKabupaten               string                  `json:"id_kabupaten"`
	IDKecamatan               string                  `json:"id_kecamatan"`
	IDKelurahan               string                  `json:"id_kelurahan"`
	Kelurahan                 string                  `json:"kelurahan"`
	Provinsi                  string                  `json:"provinsi"`
	Kabupaten                 string                  `json:"kabupaten"`
	Kecamatan                 string                  `json:"kecamatan"`
	KodePos                   string                  `json:"kode_pos"`
	NoHP                      string                  `json:"no_hp"`
	Avatar                    string                  `json:"avatar"`
	FotoKTP                   string                  `json:"foto_ktp"`
	IsEmailVerified           string                  `json:"is_email_verified"`
	Kewarganegaraan           string                  `json:"kewarganegaraan"`
	JenisIdentitas            string                  `json:"jenis_identitas"`
	NoIdentitas               string                  `json:"no_identitas"`
	TglExpiredIdentitas       string                  `json:"tgl_expired_identitas"`
	NoNPWP                    string                  `json:"no_npwp"`
	FotoNPWP                  string                  `json:"foto_npwp"`
	NoSid                     interface{}             `json:"no_sid"`
	FotoSid                   interface{}             `json:"foto_sid"`
	StatusKawin               string                  `json:"status_kawin"`
	Norek                     string                  `json:"norek"`
	Saldo                     string                  `json:"saldo"`
	AktifasiTransFinansial    string                  `json:"AktifasiTransFinansial"`
	IsDukcapilVerified        string                  `json:"is_dukcapil_verified"`
	IsOpenTe                  string                  `json:"is_open_te"`
	ReferralCode              interface{}             `json:"referral_code"`
	GoldCardApplicationNumber string                  `json:"gold_card_application_number"`
	GoldCardAccountNumber     interface{}             `json:"gold_card_account_number"`
	KodeCabang                string                  `json:"kode_cabang"`
	TabunganEmas              *CustomerTabunganEmasVO `json:"tabungan_emas"`
	IsFirstLogin              bool                    `json:"is_first_login"`
	IsForceUpdatePassword     bool                    `json:"is_force_update_password"`
}

type UserVO struct {
	UserAiid                  string `json:"user_AIID,omitempty"`
	JenisIdentitas            string `json:"jenis_identitas,omitempty"`
	NoKtp                     string `json:"no_ktp"`
	TanggalExpiredIdentitas   string `json:"tanggal_expired_identitas,omitempty"`
	Cif                       string `json:"cif,omitempty"`
	Email                     string `json:"email"`
	NoHp                      string `json:"no_hp"`
	Password                  string `json:"password"`
	NextPasswordReset         string `json:"next_password_reset"`
	Pin                       string `json:"pin"`
	Nama                      string `json:"nama"`
	NamaIbu                   string `json:"nama_ibu"`
	JenisKelamin              string `json:"jenis_kelamin,omitempty"`
	TempatLahir               string `json:"tempat_lahir"`
	Agama                     string `json:"agama"`
	TglLahir                  string `json:"tgl_lahir"`
	Alamat                    string `json:"alamat"`
	Domisili                  string `json:"domisili"`
	Kewarganegaraan           string `json:"kewarganegaraan,omitempty"`
	StatusKawin               string `json:"status_kawin,omitempty"`
	Kodepos                   string `json:"kodepos"`
	IdKelurahan               string `json:"id_kelurahan,omitempty"`
	NoNpwp                    string `json:"no_npwp,omitempty"`
	FotoNpwp                  string `json:"foto_npwp,omitempty"`
	NoSid                     string `json:"no_sid"`
	FotoSid                   string `json:"foto_sid"`
	KodeCabang                string `json:"kode_cabang,omitempty"`
	FotoUrl                   string `json:"foto_url"`
	FotoKtpUrl                string `json:"foto_ktp_url,omitempty"`
	Status                    string `json:"status"`
	IsLocked                  string `json:"is_locked"`
	LoginFailCount            string `json:"login_fail_count,omitempty"`
	EmailVerified             string `json:"email_verified,omitempty"`
	KycVerified               string `json:"kyc_verified,omitempty"`
	EmailVerificationToken    string `json:"email_verification_token,omitempty"`
	Token                     string `json:"token"`
	TokenWeb                  string `json:"token_web,omitempty"`
	FcmToken                  string `json:"fcm_token,omitempty"`
	LastUpdate                string `json:"last_update"`
	Norek                     string `json:"norek,omitempty"`
	Saldo                     string `json:"saldo,omitempty"`
	PinTemp                   string `json:"pin_temp"`
	LastUpdateDataNasabah     string `json:"last_update_data_nasabah,omitempty"`
	LastUpdateDataNpwp        string `json:"last_update_data_npwp,omitempty"`
	LastUpdateLinkCif         string `json:"last_update_link_cif,omitempty"`
	LastUpdateUnlinkCif       string `json:"last_update_unlink_cif,omitempty"`
	LastUpdatePin             string `json:"last_update_pin"`
	AktifasiTransFinansial    string `json:"aktifasiTransFinansial,omitempty"`
	TanggalAktifasiFinansial  string `json:"tanggal_aktifasi_finansial"`
	IsDukcapilVerified        string `json:"is_dukcapil_verified"`
	IsOpenTe                  string `json:"is_open_te"`
	ReferralCode              string `json:"referral_code"`
	GoldcardApplicationNumber string `json:"goldcard_application_number"`
	GoldcardAccountNumber     string `json:"goldcard_account_number"`
	TryLoginDate              string `json:"try_login_date"`
	WrongPasswordCount        string `json:"wrong_password_count,omitempty"`
	BlockedDate               string `json:"blocked_date"`
	BlockedToDate             string `json:"blocked_to_date"`
	NorekUtama                string `json:"norek_utama"`
	IsSetBiometric            string `json:"is_set_biometric"`
	DeviceIdBiometric         string `json:"device_id_biometric"`
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

type CustomerPhotoVO struct {
	Xid      string `json:"xid"`
	Filename string `json:"filename"`
	Filesize int    `json:"filesize"`
	Mimetype string `json:"mimetype"`
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
	FcmToken     string `json:"fcm_token"`
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
