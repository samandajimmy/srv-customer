package statement

import (
	"github.com/jmoiron/sqlx"
	q "github.com/nbs-go/nsql/pq/query"
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
	return &Verification{
		FindByCustomerID: db.PrepareFmtRebind(q.
			Select(q.Column("*")).
			From(VerificationSchema).
			Where(q.Equal(q.Column("customerId"))).
			Build(),
		),
		FindByEmailToken: db.PrepareFmtRebind(q.
			Select(q.Column("*")).
			From(VerificationSchema).
			Where(q.Equal(q.Column("emailVerificationToken"))).
			Build(),
		),
		Insert: db.PrepareNamedFmtRebind(q.
			Insert(VerificationSchema, "*").
			Build()),
		Update: db.PrepareNamedFmtRebind(q.
			Update(VerificationSchema, "*").
			Where(q.Equal(q.Column("customerId"))).
			Build(),
		),
		DeleteByID: db.PrepareFmtRebind(q.
			Delete(VerificationSchema).
			Where(q.Equal(q.Column(VerificationSchema.PrimaryKey()))).
			Build(),
		),
	}
}
