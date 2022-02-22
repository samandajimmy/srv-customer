package model

import "database/sql"

type UserPin struct {
	UserID         int64          `db:"user_id"`
	Cif            sql.NullString `db:"cif"`
	LastAccessTime sql.NullTime   `db:"last_access_time"`
	Counter        int64          `db:"counter"`
	IsBlocked      int64          `db:"is_blocked"`
	BlockedDate    sql.NullTime   `db:"blocked_date"`
	CreatedAt      string         `db:"created_at"`
	UpdatedAt      sql.NullTime   `db:"updated_at"`
}
