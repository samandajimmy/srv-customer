package model

import (
	"time"
)

type UserRegister struct {
	Id        string    `db:"id"`
	NoHp      string    `db:"no_hp"`
	CreatedAt time.Time `db:"created_at"`
}
