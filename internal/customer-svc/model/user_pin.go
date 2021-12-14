package model

import "database/sql"

type UserPin struct {
	UserId         int64          `db:"userId"`
	Cif            sql.NullString `db:"cif"`
	LastAccessTime sql.NullTime   `db:"lastAccessTime"`
	Counter        int64          `db:"counter"`
	IsBlocked      int64          `db:"isBlocked"`
	BlockedDate    sql.NullTime   `db:"blockedDate"`
	CreatedAt      string         `db:"createdAt"`
	UpdatedAt      sql.NullTime   `db:"updatedAt"`
}
