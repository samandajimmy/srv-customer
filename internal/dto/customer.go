package dto

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
)

type ValidatePasswordResult struct {
	IsValid bool   `json:"isValid"`
	ErrCode string `json:"errCode"`
	Message string `json:"message"`
}

type RegisterPayload struct {
	Name           string `json:"nama"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"no_hp"`
	Password       string `json:"password"`
	FcmToken       string `json:"fcm_token"`
	RegistrationID string `json:"register_id"`
	Agen           string `json:"agen"`
	Version        string `json:"version"`
}

type RegisterResult struct {
	*LoginResult
	Ekyc   *EKyc       `json:"ekyc"`
	GPoint interface{} `json:"gpoint"`
	GCash  *GCash      `json:"gcash"`
}

type NewRegisterResponse struct {
	Token  string `json:"token"`
	ReffID int64  `json:"reffId"`
}

type SendOTPPayload struct {
	Name        string `json:"nama"`
	Email       string `json:"email"`
	PhoneNumber string `json:"no_hp"`
	Agen        string `json:"agen"`
}

type RegisterResendOTPPayload struct {
	PhoneNumber string `json:"no_hp"`
}

type RegisterResendOTPResult struct {
	Action string `json:"action"`
}

type SendOTPResult struct {
	Action string `json:"action"`
}

type RegisterVerifyOTPPayload struct {
	PhoneNumber string `json:"no_hp"`
	OTP         string `json:"otp"`
}

