package model

import (
	"encoding/json"
	"time"
)

type Credential struct {
	Xid                 string          `db:"xid"`
	CustomerId          int64           `db:"customerId"`
	Password            string          `db:"password"`
	NextPasswordResetAt *time.Time      `db:"nextPasswordResetAt"`
	Pin                 string          `db:"pin"`
	PinCif              string          `db:"pinCif"`
	PinUpdatedAt        *time.Time      `db:"pinUpdatedAt"`
	PinLastAccessAt     *time.Time      `db:"pinLastAccessAt"`
	PinCounter          int64           `db:"pinCounter"`
	PinBlockedStatus    int64           `db:"pinBlockedStatus"`
	IsLocked            int64           `db:"isLocked"`
	LoginFailCount      int64           `db:"loginFailCount"`
	WrongPasswordCount  int64           `db:"wrongPasswordCount"`
	BlockedAt           *time.Time      `db:"blockedAt"`
	BlockedUntilAt      *time.Time      `db:"blockedUntilAt"`
	BiometricLogin      int64           `db:"biometricLogin"`
	BiometricDeviceId   string          `db:"biometricDeviceId"`
	Metadata            json.RawMessage `db:"metadata"`
	ItemMetadata
}
