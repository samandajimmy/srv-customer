package statement

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	q "github.com/nbs-go/nsql/pq/query"
	"github.com/nbs-go/nsql/schema"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

var customerSchema = schema.New(schema.FromModelRef(model.Customer{}))

type Customer struct {
	Insert             *sqlx.NamedStmt
	UpdateByPhone      *sqlx.NamedStmt
	FindByRefId        *sqlx.Stmt
	FindById           *sqlx.Stmt
	FindByPhoneOrCIF   *sqlx.Stmt
	FindByPhone        *sqlx.Stmt
	FindByEmail        *sqlx.Stmt
	FindByEmailOrPhone *sqlx.Stmt
	ReferralCodeExist  *sqlx.Stmt
	UpdateByCIF        *sqlx.NamedStmt
}

func NewCustomer(db *nsql.DatabaseContext) *Customer {
	return &Customer{
		Insert: db.PrepareNamedFmtRebind(fmt.Sprintf(`%s ON CONFLICT DO NOTHING RETURNING "id"`, q.
			Insert(customerSchema, "*").
			Build())),
		FindById: db.PrepareFmtRebind(q.Select(q.Column("*")).
			From(customerSchema).
			Where(q.Equal(q.Column(customerSchema.PrimaryKey()))).
			Build()),
		FindByRefId: db.PrepareFmtRebind(q.Select(q.Column("*")).
			From(customerSchema).
			Where(q.Equal(q.Column("userRefId"))).
			Build()),
		FindByPhone: db.PrepareFmtRebind(q.Select(q.Column("*")).
			From(customerSchema).
			Where(q.Equal(q.Column("phone"))).
			Build()),
		FindByPhoneOrCIF: db.PrepareFmtRebind(q.Select(q.Column("*")).
			From(customerSchema).
			Where(
				q.Or(
					q.Equal(q.Column("cif")),
					q.Equal(q.Column("phone")),
				),
			).
			Build()),
		FindByEmail: db.PrepareFmtRebind(q.Select(q.Column("*")).
			From(customerSchema).
			Where(q.Equal(q.Column("email"))).
			Build()),
		ReferralCodeExist: db.PrepareFmtRebind(q.Select(q.Column("*")).
			From(customerSchema).
			Where(q.Equal(q.Column("referralCode"))).
			Limit(1).
			Build()),
		FindByEmailOrPhone: db.PrepareFmtRebind(q.
			Select(q.Column("*")).
			From(customerSchema).
			Where(
				q.Or(
					q.Equal(q.Column("phone")),
					q.Equal(q.Column("email")),
				),
			).
			Build(),
		),
		UpdateByCIF: db.PrepareNamedFmtRebind(q.
			Update(customerSchema, "*").
			Where(q.Equal(q.Column("cif"))).
			Build()),
		UpdateByPhone: db.PrepareNamedFmtRebind(q.
			Update(customerSchema, "*").
			Where(q.Equal(q.Column("phone"))).
			Build()),
	}
}
