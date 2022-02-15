package statement

import (
	"github.com/jmoiron/sqlx"
	q "github.com/nbs-go/nsql/pq/query"
	"github.com/nbs-go/nsql/schema"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

var credentialSchema = schema.New(schema.FromModelRef(model.Credential{}))

type Credential struct {
	FindByCustomerID            *sqlx.Stmt
	Insert                      *sqlx.NamedStmt
	Update                      *sqlx.NamedStmt
	DeleteByID                  *sqlx.Stmt
	FindByPasswordAndCustomerID *sqlx.Stmt
}

func NewCredential(db *nsql.DatabaseContext) *Credential {
	return &Credential{
		FindByCustomerID: db.PrepareFmtRebind(q.
			Select(q.Column("*")).
			From(credentialSchema).
			Where(q.Equal(q.Column("customerId"))).
			Build(),
		),
		FindByPasswordAndCustomerID: db.PrepareFmtRebind(q.
			Select(q.Column("*")).
			From(credentialSchema).
			Where(
				q.Equal(q.Column("customerId")),
				q.Equal(q.Column("password")),
			).Build(),
		),
		Insert: db.PrepareNamedFmtRebind(q.
			Insert(credentialSchema, "*").
			Build()),
		Update: db.PrepareNamedFmtRebind(q.
			Update(credentialSchema, "*").
			Where(q.Equal(q.Column("customerId"))).
			Build()),
		DeleteByID: db.PrepareFmtRebind(q.
			Delete(credentialSchema).
			Where(q.Equal(q.Column(credentialSchema.PrimaryKey()))).
			Build()),
	}
}
