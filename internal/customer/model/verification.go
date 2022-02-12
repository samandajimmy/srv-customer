package model

import (
	"database/sql"
	"encoding/json"
)

type Verification struct {
	Id                              int64           `db:"id"`
	Xid                             string          `db:"xid"`
	CustomerId                      int64           `db:"customerId"`
	KycVerifiedStatus               int64           `db:"kycVerifiedStatus"`
	KycVerifiedAt                   sql.NullTime    `db:"kycVerifiedAt"`
	EmailVerificationToken          string          `db:"emailVerificationToken"`
	EmailVerifiedStatus             int64           `db:"emailVerifiedStatus"`
	EmailVerifiedAt                 sql.NullTime    `db:"emailVerifiedAt"`
	DukcapilVerifiedStatus          int64           `db:"dukcapilVerifiedStatus"`
	DukcapilVerifiedAt              sql.NullTime    `db:"dukcapilVerifiedAt"`
	FinancialTransactionStatus      int64           `db:"financialTransactionStatus"`
	FinancialTransactionActivatedAt sql.NullTime    `db:"financialTransactionActivatedAt"`
	Metadata                        json.RawMessage `db:"metadata"`
	ItemMetadata
}
