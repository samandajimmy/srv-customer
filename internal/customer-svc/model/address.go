package model

import "encoding/json"

type Address struct {
	Id              int64           `db:"id"`
	Xid             string          `db:"xid"`
	CustomerId      string          `db:"customerId"`
	Purpose         string          `db:"purpose"`
	ProvinceId      string          `db:"provinceId"`
	ProvinceName    string          `db:"provinceName"`
	CityId          string          `db:"cityId"`
	CityName        string          `db:"cityName"`
	DistrictId      string          `db:"districtId"`
	DistrictName    string          `db:"districtName"`
	SubDistrictId   string          `db:"subDistrictId"`
	SubDistrictName string          `db:"subDistrictName"`
	Metadata        json.RawMessage `db:"metadata"`
	ItemMetadata
}
