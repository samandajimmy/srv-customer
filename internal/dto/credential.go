package dto

type MetadataCredential struct {
	TryLoginAt   string `json:"tryLoginAt"`
	PinCreatedAt string `json:"pinCreatedAt"`
	PinBlockedAt string `json:"pinBlockedAt"`
}

type ValidatePinPayload struct {
	NewPin string `json:"new_pin"`
}

type CheckPinPayload struct {
	Pin       string `json:"pin"`
	UserRefID string `json:"userRefId"`
}
