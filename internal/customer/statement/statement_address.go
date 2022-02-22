package statement

import (
	"github.com/jmoiron/sqlx"
	"github.com/nbs-go/nsql/pq/query"
	"github.com/nbs-go/nsql/schema"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

var AddressSchema = schema.New(schema.FromModelRef(model.Address{}))

type Address struct {
	FindByCustomerID        *sqlx.Stmt
	Insert                  *sqlx.NamedStmt
	Update                  *sqlx.NamedStmt
	FindPrimaryByCustomerID *sqlx.Stmt
}

func NewAddress(db *nsql.DatabaseContext) *Address {
	// Init query Schema Builder
	sb := query.Schema(AddressSchema)

	// Init query
	findByCustomerID := query.Select(query.Column("*")).
		From(AddressSchema).
		Where(query.Equal(query.Column("customerId"))).
		Build()

	findPrimaryByCustomerID := query.
		Select(query.Column("*")).
		From(AddressSchema).
		Where(
			query.Equal(query.Column("isPrimary")),
			query.Equal(query.Column("customerId")),
		).
		Limit(1).
		Build()

	return &Address{
		Insert:                  db.PrepareNamedFmtRebind(sb.Insert()),
		Update:                  db.PrepareNamedFmtRebind(sb.Update()),
		FindByCustomerID:        db.PrepareFmtRebind(findByCustomerID),
		FindPrimaryByCustomerID: db.PrepareFmtRebind(findPrimaryByCustomerID),
	}
}
