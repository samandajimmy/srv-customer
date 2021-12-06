package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type AccessSessionStatement struct {
	Insert *sqlx.NamedStmt
	Update *sqlx.NamedStmt
}

func NewAccessSessionStatement(db *nsql.DB) *AccessSessionStatement {

	tableName := "AccessSession"
	columns := `"xid", "metadata", "createdAt", "updatedAt", "modifiedBy", "version", "customerId", "expiredAt", "notificationToken", "notificationProvider"`
	namedColumns := `:xid,:metadata,:createdAt,:updatedAt,:modifiedBy,:version,:customerId,:expiredAt,:notificationToken,:notificationProvider`
	updatedNamedColumns := `"xid" = :xid, "metadata" = :metadata, "updatedAt" = :updatedAt, "modifiedBy" = :modifiedBy, "version" = :version, "customerId" = :customerId, "expiredAt" = :expiredAt, "notificationToken" = :notificationToken, "notificationProvider" = :notificationProvider`
	return &AccessSessionStatement{
		Insert: db.PrepareNamedFmt(`INSERT INTO "%s" (%s) VALUES (%s)`, tableName, columns, namedColumns),
		Update: db.PrepareNamedFmt(`UPDATE "%s" SET %s WHERE id = :id`, tableName, updatedNamedColumns),
	}
}
