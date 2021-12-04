package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type CredentialStatement struct {
	Insert *sqlx.NamedStmt
}

func NewCredentialStatement(db *nsql.DB) *CredentialStatement {
	tableName := "Credential"
	columns := `id, xid, metadata, createdAt, updatedAt, modifiedBy, version, customerId, password, nextPasswordResetAt, pin, pinCif, pinUpdatedAt, pinLastAccessAt, pinCounter, pinBlockedStatus, isLocked, loginFailCount, wrongPasswordCount, blockedAt, blockedUntilAt, biometricLogin, biometricDeviceId`
	namedColumns := `:id,:xid,:metadata,:createdAt,:updatedAt,:modifiedBy,:version,:customerId,:password,:nextPasswordResetAt,:pin,:pinCif,:pinUpdatedAt,:pinLastAccessAt,:pinCounter,:pinBlockedStatus,:isLocked,:loginFailCount,:wrongPasswordCount,:blockedAt,:blockedUntilAt,:biometricLogin,:biometricDeviceId`

	return &CredentialStatement{
		Insert: db.PrepareNamedFmt("INSERT INTO %s(%s) VALUES (%s)", tableName, columns, namedColumns),
	}
}
