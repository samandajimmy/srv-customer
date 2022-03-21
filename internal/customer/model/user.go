package model

import (
	"database/sql"
	"encoding/json"
	"github.com/rs/xid"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
	"strings"
)

// User model old structure from PDS-API
type User struct {
	UserAiid                  int64                  `db:"user_AIID"`
	JenisIdentitas            string                 `db:"jenis_identitas"`
	NoKtp                     sql.NullString         `db:"no_ktp"`
	TanggalExpiredIdentitas   sql.NullString         `db:"tanggal_expired_identitas"`
	Cif                       string                 `db:"cif"`
	Email                     sql.NullString         `db:"email"`
	NoHp                      sql.NullString         `db:"no_hp"`
	Password                  sql.NullString         `db:"password"`
	NextPasswordReset         sql.NullTime           `db:"next_password_reset"`
	Pin                       sql.NullString         `db:"pin"`
	Nama                      sql.NullString         `db:"nama"`
	NamaIbu                   sql.NullString         `db:"nama_ibu"`
	JenisKelamin              string                 `db:"jenis_kelamin"`
	TempatLahir               sql.NullString         `db:"tempat_lahir"`
	Agama                     sql.NullString         `db:"agama"`
	TglLahir                  sql.NullTime           `db:"tgl_lahir"`
	Alamat                    sql.NullString         `db:"alamat"`
	Domisili                  sql.NullString         `db:"domisili"`
	Kewarganegaraan           string                 `db:"kewarganegaraan"`
	StatusKawin               string                 `db:"status_kawin"`
	Kodepos                   sql.NullString         `db:"kodepos"`
	IDKelurahan               string                 `db:"id_kelurahan"`
	NoNpwp                    string                 `db:"no_npwp"`
	FotoNpwp                  string                 `db:"foto_npwp"`
	NoSid                     sql.NullString         `db:"no_sid"`
	FotoSid                   sql.NullString         `db:"foto_sid"`
	KodeCabang                string                 `db:"kode_cabang"`
	FotoURL                   sql.NullString         `db:"foto_url"`
	FotoKtpURL                string                 `db:"foto_ktp_url"`
	Status                    sql.NullInt64          `db:"status"`
	IsLocked                  sql.NullInt64          `db:"is_locked"`
	LoginFailCount            int64                  `db:"login_fail_count"`
	EmailVerified             int64                  `db:"email_verified"`
	KycVerified               int64                  `db:"kyc_verified"`
	EmailVerificationToken    string                 `db:"email_verification_token"`
	Token                     sql.NullString         `db:"token"`
	TokenWeb                  string                 `db:"token_web"`
	FcmToken                  string                 `db:"fcm_token"`
	LastUpdate                sql.NullTime           `db:"last_update"`
	Norek                     string                 `db:"norek"`
	Saldo                     int64                  `db:"saldo"`
	PinTemp                   sql.NullString         `db:"pin_temp"`
	LastUpdateDataNasabah     string                 `db:"last_update_data_nasabah"`
	LastUpdateDataNpwp        sql.NullTime           `db:"last_update_data_npwp"`
	LastUpdateLinkCif         sql.NullTime           `db:"last_update_link_cif"`
	LastUpdateUnlinkCif       sql.NullTime           `db:"last_update_unlink_cif"`
	LastUpdatePin             sql.NullTime           `db:"last_update_pin"`
	AktifasiTransFinansial    constant.ControlStatus `db:"aktifasiTransFinansial"`
	TanggalAktifasiFinansial  sql.NullTime           `db:"tanggal_aktifasi_finansial"`
	IsDukcapilVerified        sql.NullInt64          `db:"is_dukcapil_verified"`
	IsOpenTe                  sql.NullInt64          `db:"is_open_te"`
	ReferralCode              sql.NullString         `db:"referral_code"`
	GoldcardApplicationNumber sql.NullString         `db:"goldcard_application_number"`
	GoldcardAccountNumber     sql.NullString         `db:"goldcard_account_number"`
	TryLoginDate              sql.NullTime           `db:"try_login_date"`
	WrongPasswordCount        int64                  `db:"wrong_password_count"`
	BlockedDate               sql.NullTime           `db:"blocked_date"`
	BlockedToDate             sql.NullTime           `db:"blocked_to_date"`
	NorekUtama                sql.NullString         `db:"norek_utama"`
	IsSetBiometric            sql.NullInt64          `db:"is_set_biometric"`
	DeviceIDBiometric         sql.NullString         `db:"device_id_biometric"`
}

