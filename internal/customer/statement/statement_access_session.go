package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

var AccessSessionSchema = nsql.NewSchema(
	"AccessSession",
	append(CommonColumns, "xid", "customerId", "expiredAt", "notificationToken", "notificationProvider"),
	nsql.WithAlias("as"),
)

type AccessSession struct {
	Insert *sqlx.NamedStmt
	Update *sqlx.NamedStmt
}

func NewAccessSession(db *nsql.DatabaseContext) *AccessSession {

	return &AccessSession{
		Insert: db.PrepareNamedFmt(`INSERT INTO "%s" (%s) VALUES (%s)`,
			AccessSessionSchema.TableName, AccessSessionSchema.InsertColumns(), AccessSessionSchema.InsertNamedColumns()),
		Update: db.PrepareNamedFmt(`UPDATE "%s" SET %s WHERE id = :id AND "version" = :currentVersion`, AccessSessionSchema.TableName, AccessSessionSchema.UpdateNamedColumns()),
	}
}
