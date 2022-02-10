package ncore

const (
	NumCharSet           = "0123456789"
	AlphaCharSet         = "abcdefghijklmnopqrstuvwxyz"
	AlphaUpperCharSet    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	AlphaNumCharSet      = AlphaCharSet + NumCharSet
	AlphaNumUpperCharSet = AlphaUpperCharSet + NumCharSet
	AlphaNumRandomSet    = AlphaCharSet + AlphaUpperCharSet + NumCharSet
	SlugRandomSet        = AlphaNumRandomSet + "-"
)
