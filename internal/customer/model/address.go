package model

import (
	"database/sql"
)

type Address struct {
	BaseField
	ID              int64          `db:"id"`
	Xid             string         `db:"xid"`
	CustomerID      int64          `db:"customerId"`
	Purpose         int64          `db:"purpose"`
	ProvinceID      sql.NullInt64  `db:"provinceId"`
	ProvinceName    sql.NullString `db:"provinceName"`
	CityID          sql.NullInt64  `db:"cityId"`
	CityName        sql.NullString `db:"cityName"`
	DistrictID      sql.NullInt64  `db:"districtId"`
	DistrictName    sql.NullString `db:"districtName"`
	SubDistrictID   sql.NullInt64  `db:"subDistrictId"`
	SubDistrictName sql.NullString `db:"subDistrictName"`
	PostalCode      sql.NullString `db:"postalCode"`
	Line            sql.NullString `db:"line"`
	IsPrimary       sql.NullBool   `db:"isPrimary"`
}

type AddressExternal struct {
	CustomerID  string         `db:"user_AIID"`
	Alamat      sql.NullString `db:"alamat"`
	Kodepos     sql.NullString `db:"kodepos"`
	Kelurahan   sql.NullString `db:"kelurahan"`
	Kecamatan   sql.NullString `db:"kecamatan"`
	Kabupaten   sql.NullString `db:"kabupaten"`
	Provinsi    sql.NullString `db:"provinsi"`
	IDKelurahan sql.NullInt64  `db:"id_kelurahan"`
	IDKecamatan sql.NullInt64  `db:"id_kecamatan"`
	IDKabupaten sql.NullInt64  `db:"id_kabupaten"`
	IDProvinsi  sql.NullInt64  `db:"id_provinsi"`
}
