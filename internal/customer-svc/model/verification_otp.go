package model

import "time"

type VerificationOTP struct {
	CreatedAt      time.Time `db:"createdAt"`
	Phone          string    `db:"phone"`
	RegistrationId string    `db:"registrationId"`
}
