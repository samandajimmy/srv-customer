package model

import "time"

type VerificationOTP struct {
	Id             int64     `db:"id"`
	CreatedAt      time.Time `db:"createdAt"`
	Phone          string    `db:"phone"`
	RegistrationId string    `db:"registrationId"`
}
