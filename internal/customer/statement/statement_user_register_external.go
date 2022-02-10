package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type UserRegister struct {
	Insert *sqlx.NamedStmt
}

func NewUserRegister(db *nsql.DatabaseContext) *UserRegister {
	tableName := `user_register`
	columns := `id, no_hp, created_at`
	namedColumns := `:id, :no_hp, :created_at`

	return &UserRegister{
		Insert: db.PrepareNamedFmt(`INSERT INTO %s (%s) VALUES (%s)`, tableName, columns, namedColumns),
	}
}
