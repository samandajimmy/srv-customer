package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type VerificationOTPStatement struct {
	Insert *sqlx.NamedStmt
}

func NewVerificationOTPStatement(db *nsql.DB) *VerificationOTPStatement {
	customerTable := "VerificationOTP"
	columns := "\"createdAt\", \"registrationId\", \"phone\""
	namedColumns := ":createdAt,:registrationId,:phone"
	return &VerificationOTPStatement{
		Insert: db.PrepareNamedFmt("INSERT INTO \"%s\"(%s) VALUES (%s) RETURNING id", customerTable, columns, namedColumns),
	}
}
