package statement

import (
	"github.com/jmoiron/sqlx"
	q "github.com/nbs-go/nsql/pq/query"
	opt "github.com/nbs-go/nsql/query/option"
	"github.com/nbs-go/nsql/schema"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

var AuditLoginSchema = schema.New(schema.FromModelRef(model.AuditLogin{}))

type AuditLogin struct {
	Insert             *sqlx.NamedStmt
	UpdateByCustomerId *sqlx.NamedStmt
	CountLogin         *sqlx.Stmt
}

func NewAuditLogin(db *nsql.DatabaseContext) *AuditLogin {
	return &AuditLogin{
		Insert: db.PrepareNamedFmtRebind(q.Insert(AuditLoginSchema, "*").Build()),
		UpdateByCustomerId: db.PrepareNamedFmtRebind(q.
			Update(AuditLoginSchema, "*").
			Where(q.Equal(q.Column("customerId"))).
			Build(),
		),
		CountLogin: db.PrepareFmtRebind(q.
			Select(q.Count("id", opt.Schema(AuditLoginSchema), opt.As("count"))).
			From(AuditLoginSchema).
			Where(q.Equal(q.Column("customerId"))).
			Build(),
		),
	}
}
