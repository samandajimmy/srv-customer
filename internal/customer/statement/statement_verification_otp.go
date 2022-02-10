package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type VerificationOTP struct {
	Insert               *sqlx.NamedStmt
	FindByRegistrationId *sqlx.Stmt
	Delete               *sqlx.Stmt
}

func NewVerificationOTP(db *nsql.DatabaseContext) *VerificationOTP {
	verificationOTPTable := `VerificationOTP`
	columns := `"createdAt", "registrationId", "phone"`
	namedColumns := `:createdAt,:registrationId,:phone`
	return &VerificationOTP{
		Insert:               db.PrepareNamedFmt(`INSERT INTO "%s" (%s) VALUES (%s) RETURNING id`, verificationOTPTable, columns, namedColumns),
		FindByRegistrationId: db.PrepareFmt(`SELECT "registrationId" FROM "%s" WHERE "registrationId" = $1 AND phone = $2`, verificationOTPTable),
		Delete:               db.PrepareFmt(`DELETE FROM "%s" WHERE "registrationId" = $1 AND "phone" = $2`, verificationOTPTable),
	}
}
