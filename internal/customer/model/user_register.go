package model

import (
	"time"
)

type UserRegister struct {
	ID        string    `db:"id"`
	NoHp      string    `db:"no_hp"`
	CreatedAt time.Time `db:"created_at"`
}