func UserToCustomerProfileDTO(user *User) *dto.CustomerProfileVO {
	dateOfBirth := sql.NullString{}

	if dob := user.TglLahir; dob.Valid && !dob.Time.IsZero() {
		dateOfBirth.Valid = true
		dateOfBirth.String = nval.ParseStringFallback(user.TglLahir.Time.Format("02-01-2006"), "")
	}

	return &dto.CustomerProfileVO{
		MaidenName:         user.NamaIbu.String,
		Gender:             user.JenisKelamin,
		Nationality:        user.Kewarganegaraan,
		DateOfBirth:        dateOfBirth.String,
		PlaceOfBirth:       user.TempatLahir.String,
		IdentityPhotoFile:  user.FotoKtpURL,
		MarriageStatus:     user.StatusKawin,
		NPWPNumber:         user.NoNpwp,
		NPWPPhotoFile:      user.FotoNpwp,
		SidPhotoFile:       user.FotoSid.String,
		NPWPUpdatedAt:      user.LastUpdateDataNpwp.Time.Unix(),
		ProfileUpdatedAt:   user.LastUpdate.Time.Unix(),
		CifLinkUpdatedAt:   user.LastUpdate.Time.Unix(),
		CifUnlinkUpdatedAt: user.LastUpdate.Time.Unix(),
	}
}

func UserToCustomer(user *User) (*Customer, error) {
	// Prepare profile value object
	profile := UserToCustomerProfileDTO(user)

	photo := &CustomerPhoto{
		Xid:      strings.ToUpper(xid.New().String()),
		FileName: nval.ParseStringFallback(user.FotoURL, ""),
		FileSize: 0,
		Mimetype: "",
	}

	userSnapshot, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	// Add metadata
	customerMetaData := &dto.CustomerMetadata{
		Snapshot:          string(userSnapshot),
		SnapshotSignature: "", // TODO snapshot signature
	}

	customerMetadataRaw, err := json.Marshal(customerMetaData)
	if err != nil {
		return nil, err
	}

	baseField := EmptyBaseField
	baseField.Metadata = customerMetadataRaw

	customerXID := strings.ToUpper(xid.New().String())
	customer := &Customer{
		CustomerXID:    customerXID,
		FullName:       user.Nama.String,
		Phone:          user.NoHp.String,
		Email:          user.Email.String,
		IdentityNumber: user.NoKtp.String,
		Photos:         photo,
		Cif:            user.Cif,
		Sid:            user.NoSid.String,
		UserRefID:      sql.NullString{String: nval.ParseStringFallback(user.UserAiid, "")},
		IdentityType:   nval.ParseInt64Fallback(user.JenisIdentitas, 0),
		Profile:        ToCustomerProfile(profile),
		ReferralCode:   user.ReferralCode,
		Status:         user.Status.Int64,
		BaseField:      baseField,
	}

	return customer, nil
}

