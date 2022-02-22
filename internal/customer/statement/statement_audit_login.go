package statement

import (
	"github.com/jmoiron/sqlx"
	"github.com/nbs-go/nsql/option"
	"github.com/nbs-go/nsql/pq/query"
	"github.com/nbs-go/nsql/schema"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

var AuditLoginSchema = schema.New(schema.FromModelRef(model.AuditLogin{}))

type AuditLogin struct {
	Insert             *sqlx.NamedStmt
	UpdateByCustomerID *sqlx.NamedStmt
	CountLogin         *sqlx.Stmt
}

func NewAuditLogin(db *nsql.DatabaseContext) *AuditLogin {
	// Init query Schema Builder
	sb := query.Schema(AuditLoginSchema)

	// Init Query
	countLogin := query.
		Select(query.Count("id", option.Schema(AuditLoginSchema), option.As("count"))).
		From(AuditLoginSchema).
		Where(query.Equal(query.Column("customerId"))).
		Build()

	return &AuditLogin{
		Insert: db.PrepareNamedFmtRebind(sb.Insert()),
		UpdateByCustomerID: db.PrepareNamedFmtRebind(query.
			Update(AuditLoginSchema, "*").
			Where(query.Equal(query.Column("customerId"))).
			Build(),
		),
		CountLogin: db.PrepareFmtRebind(countLogin),
	}
}
