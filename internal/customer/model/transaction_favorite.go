package model

import (
	"database/sql"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
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

func ToTransactionFavoriteDTO(m TransactionFavorite) dto.Favorite {
	return dto.Favorite{
		BaseField:       ToBaseFieldDTO(&m.BaseField),
		XID:             m.XID,
		Type:            m.Type,
		TypeTransaction: m.TypeTransaction,
		AccountName:     m.AccountName,
		AccountNumber:   m.AccountNumber,
		BankName:        m.BankName.String,
		BankCode:        m.BankCode.String,
		GroupMPO:        m.GroupMPO.String,
		ServiceCodeMPO:  m.ServiceCodeMPO.String,
	}
}

type ListTransactionFavoriteResult struct {
	Rows  []TransactionFavorite
	Count int64
}
