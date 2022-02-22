package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type Customer struct {
	BaseField
	ID             int64            `db:"id"`
	CustomerXID    string           `db:"xid"`
	FullName       string           `db:"fullName"`
	Phone          string           `db:"phone"`
	Email          string           `db:"email"`
	IdentityType   int64            `db:"identityType"`
	IdentityNumber string           `db:"identityNumber"`
	UserRefID      sql.NullString   `db:"userRefId"`
	Photos         *CustomerPhoto   `db:"photos"`
	Profile        *CustomerProfile `db:"profile"`
	Cif            string           `db:"cif"`
	Sid            string           `db:"sid"`
	ReferralCode   string           `db:"referralCode"`
	Status         int64            `db:"status"`
	Metadata       json.RawMessage  `db:"metadata"`
}

type CustomerPhoto struct {
	Xid      string `json:"xid"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	Mimetype string `json:"mime_type"`
}

func (m *CustomerPhoto) Scan(src interface{}) error {
	return nsql.ScanJSON(src, m)
}

func (m *CustomerPhoto) Value() (driver.Value, error) {
	return json.Marshal(m)
}

type ValidatePassword struct {
	CustomerID int64  `db:"customerId"`
	Password   string `db:"password"`
}

type UpdateCustomer struct {
	*Customer
	CurrentVersion int64 `db:"currentVersion"`
}

type UpdateByCIF struct {
	*Customer
	Cif string `db:"cif"`
}

type UpdateByID struct {
	*Customer
	ID int64 `db:"id"`
}

type UpdateCustomerByUserRefID struct {
	*Customer
	UserRefID string `db:"userRefId"`
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
	IdentityExpiredAt  string `json:"IdentityExpiredAt"`
	Religion           string `json:"religion"`
	MarriageStatus     string `json:"marriageStatus"`
	NPWPNumber         string `json:"npwpNumber"`
	NPWPPhotoFile      string `json:"npwpPhotoFile"`
	NPWPUpdatedAt      int64  `json:"npwpUpdatedAt,string,omitempty"`
	ProfileUpdatedAt   int64  `json:"profileUpdatedAt,string,omitempty"`
	CifLinkUpdatedAt   int64  `json:"cifLinkUpdatedAt,string,omitempty"`
	CifUnlinkUpdatedAt int64  `json:"cifUnlinkUpdatedAt,string,omitempty"`
	SidPhotoFile       string `json:"sidPhotoFile"`
}

func (m *CustomerProfile) Scan(src interface{}) error {
	return nsql.ScanJSON(src, m)
}

func (m *CustomerProfile) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func ToCustomerProfile(d *dto.CustomerProfileVO) *CustomerProfile {
	if d == nil {
		return nil
	}

	return &CustomerProfile{
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

var EmptyCustomerProfile = &CustomerProfile{
	MaidenName:         "",
	Gender:             "",
	Nationality:        "",
	DateOfBirth:        "",
	PlaceOfBirth:       "",
	IdentityPhotoFile:  "",
	IdentityExpiredAt:  "",
	Religion:           "",
	MarriageStatus:     "",
	NPWPNumber:         "",
	NPWPPhotoFile:      "",
	NPWPUpdatedAt:      0,
	ProfileUpdatedAt:   0,
	CifLinkUpdatedAt:   0,
	CifUnlinkUpdatedAt: 0,
	SidPhotoFile:       "",
}
