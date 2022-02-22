package model

import "time"

type VerificationOTP struct {
	ID             int64     `db:"id"`
	CreatedAt      time.Time `db:"createdAt"`
	Phone          string    `db:"phone"`
	RegistrationID string    `db:"registrationId"`
}
