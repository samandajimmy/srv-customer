package statement

import (
	"github.com/jmoiron/sqlx"
	"github.com/nbs-go/nsql/pq/query"
	"github.com/nbs-go/nsql/schema"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

var CredentialSchema = schema.New(schema.FromModelRef(model.Credential{}))

type Credential struct {
	FindByCustomerID            *sqlx.Stmt
	Insert                      *sqlx.NamedStmt
	Update                      *sqlx.NamedStmt
	UpdatePasswordByCustomerID  *sqlx.NamedStmt
	DeleteByID                  *sqlx.Stmt
	FindByPasswordAndCustomerID *sqlx.Stmt
}

func NewCredential(db *nsql.DatabaseContext) *Credential {
	// Init query Schema Builder
	sb := query.Schema(CredentialSchema)

	// Init query
	updateByCustomerId := query.
		Update(CredentialSchema, "*").
		Where(query.Equal(query.Column("customerId"))).
		Build()

	findByPasswordAndCustomerId := query.
		Select(query.Column("*")).
		From(CredentialSchema).
		Where(
			query.Equal(query.Column("customerId")),
			query.Equal(query.Column("password")),
		).Build()

	findByCustomerId := query.
		Select(query.Column("*")).
		From(CredentialSchema).
		Where(query.Equal(query.Column("customerId"))).
		Build()

	updatePasswordByCustomerId := query.
		Update(CredentialSchema, "password").
		Where(query.Equal(query.Column("customerId"))).
		Build()

	return &Credential{
		FindByCustomerID:            db.PrepareFmtRebind(findByCustomerId),
		FindByPasswordAndCustomerID: db.PrepareFmtRebind(findByPasswordAndCustomerId),
		Insert:                      db.PrepareNamedFmtRebind(sb.Insert()),
		Update:                      db.PrepareNamedFmtRebind(updateByCustomerId),
		DeleteByID:                  db.PrepareFmtRebind(sb.Delete()),
		UpdatePasswordByCustomerID:  db.PrepareNamedFmtRebind(updatePasswordByCustomerId),
	}
}
