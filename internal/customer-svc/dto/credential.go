package dto

type MetadataCredential struct {
	TryLoginAt   string `json:"tryLoginAt"`
	PinCreatedAt string `json:"pinCreatedAt"`
	PinBlockedAt string `json:"pinBlockedAt"`
}
