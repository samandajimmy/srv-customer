package model

import (
	"database/sql/driver"
	"encoding/json"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type Customer struct {
	Id             int64           `db:"id"`
	CustomerXID    string          `db:"xid"`
	FullName       string          `db:"fullName"`
	Phone          string          `db:"phone"`
	Email          string          `db:"email"`
	IdentityType   int64           `db:"identityType"`
	IdentityNumber string          `db:"identityNumber"`
	UserRefId      string          `db:"userRefId"`
	Photos         json.RawMessage `db:"photos"`
	Profile        CustomerProfile `db:"profile"`
	Cif            string          `db:"cif"`
	Sid            string          `db:"sid"`
	ReferralCode   string          `db:"referralCode"`
	Status         int64           `db:"status"`
	Metadata       json.RawMessage `db:"metadata"`
	ItemMetadata
}

type UpdateCustomer struct {
	*Customer
	CurrentVersion int64 `db:"currentVersion"`
}

type UpdateByCIF struct {
	*Customer
	Cif string `db:"cif"`
}

type CustomerDetail struct {
	Customer
	Address
}

type CustomerProfile struct {
	MaidenName         string `json:"maidenName"`
	Gender             string `json:"gender"`
	Nationality        string `json:"nationality"`
	DateOfBirth        string `json:"dateOfBirth"`
	PlaceOfBirth       string `json:"placeOfBirth"`
	IdentityPhotoFile  string `json:"identityPhotoFile"`
	MarriageStatus     string `json:"marriageStatus"`
	NPWPNumber         string `json:"npwpNumber"`
	NPWPPhotoFile      string `json:"npwpPhotoFile"`
	NPWPUpdatedAt      string `json:"npwpUpdatedAt"`
	ProfileUpdatedAt   string `json:"profileUpdatedAt"`
	CifLinkUpdatedAt   string `json:"cifLinkUpdatedAt"`
	CifUnlinkUpdatedAt string `json:"cifUnlinkUpdatedAt"`
	SidPhotoFile       string `json:"sidPhotoFile"`
}

func (m *CustomerProfile) Scan(src interface{}) error {
	return nsql.ScanJSON(src, m)
}

func (m *CustomerProfile) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func ToCustomerProfile(d dto.CustomerProfileVO) CustomerProfile {
	return CustomerProfile{
		MaidenName:         d.MaidenName,
		Gender:             d.Gender,
		Nationality:        d.Nationality,
		DateOfBirth:        d.DateOfBirth,
		PlaceOfBirth:       d.PlaceOfBirth,
		IdentityPhotoFile:  d.IdentityPhotoFile,
		MarriageStatus:     d.MarriageStatus,
		NPWPNumber:         d.NPWPNumber,
		NPWPPhotoFile:      d.NPWPPhotoFile,
		NPWPUpdatedAt:      d.NPWPUpdatedAt,
		ProfileUpdatedAt:   d.ProfileUpdatedAt,
		CifLinkUpdatedAt:   d.CifLinkUpdatedAt,
		CifUnlinkUpdatedAt: d.CifUnlinkUpdatedAt,
		SidPhotoFile:       d.SidPhotoFile,
	}
}

func ToDTOCustomerProfile(m CustomerProfile) dto.CustomerProfileVO {
	return dto.CustomerProfileVO{
		MaidenName:         m.MaidenName,
		Gender:             m.Gender,
		Nationality:        m.Nationality,
		DateOfBirth:        m.DateOfBirth,
		PlaceOfBirth:       m.PlaceOfBirth,
		IdentityPhotoFile:  m.IdentityPhotoFile,
		MarriageStatus:     m.MarriageStatus,
		NPWPNumber:         m.NPWPNumber,
		NPWPPhotoFile:      m.NPWPPhotoFile,
		NPWPUpdatedAt:      m.NPWPUpdatedAt,
		ProfileUpdatedAt:   m.ProfileUpdatedAt,
		CifLinkUpdatedAt:   m.CifLinkUpdatedAt,
		CifUnlinkUpdatedAt: m.CifUnlinkUpdatedAt,
		SidPhotoFile:       m.SidPhotoFile,
	}
}
