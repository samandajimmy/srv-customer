package customer

import (
  "crypto/md5"
  "database/sql"
  "encoding/json"
  "errors"
  "fmt"
  "github.com/lestrrat-go/jwx/jwt"
  "github.com/nbs-go/nlogger"
  "regexp"
  "repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
  "repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/convert"
  "repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
  "repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
  "repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
  "repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ntime"
  "repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
  "strings"
  "time"
)

func (s *Service) Login(payload dto.LoginRequest) (*dto.LoginResponse, error) {
  ctx := s.ctx
  // Check if user exists
  t := time.Now()
  customer, err := s.repo.FindCustomerByEmail(payload.Email)
  if errors.Is(err, sql.ErrNoRows) {
	// If data not found on internal database check on external database.
	user, err := s.repoExternal.FindUserExternalByEmailOrPhone(payload.Email)
	if err != nil {
	  s.log.Debug("Phone or email is not registered")
	  return nil, s.responses.GetError("E_AUTH_10")
	}

	// sync data external to internal
	customer, err = s.syncExternalToInternal(user)
	if err != nil {
	  s.log.Error("error while sync data External to Internal", nlogger.Error(err), nlogger.Context(ctx))
	  return nil, ncore.TraceError("error sync data", err)
	}

  } else if err != nil {
	s.log.Error("failed to retrieve customer", nlogger.Error(err), nlogger.Context(ctx))
	return nil, ncore.TraceError("error find customer by email", err)
  }

  // Get credential customer
  credential, err := s.repo.FindCredentialByCustomerID(customer.Id)
  if err != nil {
	if errors.Is(err, sql.ErrNoRows) {
	  s.log.Error("failed to retrieve credential not found", nlogger.Error(err), nlogger.Context(ctx))
	  return nil, s.responses.GetError("E_AUTH_8")
	}
	s.log.Errorf("failed to retrieve credential", nlogger.Error(err), nlogger.Context(ctx))
	return nil, ncore.TraceError("error", err)
  }

  // Check if account isn't blocked
  blockedUntil := ntime.ChangeTimezone(credential.BlockedUntilAt.Time, constant.WIB)
  now := ntime.ChangeTimezone(t, constant.WIB)
  if credential.BlockedUntilAt.Valid != false && blockedUntil.After(now) {
	// Set response if blocked
	message := "Akun dikunci hingga %v karena gagal login %v kali. Hubungi call center jika ini bukan kamu"
	timeBlocked := ntime.ChangeTimezone(credential.BlockedUntilAt.Time, constant.WIB).Format("02-Jan-2006 15:04:05")
	setResponse := ncore.Success.HttpStatus(401).SetMessage(message, timeBlocked, credential.WrongPasswordCount)
	return nil, &setResponse
  }

  // Counter wrong password count
  passwordRequest := fmt.Sprintf("%x", md5.Sum([]byte(payload.Password)))
  if credential.Password != passwordRequest {
	err := s.handleWrongPassword(credential, customer)
	return nil, err
  }

  // get userRefId from external DB
  if customer.UserRefId.Valid == false {
	registerPayload := &dto.CustomerSynchronizeRequest{
	  Name:        customer.FullName,
	  Email:       customer.Email,
	  PhoneNumber: customer.Phone,
	  Password:    credential.Password,
	  FcmToken:    payload.FcmToken,
	}
	// sync data from customer service to PDS API
	resultSync, err := s.syncInternalToExternal(registerPayload)
	if err != nil {
	  s.log.Error("error found when sync to pds api", nlogger.Error(err), nlogger.Context(ctx))
	  return nil, ncore.TraceError("error found when sync to pds api", err)
	}
	// set userRefId
	customer.UserRefId = sql.NullString{
	  Valid:  true,
	  String: nval.ParseStringFallback(resultSync.UserAiid, ""),
	}
	// update customer
	err = s.repo.UpdateCustomerByPhone(customer)
	if err != nil {
	  s.log.Error("failed to update userRefId", nlogger.Error(err), nlogger.Context(ctx))
	  return nil, ncore.TraceError("update customer by phone", err)
	}
  }

  // Get token from cache
  var token string
  cacheTokenKey := fmt.Sprintf("%v:%v:%v", constant.Prefix, constant.CacheTokenJWT, customer.UserRefId.String)
  token, err = s.setTokenAuthentication(customer, payload.Agen, payload.Version, cacheTokenKey)
  if err != nil {
	s.log.Error("error found when get access token from cache", nlogger.Error(err), nlogger.Context(ctx))
	return nil, ncore.TraceError("failed to get token from cache", err)
  }

  // Check account is first login or not
  countAuditLog, err := s.repo.CountAuditLogin(customer.Id)
  if err != nil {
	s.log.Error("error found when count audit login", nlogger.Error(err), nlogger.Context(ctx))
	return nil, ncore.TraceError("error when count audit login", err)
  }

  // Set is first login is true or false.
  var isFirstLogin = true
  if countAuditLog > 0 {
	isFirstLogin = false
  }

  // Prepare to insert audit login
  auditLogin := model.AuditLogin{
	CustomerId:   customer.Id,
	ChannelId:    GetChannelByAgen(payload.Agen),
	DeviceId:     payload.DeviceId,
	IP:           payload.IP,
	Latitude:     payload.Latitude,
	Longitude:    payload.Longitude,
	Timestamp:    t.Format(time.RFC3339),
	Timezone:     payload.Timezone,
	Brand:        payload.Brand,
	OsVersion:    payload.OsVersion,
	Browser:      payload.Browser,
	UseBiometric: payload.UseBiometric,
	Status:       1,
	Metadata:     []byte("{}"),
	ItemMetadata: model.NewItemMetadata(
	  convert.ModifierDTOToModel(
		dto.Modifier{ID: "", Role: "", FullName: ""},
	  ),
	),
  }

  // Persist audit login
  err = s.repo.CreateAuditLogin(&auditLogin)
  if err != nil {
	s.log.Error("error found when create audit login", nlogger.Error(err), nlogger.Context(ctx))
	return nil, s.responses.GetError("E_AUTH_1")
  }

  // TODO get tabungan emas service / get tabungan emas from financial table

  // Check is force update password
  validatePassword := s.validatePassword(payload.Password)
  isForceUpdatePassword := false
  if validatePassword.IsValid != true {
	isForceUpdatePassword = true
  }

  // Get data address
  address, err := s.repo.FindAddressByCustomerId(customer.Id)
  if errors.Is(err, sql.ErrNoRows) {
	address = &model.Address{}
  } else if err != nil {
	s.log.Error("error found when get customer address", nlogger.Error(err), nlogger.Context(ctx))
	return nil, s.responses.GetError("E_AUTH_8")
  }

  // Get data verification
  verification, err := s.repo.FindVerificationByCustomerID(customer.Id)
  if errors.Is(err, sql.ErrNoRows) {
	verification = &model.Verification{}
  } else if err != nil {
	s.log.Errorf("error found when get data verification", nlogger.Error(err), nlogger.Context(ctx))
	return nil, s.responses.GetError("E_AUTH_8")
  }

  return composeLoginResponse(dto.LoginVO{
	Customer:              customer,
	Address:               address,
	Profile:               customer.Profile,
	Verification:          verification,
	IsFirstLogin:          isFirstLogin,
	IsForceUpdatePassword: isForceUpdatePassword,
	Token:                 token,
  })
}

