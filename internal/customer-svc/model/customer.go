package model

import (
	"encoding/json"
	"time"
)

type Customer struct {
	Id             int64           `db:"id"`
	CustomerXID    string          `db:"xid"`
	FullName       string          `db:"fullName"`
	Phone          string          `db:"phone"`
	Email          string          `db:"email"`
	IdentityType   int64           `db:"identityType"`
	IdentityNumber string          `db:"identityNumber"`
	UserRefId      int64           `db:"userRefId"`
	Photos         json.RawMessage `db:"photos"`
	Profile        json.RawMessage `db:"profile"`
	Cif            string          `db:"cif"`
	Sid            string          `db:"sid"`
	ReferralCode   string          `db:"referralCode"`
	Status         int64           `db:"status"`
	Metadata       json.RawMessage `db:"metadata"`
	ItemMetadata
}

type CustomerAuthentication struct {
	CustomerId         int64           `db:"customerId"`
	FullName           string          `db:"fullName"`
	Phone              string          `db:"phone"`
	Email              string          `db:"email"`
	BlockedAt          time.Time       `db:"blockedAt"`
	BlockedUntilAt     time.Time       `db:"blockedUntilAt"`
	Metadata           json.RawMessage `db:"metadata"`
	WrongPasswordCount int64           `db:"wrongPasswordCount"`
	Password           string          `db:"password"`
}