func UserToCredential(user *User, userPin *UserPin) (*Credential, error) {
	credential := &Credential{
		Xid:                 strings.ToUpper(xid.New().String()),
		Password:            user.Password.String,
		NextPasswordResetAt: ModifierNullTime(user.NextPasswordReset),
		Pin:                 user.Pin.String,
		PinCif:              sql.NullString{},
		PinUpdatedAt:        ModifierNullTime(user.LastUpdatePin),
		PinLastAccessAt:     sql.NullTime{},
		PinCounter:          0,
		PinBlockedStatus:    0,
		IsLocked:            user.IsLocked.Int64,
		LoginFailCount:      user.LoginFailCount,
		WrongPasswordCount:  user.WrongPasswordCount,
		BlockedAt:           ModifierNullTime(user.BlockedDate),
		BlockedUntilAt:      ModifierNullTime(user.BlockedToDate),
		BiometricLogin:      user.IsSetBiometric.Int64,
		BiometricDeviceID:   user.DeviceIDBiometric.String,
		BaseField:           EmptyBaseField,
	}

	metadata := dto.MetadataCredential{
		TryLoginAt:   user.TryLoginDate.Time.String(),
		PinCreatedAt: "",
		PinBlockedAt: "",
	}

	metadataRawMessage, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	credential.Metadata = metadataRawMessage
	if userPin != nil {
		metadata.PinCreatedAt = userPin.CreatedAt
		metadata.PinBlockedAt = userPin.BlockedDate.Time.String()

		credential.PinCounter = userPin.Counter
		credential.PinLastAccessAt = userPin.LastAccessTime
		credential.PinCif = userPin.Cif
		credential.PinBlockedStatus = userPin.IsBlocked
	}

	return credential, nil
}

func UserToVerification(user *User) (*Verification, error) {
	// Verification
	verification := &Verification{
		Xid:                             strings.ToUpper(xid.New().String()),
		KycVerifiedStatus:               user.KycVerified,
		KycVerifiedAt:                   sql.NullTime{},
		EmailVerificationToken:          user.EmailVerificationToken,
		EmailVerifiedStatus:             user.EmailVerified,
		EmailVerifiedAt:                 sql.NullTime{},
		DukcapilVerifiedStatus:          user.IsDukcapilVerified.Int64,
		DukcapilVerifiedAt:              sql.NullTime{},
		FinancialTransactionStatus:      user.AktifasiTransFinansial,
		FinancialTransactionActivatedAt: ModifierNullTime(user.TanggalAktifasiFinansial),
		BaseField:                       EmptyBaseField,
	}

	return verification, nil
}

func UserToFinancialData(user *User) (*FinancialData, error) {
	financialData := &FinancialData{
		Xid:                       strings.ToUpper(xid.New().String()),
		MainAccountNumber:         user.NorekUtama.String,
		AccountNumber:             user.Norek,
		GoldSavingStatus:          user.IsOpenTe.Int64,
		GoldCardApplicationNumber: user.GoldcardApplicationNumber.String,
		GoldCardAccountNumber:     user.GoldcardApplicationNumber.String,
		Balance:                   user.Saldo,
		BaseField:                 EmptyBaseField,
	}

	return financialData, nil
}

func UserToAddress(user *User, userAddress *AddressExternal) (*Address, error) {
	var purpose int64
	purpose = constant.IdentityCard
	if user.Domisili.String == "1" {
		purpose = constant.Domicile
	}

	address := &Address{
		BaseField: EmptyBaseField,
		XID:       strings.ToUpper(xid.New().String()),
		Purpose:   purpose,
		ProvinceID: sql.NullInt64{
			Int64: userAddress.IDProvinsi.Int64,
			Valid: true,
		},
		ProvinceName: sql.NullString{
			String: userAddress.Provinsi.String,
			Valid:  true,
		},
		CityID: sql.NullInt64{
			Int64: userAddress.IDKabupaten.Int64,
			Valid: true,
		},
		CityName: sql.NullString{
			String: userAddress.Kabupaten.String,
			Valid:  true,
		},
		DistrictID: sql.NullInt64{
			Int64: userAddress.IDKecamatan.Int64,
			Valid: true,
		},
		DistrictName: sql.NullString{
			String: userAddress.Kecamatan.String,
			Valid:  true,
		},
		SubDistrictID: sql.NullInt64{
			Int64: userAddress.IDKelurahan.Int64,
			Valid: true,
		},
		SubDistrictName: sql.NullString{
			String: userAddress.Kelurahan.String,
			Valid:  true,
		},
		PostalCode: userAddress.Kodepos,
		Line:       sql.NullString{String: "", Valid: false},
		IsPrimary:  sql.NullBool{Bool: false, Valid: false},
	}

	return address, nil
}
