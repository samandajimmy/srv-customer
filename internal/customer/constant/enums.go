package constant

import "github.com/lestrrat-go/jwx/jwa"

/// Control Status

type ControlStatus = int8

const (
	Enabled = ControlStatus(iota + 1)
	Disabled
)

const (
	Prefix = "PDSAPI"
)

const (
	AgenAndroid = "android"
	AgenMobile  = "mobile"
	AgenWeb     = "web"
)

const (
	ChannelAndroid = "6017"
	ChannelMobile  = "6017"
	ChannelWeb     = "6014"
)

const (
	Warn2XWrongPassword = 2
	Warn4XWrongPassword = 4
	MinWrongPassword    = 3
	MaxWrongPassword    = 5
)

const (
	WIB  = "Asia/Jakarta"
	WITA = "Asia/Makassar"
	WIT  = "Asia/Jayapura"
)

const (
	TypeProfile = "profile"
)

const (
	Domicile     int64 = 1
	IdentityCard int64 = 2
)

// Identity

type IdentityType int64

const (
	IdentityTypeKTP      = IdentityType(10)
	IdentityTypeSIM      = IdentityType(11)
	IdentityTypePassport = IdentityType(12)
)

// Marriage Status

type MarriageStatus int64

const (
	Married = MarriageStatus(1)
	Single  = MarriageStatus(2)
	Widower = MarriageStatus(4)
)

// Request Type
const (
	RequestTypeRegister     = "register"
	RequestTypeBlockOneHour = "block-login-hour"
	RequestTypeBlockOneDay  = "block-login-day"
)

const (
	NotificationProviderFCM = iota + 1
)

const (
	CacheTokenSwitching = "token_switching"
	CacheTokenJWT       = "token_jwt"
	CacheGoldSavings    = "tabemas_listaccount"
)

const (
	RestSwitchingInvalidToken = "invalid_token"
)

const (
	JWTSignature = jwa.HS256
	JWTIssuer    = "https://www.pegadaian.co.id"
)

const (
	Unblocked = 0
	Blocked   = 1
)

const (
	WrongPIN    = 1
	WrongPIN2   = 2
	MaxWrongPIN = 3
)
