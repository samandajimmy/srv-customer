package constant

const (
	AlphaNumUpperCaseRandomSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	AlphaNumRandomSet          = AlphaNumUpperCaseRandomSet + "abcdefghijklmnopqrstuvwxyz0123456789"
	SlugRandomSet              = AlphaNumRandomSet + "_-"
)

const (
	SubjectKey   = "subject"
	BuildHashKey = "build_hash"
	UserRefID    = "user_ref_id"
)

const (
	SubjectIDHeader   = "x-subject-id"
	SubjectNameHeader = "x-subject-name"
	SubjectRoleHeader = "x-subject-role"
)

const (
	_ = iota
	AssetAvatarProfile
	AssetNPWP
	AssetKTP
)

var AssetDirs = map[int]string{
	AssetAvatarProfile: "user/avatar",
	AssetNPWP:          "user/npwp",
	AssetKTP:           "user/ktp",
}
