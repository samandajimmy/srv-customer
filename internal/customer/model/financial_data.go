package model

import (
	"github.com/rs/xid"
	"strings"
)

var EmptyFinancialData = &FinancialData{
	Xid:                       strings.ToUpper(xid.New().String()),
	MainAccountNumber:         "",
	AccountNumber:             "",
	GoldSavingStatus:          0,
	GoldCardApplicationNumber: "",
	GoldCardAccountNumber:     "",
	Balance:                   0,
	BaseField:                 EmptyBaseField,
}

type FinancialData struct {
	BaseField
	ID                        int64  `db:"id"`
	Xid                       string `db:"xid"`
	CustomerID                int64  `db:"customerId"`
	MainAccountNumber         string `db:"mainAccountNumber"`
	AccountNumber             string `db:"accountNumber"`
	GoldSavingStatus          int64  `db:"goldSavingStatus"`
	GoldCardApplicationNumber string `db:"goldCardApplicationNumber"`
	GoldCardAccountNumber     string `db:"goldCardAccountNumber"`
	Balance                   int64  `db:"balance"`
}
