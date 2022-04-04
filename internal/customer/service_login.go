package customer

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/imdario/mergo"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/nbs-go/errx"
	logOption "github.com/nbs-go/nlogger/v2/option"
	"regexp"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ntime"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
	"strings"
	"time"
)

func (s *Service) Login(payload dto.LoginPayload) (*dto.LoginResult, error) {
	// Check if user exists
	t := time.Now()
	customer, err := s.repo.FindCustomerByEmailOrPhone(payload.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("failed to retrieve customer", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	// Handle null Profile
	if customer.Profile == nil {
		customer.Profile = model.EmptyCustomerProfile
		// Update profile json
		err = s.repo.UpdateCustomerByPhone(customer)
		if err != nil {
			s.log.Error("error when update customer by phone", logOption.Error(err))
			return nil, errx.Trace(err)
		}
	}

	// Check on external database
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		// If data not found on internal database check on external database.
		user, errExternal := s.repoExternal.FindUserExternalByEmailOrPhone(payload.Email)
		if errExternal != nil && !errors.Is(errExternal, sql.ErrNoRows) {
			s.log.Error("error found when query find by email or phone", logOption.Error(errExternal))
			return nil, errx.Trace(errExternal)
		}

		if errExternal != nil && errors.Is(errExternal, sql.ErrNoRows) {
			s.log.Debug("Phone or email is not registered")
			return nil, constant.NoPhoneEmailError.Trace(errx.Source(errExternal))
		}

		// sync data external to internal
		customer, err = s.syncExternalToInternal(user)
		if err != nil {
			s.log.Error("error while sync data External to Internal", logOption.Error(err))
			return nil, errx.Trace(err)
		}
	}

	// Get credential customer
	credential, err := s.repo.FindCredentialByCustomerID(customer.ID)
	if err != nil {
		err = handleErrorRepository(err, constant.InvalidEmailPassInputError.Trace())
		return nil, errx.Trace(err)
	}

	// Check if account isn't blocked
	blockedUntil := ntime.ChangeTimezone(credential.BlockedUntilAt.Time, constant.WIB)
	now := ntime.ChangeTimezone(t, constant.WIB)
	if credential.BlockedUntilAt.Valid && blockedUntil.After(now) {
		// Compose custom message response
		message := "Akun dikunci hingga %v karena gagal login %v kali. Hubungi call center jika ini bukan kamu"
		timeBlocked := ntime.ChangeTimezone(credential.BlockedUntilAt.Time, constant.WIB).Format("02-Jan-2006 15:04:05")
		return nil, nhttp.UnauthorizedError.
			AddMetadata(nhttp.OverrideMessageMetadata, fmt.Sprintf(message, timeBlocked, credential.WrongPasswordCount))
	}

	// Counter wrong password count
	passwordRequest := stringToMD5(payload.Password)
	if credential.Password != passwordRequest {
		errB := s.HandleWrongPassword(credential, customer)
		return nil, errB
	}

	// If password verified remove counter
	err = s.UnblockPassword(credential)
	if err != nil {
		s.log.Error("error found when unblock password", logOption.Error(err))
		return nil, err
	}

	// get userRefId from external DB
	if !customer.UserRefID.Valid {
		registerPayload := &dto.CustomerSynchronizePayload{
			Name:        customer.FullName,
			Email:       customer.Email,
			PhoneNumber: customer.Phone,
			Password:    credential.Password,
			FcmToken:    payload.FcmToken,
		}
		// sync data from customer service to PDS API
		resultSync, errInternal := s.syncInternalToExternal(registerPayload)
		if errInternal != nil {
			s.log.Error("error found when sync to pds api", logOption.Error(errInternal))
			return nil, errx.Trace(errInternal)
		}
		// set userRefId
		customer.UserRefID = sql.NullString{
			Valid:  true,
			String: resultSync.UserAiid,
		}
		// update customer
		err = s.repo.UpdateCustomerByPhone(customer)
		if err != nil {
			s.log.Error("failed to update userRefId", logOption.Error(err))
			return nil, errx.Trace(err)
		}
	}

	// Get token from cache
	var token string
	cacheTokenKey := fmt.Sprintf("%v:%v:%v", constant.Prefix, constant.CacheTokenJWT, customer.UserRefID.String)
	token, err = s.setTokenAuthentication(customer, payload.Agen, payload.Version, cacheTokenKey)
	if err != nil {
		s.log.Error("error found when get access token from cache", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	// Check account is first login or not
	isFirstLogin, err := s.isFirstLogin(customer)
	if err != nil {
		return nil, errx.Trace(err)
	}

	// Prepare to insert audit login
	auditLogin := model.NewAuditLogin(customer, t, payload, GetChannelByAgen(payload.Agen))

	// Persist audit login
	err = s.repo.CreateAuditLogin(&auditLogin)
	if err != nil {
		s.log.Error("error found when create audit login", logOption.Error(err))
		return nil, constant.InvalidCredentialError.Trace()
	}

	// Check is force update password
	validatePassword := s.ValidatePassword(payload.Password)
	isForceUpdatePassword := false
	if !validatePassword.IsValid {
		isForceUpdatePassword = true
	}

	// Get data address
	address, err := s.repo.FindAddressByCustomerId(customer.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("error found when get customer address", logOption.Error(err))
		return nil, constant.InvalidEmailPassInputError.Trace(errx.Source(err))
	}

	// Get data verification
	verification, err := s.repo.FindVerificationByCustomerID(customer.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("error found when get data verification", logOption.Error(err))
		return nil, constant.InvalidEmailPassInputError.Trace(errx.Source(err))
	}

	// Get financial data
	financial, err := s.repo.FindFinancialDataByCustomerID(customer.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("error found when get financial data", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	var gs interface{}
	if len(customer.Cif) == 0 {
		gs = false
	} else {
		// Get gold saving account
		goldSaving, errGs := s.getListAccountNumber(customer.Cif, customer.UserRefID.String)
		if errGs != nil {
			return nil, errx.Trace(errGs)
		}

		gs = &dto.GoldSavingVO{
			TotalSaldoBlokir:  goldSaving.TotalSaldoBlokir,
			TotalSaldoSeluruh: goldSaving.TotalSaldoSeluruh,
			TotalSaldoEfektif: goldSaving.TotalSaldoEfektif,
			ListTabungan:      goldSaving.ListTabungan,
			PrimaryRekening:   goldSaving.PrimaryRekening,
		}
	}

	err = s.synchronizeWhenAuthenticated(customer, financial, verification, credential, address)
	if err != nil {
		s.log.Error("error found when synchronize", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	resp := dto.LoginVO{
		Customer:              customer,
		Address:               address,
		Profile:               customer.Profile,
		Verification:          verification,
		Financial:             financial,
		IsFirstLogin:          isFirstLogin,
		IsForceUpdatePassword: isForceUpdatePassword,
		GoldSaving:            gs,
		Token:                 token,
	}

	return s.composeLoginResponse(resp)
}

func (s *Service) syncInternalToExternal(payload *dto.CustomerSynchronizePayload) (*dto.UserVO, error) {
	// Set body
	reqBody := map[string]interface{}{
		"nama":      payload.Name,
		"email":     payload.Email,
		"no_hp":     payload.PhoneNumber,
		"password":  payload.Password,
		"fcm_token": payload.FcmToken,
	}

	// Sync customer
	resp, err := s.SynchronizeCustomer(reqBody)
	if err != nil {
		s.log.Error("error found when sync data customer via API PDS", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	// handle status error
	if resp.Status != constant.ResponseSuccess {
		s.log.Error("Get Error from SynchronizeCustomer")
		return nil, nhttp.InternalError.Trace(errx.Errorf(resp.Message))
	}

	// parsing response
	var user dto.UserVO
	err = json.Unmarshal(resp.Data, &user)
	if err != nil {
		s.log.Error("Cannot unmarshall data login pds", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	// set result
	result := &user

	return result, nil
}

func (s *Service) syncExternalToInternal(user *model.User) (*model.Customer, error) {
	// prepare customer
	customer, err := model.UserToCustomer(user)
	if err != nil {
		s.log.Error("failed to convert to model customer", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	// Check has userPin or not
	userPin, err := s.repoExternal.FindUserPINByCustomerID(customer.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("failed retrieve user pin from external database", logOption.Error(err))
		return nil, errx.Trace(err)
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		userPin = &model.UserPin{}
	}

	// Check user address on external database
	addressExternal, err := s.repoExternal.FindUserExternalAddressByCustomerID(user.UserAiid)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("failed retrieve address from external database", logOption.Error(err))
		return nil, errx.Trace(err)
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		addressExternal = &model.AddressExternal{}
	}

	// Prepare credential
	credential, err := model.UserToCredential(user, userPin)
	if err != nil {
		s.log.Error("failed convert to credential model", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	// Prepare financial data
	financialData, err := model.UserToFinancialData(user)
	if err != nil {
		s.log.Error("failed convert to financial data", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	// Prepare verification
	verification, err := model.UserToVerification(user)
	if err != nil {
		s.log.Error("failed convert to verification", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	// Prepare address
	address, err := model.UserToAddress(user, addressExternal)
	if err != nil {
		s.log.Error("failed convert address", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	// Persist customer data
	customerID, err := s.repo.CreateCustomer(customer)
	if err != nil {
		s.log.Error("failed to persist customer", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	// Set credential customer id to last inserted
	credential.CustomerID = customerID
	financialData.CustomerID = customerID
	verification.CustomerID = customerID
	address.CustomerID = customerID

	// persist credential
	err = s.repo.InsertOrUpdateCredential(credential)
	if err != nil {
		s.log.Error("failed persist to credential", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	// persist financial data
	err = s.repo.InsertOrUpdateFinancialData(financialData)
	if err != nil {
		s.log.Error("failed persist to financial data", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	// persist verification
	err = s.repo.InsertOrUpdateVerification(verification)
	if err != nil {
		s.log.Error("failed persist verification", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	// persist address
	err = s.repo.InsertOrUpdateAddress(address)
	if err != nil {
		s.log.Error("failed persist address", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	customer, err = s.repo.FindCustomerByID(customerID)
	if err != nil {
		s.log.Error("failed to retrieve customer not found", logOption.Error(err))
		return nil, constant.ResourceNotFoundError.Trace()
	}

	return customer, nil
}

func (s *Service) composeLoginResponse(data dto.LoginVO) (*dto.LoginResult, error) {
	// Cast to model
	customer := data.Customer.(*model.Customer)
	profile := data.Profile.(*model.CustomerProfile)
	verification := data.Verification.(*model.Verification)
	address := data.Address.(*model.Address)
	financial := data.Financial.(*model.FinancialData)
	gs := data.GoldSaving

	// Get asset url
	// -- Avatar URL
	avatarURL := s.AssetGetPublicURL(constant.AssetAvatarProfile, customer.Photos.FileName)
	// -- NPWP URL
	npwpURL := s.AssetGetPublicURL(constant.AssetNPWP, customer.Profile.NPWPPhotoFile)
	// -- KTP URL
	ktpURL := s.AssetGetPublicURL(constant.AssetKTP, customer.Profile.IdentityPhotoFile)
	// -- SID URL
	sidURL := s.AssetGetPublicURL(constant.AssetNPWP, customer.Profile.SidPhotoFile)

	return &dto.LoginResult{
		User: &dto.LoginUserVO{
			CustomerVO: dto.CustomerVO{
				ID:                        customer.UserRefID.String,
				IsKYC:                     nval.ParseStringFallback(verification.KycVerifiedStatus, "0"),
				Cif:                       customer.Cif,
				Nama:                      customer.FullName,
				NamaIbu:                   profile.MaidenName,
				NoKTP:                     customer.IdentityNumber,
				Email:                     customer.Email,
				JenisKelamin:              profile.Gender,
				TempatLahir:               profile.PlaceOfBirth,
				TglLahir:                  profile.DateOfBirth,
				Alamat:                    address.Line.String,
				IDProvinsi:                address.ProvinceID.Int64,
				IDKabupaten:               address.CityID.Int64,
				IDKecamatan:               address.DistrictID.Int64,
				IDKelurahan:               address.SubDistrictID.Int64,
				Kelurahan:                 address.DistrictName.String,
				Provinsi:                  address.ProvinceName.String,
				Kabupaten:                 address.CityName.String,
				Kecamatan:                 address.DistrictName.String,
				KodePos:                   address.PostalCode.String,
				NoHP:                      customer.Phone,
				FotoKTP:                   ktpURL,
				IsEmailVerified:           nval.ParseStringFallback(verification.EmailVerifiedStatus, "0"),
				Kewarganegaraan:           profile.Nationality,
				JenisIdentitas:            fmt.Sprintf("%v", customer.IdentityType),
				NoIdentitas:               customer.IdentityNumber,
				TglExpiredIdentitas:       profile.IdentityExpiredAt,
				NoNPWP:                    profile.NPWPNumber,
				Avatar:                    avatarURL,
				FotoNPWP:                  npwpURL,
				FotoSid:                   sidURL,
				NoSid:                     customer.Sid,
				StatusKawin:               profile.MarriageStatus,
				Norek:                     financial.AccountNumber,
				Saldo:                     nval.ParseStringFallback(financial.Balance, "0"),
				AktifasiTransFinansial:    nval.ParseStringFallback(verification.FinancialTransactionStatus, ""),
				IsDukcapilVerified:        nval.ParseStringFallback(verification.DukcapilVerifiedStatus, "0"),
				IsOpenTe:                  nval.ParseStringFallback(financial.GoldSavingStatus, "0"),
				ReferralCode:              customer.ReferralCode.String,
				GoldCardApplicationNumber: financial.GoldCardApplicationNumber,
				GoldCardAccountNumber:     financial.GoldCardAccountNumber,
				KodeCabang:                "", // TODO: Branch Code
				TabunganEmas:              gs,
			},
			IsFirstLogin:          data.IsFirstLogin,
			IsForceUpdatePassword: data.IsForceUpdatePassword,
		},
		JwtToken: data.Token,
	}, nil
}

func (s *Service) setTokenAuthentication(customer *model.Customer, agen string, version string, cacheTokenKey string) (string, error) {
	var accessToken string
	accessToken, err := s.CacheGet(cacheTokenKey)
	if err != nil {
		s.log.Error("error found when get cache", logOption.Error(err))
		return "", err
	}

	accessToken, err = s.setNewAccessToken(accessToken, cacheTokenKey)
	if err != nil {
		return "", errx.Trace(err)
	}

	channelID := GetChannelByAgen(agen)
	now := time.Now()

	// Generate JWT
	token, err := jwt.NewBuilder().
		Claim("id", customer.UserRefID.String).
		Claim("email", customer.Email).
		Claim("nama", customer.FullName).
		Claim("no_hp", customer.Phone).
		Claim("access_token", accessToken).
		Claim("agen", agen).
		Claim("channelId", channelID).
		Claim("version", version).
		IssuedAt(now).
		Expiration(now.Add(time.Second * time.Duration(s.config.ClientConfig.JWTExpiry))).
		Issuer(constant.JWTIssuer).
		Build()
	if err != nil {
		s.log.Error("error found when generate JWT", logOption.Error(err))
		return "", errx.Trace(err)
	}

	jwtKey := s.config.ClientConfig.JWTKey
	jwtKeyBytes := []byte(jwtKey)

	// sign token
	signed, err := jwt.Sign(token, constant.JWTSignature, jwtKeyBytes)
	if err != nil {
		s.log.Error("failed to sign token", logOption.Error(err))
		return "", errx.Trace(err)
	}
	tokenString := string(signed)

	return tokenString, nil
}

func (s *Service) setNewAccessToken(accessToken string, cacheTokenKey string) (string, error) {
	if accessToken != "" {
		return accessToken, nil
	}
	// Generate access token
	newAccessToken := nval.Bin2Hex(nval.RandStringBytes(78))
	// Set token to cache
	cacheToken, cErr := s.CacheSetThenGet(cacheTokenKey, newAccessToken, s.config.ClientConfig.JWTExpiry)
	if cErr != nil {
		s.log.Error("error found when set token to cache", logOption.Error(cErr))
		return "", errx.Trace(cErr)
	}
	return cacheToken, nil
}

func (s *Service) ValidatePassword(password string) *dto.ValidatePasswordResult {
	var validation dto.ValidatePasswordResult

	lowerCase := regexp.MustCompile(`[a-z]+`)
	upperCase := regexp.MustCompile(`[A-Z]+`)
	allNumber := regexp.MustCompile(`[0-9]+`)

	if len(lowerCase.FindStringSubmatch(password)) < 1 {
		validation.IsValid = false
		validation.ErrCode = "isLower"
		validation.Message = "Password harus terdapat satu huruf kecil."
		return &validation
	}

	if len(upperCase.FindStringSubmatch(password)) < 1 {
		validation.IsValid = false
		validation.ErrCode = "isUpper"
		validation.Message = "Password harus terdapat satu huruf kapital.."
		return &validation
	}

	if len(allNumber.FindStringSubmatch(password)) < 1 {
		validation.IsValid = false
		validation.ErrCode = "isNumber"
		validation.Message = "Password harus terdapat angka."
		return &validation
	}

	if len(password) < 8 {
		validation.IsValid = false
		validation.ErrCode = "length"
		validation.Message = "Pasword harus terdapat minimal 8 karakter."
		return &validation
	}

	if strings.Contains(strings.ToLower(password), "gadai") {
		validation.IsValid = false
		validation.ErrCode = "containsGadai"
		validation.Message = "Hindari menggunakan kata gadai."
		return &validation
	}

	validation.IsValid = true
	validation.ErrCode = ""

	return &validation
}

func (s *Service) HandleWrongPassword(credential *model.Credential, customer *model.Customer) error {
	t := time.Now()

	// Unmarshalling metadata credential to get tryLoginAt
	var Metadata dto.MetadataCredential
	err := json.Unmarshal(credential.Metadata, &Metadata)
	if err != nil {
		s.log.Error("error found when unmarshal metadata credential", logOption.Error(err))
		return errx.Trace(err)
	}
	// Parse time from metadata string to time
	tryLoginAt, _ := time.Parse(time.RFC3339, Metadata.TryLoginAt)
	now := ntime.ChangeTimezone(t, constant.WIB)

	// If user is not trying to login after 1 day, set wrongPassword to 0
	if now.After(tryLoginAt.Add(time.Hour * time.Duration(24))) {
		credential.WrongPasswordCount = 0
	}

	// If now is after than blockedUntilAt set wrong password to 0 and unblock account
	blockedUntil := ntime.ChangeTimezone(credential.BlockedUntilAt.Time, constant.WIB)
	if credential.WrongPasswordCount == constant.MaxWrongPassword && now.After(blockedUntil) {
		credential.BlockedAt = sql.NullTime{}
		credential.BlockedUntilAt = sql.NullTime{}
		credential.WrongPasswordCount = 0
	}

	wrongCount := credential.WrongPasswordCount + 1

	switch wrongCount {
	case constant.Warn2XWrongPassword:
		err = constant.InvalidPhoneInput1Error.Trace()
		credential.WrongPasswordCount = wrongCount
	case constant.Warn4XWrongPassword:
		err = constant.InvalidPhoneInput2Error.Trace()
		credential.WrongPasswordCount = wrongCount
	case constant.MinWrongPassword:
		// Set block account
		hour := 1 // Block for 1 hours
		duration := time.Hour * time.Duration(hour)
		addDuration := t.Add(duration)
		credential.BlockedAt = sql.NullTime{
			Time:  t,
			Valid: true,
		}
		credential.BlockedUntilAt = sql.NullTime{
			Time:  addDuration,
			Valid: true,
		}
		credential.WrongPasswordCount = wrongCount

		// Set response if blocked for 1 hour
		message := "Akun dikunci hingga %v WIB karena gagal login %v kali. Hubungi call center jika ini bukan kamu"
		timeBlocked := ntime.ChangeTimezone(credential.BlockedUntilAt.Time, constant.WIB).Format("02-Jan-2006 15:04:05")
		err = nhttp.UnauthorizedError.
			AddMetadata(nhttp.OverrideMessageMetadata, fmt.Sprintf(message, timeBlocked, credential.WrongPasswordCount))

		// Send OTP To Phone Number
		request := dto.SendOTPRequest{
			PhoneNumber: customer.Phone,
			RequestType: constant.RequestTypeBlockOneHour,
		}
		_, errOTP := s.SendOTP(request)
		if errOTP != nil {
			s.log.Debug("error found when sending otp block one hour", logOption.Error(errOTP))
		}

		// Send Notification Blocked Login One Hour
		// TODO Refactor
		_ = s.SendNotificationBlock(dto.NotificationBlock{
			Customer:     customer,
			Message:      fmt.Sprintf(message, timeBlocked, credential.WrongPasswordCount),
			LastTryLogin: ntime.NewTimeWIB(tryLoginAt).Format("02-Jan-2006 15:04:05"),
		})
	case constant.MaxWrongPassword:
		// Set block account
		hour := 24 // Block for 24 hours
		duration := time.Hour * time.Duration(hour)
		addDuration := t.Add(duration)
		credential.BlockedAt = sql.NullTime{
			Time:  t,
			Valid: true,
		}
		credential.BlockedUntilAt = sql.NullTime{
			Time:  addDuration,
			Valid: true,
		}
		credential.WrongPasswordCount = wrongCount

		// Set response if blocked for 24 hour
		message := "Akun dikunci hingga %v WIB karena gagal login %v kali. Hubungi call center jika ini bukan kamu"
		timeBlocked := ntime.ChangeTimezone(credential.BlockedUntilAt.Time, constant.WIB).Format("02-Jan-2006 15:04:05")
		err = nhttp.UnauthorizedError.
			AddMetadata(nhttp.OverrideMessageMetadata, fmt.Sprintf(message, timeBlocked, credential.WrongPasswordCount))

		// Send OTP To Phone Number
		request := dto.SendOTPRequest{
			PhoneNumber: customer.Phone,
			RequestType: constant.RequestTypeBlockOneDay,
		}
		_, errOTP := s.SendOTP(request)
		if errOTP != nil {
			s.log.Debug("Error when sending otp block one hour", logOption.Error(errOTP))
		}

		// Send Notification Blocked Login One Day
		_ = s.SendNotificationBlock(dto.NotificationBlock{
			Customer:     customer,
			Message:      fmt.Sprintf(message, timeBlocked, credential.WrongPasswordCount),
			LastTryLogin: ntime.NewTimeWIB(tryLoginAt).Format("02-Jan-2006 15:04:05"),
		})
	default:
		err = constant.InvalidEmailPassInputError.Trace()
		credential.WrongPasswordCount = wrongCount
	}

	// Handle notification error
	if err != nil {
		s.log.Debug("Error when sending notification block", logOption.Error(err))
	}

	// Set trying login at to metadata
	var Format dto.MetadataCredential
	Format.TryLoginAt = t.Format(time.RFC3339)
	Format.PinCreatedAt = Metadata.PinCreatedAt
	Format.PinBlockedAt = Metadata.PinBlockedAt
	MetadataCredential, _ := json.Marshal(&Format)
	credential.Metadata = MetadataCredential

	uErr := s.repo.UpdateCredential(credential)
	if uErr != nil {
		s.log.Error("error when update credential.", logOption.Error(err))
		return errx.Trace(uErr)
	}

	return err
}

func (s *Service) UnblockPassword(credential *model.Credential) error {
	credential.BlockedAt = sql.NullTime{}
	credential.BlockedUntilAt = sql.NullTime{}
	credential.WrongPasswordCount = 0

	err := s.repo.UpdateCredential(credential)
	if err != nil {
		s.log.Error("error when update credential.", logOption.Error(err))
		return errx.Trace(err)
	}

	return nil
}

func (s *Service) isFirstLogin(customer *model.Customer) (bool, error) {
	// Check account is first login or not
	countAuditLog, err := s.repo.CountAuditLogin(customer.ID)
	if err != nil {
		s.log.Error("error found when count audit login", logOption.Error(err))
		return false, errx.Trace(err)
	}

	// Set is first login is true or false.
	var isFirstLogin = true
	if countAuditLog > 0 {
		isFirstLogin = false
	}

	return isFirstLogin, nil
}

func prepareBodyCustomer(customer *model.Customer, dest map[string]interface{}) error {
	profile := customer.Profile
	parseDateOfBirth, _ := time.Parse("02-01-2006", profile.DateOfBirth)
	data := map[string]interface{}{
		// `user` table
		"nama":                      customer.FullName,
		"email":                     customer.Email,
		"no_hp":                     customer.Phone,
		"foto_url":                  customer.Photos.FileName,
		"cif":                       customer.Cif,
		"no_sid":                    customer.Sid,
		"referral_code":             customer.ReferralCode.String,
		"status":                    customer.Status,
		"jenis_identitas":           customer.IdentityType,
		"no_ktp":                    customer.IdentityNumber,
		"last_update":               customer.UpdatedAt.Format(constant.DateTimeLayout),
		"last_update_data_npwp":     profile.NPWPUpdatedAt,
		"last_update_data_nasabah":  profile.ProfileUpdatedAt,
		"last_update_link_cif":      profile.CifLinkUpdatedAt,
		"last_update_unlink_cif":    profile.CifUnlinkUpdatedAt,
		"nama_ibu":                  profile.MaidenName,
		"jenis_kelamin":             profile.Gender,
		"kewarganegaraan":           profile.Nationality,
		"tgl_lahir":                 parseDateOfBirth.Format(constant.DateLayout),
		"tempat_lahir":              profile.PlaceOfBirth,
		"foto_ktp_url":              profile.IdentityPhotoFile,
		"tanggal_expired_identitas": profile.IdentityExpiredAt,
		"agama":                     profile.Religion,
		"status_kawin":              profile.MarriageStatus,
		"no_npwp":                   profile.NPWPNumber,
		"foto_npwp":                 profile.NPWPPhotoFile,
		"foto_sid":                  profile.SidPhotoFile,
	}
	err := mergo.Merge(&dest, data)
	if err != nil {
		return err
	}

	return nil
}

func prepareBodyFinancial(financial *model.FinancialData, dest map[string]interface{}) error {
	data := map[string]interface{}{
		"norek_utama":                 financial.MainAccountNumber,
		"norek":                       financial.AccountNumber,
		"is_open_te":                  financial.GoldSavingStatus,
		"goldcard_application_number": financial.GoldCardApplicationNumber,
		"goldcard_account_number":     financial.GoldCardAccountNumber,
		"saldo":                       financial.Balance,
	}
	err := mergo.Merge(&dest, data)
	if err != nil {
		return err
	}
	return nil
}

func prepareBodyVerification(verification *model.Verification, dest map[string]interface{}) error {
	data := map[string]interface{}{
		"email_verification_token":   verification.EmailVerificationToken,
		"email_verified":             verification.EmailVerifiedStatus,
		"is_dukcapil_verified":       verification.DukcapilVerifiedStatus,
		"aktifasiTransFinansial":     verification.DukcapilVerifiedStatus,
		"tanggal_aktifasi_finansial": verification.FinancialTransactionActivatedAt,
	}
	err := mergo.Merge(&dest, data)
	if err != nil {
		return err
	}
	return nil
}

func prepareBodyCredential(credential *model.Credential, dest map[string]interface{}) error {
	data := map[string]interface{}{
		"password":             credential.Password,
		"next_password_reset":  credential.NextPasswordResetAt,
		"pin":                  credential.Pin,
		"last_update_pin":      credential.PinUpdatedAt,
		"blocked_date":         credential.BlockedAt,
		"blocked_to_date":      credential.BlockedUntilAt,
		"login_fail_count":     credential.LoginFailCount,
		"wrong_password_count": credential.WrongPasswordCount,
		"is_set_biometric":     credential.BiometricLogin,
		"device_id_biometric":  credential.BiometricDeviceID,
	}
	err := mergo.Merge(&dest, data)
	if err != nil {
		return err
	}
	return nil
}

func prepareBodyAddress(address *model.Address, dest map[string]interface{}) error {
	data := map[string]interface{}{
		"id_kelurahan": address.SubDistrictID.Int64,
		"alamat":       address.Line.String,
		"kodepos":      address.PostalCode.String,
	}
	err := mergo.Merge(&dest, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) synchronizeWhenAuthenticated(c *model.Customer, f *model.FinancialData, v *model.Verification,
	cr *model.Credential, a *model.Address) error {
	// Prepare Data
	requestBodySyncCustomer := map[string]interface{}{}
	err := prepareBodyCustomer(c, requestBodySyncCustomer)
	if err != nil {
		return errx.Trace(err)
	}
	err = prepareBodyFinancial(f, requestBodySyncCustomer)
	if err != nil {
		return errx.Trace(err)
	}
	err = prepareBodyVerification(v, requestBodySyncCustomer)
	if err != nil {
		return errx.Trace(err)
	}
	err = prepareBodyCredential(cr, requestBodySyncCustomer)
	if err != nil {
		return errx.Trace(err)
	}
	err = prepareBodyAddress(a, requestBodySyncCustomer)
	if err != nil {
		return errx.Trace(err)
	}

	// Execute synchronize to PDS API
	sync, err := s.SynchronizeCustomer(requestBodySyncCustomer)
	if err != nil {
		s.log.Error("error found when sync customer", logOption.Error(err))
		return errx.Trace(err)
	}

	// Handle status error
	if sync.Status != constant.ResponseSuccess {
		s.log.Error("Get Error from SynchronizeCustomer")
		return nhttp.InternalError.Trace(errx.Errorf(sync.Message))
	}

	return nil
}

func GetChannelByAgen(agen string) string {
	// Generalize agen
	agen = strings.ToLower(agen)

	if agen == constant.AgenAndroid {
		return constant.ChannelAndroid
	}

	if agen == constant.AgenMobile {
		return constant.ChannelMobile
	}

	if agen == constant.AgenWeb {
		return constant.ChannelWeb
	}

	return ""
}
