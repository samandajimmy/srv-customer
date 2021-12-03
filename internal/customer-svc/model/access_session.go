package model

import "time"

type AccessSession struct {
	CustomerId           int64     `db:"customerId"`
	ExpiredAt            time.Time `db:"expiredAt"`
	NotificationToken    string    `db:"notificationToken"`
	NotificationProvider int64     `db:"notificationProvider"`
	ItemMetadata
}
