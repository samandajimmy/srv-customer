package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type UserPin struct {
	FindByCustomerID *sqlx.Stmt
}

func NewUserPin(db *nsql.DatabaseContext) *UserPin {
	// TODO: Update query build nsql (mysql ver)
	tableName := `user_pin`
	getColumns := `user_id, cif, last_access_time, counter, is_blocked, blocked_date, created_at, updated_at`

	return &UserPin{
		FindByCustomerID: db.PrepareFmt(`SELECT %s FROM %s WHERE user_id = ?`, getColumns, tableName),
	}
}
