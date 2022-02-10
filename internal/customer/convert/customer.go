package convert

import (
	"database/sql"
	"encoding/json"
	"github.com/rs/xid"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
	"strings"
)

func ModelUserToCustomer(user *model.User) (*model.Customer, error) {

	// prepare profile value object
	profileVO := dto.CustomerProfileVO{
		MaidenName:         user.NamaIbu.String,
		Gender:             user.JenisKelamin,
		Nationality:        user.Kewarganegaraan,
		DateOfBirth:        nval.ParseStringFallback(user.TglLahir.Time.Format("02-01-2006"), ""),
		PlaceOfBirth:       user.TempatLahir.String,
		IdentityPhotoFile:  user.FotoKtpUrl,
		MarriageStatus:     user.StatusKawin,
		NPWPNumber:         user.NoNpwp,
		NPWPPhotoFile:      user.FotoNpwp,
		NPWPUpdatedAt:      user.LastUpdateDataNpwp,
		ProfileUpdatedAt:   nval.ParseStringFallback(user.LastUpdate.Time, ""),
		CifLinkUpdatedAt:   user.LastUpdateLinkCif,
		CifUnlinkUpdatedAt: nval.ParseStringFallback(user.LastUpdate.Time, ""),
		SidPhotoFile:       user.FotoSid.String,
	}
	profile, err := json.Marshal(profileVO)
	if err != nil {
		return nil, err
	}
	photosVO := dto.CustomerPhotoVO{
		Xid:      strings.ToUpper(xid.New().String()),
		Filename: nval.ParseStringFallback(user.FotoUrl, ""),
		Filesize: 0,
		Mimetype: "",
	}

	photoRawMessage, err := json.Marshal(photosVO)
	if err != nil {
		return nil, err
	}

	userSnapshot, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	customerMetaData := &dto.CustomerMetadata{
		Snapshot:          string(userSnapshot),
		SnapshotSignature: "", // TODO snapshot signature
	}

	customerMetadataRaw, err := json.Marshal(customerMetaData)
	if err != nil {
		return nil, err
	}

	itemMetaData := model.NewItemMetadata(ModifierDTOToModel(dto.Modifier{ID: "", Role: "", FullName: ""}))
	//customer
	customerXID := strings.ToUpper(xid.New().String())
	customer := &model.Customer{
		CustomerXID:    customerXID,
		FullName:       user.Nama.String,
		Phone:          user.NoHp.String,
		Email:          user.Email.String,
		IdentityType:   nval.ParseInt64Fallback(user.JenisIdentitas, 0),
		IdentityNumber: user.NoKtp.String,
		UserRefId:      nval.ParseStringFallback(user.UserAiid, ""),
		Photos:         photoRawMessage,
		Profile:        profile,
		Cif:            user.Cif,
		Sid:            user.NoSid.String,
		ReferralCode:   nval.ParseStringFallback(user.ReferralCode, ""),
		Status:         user.Status.Int64,
		Metadata:       customerMetadataRaw,
		ItemMetadata:   itemMetaData,
	}

	return customer, nil
}

