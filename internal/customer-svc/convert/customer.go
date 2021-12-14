package convert

import (
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/rs/xid"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
)

func ModelUserToCustomer(user *model.User) (*model.Customer, error) {

	customerVO := dto.CustomerProfileVO{
		MaidenName:         user.NamaIbu.String,
		Nationality:        user.Kewarganegaraan,
		DateOfBirth:        nval.ParseStringFallback(user.TglLahir.Time, ""),
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
	profile, err := json.Marshal(customerVO)
	if err != nil {
		return nil, err
	}
	// TODO adjust foto_url
	photos := map[string]string{
		"foto_url": nval.ParseStringFallback(user.FotoUrl, ""),
	}

	photoRawMessage, err := json.Marshal(photos)
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
		UserRefId:      nval.ParseInt64Fallback(user.UserAiid, 0),
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

func ModelUserToCredential(user model.User, userPin *model.UserPin) (*model.Credential, error) {
	itemMetaData := model.NewItemMetadata(ModifierDTOToModel(dto.Modifier{ID: "", Role: "", FullName: ""}))

	//credential
	credential := &model.Credential{
		Xid:      strings.ToUpper(xid.New().String()),
		Password: user.Password.String,
		NextPasswordResetAt: sql.NullTime{
			Time:  user.NextPasswordReset.Time,
			Valid: user.NextPasswordReset.Valid,
		},
		Pin:    user.Pin.String,
		PinCif: "", // TODO insert from user_pin.cif
		PinUpdatedAt: sql.NullTime{
			Time:  user.LastUpdatePin.Time,
			Valid: user.LastUpdatePin.Valid,
		},
		PinLastAccessAt:    sql.NullTime{}, // TODO insert from user_pin.last_access_time
		PinCounter:         0,              // TODO insert from user_pin.counter
		PinBlockedStatus:   0,              // TODO insert from user_pin.is_blocked
		IsLocked:           user.IsLocked.Int64,
		LoginFailCount:     user.LoginFailCount,
		WrongPasswordCount: user.WrongPasswordCount,
		BlockedAt: sql.NullTime{
			Time:  user.BlockedDate.Time,
			Valid: user.BlockedDate.Valid,
		},
		BlockedUntilAt: sql.NullTime{
			Time:  user.BlockedToDate.Time,
			Valid: user.BlockedToDate.Valid,
		},
		BiometricLogin:    user.IsSetBiometric.Int64,
		BiometricDeviceId: user.DeviceIdBiometric.String,
		ItemMetadata:      itemMetaData,
	}

	metadata := dto.MetadataCredential{
		TryLoginAt:   user.TryLoginDate.Time.String(),
		PinCreatedAt: "", // TODO insert from user_pin.created_at
		PinBlockedAt: "", // TODO insert from user_pin.blocked_date
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
		credential.PinCif = userPin.Cif.String
		credential.PinBlockedStatus = userPin.IsBlocked
	}

	return credential, nil

}
