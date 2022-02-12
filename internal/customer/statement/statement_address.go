package statement

import (
	"github.com/jmoiron/sqlx"
	q "github.com/nbs-go/nsql/pq/query"
	"github.com/nbs-go/nsql/schema"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

var AddressSchema = schema.New(schema.FromModelRef(model.Address{}))

type Address struct {
	FindByCustomerId        *sqlx.Stmt
	Insert                  *sqlx.NamedStmt
	Update                  *sqlx.NamedStmt
	FindPrimaryByCustomerId *sqlx.Stmt
}

func NewAddress(db *nsql.DatabaseContext) *Address {
	return &Address{
		FindByCustomerId: db.PrepareFmtRebind(q.
			Select(q.Column("*")).
			From(AddressSchema).
			Where(q.Equal(q.Column("customerId"))).
			Build(),
		),
		Insert: db.PrepareNamedFmtRebind(q.Insert(AddressSchema, "*").Build()),
		Update: db.PrepareNamedFmtRebind(q.Update(AddressSchema, "*").Build()),
		FindPrimaryByCustomerId: db.PrepareFmtRebind(q.
			Select(q.Column("*")).
			From(AddressSchema).
			Where(
				q.Equal(q.Column("isPrimary")),
				q.Equal(q.Column("customerId")),
			).
			Limit(1).
			Build(),
		),
	}
}
