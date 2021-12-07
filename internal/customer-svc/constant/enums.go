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
	PREFIX = "MSCUSTOMER"
)

const (
	AGEN_ANDROID = "android"
	AGEN_MOBILE  = "mobile"
	AGEN_WEB     = "web"
)

const (
	CHANNEL_ANDROID = "6017"
	CHANNEL_MOBILE  = "6017"
	CHANNEL_WEB     = "6014"
)

const (
	WARN_2X_WRONG_PASSWORD = 2
	WARN_4X_WRONG_PASSWORD = 4
	MIN_WRONG_PASSWORD     = 3
	MAX_WRONG_PASSWORD     = 5
)

const (
	WIB  = "Asia/Jakarta"
	WITA = "Asia/Makassar"
	WIT  = "Asia/Jayapura"
)