func ModelUserToCredential(user *model.User, userPin *model.UserPin) (*model.Credential, error) {
	itemMetaData := model.NewItemMetadata(ModifierDTOToModel(dto.Modifier{ID: "", Role: "", FullName: ""}))

	credential := &model.Credential{
		Xid:                 strings.ToUpper(xid.New().String()),
		Password:            user.Password.String,
		NextPasswordResetAt: ModifierNullTime(user.NextPasswordReset),
		Pin:                 user.Pin.String,
		PinCif:              sql.NullString{},
		PinUpdatedAt:        ModifierNullTime(user.LastUpdatePin),
		PinLastAccessAt:     sql.NullTime{},
		PinCounter:          0,
		PinBlockedStatus:    0,
		IsLocked:            user.IsLocked.Int64,
		LoginFailCount:      user.LoginFailCount,
		WrongPasswordCount:  user.WrongPasswordCount,
		BlockedAt:           ModifierNullTime(user.BlockedDate),
		BlockedUntilAt:      ModifierNullTime(user.BlockedToDate),
		BiometricLogin:      user.IsSetBiometric.Int64,
		BiometricDeviceId:   user.DeviceIdBiometric.String,
		ItemMetadata:        itemMetaData,
	}

	metadata := dto.MetadataCredential{
		TryLoginAt:   user.TryLoginDate.Time.String(),
		PinCreatedAt: "",
		PinBlockedAt: "",
	}

	metadataRawMessage, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	credential.Metadata = metadataRawMessage
	if userPin != nil {
		metadata.PinCreatedAt = userPin.CreatedAt
		metadata.PinBlockedAt = userPin.BlockedDate.Time.String()

		credential.PinCounter = userPin.Counter
		credential.PinLastAccessAt = userPin.LastAccessTime
		credential.PinCif = userPin.Cif
		credential.PinBlockedStatus = userPin.IsBlocked
	}

	return credential, nil
}

func ModelUserToVerification(user *model.User) (*model.Verification, error) {
	itemMetaData := model.NewItemMetadata(ModifierDTOToModel(dto.Modifier{ID: "", Role: "", FullName: ""}))

	//verification
	verification := &model.Verification{
		Xid:                             strings.ToUpper(xid.New().String()),
		KycVerifiedStatus:               user.KycVerified,
		KycVerifiedAt:                   sql.NullTime{},
		EmailVerificationToken:          user.EmailVerificationToken,
		EmailVerifiedStatus:             user.EmailVerified,
		EmailVerifiedAt:                 sql.NullTime{},
		DukcapilVerifiedStatus:          user.IsDukcapilVerified.Int64,
		DukcapilVerifiedAt:              sql.NullTime{},
		FinancialTransactionStatus:      user.AktifasiTransFinansial,
		FinancialTransactionActivatedAt: ModifierNullTime(user.TanggalAktifasiFinansial),
		Metadata:                        []byte("{}"),
		ItemMetadata:                    itemMetaData,
	}

	return verification, nil
}

func ModelUserToFinancialData(user *model.User) (*model.FinancialData, error) {
	itemMetaData := model.NewItemMetadata(ModifierDTOToModel(dto.Modifier{ID: "", Role: "", FullName: ""}))

	financialData := &model.FinancialData{
		Xid:                       strings.ToUpper(xid.New().String()),
		MainAccountNumber:         user.NorekUtama.String,
		AccountNumber:             user.Norek,
		GoldSavingStatus:          user.IsOpenTe.Int64,
		GoldCardApplicationNumber: user.GoldcardApplicationNumber.String,
		GoldCardAccountNumber:     user.GoldcardApplicationNumber.String,
		Balance:                   user.Saldo,
		Metadata:                  []byte("{}"),
		ItemMetadata:              itemMetaData,
	}

	return financialData, nil
}

func ModelUserToAddress(user *model.User, userAddress *model.AddressExternal) (*model.Address, error) {
	itemMetaData := model.NewItemMetadata(ModifierDTOToModel(dto.Modifier{ID: "", Role: "", FullName: ""}))

	var purpose int64
	if user.Domisili.String == "1" {
		purpose = constant.Domicile
	} else {
		purpose = constant.IdentityCard
	}

	address := &model.Address{
		Xid:             strings.ToUpper(xid.New().String()),
		Purpose:         purpose,
		ProvinceId:      userAddress.IdProvinsi,
		ProvinceName:    userAddress.Provinsi,
		CityId:          userAddress.IdKabupaten,
		CityName:        userAddress.Kabupaten,
		DistrictId:      userAddress.IdKecamatan,
		DistrictName:    userAddress.Kecamatan,
		SubDistrictId:   userAddress.IdKelurahan,
		SubDistrictName: userAddress.Kelurahan,
		PostalCode:      userAddress.Kodepos,
		Metadata:        []byte("{}"),
		ItemMetadata:    itemMetaData,
	}

	return address, nil
}
