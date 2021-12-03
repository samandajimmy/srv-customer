package model

import "time"

type Verification struct {
	CustomerXID                     string    `db:"xid"`
	CustomerId                      int64     `db:"customerId"`
	KycVerifiedStatus               int64     `db:"kycVerifiedStatus"`
	KycVerifiedAt                   time.Time `db:"kycVerifiedAt"`
	EmailVerifiedStatus             int64     `db:"emailVerifiedStatus"`
	EmailVerifiedAt                 time.Time `db:"emailVerifiedAt"`
	DukcapilVerifiedStatus          int64     `db:"dukcapilVerifiedStatus"`
	DukcapilVerifiedAt              time.Time `db:"dukcapilVerifiedAt"`
	FinancialTransactionStatus      int64     `db:"financialTransactionStatus"`
	FinancialTransactionActivatedAt time.Time `db:"financialTransactionActivatedAt"`
	ItemMetadata
}
