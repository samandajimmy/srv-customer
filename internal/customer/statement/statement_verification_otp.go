package statement

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/nbs-go/nsql/pq/query"
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
	// Init query Schema Builder
	sb := query.Schema(VerificationSchema)

	insert := fmt.Sprintf(`%s ON CONFLICT DO NOTHING RETURNING "id"`, sb.Insert())

	findByRegistrationId := query.
		Select(query.Column("*")).
		From(VerificationOTPSchema).
		Where(query.Equal(query.Column("registrationId")), query.Equal(query.Column("phone"))).
		Build()

	deleteByRegistrationIdAndPhone := query.
		Delete(VerificationOTPSchema).
		Where(
			query.And(
				query.Equal(query.Column("registrationId")),
				query.Equal(query.Column("phone")),
			),
		).
		Build()

	return &VerificationOTP{
		Insert:               db.PrepareNamedFmtRebind(insert),
		FindByRegistrationId: db.PrepareFmtRebind(findByRegistrationId),
		Delete:               db.PrepareFmtRebind(deleteByRegistrationIdAndPhone),
	}
}
