package statement

import (
	"github.com/jmoiron/sqlx"
	q "github.com/nbs-go/nsql/pq/query"
	"github.com/nbs-go/nsql/schema"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

var OTPSchema = schema.New(schema.FromModelRef(model.OTP{}))

type OTP struct {
	Insert *sqlx.NamedStmt
}

func NewOTP(db *nsql.DatabaseContext) *OTP {
	return &OTP{
		Insert: db.PrepareNamedFmtRebind(q.
			Insert(OTPSchema, "*").
			Build(),
		),
	}
}
