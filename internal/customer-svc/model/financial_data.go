package model

type FinancialData struct {
	Xid                       string `db:"xid"`
	CustomerId                int64  `db:"customerId"`
	MainAccountNumber         string `db:"mainAccountNumber"`
	AccountNumber             string `db:"accountNumber"`
	GoldSavingStatus          int64  `db:"goldSavingStatus"`
	GoldCardApplicationNumber string `db:"goldCardApplicationNumber"`
	GoldCardAccountNumber     string `db:"goldCardAccountNumber"`
	Balance                   int64  `db:"balance"`
	ItemMetadata
}
