package constant

const (
	AlphaNumUpperCaseRandomSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	AlphaNumRandomSet          = AlphaNumUpperCaseRandomSet + "abcdefghijklmnopqrstuvwxyz0123456789"
	SlugRandomSet              = AlphaNumRandomSet + "_-"
)

const (
	SubjectKey   = "subject"
	BuildHashKey = "build_hash"
)

const (
	SubjectIDHeader   = "x-subject-id"
	SubjectNameHeader = "x-subject-name"
	SubjectRoleHeader = "x-subject-role"
)
