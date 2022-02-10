package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type AuditLogin struct {
	Insert             *sqlx.NamedStmt
	UpdateByCustomerId *sqlx.NamedStmt
	CountLogin         *sqlx.Stmt
}

func NewAuditLogin(db *nsql.DatabaseContext) *AuditLogin {
	tableName := "AuditLogin"
	columns := `"customerId", "channelId", "deviceId", ip, latitude, longitude, timestamp, timezone, "brand", "osVersion", browser, "useBiometric","metadata","createdAt","updatedAt","modifiedBy","version"`
	namedColumns := `:customerId, :channelId, :deviceId, :ip, :latitude, :longitude, :timestamp, :timezone, :brand, :osVersion, :browser, :useBiometric,:metadata,:createdAt,:updatedAt,:modifiedBy,:version`
	updateColumns := `"customerId" = :customerId, "channelId" = :channelId, "deviceId" = :deviceId, "ip" = :id, "latitude" = :latitude, "longitude" = :longitude, "timestamp" = :timestamp, "timezone" = :timezone, "brand" = :brand, "osVersion" = :osVersion, "browser" = :browser, "useBiometric" = :useBiometric`

	return &AuditLogin{
		Insert:             db.PrepareNamedFmt(`INSERT INTO "%s" (%s) VALUES (%s)`, tableName, columns, namedColumns),
		UpdateByCustomerId: db.PrepareNamedFmt(`UPDATE "%s" SET %s WHERE "customerId" = :customerId`, tableName, updateColumns),
		CountLogin:         db.PrepareFmt(`SELECT COUNT(id) as count FROM "%s" WHERE "customerId" = $1`, tableName),
	}
}
