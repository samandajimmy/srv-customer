package model

import (
	"encoding/json"
	"time"
)

type AccessSession struct {
	Xid                  string          `db:"xid"`
	CustomerId           int64           `db:"customerId"`
	ExpiredAt            time.Time       `db:"expiredAt"`
	NotificationToken    string          `db:"notificationToken"`
	NotificationProvider int64           `db:"notificationProvider"`
	Metadata             json.RawMessage `db:"metadata"`
	ItemMetadata
}

type UpdateAccessSession struct {
	*AccessSession
	CurrentVersion int64 `db:"currentVersion"`
}
