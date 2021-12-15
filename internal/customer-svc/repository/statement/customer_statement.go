package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type CustomerStatement struct {
	Insert             *sqlx.NamedStmt
	UpdateByPhone      *sqlx.NamedStmt
	FindById           *sqlx.Stmt
	FindByPhone        *sqlx.Stmt
	FindByEmail        *sqlx.Stmt
	FindByEmailOrPhone *sqlx.Stmt
}

func NewCustomerStatement(db *nsql.DB) *CustomerStatement {
	tableName := `Customer`
	columns := `"xid","metadata","createdAt","updatedAt","modifiedBy","version","fullName","phone","email","identityType","identityNumber","userRefId","photos","profile","cif","sid","referralCode","status"`
	namedColumns := `:xid,:metadata,:createdAt,:updatedAt,:modifiedBy,:version,:fullName,:phone,:email,:identityType,:identityNumber,:userRefId,:photos,:profile,:cif,:sid,:referralCode,:status`
	updateColumns := `"xid" = :xid, "metadata" = :metadata, "createdAt" = :createdAt, "updatedAt" = :updatedAt, "modifiedBy" = :modifiedBy, "version" = :version, "fullName" = :fullName, "phone" = :phone, "email" = :email, "identityType" = :identityType, "identityNumber" = :identityNumber, "userRefId" = :userRefId, "photos" = :photos, "profile" = :profile, "cif" = :cif, "sid" = :sid, "referralCode" = :referralCode, "status" = :status`

	columnsWithId := `"id",` + columns

	return &CustomerStatement{
		Insert:      db.PrepareNamedFmt(`INSERT INTO "%s" (%s) VALUES (%s) RETURNING id`, tableName, columns, namedColumns),
		FindById:    db.PrepareFmt(`SELECT "id", %s FROM "%s" WHERE "id" = $1`, columns, tableName),
		FindByPhone: db.PrepareFmt(`SELECT "id", %s FROM "%s" WHERE "phone" = $1`, columns, tableName),
		FindByEmail: db.PrepareFmt(`SELECT "id", %s FROM "%s" WHERE "email" = $1`, columns, tableName),
		FindByEmailOrPhone: db.PrepareFmt(
			`SELECT %s FROM "%s" WHERE (phone = $1) OR (email = $1)`, columnsWithId, tableName,
		),
		UpdateByPhone: db.PrepareNamedFmt(`UPDATE "%s" SET %s WHERE "phone" = :phone`, tableName, updateColumns),
	}
}
