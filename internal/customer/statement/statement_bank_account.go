package statement

import (
	"github.com/jmoiron/sqlx"
	"github.com/nbs-go/nsql/pq/query"
	"github.com/nbs-go/nsql/schema"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

// Define schema

var BankAccountSchema = schema.New(schema.FromModelRef(model.BankAccount{}))

type BankAccount struct {
	Insert                           *sqlx.NamedStmt
	Update                           *sqlx.NamedStmt
	DeleteByXID                      *sqlx.Stmt
	FindByXID                        *sqlx.Stmt
	FindAllByCustomerID              *sqlx.Stmt
	FindByAccountNumberAndCustomerID *sqlx.Stmt
}

func NewBankAccount(db *nsql.DatabaseContext) *BankAccount {
	// Init query Schema Builder
	sb := query.Schema(BankAccountSchema)

	// Init query
	findByXID := query.Select(query.Column("*")).
		From(BankAccountSchema).
		Where(query.Equal(query.Column("xid"))).
		Limit(1).
		Build()

	findAllByCustomerID := query.Select(query.Column("*")).
		From(BankAccountSchema).
		Where(query.Equal(query.Column("customerId"))).
		Build()

	update := query.
		Update(
			BankAccountSchema,
			"version", "modifiedBy", "updatedAt", "accountName", "accountNumber", "bank", "status").
		Build()

	deleteByXID := query.Delete(BankAccountSchema).
		Where(query.Equal(query.Column("xid"))).
		Build()

	findByAccountNumberAndCustomerID := query.
		Select(query.Column("*")).
		From(BankAccountSchema).
		Where(
			query.Equal(query.Column("accountNumber")),
			query.Equal(query.Column("customerId"))).
		Limit(1).
		Build()

	return &BankAccount{
		Insert:                           db.PrepareNamed(sb.Insert()),
		Update:                           db.PrepareNamed(update),
		DeleteByXID:                      db.PrepareFmtRebind(deleteByXID),
		FindByXID:                        db.PrepareFmtRebind(findByXID),
		FindAllByCustomerID:              db.PrepareFmtRebind(findAllByCustomerID),
		FindByAccountNumberAndCustomerID: db.PrepareFmtRebind(findByAccountNumberAndCustomerID),
	}
}
