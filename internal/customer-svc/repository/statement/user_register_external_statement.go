package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type UserRegisterStatement struct {
	Insert *sqlx.NamedStmt
}

func NewUserRegisterStatement(db *nsql.DB) *UserRegisterStatement {
	tableName := `user_register`
	columns := `id, no_hp, created_at`
	namedColumns := `:id, :no_hp, :created_at`

	return &UserRegisterStatement{
		Insert: db.PrepareNamedFmt(`INSERT INTO %s (%s) VALUES (%s)`, tableName, columns, namedColumns),
	}
}
