package statement

import (
	"github.com/jmoiron/sqlx"
	q "github.com/nbs-go/nsql/pq/query"
	"github.com/nbs-go/nsql/schema"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

var AccessSessionSchema = schema.New(schema.FromModelRef(model.AccessSession{}))

type AccessSession struct {
	Insert *sqlx.NamedStmt
	Update *sqlx.NamedStmt
}

func NewAccessSession(db *nsql.DatabaseContext) *AccessSession {
	s := AccessSessionSchema
	return &AccessSession{
		Insert: db.PrepareNamedFmt(q.Insert(s, "*").Build()),
		Update: db.PrepareNamedFmt(q.Update(s, "*").Build()),
	}
}
