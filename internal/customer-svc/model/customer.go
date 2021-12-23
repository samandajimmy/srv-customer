package model

import (
	"encoding/json"
)

type Customer struct {
	Id             int64           `db:"id,omitempty"`
	CustomerXID    string          `db:"xid"`
	FullName       string          `db:"fullName"`
	Phone          string          `db:"phone"`
	Email          string          `db:"email"`
	IdentityType   int64           `db:"identityType"`
	IdentityNumber string          `db:"identityNumber"`
	UserRefId      string          `db:"userRefId"`
	Photos         json.RawMessage `db:"photos"`
	Profile        json.RawMessage `db:"profile"`
	Cif            string          `db:"cif"`
	Sid            string          `db:"sid"`
	ReferralCode   string          `db:"referralCode"`
	Status         int64           `db:"status"`
	Metadata       json.RawMessage `db:"metadata"`
	ItemMetadata
}
