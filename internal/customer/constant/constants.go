package constant

const (
	AlphaNumUpperCaseRandomSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	AlphaNumRandomSet          = AlphaNumUpperCaseRandomSet + "abcdefghijklmnopqrstuvwxyz0123456789"
	SlugRandomSet              = AlphaNumRandomSet + "_-"
)

const (
	SubjectIDHeader   = "x-subject-id"
	SubjectNameHeader = "x-subject-name"
	SubjectRoleHeader = "x-subject-role"
)

type AssetType = int

const (
	_ = AssetType(iota)
	AssetAvatarProfile
	AssetNPWP
	AssetKTP
)

var AssetDirs = map[AssetType]string{
	AssetAvatarProfile: "user/avatar",
	AssetNPWP:          "user/npwp",
	AssetKTP:           "user/ktp",
}

const (
	KeyUserFile = "userfile"
)

// Context Key

type ContextKey uint8

const (
	_ ContextKey = iota
	RequestIDContextKey
	SubjectContextKey
	UserRefIDContextKey
)
