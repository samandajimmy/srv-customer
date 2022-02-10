package model

import "database/sql"

// User model old structure from PDS-API
type User struct {
	UserAiid                  int64          `db:"user_AIID"`
	JenisIdentitas            string         `db:"jenis_identitas"`
	NoKtp                     sql.NullString `db:"no_ktp"`
	TanggalExpiredIdentitas   string         `db:"tanggal_expired_identitas"`
	Cif                       string         `db:"cif"`
	Email                     sql.NullString `db:"email"`
	NoHp                      sql.NullString `db:"no_hp"`
	Password                  sql.NullString `db:"password"`
	NextPasswordReset         sql.NullTime   `db:"next_password_reset"`
	Pin                       sql.NullString `db:"pin"`
	Nama                      sql.NullString `db:"nama"`
	NamaIbu                   sql.NullString `db:"nama_ibu"`
	JenisKelamin              string         `db:"jenis_kelamin"`
	TempatLahir               sql.NullString `db:"tempat_lahir"`
	Agama                     sql.NullString `db:"agama"`
	TglLahir                  sql.NullTime   `db:"tgl_lahir"`
	Alamat                    sql.NullString `db:"alamat"`
	Domisili                  sql.NullString `db:"domisili"`
	Kewarganegaraan           string         `db:"kewarganegaraan"`
	StatusKawin               string         `db:"status_kawin"`
	Kodepos                   sql.NullString `db:"kodepos"`
	IdKelurahan               string         `db:"id_kelurahan"`
	NoNpwp                    string         `db:"no_npwp"`
	FotoNpwp                  string         `db:"foto_npwp"`
	NoSid                     sql.NullString `db:"no_sid"`
	FotoSid                   sql.NullString `db:"foto_sid"`
	KodeCabang                string         `db:"kode_cabang"`
	FotoUrl                   sql.NullString `db:"foto_url"`
	FotoKtpUrl                string         `db:"foto_ktp_url"`
	Status                    sql.NullInt64  `db:"status"`
	IsLocked                  sql.NullInt64  `db:"is_locked"`
	LoginFailCount            int64          `db:"login_fail_count"`
	EmailVerified             int64          `db:"email_verified"`
	KycVerified               int64          `db:"kyc_verified"`
	EmailVerificationToken    string         `db:"email_verification_token"`
	Token                     sql.NullString `db:"token"`
	TokenWeb                  string         `db:"token_web"`
	FcmToken                  string         `db:"fcm_token"`
	LastUpdate                sql.NullTime   `db:"last_update"`
	Norek                     string         `db:"norek"`
	Saldo                     int64          `db:"saldo"`
	PinTemp                   sql.NullString `db:"pin_temp"`
	LastUpdateDataNasabah     string         `db:"last_update_data_nasabah"`
	LastUpdateDataNpwp        string         `db:"last_update_data_npwp"`
	LastUpdateLinkCif         string         `db:"last_update_link_cif"`
	LastUpdateUnlinkCif       string         `db:"last_update_unlink_cif"`
	LastUpdatePin             sql.NullTime   `db:"last_update_pin"`
	AktifasiTransFinansial    int64          `db:"aktifasiTransFinansial"`
	TanggalAktifasiFinansial  sql.NullTime   `db:"tanggal_aktifasi_finansial"`
	IsDukcapilVerified        sql.NullInt64  `db:"is_dukcapil_verified"`
	IsOpenTe                  sql.NullInt64  `db:"is_open_te"`
	ReferralCode              sql.NullString `db:"referral_code"`
	GoldcardApplicationNumber sql.NullString `db:"goldcard_application_number"`
	GoldcardAccountNumber     sql.NullString `db:"goldcard_account_number"`
	TryLoginDate              sql.NullTime   `db:"try_login_date"`
	WrongPasswordCount        int64          `db:"wrong_password_count"`
	BlockedDate               sql.NullTime   `db:"blocked_date"`
	BlockedToDate             sql.NullTime   `db:"blocked_to_date"`
	NorekUtama                sql.NullString `db:"norek_utama"`
	IsSetBiometric            sql.NullInt64  `db:"is_set_biometric"`
	DeviceIdBiometric         sql.NullString `db:"device_id_biometric"`
}
