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

type UpdatePinPayload struct {
	PIN                string `json:"pin"`
	NewPIN             string `json:"new_pin"`
	NewPINConfirmation string `json:"new_pin_confirmation"`
	UserRefID          string `json:"userRefId"`
	CheckPIN           bool   `json:"check_pin"`
}

type UpdatePinResult struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}
