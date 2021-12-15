package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type FinancialDataStatement struct {
	FindByCustomerID *sqlx.Stmt
	Insert           *sqlx.NamedStmt
	Update           *sqlx.NamedStmt
	DeleteByID       *sqlx.Stmt
}

func NewFinancialDataStatement(db *nsql.DB) *FinancialDataStatement {

	tableName := `FinancialData`
	getColumns := `"xid", "metadata", "createdAt", "updatedAt", "modifiedBy", "version", "customerId", "mainAccountNumber", "accountNumber", "goldSavingStatus", "goldCardApplicationNumber", "goldCardAccountNumber", "balance"`
	columns := `"xid", "metadata", "createdAt", "updatedAt", "modifiedBy", "version", "customerId", "mainAccountNumber", "accountNumber", "goldSavingStatus", "goldCardApplicationNumber", "goldCardAccountNumber", "balance" `
	namedColumns := `:xid, :metadata, :createdAt, :updatedAt, :modifiedBy, :version, :customerId, :mainAccountNumber, :accountNumber, :goldSavingStatus, :goldCardApplicationNumber, :goldCardAccountNumber, :balance`
	updatedNamedColumns := `"xid" = :xid, "metadata" = :metadata, "updatedAt" = :updatedAt, "modifiedBy" = :modifiedBy, "version" = :version, "mainAccountNumber" = :mainAccountNumber, "accountNumber" = :accountNumber, "goldSavingStatus" = :goldSavingStatus, "goldCardApplicationNumber" = :goldCardApplicationNumber, "goldCardAccountNumber" = :goldCardAccountNumber, "balance" = :balance`

	return &FinancialDataStatement{
		FindByCustomerID: db.PrepareFmt(`SELECT %s FROM "%s" WHERE "customerId" = $1`, getColumns, tableName),
		Insert:           db.PrepareNamedFmt(`INSERT INTO "%s" (%s) VALUES (%s)`, tableName, columns, namedColumns),
		Update:           db.PrepareNamedFmt(`UPDATE "%s" SET %s WHERE "customerId" = :customerId`, tableName, updatedNamedColumns),
		DeleteByID:       db.PrepareFmt(`DELETE FROM "%s" WHERE "id" = $1`, tableName),
	}
}
