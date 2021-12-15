package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type AddressStatement struct {
	FindByCustomerId        *sqlx.Stmt
	Insert                  *sqlx.NamedStmt
	Update                  *sqlx.NamedStmt
	FindPrimaryByCustomerId *sqlx.Stmt
}

func NewAddressStatement(db *nsql.DB) *AddressStatement {

	tableName := `Address`
	getColumns := `"xid", "metadata", "createdAt", "updatedAt", "modifiedBy", "version", "customerId", "purpose", "provinceId", "provinceName", "cityId", "cityName", "districtId", "districtName", "subDistrictId", "subDistrictName", "line", "postalCode", "isPrimary"`
	columns := `"xid", "metadata", "createdAt", "updatedAt", "modifiedBy", "version", "customerId", "purpose", "provinceId", "provinceName", "cityId", "cityName", "districtId", "districtName", "subDistrictId", "subDistrictName", "line", "postalCode", "isPrimary"`
	namedColumns := `:xid, :metadata, :createdAt, :updatedAt, :modifiedBy, :version, :customerId, :purpose, :provinceId, :provinceName, :cityId, :cityName, :districtId, :districtName, :subDistrictId, :subDistrictName, :line, :postalCode, :isPrimary`
	updatedNamedColumns := `"xid" = :xid, "metadata" = :metadata, "updatedAt" = :updatedAt, "modifiedBy" = :modifiedBy, "version" = :version, "customerId" = :customerId, "purpose" = :purpose, "provinceId" = :provinceId, "provinceName" = :provinceName, "cityId" = :cityId, "cityName" = :cityName, "districtId" = :districtId, "districtName" = :districtName, "subDistrictId" = :subDistrictId, "subDistrictName" = :subDistrictName, "line" = :line, "postalCode" = :postalCode, "isPrimary" = :isPrimary`

	return &AddressStatement{
		FindByCustomerId:        db.PrepareFmt(`SELECT %s FROM "%s" WHERE "customerId" = $1`, getColumns, tableName),
		Insert:                  db.PrepareNamedFmt(`INSERT INTO "%s" (%s) VALUES (%s)`, tableName, columns, namedColumns),
		Update:                  db.PrepareNamedFmt(`UPDATE "%s" SET %s WHERE id = :id`, tableName, updatedNamedColumns),
		FindPrimaryByCustomerId: db.PrepareFmt(`SELECT %s FROM "%s" WHERE "isPrimary" = true AND "customerId" = $1 LIMIT 1`, columns, tableName),
	}
}
