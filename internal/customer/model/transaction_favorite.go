package model

import (
	"database/sql"
)

type TransactionFavorite struct {
	BaseField
	ID              int64          `db:"id"`
	XID             string         `db:"xid"`
	CustomerID      int64          `db:"customerId"`
	Type            string         `db:"type"`
	TypeTransaction string         `db:"typeTransaction"`
	AccountNumber   string         `db:"accountNumber"`
	AccountName     string         `db:"accountName"`
	BankName        sql.NullString `db:"bankName"`
	BankCode        sql.NullString `db:"bankCode"`
	GroupMPO        sql.NullString `db:"groupMpo"`
	ServiceCodeMPO  sql.NullString `db:"serviceCodeMpo"`
}
