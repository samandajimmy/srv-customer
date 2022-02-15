package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
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
	Customer              interface{}
	Address               interface{}
	Profile               interface{}
	Verification          interface{}
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

type LoginUserVO struct {
	CustomerVO
	IsFirstLogin          bool `json:"isFirstLogin"`
	IsForceUpdatePassword bool `json:"isForceUpdatePassword"`
}

type LoginResponse struct {
	User     *LoginUserVO `json:"user"`
	JwtToken string       `json:"token"`
}

type ProfileResponse struct {
	CustomerVO
}

type CustomerVO struct {
	ID                        string        `json:"id"`
	Cif                       string        `json:"cif"`
	IsKYC                     string        `json:"isKyc"`
	Nama                      string        `json:"nama"`
	NamaIbu                   string        `json:"namaIbu"`
	NoKTP                     string        `json:"noKtp"`
	Email                     string        `json:"email"`
	JenisKelamin              string        `json:"jenisKelamin"`
	TempatLahir               string        `json:"tempatLahir"`
	TglLahir                  string        `json:"tglLahir"`
	Alamat                    string        `json:"alamat"`
	IDProvinsi                interface{}   `json:"idProvinsi"`
	IDKabupaten               interface{}   `json:"idKabupaten"`
	IDKecamatan               interface{}   `json:"idKecamatan"`
	IDKelurahan               interface{}   `json:"idKelurahan"`
	Kelurahan                 string        `json:"kelurahan"`
	Provinsi                  string        `json:"provinsi"`
	Kabupaten                 string        `json:"kabupaten"`
	Kecamatan                 string        `json:"kecamatan"`
	KodePos                   string        `json:"kodePos"`
	NoHP                      string        `json:"noHp"`
	Avatar                    string        `json:"avatar"`
	FotoKTP                   string        `json:"fotoKTP"`
	IsEmailVerified           string        `json:"isEmailVerified"`
	Kewarganegaraan           string        `json:"kewarganegaraan"`
	JenisIdentitas            string        `json:"jenisIdentitas"`
	NoIdentitas               string        `json:"noIdentitas"`
	TglExpiredIdentitas       string        `json:"tglExpiredIdentitas"`
	NoNPWP                    string        `json:"noNPWP"`
	FotoNPWP                  string        `json:"fotoNPWP"`
	NoSid                     interface{}   `json:"noSid"`
	FotoSid                   interface{}   `json:"fotoSid"`
	StatusKawin               string        `json:"statusKawin"`
	Norek                     string        `json:"norek"`
	Saldo                     string        `json:"saldo"`
	AktifasiTransFinansial    string        `json:"aktifasiTransFinansial"`
	IsDukcapilVerified        string        `json:"isDukcapilVerified"`
	IsOpenTe                  string        `json:"isOpenTe"`
	ReferralCode              interface{}   `json:"referralCode"`
	GoldCardApplicationNumber string        `json:"goldCardApplicationNumber"`
	GoldCardAccountNumber     interface{}   `json:"goldCardAccountNumber"`
	KodeCabang                string        `json:"kodeCabang"`
	TabunganEmas              *GoldSavingVO `json:"tabunganEmas"`
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

type GoldSavingVO struct {
	TotalSaldoBlokir  string            `json:"totalSaldoBlokir"`
	TotalSaldoSeluruh string            `json:"totalSaldoSeluruh"`
	TotalSaldoEfektif string            `json:"totalSaldoEfektif"`
	ListTabungan      []AccountSavingVO `json:"listTabungan"`
	PrimaryRekening   *AccountSavingVO  `json:"primaryRekening"`
}

type AccountSavingVO struct {
	Cif          string `json:"cif"`
	KodeCabang   string `json:"kodeCabang"`
	NamaNasabah  string `json:"namaNasabah"`
	NoBuku       string `json:"noBuku"`
	Norek        string `json:"norek"`
	SaldoBlokir  string `json:"saldoBlokir"`
	SaldoEmas    string `json:"saldoEmas"`
	TglBuka      string `json:"tglBuka"`
	SaldoEfektif string `json:"saldoEfektif"`
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

type UpdateProfileRequest struct {
	Nama                    string `json:"nama"`
	Alamat                  string `json:"alamat"`
	NamaIbu                 string `json:"nama_ibu"`
	Agama                   string `json:"agama"`
	TempatLahir             string `json:"tempat_lahir"`
	TglLahir                string `json:"tgl_lahir"`
	JenisIdentitas          string `json:"jenis_identitas"`
	NoKtp                   string `json:"no_ktp"`
	IdProvinsi              string `json:"id_provinsi"`
	NamaProvinsi            string `json:"nama_provinsi"`
	IdKabupaten             string `json:"id_kabupaten"`
	NamaKabupaten           string `json:"nama_kabupaten"`
	IdKecamatan             string `json:"id_kecamatan"`
	NamaKecamatan           string `json:"nama_kecamatan"`
	IdKelurahan             string `json:"id_kelurahan"`
	NamaKelurahan           string `json:"nama_kelurahan"`
	KodePos                 string `json:"kode_pos"`
	JenisKelamin            string `json:"jenis_kelamin"`
	Kewarganegaraan         string `json:"kewarganegaraan"`
	TanggalExpiredIdentitas string `json:"tanggal_expired_identitas"`
	StatusKawin             string `json:"status_kawin"`
}

func (d UpdateProfileRequest) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Nama, validation.Required),
		validation.Field(&d.Alamat, validation.Required),
		validation.Field(&d.NamaIbu, validation.Required),
		validation.Field(&d.TempatLahir, validation.Required),
		validation.Field(&d.TglLahir, validation.Required),
		validation.Field(&d.NoKtp, validation.Required),
		validation.Field(&d.IdKelurahan, validation.Required),
		validation.Field(&d.JenisKelamin, validation.Required),
		validation.Field(&d.JenisIdentitas, validation.Required),
		validation.Field(&d.Kewarganegaraan, validation.Required),
		validation.Field(&d.Agama, validation.Required),
	)
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
