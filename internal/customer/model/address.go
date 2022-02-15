package model

import (
	"database/sql"
	"encoding/json"
)

type Address struct {
	Id              int64           `db:"id"`
	Xid             string          `db:"xid"`
	CustomerId      int64           `db:"customerId"`
	Purpose         int64           `db:"purpose"`
	ProvinceId      sql.NullInt64   `db:"provinceId"`
	ProvinceName    sql.NullString  `db:"provinceName"`
	CityId          sql.NullInt64   `db:"cityId"`
	CityName        sql.NullString  `db:"cityName"`
	DistrictId      sql.NullInt64   `db:"districtId"`
	DistrictName    sql.NullString  `db:"districtName"`
	SubDistrictId   sql.NullInt64   `db:"subDistrictId"`
	SubDistrictName sql.NullString  `db:"subDistrictName"`
	PostalCode      sql.NullString  `db:"postalCode"`
	Line            sql.NullString  `db:"line"`
	IsPrimary       sql.NullBool    `db:"isPrimary"`
	Metadata        json.RawMessage `db:"metadata"`
	ItemMetadata
}

type AddressExternal struct {
	CustomerId  string         `db:"user_AIID"`
	Alamat      sql.NullString `db:"alamat"`
	Kodepos     sql.NullString `db:"kodepos"`
	Kelurahan   sql.NullString `db:"kelurahan"`
	Kecamatan   sql.NullString `db:"kecamatan"`
	Kabupaten   sql.NullString `db:"kabupaten"`
	Provinsi    sql.NullString `db:"provinsi"`
	IdKelurahan sql.NullInt64  `db:"id_kelurahan"`
	IdKecamatan sql.NullInt64  `db:"id_kecamatan"`
	IdKabupaten sql.NullInt64  `db:"id_kabupaten"`
	IdProvinsi  sql.NullInt64  `db:"id_provinsi"`
}