type RegisterVerifyOTPResult struct {
	RegisterID string `json:"register_id"`
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
	Financial             interface{}
	GoldSaving            interface{}
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

type LoginResult struct {
	User     *LoginUserVO `json:"user"`
	JwtToken string       `json:"token"`
}

type ProfileResponse struct {
	CustomerVO
}

type CustomerVO struct {
	ID                        string      `json:"id,omitempty"`
	Cif                       string      `json:"cif"`
	IsKYC                     string      `json:"isKyc"`
	Nama                      string      `json:"nama"`
	NamaIbu                   string      `json:"namaIbu"`
	NoKTP                     string      `json:"noKtp"`
	Email                     string      `json:"email"`
	JenisKelamin              string      `json:"jenisKelamin"`
	TempatLahir               string      `json:"tempatLahir"`
	TglLahir                  string      `json:"tglLahir"`
	Alamat                    string      `json:"alamat"`
	IDProvinsi                interface{} `json:"idProvinsi"`
	IDKabupaten               interface{} `json:"idKabupaten"`
	IDKecamatan               interface{} `json:"idKecamatan"`
	IDKelurahan               interface{} `json:"idKelurahan"`
	Kelurahan                 string      `json:"kelurahan"`
	Provinsi                  string      `json:"provinsi"`
	Kabupaten                 string      `json:"kabupaten"`
	Kecamatan                 string      `json:"kecamatan"`
	KodePos                   string      `json:"kodePos"`
	NoHP                      string      `json:"noHp"`
	Avatar                    string      `json:"avatar"`
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
	IsDukcapilVerified        string      `json:"isDukcapilVerified"`
	IsOpenTe                  string      `json:"isOpenTe"`
	ReferralCode              interface{} `json:"referralCode"`
	GoldCardApplicationNumber string      `json:"goldCardApplicationNumber"`
	GoldCardAccountNumber     interface{} `json:"goldCardAccountNumber"`
	KodeCabang                string      `json:"kodeCabang"`
	TabunganEmas              interface{} `json:"tabunganEmas"`
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
	IDKelurahan               string `json:"id_kelurahan,omitempty"`
	NoNpwp                    string `json:"no_npwp,omitempty"`
	FotoNpwp                  string `json:"foto_npwp,omitempty"`
	NoSid                     string `json:"no_sid"`
	FotoSid                   string `json:"foto_sid"`
	KodeCabang                string `json:"kode_cabang,omitempty"`
	FotoURL                   string `json:"foto_url"`
	FotoKtpURL                string `json:"foto_ktp_url,omitempty"`
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
	DeviceIDBiometric         string `json:"device_id_biometric"`
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
	NPWPUpdatedAt      int64  `json:"npwpUpdatedAt,string,omitempty"`
	ProfileUpdatedAt   int64  `json:"profileUpdatedAt,string,omitempty"`
	CifLinkUpdatedAt   int64  `json:"cifLinkUpdatedAt,string,omitempty"`
	CifUnlinkUpdatedAt int64  `json:"cifUnlinkUpdatedAt,string,omitempty"`
	SidPhotoFile       string `json:"sidPhotoFile"`
	Religion           string `json:"religion"`
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
	ListTabungan      []AccountSavingVO `json:"listTabungan,omitempty"`
	PrimaryRekening   *AccountSavingVO  `json:"primaryRekening,omitempty"`
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

type LoginPayload struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	Agen         string `json:"agen"`
	Version      string `json:"version"`
	DeviceID     string `json:"device_id"`
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

type UpdateProfilePayload struct {
	Nama                    string `json:"nama"`
	Alamat                  string `json:"alamat"`
	NamaIbu                 string `json:"nama_ibu"`
	Agama                   string `json:"agama"`
	TempatLahir             string `json:"tempat_lahir"`
	TglLahir                string `json:"tgl_lahir"`
	JenisIdentitas          string `json:"jenis_identitas"`
	NoKtp                   string `json:"no_ktp"`
	IDProvinsi              string `json:"id_provinsi"`
	NamaProvinsi            string `json:"nama_provinsi"`
	IDKabupaten             string `json:"id_kabupaten"`
	NamaKabupaten           string `json:"nama_kabupaten"`
	IDKecamatan             string `json:"id_kecamatan"`
	NamaKecamatan           string `json:"nama_kecamatan"`
	IDKelurahan             string `json:"id_kelurahan"`
	NamaKelurahan           string `json:"nama_kelurahan"`
	KodePos                 string `json:"kode_pos"`
	JenisKelamin            string `json:"jenis_kelamin"`
	Kewarganegaraan         string `json:"kewarganegaraan"`
	TanggalExpiredIdentitas string `json:"tanggal_expired_identitas"`
	StatusKawin             string `json:"status_kawin"`
}

type UpdatePasswordCheckPayload struct {
	CurrentPassword string `json:"current_password"`
}

type UpdatePasswordPayload struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type UpdateNPWPPayload struct {
	NoNPWP    string `json:"no_npwp"`
	UserRefID string
	FileName  string
}

type UpdateSIDPayload struct {
	NoSID     string `json:"no_sid"`
	UserRefID string
	FileName  string
}

type UpdateAvatarPayload struct {
	UpdateUserFile
	FileSize int64
	MimeType string
}

type UpdateUserFile struct {
	FileName  string
	UserRefID string
	AssetType constant.AssetType
}

type UploadUserFilePayload struct {
	File      nhttp.MultipartFile
	AssetType constant.AssetType
}

type CheckStatusResponse struct {
	Cif                    string `json:"cif"`
	EmailVerified          bool   `json:"emailVerified"`
	KycVerified            bool   `json:"kycVerified"`
	PinAvailable           bool   `json:"pinAvailable"`
	AktifasiTransFinansial string `json:"aktifasiTransFinansial"`
}

type CustomerInquiryVO struct {
	Cif                string `json:"cif"`
	Hp                 string `json:"hp"`
	IbuKandung         string `json:"ibuKandung"`
	JenisKelamin       string `json:"jenisKelamin"`
	KodeCabang         string `json:"kodeCabang"`
	NamaNasabah        string `json:"namaNasabah"`
	NoIdentitas        string `json:"noIdentitas"`
	Telp               string `json:"telp"`
	TempatLahir        string `json:"tempatLahir"`
	TglLahir           string `json:"tglLahir"`
	TglKyc             string `json:"tglKyc"`
	StatusKyc          string `json:"statusKyc"`
	StatusKawin        string `json:"statusKawin"`
	Kewarganegaraan    string `json:"kewarganegaraan"`
	TipeIdentitas      string `json:"tipeIdentitas"`
	IDKelurahan        string `json:"idKelurahan"`
	Jalan              string `json:"jalan"`
	IsDukcapilVerified string `json:"isDukcapilVerified"`
	Agama              string `json:"agama"`
}

type PINConfirmation struct {
	NewPIN             string `json:"new_pin"`
	NewPINConfirmation string `json:"new_pin_confirmation"`
}

type OTPResetPasswordPayload struct {
	Email string `json:"email"`
}

type VerifyOTPResetPasswordPayload struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

type ResetPasswordByOTPPayload struct {
	Email    string `json:"email"`
	OTP      string `json:"otp"`
	Password string `json:"password"`
}

type ChangePhoneNumberPayload struct {
	UserRefID          string `json:"userRefID"`
	MaidenName         string `json:"maidenName"`
	FullName           string `json:"fullName"`
	DateOfBirth        string `json:"dateOfBirth"`
	CurrentPhoneNumber string `json:"currentPhoneNumber"`
	NewPhoneNumber     string `json:"newPhoneNumber"`
}

type ChangePhoneNumberResult struct {
	PhoneNumber string `json:"currentPhoneNumber"`
}

type ChangePhoneNumberRequestCore struct {
	CurrentPhoneNumber string
	NewPhoneNumber     string
	FullName           string
	MaidenName         string
	Cif                string
	DateOfBirth        string
}

type UpdateLinkCifPayload struct {
	Cif         string `json:"cif"`
	PhoneNumber string `json:"no_tlp"`
}

type UnlinkCifPayload struct {
	Cif string `json:"cif"`
}

type PutSynchronizeCustomerPayload struct {
	Customer     *CustomerSyncVO     `json:"customer"`
	Financial    *FinancialSyncVO    `json:"financial"`
	Credential   *CredentialSyncVO   `json:"credential"`
	Verification *VerificationSyncVO `json:"verification"`
	Address      *AddressSyncVO      `json:"address"`
}

type PutSynchronizeCustomerResult struct {
	Customer     *CustomerSyncVO     `json:"customer"`
	Financial    *FinancialSyncVO    `json:"financial"`
	Credential   *CredentialSyncVO   `json:"credential"`
	Verification *VerificationSyncVO `json:"verification"`
	Address      *AddressSyncVO      `json:"adddress"`
}

type CustomerSyncVO struct {
	Photos         PhotosVO               `json:"photos"`
	Profile        CustomerProfileVO      `json:"profile"`
	FullName       string                 `json:"fullName"`
	Phone          string                 `json:"phone"`
	Email          string                 `json:"email"`
	IdentityType   int64                  `json:"identityType"`
	IdentityNumber string                 `json:"identityNumber"`
	Cif            string                 `json:"cif"`
	Sid            string                 `json:"sid"`
	ReferralCode   string                 `json:"referralCode"`
	Status         constant.ControlStatus `json:"status"`
}

type PhotosVO struct {
	FileName string `json:"fileName"`
	FileSize int64  `json:"fileSize"`
	MimeType string `json:"mimeType"`
}

type FinancialSyncVO struct {
	MainAccountNumber         string `json:"mainAccountNumber"`
	AccountNumber             string `json:"accountNumber"`
	GoldSavingStatus          int64  `json:"goldSavingStatus"`
	GoldCardApplicationNumber string `json:"goldCardApplicationNumber"`
	GoldCardAccountNumber     string `json:"goldCardAccountNumber"`
	Balance                   int64  `json:"balance"`
}

type CredentialSyncVO struct {
	Password            string                 `json:"password"`
	NextPasswordResetAt int64                  `json:"nextPasswordResetAt"`
	Pin                 string                 `json:"pin"`
	PinUpdatedAt        int64                  `json:"pinUpdatedAt"`
	PinLastAccessAt     int64                  `json:"pinLastAccessAt"`
	PinCounter          int64                  `json:"pinCounter"`
	PinBlockedStatus    int64                  `json:"pinBlockedStatus"`
	IsLocked            int64                  `json:"isLocked"`
	LoginFailCount      int64                  `json:"loginFailCount"`
	WrongPasswordCount  int64                  `json:"wrongPasswordCount"`
	BlockedAt           int64                  `json:"blockedAt"`
	BlockedUntilAt      int64                  `json:"blockedUntilAt"`
	BiometricLogin      constant.ControlStatus `json:"biometricLogin"`
	BiometricDeviceID   string                 `json:"biometricDeviceId"`
	Metadata            MetadataCredentialVO   `json:"metadata"`
}

type MetadataCredentialVO struct {
	TryLoginAt   string `json:"tryLoginAt"`
	PinCreatedAt string `json:"pinCreatedAt"`
	PinBlockedAt string `json:"pinBlockedAt"`
}

type VerificationSyncVO struct {
	KycVerifiedStatus               int64                  `json:"kycVerifiedStatus"`
	EmailVerificationToken          string                 `json:"emailVerificationToken"`
	EmailVerifiedStatus             int64                  `json:"emailVerifiedStatus"`
	DukcapilVerifiedStatus          int64                  `json:"dukcapilVerifiedStatus"`
	FinancialTransactionStatus      constant.ControlStatus `json:"financialTransactionStatus"`
	FinancialTransactionActivatedAt int64                  `json:"financialTransactionActivatedAt"`
}

type AddressSyncVO struct {
	Purpose         int64  `json:"purpose"`
	ProvinceID      int64  `json:"provinceId"`
	ProvinceName    string `json:"provinceName"`
	CityID          int64  `json:"cityId"`
	CityName        string `json:"cityName"`
	DistrictID      int64  `json:"districtId"`
	DistrictName    string `json:"districtName"`
	SubDistrictID   int64  `json:"subDistrictId"`
	SubDistrictName string `json:"subDistrictName"`
	Line            string `json:"line"`
	PostalCode      string `json:"postalCode"`
	IsPrimary       bool   `json:"isPrimary"`
}
