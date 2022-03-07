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

type CheckOTPPinPayload struct {
	OTP       string `json:"otp"`
	UserRefID string `json:"userRefId"`
}

type RestSwitchingOTPPinCreate struct {
	Cif  string `json:"cif"`
	OTP  string `json:"otp"`
	NoHp string `json:"no_hp"`
}

type PostCreatePinPayload struct {
	NewPIN             string `json:"new_pin"`
	NewPINConfirmation string `json:"new_pin_confirmation"`
	UserRefID          string `json:"userRefId"`
	OTP                string `json:"otp"`
}

type RestSwitchingOTPForgetPin struct {
	Cif         string `json:"cif"`
	Flag        string `json:"flag"`
	NoHp        string `json:"noHp"`
	NoRek       string `json:"noRek"`
	RequestType string `json:"requestType"`
	OTP         string `json:"otp"`
}
