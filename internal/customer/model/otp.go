package model

import "time"

type OTP struct {
	ID         int64     `db:"id"`
	CustomerID int64     `db:"customerId"`
	Content    string    `db:"content"`
	Type       string    `db:"type"`
	Data       string    `db:"data"`
	Status     string    `db:"status"`
	UpdatedAt  time.Time `db:"updatedAt"`
}
