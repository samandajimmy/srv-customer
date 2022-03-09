package statement

import (
	"github.com/jmoiron/sqlx"
	"github.com/nbs-go/nsql/pq/query"
	"github.com/nbs-go/nsql/schema"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

var VerificationSchema = schema.New(schema.FromModelRef(model.Verification{}))

type Verification struct {
	FindByCustomerID *sqlx.Stmt
	FindByEmailToken *sqlx.Stmt
	Insert           *sqlx.NamedStmt
	Update           *sqlx.NamedStmt
	DeleteByID       *sqlx.Stmt
}

func NewVerification(db *nsql.DatabaseContext) *Verification {
	// Init query Schema Builder
	sb := query.Schema(VerificationSchema)

	// Init query
	findByCustomerID := query.
		Select(query.Column("*")).
		From(VerificationSchema).
		Where(query.Equal(query.Column("customerId"))).
		Build()

	findByEmailToken := query.
		Select(query.Column("*")).
		From(VerificationSchema).
		Where(query.Equal(query.Column("emailVerificationToken"))).
		Build()

	updateByCustomerID := query.
		Update(VerificationSchema, "*").
		Where(query.Equal(query.Column("customerId"))).
		Build()

	return &Verification{

		FindByCustomerID: db.PrepareFmtRebind(findByCustomerID),
		FindByEmailToken: db.PrepareFmtRebind(findByEmailToken),
		Insert:           db.PrepareNamedFmtRebind(sb.Insert()),
		Update:           db.PrepareNamedFmtRebind(updateByCustomerID),
		DeleteByID:       db.PrepareFmtRebind(sb.Delete()),
	}
}
