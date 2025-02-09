package statement

import (
	"github.com/jmoiron/sqlx"
	"github.com/nbs-go/nsql/pq/query"
	"github.com/nbs-go/nsql/schema"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

var FinancialDataSchema = schema.New(schema.FromModelRef(model.FinancialData{}))

type FinancialData struct {
	FindByCustomerID *sqlx.Stmt
	Insert           *sqlx.NamedStmt
	Update           *sqlx.NamedStmt
	DeleteByID       *sqlx.Stmt
}

func NewFinancialData(db *nsql.DatabaseContext) *FinancialData {
	// Init query Schema Builder
	sb := query.Schema(FinancialDataSchema)

	// Init query
	findByCustomerID := query.Select(query.Column("*")).
		Where(query.Equal(query.Column("customerId"))).
		From(FinancialDataSchema).
		Build()

	updateByCustomerID := query.Update(FinancialDataSchema, "*").
		Where(query.Equal(query.Column("customerId"))).
		Build()

	return &FinancialData{
		FindByCustomerID: db.PrepareFmtRebind(findByCustomerID),
		Insert:           db.PrepareNamedFmtRebind(sb.Insert()),
		Update:           db.PrepareNamedFmtRebind(updateByCustomerID),
		DeleteByID:       db.PrepareFmtRebind(sb.Delete()),
	}
}
