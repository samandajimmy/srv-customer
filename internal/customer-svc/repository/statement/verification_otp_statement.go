package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type VerificationOTPStatement struct {
	Insert               *sqlx.NamedStmt
	FindByRegistrationId *sqlx.Stmt
	Delete               *sqlx.Stmt
}

func NewVerificationOTPStatement(db *nsql.DB) *VerificationOTPStatement {
	verificationOTPTable := `VerificationOTP`
	columns := `"createdAt", "registrationId", "phone"`
	namedColumns := `:createdAt,:registrationId,:phone`
	return &VerificationOTPStatement{
		Insert:               db.PrepareNamedFmt(`INSERT INTO "%s" (%s) VALUES (%s) RETURNING id`, verificationOTPTable, columns, namedColumns),
		FindByRegistrationId: db.PrepareFmt(`SELECT "registrationId" FROM "%s" WHERE "registrationId" = $1 AND phone = $2`, verificationOTPTable),
		Delete:               db.PrepareFmt(`DELETE FROM "%s" WHERE "registrationId" = $1 AND "phone" = $2`, verificationOTPTable),
	}
}
