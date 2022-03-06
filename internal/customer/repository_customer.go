package customer

import (
	"database/sql"
	"errors"
	"github.com/nbs-go/errx"
	"github.com/nbs-go/nlogger/v2"
	"github.com/rs/xid"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
	"strings"
	"time"
)

func (rc *RepositoryContext) CreateCustomer(row *model.Customer) (int64, error) {
	var lastInsertID int64
	err := rc.stmt.Customer.Insert.QueryRowContext(rc.ctx, &row).Scan(&lastInsertID)
	if err != nil {
		rc.log.Error("error when insert customer", nlogger.Error(err), nlogger.Context(rc.ctx))
		return 0, err
	}

	return lastInsertID, nil
}

func (rc *RepositoryContext) FindCustomerByID(id int64) (*model.Customer, error) {
	var row model.Customer
	err := rc.stmt.Customer.FindByID.GetContext(rc.ctx, &row, id)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (rc *RepositoryContext) FindCustomerByUserRefID(id string) (*model.Customer, error) {
	var row model.Customer
	err := rc.stmt.Customer.FindByRefID.GetContext(rc.ctx, &row, id)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (rc *RepositoryContext) FindCustomerByPhone(phone string) (*model.Customer, error) {
	var row model.Customer
	err := rc.stmt.Customer.FindByPhone.GetContext(rc.ctx, &row, phone)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (rc *RepositoryContext) FindCustomerByEmail(email string) (*model.Customer, error) {
	var row model.Customer
	err := rc.stmt.Customer.FindByEmail.GetContext(rc.ctx, &row, email)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (rc *RepositoryContext) FindCustomerByEmailOrPhone(email string) (*model.Customer, error) {
	var row model.Customer
	err := rc.stmt.Customer.FindByEmailOrPhone.GetContext(rc.ctx, &row, email, email)
	return &row, err
}

func (rc *RepositoryContext) FindCustomerByPhoneOrCIF(cif string) (*model.Customer, error) {
	var row model.Customer
	err := rc.stmt.Customer.FindByPhoneOrCIF.GetContext(rc.ctx, &row, cif, cif)
	return &row, err
}

func (rc *RepositoryContext) ReferralCodeExists(referralCode string) (*model.Customer, error) {
	var row model.Customer
	err := rc.stmt.Customer.ReferralCodeExist.GetContext(rc.ctx, &row, referralCode)
	return &row, err
}

func (rc *RepositoryContext) UpdateCustomerByCIF(customer *model.Customer, cif string) error {
	result, err := rc.stmt.Customer.UpdateByCIF.ExecContext(rc.ctx, &model.UpdateByCIF{
		Customer: customer,
		Cif:      cif,
	})
	if err != nil {
		return errx.Trace(err)
	}
	if !nsql.IsUpdated(result) {
		return constant.ResourceNotFoundError
	}
	return nil
}

func (rc *RepositoryContext) UpdateCustomerProfile(customer *model.Customer, payload dto.UpdateProfileRequest) error {
	tx, err := rc.conn.BeginTxx(rc.ctx, nil)
	if err != nil {
		return errx.Trace(err)
	}
	defer rc.ReleaseTx(tx, &err)

	// update customer model
	customer.FullName = payload.Nama
	customer.Profile.MaidenName = payload.NamaIbu
	customer.Profile.PlaceOfBirth = payload.TempatLahir
	customer.Profile.DateOfBirth = payload.TglLahir
	customer.IdentityType = nval.ParseInt64Fallback(payload.JenisIdentitas, 10)
	customer.IdentityNumber = payload.NoKtp
	customer.Profile.MarriageStatus = payload.StatusKawin
	customer.Profile.Gender = payload.JenisKelamin
	customer.Profile.Nationality = payload.Kewarganegaraan
	customer.Profile.IdentityExpiredAt = payload.TanggalExpiredIdentitas
	customer.Profile.Religion = payload.Agama
	customer.Profile.ProfileUpdatedAt = time.Now().Unix()

	// Get current address data
	address, errAddress := rc.FindAddressByCustomerId(customer.ID)
	if errAddress != nil && !errors.Is(errAddress, sql.ErrNoRows) {
		rc.log.Error("error when get customer data", nlogger.Error(err))
		return errx.Trace(err)
	}

	// Update address model
	address.Line = nsql.NewNullString(payload.Alamat)
	address.ProvinceID = sql.NullInt64{Int64: nval.ParseInt64Fallback(payload.IDProvinsi, 0), Valid: true}
	address.CityID = sql.NullInt64{Int64: nval.ParseInt64Fallback(payload.IDKabupaten, 0), Valid: true}
	address.DistrictID = sql.NullInt64{Int64: nval.ParseInt64Fallback(payload.IDKecamatan, 0), Valid: true}
	address.SubDistrictID = sql.NullInt64{Int64: nval.ParseInt64Fallback(payload.IDKelurahan, 0), Valid: true}
	address.ProvinceName = sql.NullString{String: payload.NamaProvinsi, Valid: true}
	address.CityName = sql.NullString{String: payload.NamaKabupaten, Valid: true}
	address.DistrictName = sql.NullString{String: payload.NamaKecamatan, Valid: true}
	address.SubDistrictName = sql.NullString{String: payload.NamaKelurahan, Valid: true}
	address.PostalCode = sql.NullString{String: payload.KodePos, Valid: true}

	// if empty create new address
	if errors.Is(errAddress, sql.ErrNoRows) {
		address.CustomerID = customer.ID
		address.Xid = strings.ToUpper(xid.New().String())
		address.Metadata = nsql.EmptyObjectJSON
		address.Purpose = constant.IdentityCard
		address.IsPrimary = sql.NullBool{Bool: true, Valid: true}
		address.BaseField = model.EmptyBaseField
	}

	// Update the data to repositories
	customerUpdate, err := rc.stmt.Customer.UpdateByUserRefID.ExecContext(rc.ctx, &model.UpdateByID{
		Customer: customer,
		ID:       customer.ID,
	})
	if err != nil {
		return errx.Trace(err)
	}
	if !nsql.IsUpdated(customerUpdate) {
		return constant.ResourceNotFoundError
	}

	err = rc.InsertOrUpdateAddress(address)
	if err != nil {
		return errx.Trace(err)
	}

	return nil
}

func (rc *RepositoryContext) UpdateCustomerByPhone(customer *model.Customer) error {
	result, err := rc.stmt.Customer.UpdateByPhone.ExecContext(rc.ctx, customer)
	if err != nil {
		return errx.Trace(err)
	}
	if !nsql.IsUpdated(result) {
		return constant.ResourceNotFoundError
	}
	return nil
}

func (rc *RepositoryContext) UpdateCustomerByUserRefID(customer *model.Customer, userRefID string) error {
	result, err := rc.stmt.Customer.UpdateByPhone.ExecContext(rc.ctx, &model.UpdateCustomerByUserRefID{
		Customer:  customer,
		UserRefID: userRefID,
	})
	if err != nil {
		rc.log.Error("error found when update customer by UserRefID", nlogger.Error(err))
		return errx.Trace(err)
	}
	if !nsql.IsUpdated(result) {
		return constant.ResourceNotFoundError
	}

	return nil
}

func (rc *RepositoryContext) FindCombineCustomerDataByUserRefID(userRefID string) (*model.Customer, *model.Verification, *model.Credential, error) {
	// Find Customer
	customer, err := rc.FindCustomerByUserRefID(userRefID)
	if err != nil {
		rc.log.Error("error when find current customer", nlogger.Error(err))
		return nil, nil, nil, errx.Trace(err)
	}
	// Find Verification
	verification, err := rc.FindVerificationByCustomerID(customer.ID)
	if err != nil {
		rc.log.Error("error when get verification", nlogger.Error(err))
		return nil, nil, nil, errx.Trace(err)
	}
	// Find Credential
	credential, err := rc.FindCredentialByCustomerID(customer.ID)
	if err != nil {
		rc.log.Error("error when get credential", nlogger.Error(err))
		return nil, nil, nil, errx.Trace(err)
	}

	return customer, verification, credential, nil
}
