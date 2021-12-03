package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type OTPStatement struct {
	Insert *sqlx.NamedStmt
}

func NewOTPStatement(db *nsql.DB) *OTPStatement {
	tableName := "OTP"
	columns := `updatedAt,customerId,content,type,data,status`
	namedColumns := `:updatedAt, :customerId, :content, :type, :data, :status`

	return &OTPStatement{
		Insert: db.PrepareNamedFmt("INSERT INTO %s(%s) VALUES (%s)", tableName, columns, namedColumns),
	}
}
