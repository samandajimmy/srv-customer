package statement

import (
	"github.com/jmoiron/sqlx"
	"github.com/nbs-go/nsql/pq/query"
	"github.com/nbs-go/nsql/schema"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

// Define schema

var TransactionFavoriteSchema = schema.New(schema.FromModelRef(model.TransactionFavorite{}))

type TransactionFavorite struct {
	Insert                           *sqlx.NamedStmt
	DeleteByXID                      *sqlx.Stmt
	FindByXID                        *sqlx.Stmt
	FindByAccountNumberAndCustomerID *sqlx.Stmt
	FindAllByCustomerID              *sqlx.Stmt
}

func NewTransactionFavorite(db *nsql.DatabaseContext) *TransactionFavorite {
	// Init query Schema Builder
	sb := query.Schema(TransactionFavoriteSchema)

	// Init query
	findByXID := query.Select(query.Column("*")).
		From(TransactionFavoriteSchema).
		Where(query.Equal(query.Column("xid"))).
		Limit(1).
		Build()

	findAllByCustomerID := query.Select(query.Column("*")).
		From(TransactionFavoriteSchema).
		Where(query.Equal(query.Column("customerId"))).
		Build()

	deleteByXID := query.Delete(TransactionFavoriteSchema).
		Where(query.Equal(query.Column("xid"))).
		Build()

	findByAccountNumberAndCustomerID := query.
		Select(query.Column("*")).
		From(TransactionFavoriteSchema).
		Where(
			query.Equal(query.Column("accountNumber")),
			query.Equal(query.Column("customerId"))).
		Limit(1).
		Build()

	return &TransactionFavorite{
		Insert:                           db.PrepareNamed(sb.Insert()),
		DeleteByXID:                      db.PrepareFmtRebind(deleteByXID),
		FindByXID:                        db.PrepareFmtRebind(findByXID),
		FindByAccountNumberAndCustomerID: db.PrepareFmtRebind(findByAccountNumberAndCustomerID),
		FindAllByCustomerID:              db.PrepareFmtRebind(findAllByCustomerID),
	}
}