func (s *Service) syncInternalToExternal(payload *dto.CustomerSynchronizeRequest) (*dto.UserVO, error) {
  // Get context
  ctx := s.ctx

  // call register pds api
  registerCustomer := dto.RegisterNewCustomer{
	Name:        payload.Name,
	Email:       payload.Email,
	PhoneNumber: payload.PhoneNumber,
	Password:    payload.Password,
	FcmToken:    payload.FcmToken,
  }
  // sync
  resp, err := s.SynchronizeCustomer(registerCustomer)
  if err != nil {
	log.Error("error found when sync data customer via API PDS", nlogger.Error(err), nlogger.Context(ctx))
	return nil, ncore.TraceError("failed to sync data via API PDS", err)
  }

  // handle status error
  if resp.Status != "success" {
	log.Error("Get Error from SynchronizeCustomer.", nlogger.Error(err))
	return nil, ncore.NewError(resp.Message)
  }

  // parsing response
  var user dto.UserVO
  err = json.Unmarshal(resp.Data, &user)
  if err != nil {
	log.Errorf("Cannot unmarshall data login pds. err: %v", err)
	return nil, ncore.TraceError("error", err)
  }

  // set result
  result := &user

  return result, nil
}

func (s *Service) syncExternalToInternal(user *model.User) (*model.Customer, error) {
  // prepare customer
  customer, err := convert.ModelUserToCustomer(user)
  if err != nil {
	s.log.Error("failed to convert to model customer", nlogger.Error(err), nlogger.Context(s.ctx))
	return nil, err
  }

  // Check has userPin or not
  userPin, err := s.repoExternal.FindUserPINByCustomerID(customer.Id)
  if errors.Is(err, sql.ErrNoRows) {
	userPin = &model.UserPin{}
  } else if err != nil {
	s.log.Error("failed retrieve user pin from external database", nlogger.Error(err), nlogger.Context(s.ctx))
	return nil, err
  }

  // Check user address on external database
  addressExternal, err := s.repoExternal.FindUserExternalAddressByCustomerID(user.UserAiid)
  if errors.Is(err, sql.ErrNoRows) {
	addressExternal = &model.AddressExternal{}
  } else if err != nil {
	s.log.Error("failed retrieve address from external database", nlogger.Error(err), nlogger.Context(s.ctx))
	return nil, err
  }

  // Prepare credential
  credential, err := convert.ModelUserToCredential(user, userPin)
  if err != nil {
	s.log.Error("failed convert to credential model", nlogger.Error(err), nlogger.Context(s.ctx))
	return nil, err
  }

  // Prepare financial data
  financialData, err := convert.ModelUserToFinancialData(user)
  if err != nil {
	s.log.Error("failed convert to financial data", nlogger.Error(err), nlogger.Context(s.ctx))
	return nil, err
  }

  // Prepare verification
  verification, err := convert.ModelUserToVerification(user)
  if err != nil {
	s.log.Error("failed convert to verification", nlogger.Error(err), nlogger.Context(s.ctx))
	return nil, err
  }

  // Prepare address
  address, err := convert.ModelUserToAddress(user, addressExternal)
  if err != nil {
	s.log.Error("failed convert address", nlogger.Error(err), nlogger.Context(s.ctx))
	return nil, err
  }

  // Persist customer data
  customerId, err := s.repo.CreateCustomer(customer)
  if err != nil {
	s.log.Error("failed to persist customer", nlogger.Error(err), nlogger.Context(s.ctx))
	return nil, err
  }

  // Set credential customer id to last inserted
  credential.CustomerId = customerId
  financialData.CustomerId = customerId
  verification.CustomerId = customerId
  address.CustomerId = customerId

  // persist credential
  err = s.repo.InsertOrUpdateCredential(credential)
  if err != nil {
	s.log.Error("failed persist to credential", nlogger.Error(err), nlogger.Context(s.ctx))
	return nil, err
  }

  // persist financial data
  err = s.repo.InsertOrUpdateFinancialData(financialData)
  if err != nil {
	s.log.Error("failed persist to financial data", nlogger.Error(err), nlogger.Context(s.ctx))
	return nil, err
  }

  // persist verification
  err = s.repo.InsertOrUpdateVerification(verification)
  if err != nil {
	s.log.Error("failed persist verification.", nlogger.Error(err), nlogger.Context(s.ctx))
	return nil, err
  }

  // persist address
  err = s.repo.InsertOrUpdateAddress(address)
  if err != nil {
	s.log.Error("failed persist address", nlogger.Error(err), nlogger.Context(s.ctx))
	return nil, err
  }

  customer, err = s.repo.FindCustomerByID(customerId)
  if err != nil {
	s.log.Error("failed to retrieve customer not found", nlogger.Error(err), nlogger.Context(s.ctx))
	return nil, s.responses.GetError("E_RES_1")
  }

  return customer, nil
}

