package model

import (
	"time"
)

type AccessSession struct {
	BaseField
	ID                   int64     `db:"id"`
	Xid                  string    `db:"xid"`
	CustomerID           int64     `db:"customerId"`
	ExpiredAt            time.Time `db:"expiredAt"`
	NotificationToken    string    `db:"notificationToken"`
	NotificationProvider int64     `db:"notificationProvider"`
}

type UpdateAccessSession struct {
	*AccessSession
	CurrentVersion int64 `db:"currentVersion"`
}
