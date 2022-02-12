package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

var UserSchema = nsql.NewSchema(
	"user",
	[]string{
		"user_AIID", "jenis_identitas", "no_ktp", "tanggal_expired_identitas", "cif", "email", "no_hp", "password",
		"next_password_reset", "pin", "nama", "nama_ibu", "jenis_kelamin", "tempat_lahir", "agama", "tgl_lahir", "alamat",
		"domisili", "kewarganegaraan", "status_kawin", "kodepos", "id_kelurahan", "no_npwp", "no_sid", "foto_sid",
		"kode_cabang", "foto_url", "foto_ktp_url", "status", "is_locked", "login_fail_count", "email_verified", "kyc_verified",
		"email_verification_token", "token", "token_web", "fcm_token", "last_update", "norek", "saldo", "pin_temp",
		"last_update_data_nasabah", "last_update_data_npwp", "last_update_link_cif", "last_update_unlink_cif",
		"last_update_pin", "aktifasiTransFinansial", "tanggal_aktifasi_finansial", "is_dukcapil_verified", "is_open_te",
		"referral_code", "goldcard_application_number", "goldcard_account_number", "try_login_date", "wrong_password_count",
		"blocked_date", "blocked_to_date", "norek_utama", "is_set_biometric", "device_id_biometric",
	},
	nsql.WithAlias("ue"),
	nsql.WithPrimaryKey("user_AIID"),
	nsql.WithAutoIncrement(true),
)

type User struct {
	FindByEmailOrPhone *sqlx.Stmt
}

func NewUser(db *nsql.DatabaseContext) *User {
	return &User{
		FindByEmailOrPhone: db.PrepareFmtRebind(
			`SELECT %s FROM %s WHERE (no_hp = (?)) OR (email = (?))`,
			UserSchema.SelectAllColumns(), UserSchema.TableName,
		),
	}
}
