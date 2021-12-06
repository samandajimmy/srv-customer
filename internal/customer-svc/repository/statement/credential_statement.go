package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type CredentialStatement struct {
	FindByCustomerID *sqlx.Stmt
	Insert           *sqlx.NamedStmt
	Update           *sqlx.NamedStmt
	DeleteByID       *sqlx.Stmt
}

func NewCredentialStatement(db *nsql.DB) *CredentialStatement {
	tableName := "Credential"
	getColumns := `"xid","metadata","createdAt","updatedAt","modifiedBy","version","customerId","password","nextPasswordResetAt","pin","pinCif","pinUpdatedAt","pinLastAccessAt","pinCounter","pinBlockedStatus","isLocked","loginFailCount","wrongPasswordCount","blockedAt","blockedUntilAt","biometricLogin","biometricDeviceId"`
	columns := `"xid", "metadata", "createdAt", "updatedAt", "modifiedBy", "version", "customerId", "password", "nextPasswordResetAt", "pin", "pinCif", "pinUpdatedAt", "pinLastAccessAt", "pinCounter", "pinBlockedStatus", "isLocked", "loginFailCount", "wrongPasswordCount", "blockedAt", "blockedUntilAt", "biometricLogin", "biometricDeviceId"`
	namedColumns := ":xid,:metadata,:createdAt,:updatedAt,:modifiedBy,:version,:customerId,:password,:nextPasswordResetAt,:pin,:pinCif,:pinUpdatedAt,:pinLastAccessAt,:pinCounter,:pinBlockedStatus,:isLocked,:loginFailCount,:wrongPasswordCount,:blockedAt,:blockedUntilAt,:biometricLogin,:biometricDeviceId"
	updatedNamedColumns := `"xid" = :xid, "metadata" = :metadata, "updatedAt" = :updatedAt, "modifiedBy" = :modifiedBy, "version" = :version, "customerId" = :customerId, "password" = :password, "nextPasswordResetAt" = :nextPasswordResetAt, "pin" = :pin, "pinCif" = :pinCif, "pinUpdatedAt" = :pinUpdatedAt, "pinLastAccessAt" = :pinLastAccessAt, "pinCounter" = :pinCounter, "pinBlockedStatus" = :pinBlockedStatus, "isLocked" = :isLocked, "loginFailCount" = :loginFailCount, "wrongPasswordCount" = :wrongPasswordCount, "blockedAt" = :blockedAt, "blockedUntilAt" = :blockedUntilAt, "biometricLogin" = :biometricLogin, "biometricDeviceId" = :biometricDeviceId`

	return &CredentialStatement{
		FindByCustomerID: db.PrepareFmt(`SELECT %s FROM "%s" WHERE "customerId" = $1`, getColumns, tableName),
		Insert:           db.PrepareNamedFmt(`INSERT INTO "%s" (%s) VALUES (%s)`, tableName, columns, namedColumns),
		Update:           db.PrepareNamedFmt(`UPDATE "%s" SET %s WHERE "customerId" = :customerId`, tableName, updatedNamedColumns),
		DeleteByID:       db.PrepareFmt(`DELETE FROM "%s" WHERE "id" = $1`, tableName),
	}
}
