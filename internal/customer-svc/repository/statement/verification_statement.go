package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type VerificationStatement struct {
	FindByCustomerID *sqlx.Stmt
	Insert           *sqlx.NamedStmt
	Update           *sqlx.NamedStmt
	DeleteByID       *sqlx.Stmt
}

func NewVerificationStatement(db *nsql.DB) *VerificationStatement {
	tableName := "Verification"
	getColumns := `"xid","metadata","createdAt","updatedAt","modifiedBy","version","customerId","kycVerifiedStatus","kycVerifiedAt","emailVerificationToken","emailVerifiedStatus","emailVerifiedAt","dukcapilVerifiedStatus","dukcapilVerifiedAt","financialTransactionStatus","financialTransactionActivatedAt"`
	columns := `"xid","metadata","createdAt","updatedAt","modifiedBy","version","customerId","kycVerifiedStatus","kycVerifiedAt","emailVerificationToken","emailVerifiedStatus","emailVerifiedAt","dukcapilVerifiedStatus","dukcapilVerifiedAt","financialTransactionStatus","financialTransactionActivatedAt"`
	namedColumns := `:xid,:metadata,:createdAt,:updatedAt,:modifiedBy,:version,:customerId,:kycVerifiedStatus,:kycVerifiedAt,:emailVerificationToken,:emailVerifiedStatus,:emailVerifiedAt,:dukcapilVerifiedStatus,:dukcapilVerifiedAt,:financialTransactionStatus,:financialTransactionActivatedAt`
	updatedNamedColumns := `"xid" = :xid, "metadata" = :metadata, "updatedAt" = :updatedAt, "modifiedBy" = :modifiedBy, "version" = :version, "customerId" = :customerId, "kycVerifiedStatus" = :kycVerifiedStatus, "kycVerifiedAt" = :kycVerifiedAt, "emailVerificationToken" = :emailVerificationToken, "emailVerifiedStatus" = :emailVerifiedStatus, "emailVerifiedAt" = :emailVerifiedAt, "dukcapilVerifiedStatus" = :dukcapilVerifiedStatus, "dukcapilVerifiedAt" = :dukcapilVerifiedAt, "financialTransactionStatus" = :financialTransactionStatus, "financialTransactionActivatedAt" = :financialTransactionActivatedAt`

	return &VerificationStatement{
		FindByCustomerID: db.PrepareFmt(`SELECT %s FROM "%s" WHERE "customerId" = $1`, getColumns, tableName),
		Insert:           db.PrepareNamedFmt(`INSERT INTO "%s" (%s) VALUES (%s)`, tableName, columns, namedColumns),
		Update:           db.PrepareNamedFmt(`UPDATE "%s" SET %s WHERE "customerId" = :customerId`, tableName, updatedNamedColumns),
		DeleteByID:       db.PrepareFmt(`DELETE FROM "%s" WHERE "id" = $1`, tableName),
	}
}