func composeLoginResponse(data dto.LoginVO) (*dto.LoginResponse, error) {
  // Cast to model
  customer := data.Customer.(*model.Customer)
  profile := data.Profile.(*model.CustomerProfile)
  verification := data.Verification.(*model.Verification)
  address := data.Address.(*model.Address)

  return &dto.LoginResponse{
	User: &dto.LoginUserVO{
	  CustomerVO: dto.CustomerVO{
		ID:                        nval.ParseStringFallback(customer.Id, ""),
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
		IDProvinsi:                address.ProvinceId.String,
		IDKabupaten:               address.CityId.String,
		IDKecamatan:               address.DistrictId.String,
		IDKelurahan:               address.SubDistrictId.String,
		Kelurahan:                 address.DistrictName.String,
		Provinsi:                  address.ProvinceName.String,
		Kabupaten:                 address.CityName.String,
		Kecamatan:                 address.DistrictName.String,
		KodePos:                   address.PostalCode.String,
		NoHP:                      customer.Phone,
		FotoKTP:                   profile.IdentityPhotoFile,
		IsEmailVerified:           nval.ParseStringFallback(verification.EmailVerifiedStatus, "0"),
		Kewarganegaraan:           profile.Nationality,
		JenisIdentitas:            fmt.Sprintf("%v", customer.IdentityType),
		NoIdentitas:               customer.IdentityNumber,
		Avatar:                    "",
		TglExpiredIdentitas:       "",
		NoNPWP:                    profile.NPWPNumber,
		FotoNPWP:                  profile.NPWPPhotoFile,
		NoSid:                     customer.Sid,
		FotoSid:                   profile.SidPhotoFile,
		StatusKawin:               profile.MarriageStatus,
		Norek:                     "",
		Saldo:                     "",
		AktifasiTransFinansial:    nval.ParseStringFallback(verification.FinancialTransactionStatus, ""),
		IsDukcapilVerified:        nval.ParseStringFallback(verification.DukcapilVerifiedStatus, "0"),
		IsOpenTe:                  "",
		ReferralCode:              "",
		GoldCardApplicationNumber: "",
		GoldCardAccountNumber:     "",
		KodeCabang:                "",
		TabunganEmas: &dto.GoldSavingVO{
		  TotalSaldoBlokir:  "",
		  TotalSaldoSeluruh: "",
		  TotalSaldoEfektif: "",
		  PrimaryRekening:   nil,
		},
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
	s.log.Error("error found when get cache", nlogger.Error(err), nlogger.Context(s.ctx))
	return "", err
  }

  if accessToken == "" {
	// Generate access token
	newAccessToken := nval.Bin2Hex(nval.RandStringBytes(78))
	// Set token to cache
	cacheToken, err := s.CacheSetThenGet(cacheTokenKey, newAccessToken, s.config.ClientConfig.JWTExpiry)
	if err != nil {
	  s.log.Error("error found when set token to cache", nlogger.Error(err), nlogger.Context(s.ctx))
	  return "", err
	}
	// Set access token
	accessToken = cacheToken
  }

  channelId := GetChannelByAgen(agen)
  now := time.Now()

  // Generate JWT
  token, err := jwt.NewBuilder().
	Claim("id", customer.UserRefId.String).
	Claim("email", customer.Email).
	Claim("nama", customer.FullName).
	Claim("no_hp", customer.Phone).
	Claim("access_token", accessToken).
	Claim("agen", agen).
	Claim("channelId", channelId).
	Claim("version", version).
	IssuedAt(now).
	Expiration(now.Add(time.Second * time.Duration(s.config.ClientConfig.JWTExpiry))).
	Issuer(constant.JWTIssuer).
	Build()
  if err != nil {
	s.log.Error("error found when generate JWT", nlogger.Error(err), nlogger.Context(s.ctx))
	return "", err
  }

  jwtKey := s.config.ClientConfig.JWTKey
  jwtKeyBytes := []byte(jwtKey)

  // sign token
  signed, err := jwt.Sign(token, constant.JWTSignature, jwtKeyBytes)
  if err != nil {
	s.log.Error("failed to sign token", nlogger.Error(err), nlogger.Context(s.ctx))
	return "", ncore.TraceError("failed to sign token", err)
  }
  tokenString := string(signed)

  return tokenString, nil
}

func (s *Service) validatePassword(password string) *dto.ValidatePassword {
  var validation dto.ValidatePassword

  lowerCase, _ := regexp.Compile(`[a-z]+`)
  upperCase, _ := regexp.Compile(`[A-Z]+`)
  allNumber, _ := regexp.Compile(`[0-9]+`)

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

func (s *Service) handleWrongPassword(credential *model.Credential, customer *model.Customer) error {
  var resp error
  t := time.Now()

  // Unmarshalling metadata credential to get tryLoginAt
  var Metadata dto.MetadataCredential
  err := json.Unmarshal(credential.Metadata, &Metadata)
  if err != nil {
	s.log.Error("error found when unmarshal metadata credential", nlogger.Error(err), nlogger.Context(s.ctx))
	return ncore.TraceError("error found when unmarshal metadata credential", err)
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
	resp = s.responses.GetError("E_AUTH_6")
	credential.WrongPasswordCount = wrongCount
  case constant.Warn4XWrongPassword:
	resp = s.responses.GetError("E_AUTH_7")
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
	setResponse := ncore.Success.HttpStatus(401).SetMessage(message, timeBlocked, credential.WrongPasswordCount)
	resp = &setResponse

	// Send OTP To Phone Number
	request := dto.SendOTPRequest{
	  PhoneNumber: customer.Phone,
	  RequestType: constant.RequestTypeBlockOneHour,
	}
	_, err := s.SendOTP(request)
	if err != nil {
	  s.log.Debug("error found when sending otp block one hour", nlogger.Error(err), nlogger.Context(s.ctx))
	}

	// Send Notification Blocked Login One Hour
	// TODO Refactor
	err = s.SendNotificationBlock(dto.NotificationBlock{
	  Customer:     customer,
	  Message:      fmt.Sprintf(message, timeBlocked, credential.WrongPasswordCount),
	  LastTryLogin: ntime.NewTimeWIB(tryLoginAt).Format("02-Jan-2006 15:04:05"),
	})

	break
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
	setResponse := ncore.Success.HttpStatus(401).SetMessage(message, timeBlocked, credential.WrongPasswordCount)
	resp = &setResponse

	// Send OTP To Phone Number
	request := dto.SendOTPRequest{
	  PhoneNumber: customer.Phone,
	  RequestType: constant.RequestTypeBlockOneDay,
	}
	_, err := s.SendOTP(request)
	if err != nil {
	  s.log.Debug("Error when sending otp block one hour", nlogger.Error(err), nlogger.Context(s.ctx))
	}

	// Send Notification Blocked Login One Day
	err = s.SendNotificationBlock(dto.NotificationBlock{
	  Customer:     customer,
	  Message:      fmt.Sprintf(message, timeBlocked, credential.WrongPasswordCount),
	  LastTryLogin: ntime.NewTimeWIB(tryLoginAt).Format("02-Jan-2006 15:04:05"),
	})
	break
  default:
	resp = s.responses.GetError("E_AUTH_8")
	credential.WrongPasswordCount = wrongCount
	break
  }

  // Handle notification error
  if err != nil {
	s.log.Debug("Error when sending notification block", nlogger.Error(err), nlogger.Context(s.ctx))
  }

  // Set trying login at to metadata
  var Format dto.MetadataCredential
  Format.TryLoginAt = t.Format(time.RFC3339)
  Format.PinCreatedAt = Metadata.PinCreatedAt
  Format.PinBlockedAt = Metadata.PinBlockedAt
  MetadataCredential, _ := json.Marshal(&Format)
  credential.Metadata = MetadataCredential

  err = s.repo.UpdateCredential(credential)
  if err != nil {
	s.log.Error("error when update credential.", nlogger.Error(err))
	return ncore.TraceError("error when update credential.", err)
  }

  return resp
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
