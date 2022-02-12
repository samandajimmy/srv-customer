package statement

import (
	"github.com/jmoiron/sqlx"
	q "github.com/nbs-go/nsql/pq/query"
	"github.com/nbs-go/nsql/schema"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

var CredentialSchema = schema.New(schema.FromModelRef(model.Credential{}))

type Credential struct {
	FindByCustomerID *sqlx.Stmt
	Insert           *sqlx.NamedStmt
	Update           *sqlx.NamedStmt
	DeleteByID       *sqlx.Stmt
}

func NewCredential(db *nsql.DatabaseContext) *Credential {
	return &Credential{
		FindByCustomerID: db.PrepareFmtRebind(q.
			Select(q.Column("*")).
			From(CredentialSchema).
			Where(q.Equal(q.Column("customerId"))).
			Build(),
		),
		Insert: db.PrepareNamedFmtRebind(q.
			Insert(CredentialSchema, "*").
			Build()),
		Update: db.PrepareNamedFmtRebind(q.
			Update(CredentialSchema, "*").
			Where(q.Equal(q.Column("customerId"))).
			Build()),
		DeleteByID: db.PrepareFmtRebind(q.
			Delete(CredentialSchema).
			Where(q.Equal(q.Column(CredentialSchema.PrimaryKey()))).
			Build()),
	}
}
