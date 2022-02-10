package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type OTP struct {
	Insert *sqlx.NamedStmt
}

func NewOTP(db *nsql.DatabaseContext) *OTP {
	tableName := "OTP"
	columns := `"updatedAt","customerId","content","type","data","status"`
	namedColumns := `:updatedAt, :customerId, :content, :type, :data, :status`

	return &OTP{
		Insert: db.PrepareNamedFmt(`INSERT INTO "%s" (%s) VALUES (%s)`, tableName, columns, namedColumns),
	}
}
