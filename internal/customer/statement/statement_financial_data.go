package statement

import (
	"github.com/jmoiron/sqlx"
	q "github.com/nbs-go/nsql/pq/query"
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
	return &FinancialData{
		FindByCustomerID: db.PrepareFmtRebind(q.Select(q.Column("*")).
			Where(q.Equal(q.Column("customerId"))).
			From(FinancialDataSchema).
			Build()),
		Insert: db.PrepareNamedFmtRebind(q.
			Insert(FinancialDataSchema, "*").
			Build()),
		Update: db.PrepareNamedFmtRebind(q.Update(FinancialDataSchema, "*").
			Where(q.Equal(q.Column("customerId"))).
			Build()),
		DeleteByID: db.PrepareFmtRebind(q.
			Delete(FinancialDataSchema).
			Where(q.Equal(q.Column(FinancialDataSchema.PrimaryKey()))).
			Build(),
		),
	}
}
