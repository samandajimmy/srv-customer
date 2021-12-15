package model

import (
	"database/sql"
	"encoding/json"
)

type Address struct {
	Id              int64           `db:"id"`
	Xid             string          `db:"xid"`
	CustomerId      int64           `db:"customerId"`
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

type AddressExternal struct {
	CustomerId  string         `db:"user_AIID"`
	Alamat      sql.NullString `db:"alamat"`
	Kodepos     sql.NullString `db:"kodepos"`
	Kelurahan   sql.NullString `db:"kelurahan"`
	Kecamatan   sql.NullString `db:"kecamatan"`
	Kabupaten   sql.NullString `db:"kabupaten"`
	Provinsi    sql.NullString `db:"provinsi"`
	IdKelurahan sql.NullString `db:"id_kelurahan"`
	IdKecamatan sql.NullString `db:"id_kecamatan"`
	IdKabupaten sql.NullString `db:"id_kabupaten"`
	IdProvinsi  sql.NullString `db:"id_provinsi"`
}
