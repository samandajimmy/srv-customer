package model

import "database/sql"

type UserRegister struct {
	Id        string       `db:"id"`
	NoHp      string       `db:"no_hp"`
	CreatedAt sql.NullTime `db:"createdAt"`
}
