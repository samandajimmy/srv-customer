package model

import (
	"database/sql/driver"
	"encoding/json"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type BankAccount struct {
	BaseField
	ID            int64  `db:"id"`
	XID           string `db:"xid"`
	CustomerID    int64  `db:"customerId"`
	Bank          *Bank  `db:"bank"`
	AccountName   string `db:"accountName"`
	AccountNumber string `db:"accountNumber"`
	Status        int8   `db:"status"`
}

type BankAccountStatus int64

type ListBankAccountResult struct {
	Rows  []BankAccount
	Count int64
}

type Bank struct {
	ID        int32  `json:"id"`
	Name      string `json:"name"`
	Title     string `json:"title"`
	Code      string `json:"code"`
	Thumbnail string `json:"thumbnail"`
}

func (m *Bank) Scan(src interface{}) error {
	return nsql.ScanJSON(src, m)
}

func (m *Bank) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func ToBankAccountDTO(m BankAccount) dto.BankAccount {
	return dto.BankAccount{
		XID:           m.XID,
		CustomerID:    m.CustomerID,
		AccountNumber: m.AccountNumber,
		AccountName:   m.AccountName,
		Bank:          ToBankDTO(m.Bank),
		BaseField:     ToBaseFieldDTO(&m.BaseField),
	}
}

func ToBankDTO(m *Bank) *dto.Bank {
	return &dto.Bank{
		ID:        m.ID,
		Name:      m.Name,
		Title:     m.Title,
		Code:      m.Code,
		Thumbnail: m.Thumbnail,
	}
}

func ToBank(d *dto.Bank) *Bank {
	return &Bank{
		ID:        d.ID,
		Name:      d.Name,
		Title:     d.Title,
		Code:      d.Code,
		Thumbnail: d.Thumbnail,
	}
}
