package model

import (
	"database/sql"
	"encoding/json"
)

type Credential struct {
	Xid                 string          `db:"xid"`
	CustomerId          int64           `db:"customerId"`
	Password            string          `db:"password"`
	NextPasswordResetAt sql.NullTime    `db:"nextPasswordResetAt"`
	Pin                 string          `db:"pin"`
	PinCif              string          `db:"pinCif"`
	PinUpdatedAt        sql.NullTime    `db:"pinUpdatedAt"`
	PinLastAccessAt     sql.NullTime    `db:"pinLastAccessAt"`
	PinCounter          int64           `db:"pinCounter"`
	PinBlockedStatus    int64           `db:"pinBlockedStatus"`
	IsLocked            int64           `db:"isLocked"`
	LoginFailCount      int64           `db:"loginFailCount"`
	WrongPasswordCount  int64           `db:"wrongPasswordCount"`
	BlockedAt           sql.NullTime    `db:"blockedAt"`
	BlockedUntilAt      sql.NullTime    `db:"blockedUntilAt"`
	BiometricLogin      int64           `db:"biometricLogin"`
	BiometricDeviceId   string          `db:"biometricDeviceId"`
	Metadata            json.RawMessage `db:"metadata"`
	ItemMetadata
}
