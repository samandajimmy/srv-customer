package dto

import "repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"

type MetadataCredential struct {
	TryLoginAt   string `json:"tryLoginAt"`
	PinCreatedAt string `json:"pinCreatedAt"`
	PinBlockedAt string `json:"pinBlockedAt"`
}

type ValidatePinPayload struct {
	NewPin string `json:"new_pin"`
}

type CheckPinPayload struct {
	UserRefID string `json:"-"`
	Pin       string `json:"pin"`
	CheckPIN  bool   `json:"checkPIN"`
}

type UpdatePinPayload struct {
	UserRefID          string `json:"-"`
	PIN                string `json:"pin"`
	NewPIN             string `json:"new_pin"`
	NewPINConfirmation string `json:"new_pin_confirmation"`
	CheckPIN           bool   `json:"check_pin"`
}

type UpdatePinResult struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type CheckOTPPinPayload struct {
	UserRefID string `json:"-"`
	OTP       string `json:"otp"`
}

type RestSwitchingOTPPinCreate struct {
	Cif  string `json:"cif"`
	OTP  string `json:"otp"`
	NoHp string `json:"no_hp"`
}

type PostCreatePinPayload struct {
	UserRefID          string `json:"-"`
	NewPIN             string `json:"new_pin"`
	NewPINConfirmation string `json:"new_pin_confirmation"`
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

type ForgetPinPayload struct {
	UserRefID          string `json:"-"`
	OTP                string `json:"otp"`
	NewPIN             string `json:"new_pin"`
	NewPINConfirmation string `json:"new_pin_confirmation"`
}

type UpdateSmartAccessPayload struct {
	UserRefID    string                 `json:"-"`
	DeviceID     string                 `json:"device_id"`
	UseBiometric constant.ControlStatus `json:"use_biometric"`
}
