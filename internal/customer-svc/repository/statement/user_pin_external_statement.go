package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type UserPinStatement struct {
	FindByCustomerID *sqlx.Stmt
}

func NewUserPinStatement(db *nsql.DB) *UserPinStatement {
	tableName := `user_pin`
	getColumns := `user_id, cif, last_access_time, counter, is_blocked, blocked_date, created_at, updated_at`

	return &UserPinStatement{
		FindByCustomerID: db.PrepareFmt(`SELECT %s FROM "%s" WHERE "user_id" = $1`, getColumns, tableName),
	}
}
