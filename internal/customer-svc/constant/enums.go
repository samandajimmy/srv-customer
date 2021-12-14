package constant

/// Control Status

type ControlStatus = int8

const (
	Enabled = ControlStatus(iota + 1)
	Disabled
)

/// Asset Type

type AssetType = int8

const (
	ThumbnailAsset = AssetType(iota + 1)
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
	Domicile     = "DOMICILE"
	IdentityCard = "IDENTITY_CARD"
)
