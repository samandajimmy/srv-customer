package statement

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	q "github.com/nbs-go/nsql/pq/query"
	"github.com/nbs-go/nsql/schema"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

var VerificationOTPSchema = schema.New(schema.FromModelRef(model.VerificationOTP{}))

type VerificationOTP struct {
	Insert               *sqlx.NamedStmt
	FindByRegistrationId *sqlx.Stmt
	Delete               *sqlx.Stmt
}

func NewVerificationOTP(db *nsql.DatabaseContext) *VerificationOTP {
	return &VerificationOTP{
		Insert: db.PrepareNamedFmtRebind(
			fmt.Sprintf(`%s ON CONFLICT DO NOTHING RETURNING "id"`, q.
				Insert(VerificationOTPSchema, "*").Build(),
			),
		),
		FindByRegistrationId: db.PrepareFmtRebind(q.
			Select(q.Column("*")).
			From(VerificationOTPSchema).
			Where(q.Equal(q.Column("registrationId")), q.Equal(q.Column("phone"))).
			Build(),
		),
		Delete: db.PrepareFmtRebind(q.
			Delete(VerificationOTPSchema).
			Where(q.Equal(q.Column("registrationId"))).
			Where(q.Equal(q.Column("phone"))).
			Build(),
		),
	}
}
