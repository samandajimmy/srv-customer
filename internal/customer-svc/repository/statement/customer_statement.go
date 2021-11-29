package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type CustomerStatement struct {
	Insert      *sqlx.NamedStmt
	FindByPhone *sqlx.Stmt
}

func NewCustomerStatement(db *nsql.DB) *CustomerStatement {
	tableName := "Customer"
	columns := "\"xid\", \"fullName\", \"phone\", \"email\", \"identityType\", \"identityNumber\", \"userRefId\", \"photos\", \"profile\", \"cif\", \"sid\", \"referralCode\", \"status\", \"metadata\", \"createdAt\", \"updatedAt\", \"modifiedBy\", \"version\""
	namedColumns := ":xid,:fullName,:phone,:email,:identityType,:identityNumber,:userRefId,:photos,:profile,:cif,:sid,:referralCode,:status,:metadata,:createdAt,:updatedAt,:modifiedBy,:version"
	return &CustomerStatement{
		Insert:      db.PrepareNamedFmt("INSERT INTO \"%s\"(%s) VALUES (%s) RETURNING id", tableName, columns, namedColumns),
		FindByPhone: db.PrepareFmt("SELECT %s FROM \"%s\" WHERE phone = $1", columns, tableName),
	}
}
