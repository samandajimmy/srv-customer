package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"github.com/nbs-go/errx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type Customer struct {
	BaseField
	ID             int64                  `db:"id"`
	CustomerXID    string                 `db:"xid"`
	FullName       string                 `db:"fullName"`
	Phone          string                 `db:"phone"`
	Email          string                 `db:"email"`
	IdentityType   int64                  `db:"identityType"`
	IdentityNumber string                 `db:"identityNumber"`
	UserRefID      sql.NullString         `db:"userRefId"`
	Photos         *CustomerPhoto         `db:"photos"`
	Profile        *CustomerProfile       `db:"profile"`
	Cif            string                 `db:"cif"`
	Sid            string                 `db:"sid"`
	ReferralCode   sql.NullString         `db:"referralCode"`
	Status         constant.ControlStatus `db:"status"`
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

type PostSynchronizeCustomerModel struct {
	Customer     *Customer      `json:"customer"`
	Credential   *Credential    `json:"credential"`
	Financial    *FinancialData `json:"financial"`
	Verification *Verification  `json:"verification"`
	Address      *Address       `json:"address"`
}

func ToCustomerSyncVO(customer *Customer) *dto.CustomerSyncVO {

	if customer == nil {
		return nil
	}

	profile := customer.Profile

	return &dto.CustomerSyncVO{
		Photos: dto.PhotosVO{
			FileName: customer.Photos.FileName,
			FileSize: customer.Photos.FileSize,
			MimeType: customer.Photos.Mimetype,
		},
		Profile: dto.CustomerProfileVO{
			MaidenName:         profile.MaidenName,
			Gender:             profile.Gender,
			Nationality:        profile.Nationality,
			DateOfBirth:        profile.DateOfBirth,
			PlaceOfBirth:       profile.PlaceOfBirth,
			IdentityPhotoFile:  profile.IdentityPhotoFile,
			MarriageStatus:     profile.MarriageStatus,
			NPWPNumber:         profile.NPWPNumber,
			NPWPPhotoFile:      profile.NPWPPhotoFile,
			NPWPUpdatedAt:      profile.NPWPUpdatedAt,
			ProfileUpdatedAt:   profile.ProfileUpdatedAt,
			CifLinkUpdatedAt:   profile.CifLinkUpdatedAt,
			CifUnlinkUpdatedAt: profile.CifUnlinkUpdatedAt,
			SidPhotoFile:       profile.SidPhotoFile,
			Religion:           profile.Religion,
		},
		FullName:       customer.FullName,
		Phone:          customer.Phone,
		Email:          customer.Email,
		IdentityType:   customer.IdentityType,
		IdentityNumber: customer.IdentityNumber,
		Cif:            customer.Cif,
		Sid:            customer.Sid,
		ReferralCode:   customer.ReferralCode.String,
		Status:         customer.Status,
	}
}

func ToFinancialSyncVO(financial *FinancialData) *dto.FinancialSyncVO {

	if financial == nil {
		return nil
	}

	return &dto.FinancialSyncVO{
		MainAccountNumber:         financial.MainAccountNumber,
		AccountNumber:             financial.AccountNumber,
		GoldSavingStatus:          financial.GoldSavingStatus,
		GoldCardApplicationNumber: financial.GoldCardApplicationNumber,
		GoldCardAccountNumber:     financial.GoldCardAccountNumber,
		Balance:                   financial.Balance,
	}
}

func ToCredentialSyncVO(credential *Credential) (*dto.CredentialSyncVO, error) {

	if credential == nil {
		return &dto.CredentialSyncVO{}, nil
	}

	var credentialMetadata dto.MetadataCredentialVO
	err := json.Unmarshal(credential.Metadata, &credentialMetadata)
	if err != nil {
		return nil, errx.Trace(err)
	}

	return &dto.CredentialSyncVO{
		Password:            "",
		NextPasswordResetAt: credential.NextPasswordResetAt.Time.Unix(),
		Pin:                 credential.Pin,
		PinUpdatedAt:        credential.PinUpdatedAt.Time.Unix(),
		PinLastAccessAt:     credential.PinLastAccessAt.Time.Unix(),
		PinCounter:          credential.PinCounter,
		PinBlockedStatus:    credential.PinBlockedStatus,
		IsLocked:            credential.IsLocked,
		LoginFailCount:      credential.LoginFailCount,
		WrongPasswordCount:  credential.WrongPasswordCount,
		BlockedAt:           credential.BlockedAt.Time.Unix(),
		BlockedUntilAt:      credential.BlockedUntilAt.Time.Unix(),
		BiometricLogin:      credential.BiometricLogin,
		BiometricDeviceID:   credential.BiometricDeviceID,
		Metadata:            credentialMetadata,
	}, nil
}

func ToVerificationSyncVO(verification *Verification) *dto.VerificationSyncVO {

	if verification == nil {
		return &dto.VerificationSyncVO{}
	}

	return &dto.VerificationSyncVO{
		KycVerifiedStatus:               verification.KycVerifiedStatus,
		EmailVerificationToken:          verification.EmailVerificationToken,
		EmailVerifiedStatus:             verification.EmailVerifiedStatus,
		DukcapilVerifiedStatus:          verification.DukcapilVerifiedStatus,
		FinancialTransactionStatus:      verification.FinancialTransactionStatus,
		FinancialTransactionActivatedAt: verification.FinancialTransactionActivatedAt.Time.Unix(),
	}
}

func ToAddressSyncVO(address *Address) *dto.AddressSyncVO {

	if address == nil {
		return &dto.AddressSyncVO{}
	}

	return &dto.AddressSyncVO{
		Purpose:         address.Purpose,
		ProvinceID:      address.ProvinceID.Int64,
		ProvinceName:    address.ProvinceName.String,
		CityID:          address.CityID.Int64,
		CityName:        address.CityName.String,
		DistrictID:      address.DistrictID.Int64,
		DistrictName:    address.DistrictName.String,
		SubDistrictID:   address.SubDistrictID.Int64,
		SubDistrictName: address.SubDistrictName.String,
		Line:            address.Line.String,
		PostalCode:      address.PostalCode.String,
		IsPrimary:       address.IsPrimary.Bool,
	}
}
