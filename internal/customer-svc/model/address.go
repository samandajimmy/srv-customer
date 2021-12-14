package model

import (
	"database/sql"
	"encoding/json"
)

type Address struct {
	Id              int64           `db:"id"`
	Xid             string          `db:"xid"`
	CustomerId      string          `db:"customerId"`
	Purpose         string          `db:"purpose"`
	ProvinceId      sql.NullString  `db:"provinceId"`
	ProvinceName    sql.NullString  `db:"provinceName"`
	CityId          sql.NullString  `db:"cityId"`
	CityName        sql.NullString  `db:"cityName"`
	DistrictId      sql.NullString  `db:"districtId"`
	DistrictName    sql.NullString  `db:"districtName"`
	SubDistrictId   sql.NullString  `db:"subDistrictId"`
	SubDistrictName sql.NullString  `db:"subDistrictName"`
	Line            sql.NullString  `db:"line"`
	IsPrimary       sql.NullBool    `db:"isPrimary"`
	Metadata        json.RawMessage `db:"metadata"`
	ItemMetadata
}
