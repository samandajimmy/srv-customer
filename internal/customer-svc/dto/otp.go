package dto

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type SendOTPRequest struct {
	PhoneNumber string
	RequestType string
}
