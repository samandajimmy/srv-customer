package statement

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	q "github.com/nbs-go/nsql/pq/query"
	"github.com/nbs-go/nsql/schema"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

var CustomerSchema = schema.New(schema.FromModelRef(model.Customer{}))

type Customer struct {
	Insert             *sqlx.NamedStmt
	UpdateByPhone      *sqlx.NamedStmt
	FindById           *sqlx.Stmt
	FindByPhone        *sqlx.Stmt
	FindByEmail        *sqlx.Stmt
	FindByEmailOrPhone *sqlx.Stmt
}

func NewCustomer(db *nsql.DatabaseContext) *Customer {
	return &Customer{
		Insert: db.PrepareNamedFmtRebind(fmt.Sprintf(`%s ON CONFLICT DO NOTHING RETURNING "id"`, q.
			Insert(CustomerSchema, "*").
			Build())),
		FindById: db.PrepareFmtRebind(q.Select(q.Column("*")).
			From(CustomerSchema).
			Where(q.Equal(q.Column(CustomerSchema.PrimaryKey()))).
			Build()),
		FindByPhone: db.PrepareFmtRebind(q.Select(q.Column("*")).
			From(CustomerSchema).
			Where(q.Equal(q.Column("phone"))).
			Build()),
		FindByEmail: db.PrepareFmtRebind(q.Select(q.Column("*")).
			From(CustomerSchema).
			Where(q.Equal(q.Column("email"))).
			Build()),
		FindByEmailOrPhone: db.PrepareFmtRebind(q.
			Select(q.Column("*")).
			From(CustomerSchema).
			Where(
				q.Or(
					q.Equal(q.Column("phone")),
					q.Equal(q.Column("email")),
				),
			).
			Build(),
		),
		UpdateByPhone: db.PrepareNamedFmtRebind(q.
			Update(CustomerSchema, "*").
			Where(q.Equal(q.Column("phone"))).
			Build()),
	}
}
