package model

import "time"

type OTP struct {
	Id         int64     `db:"id"`
	CustomerId int64     `db:"customerId"`
	Content    string    `db:"content"`
	Type       string    `db:"type"`
	Data       string    `db:"data"`
	Status     string    `db:"status"`
	UpdatedAt  time.Time `db:"updatedAt"`
}
