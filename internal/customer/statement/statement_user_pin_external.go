package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

var UserPinSchema = nsql.NewSchema(
	"user_pin",
	[]string{"user_id", "cif", "last_access_time", "counter", "is_blocked", "blocked_date", "created_at", "updated_at"},
	nsql.WithAlias("up"),
	nsql.WithPrimaryKey("user_id"),
	nsql.WithAutoIncrement(false),
)

type UserPin struct {
	FindByCustomerID *sqlx.Stmt
}

func NewUserPin(db *nsql.DatabaseContext) *UserPin {
	return &UserPin{
		FindByCustomerID: db.PrepareFmtRebind(`SELECT %s FROM %s WHERE user_id = ?`,
			UserPinSchema.SelectAllColumns(), UserPinSchema.TableName),
	}
}
