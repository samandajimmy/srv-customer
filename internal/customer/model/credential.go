package model

import (
	"database/sql"
)

type Credential struct {
	BaseField
	ID                  string         `db:"id"`
	Xid                 string         `db:"xid"`
	CustomerID          int64          `db:"customerId"`
	Password            string         `db:"password"`
	NextPasswordResetAt sql.NullTime   `db:"nextPasswordResetAt"`
	Pin                 string         `db:"pin"`
	PinCif              sql.NullString `db:"pinCif"`
	PinUpdatedAt        sql.NullTime   `db:"pinUpdatedAt"`
	PinLastAccessAt     sql.NullTime   `db:"pinLastAccessAt"`
	PinCounter          int64          `db:"pinCounter"`
	PinBlockedStatus    int64          `db:"pinBlockedStatus"`
	IsLocked            int64          `db:"isLocked"`
	LoginFailCount      int64          `db:"loginFailCount"`
	WrongPasswordCount  int64          `db:"wrongPasswordCount"`
	BlockedAt           sql.NullTime   `db:"blockedAt"`
	BlockedUntilAt      sql.NullTime   `db:"blockedUntilAt"`
	BiometricLogin      int64          `db:"biometricLogin"`
	BiometricDeviceID   string         `db:"biometricDeviceId"`
}

type UpdatePassword struct {
	CustomerID int64  `db:"customerId"`
	Password   string `db:"password"`
}
