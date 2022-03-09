package dto

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	JTI         string `json:"jti"`
}

type SendOTPRequest struct {
	PhoneNumber string
	RequestType string
}

type VerifyOTPRequest struct {
	PhoneNumber string
	RequestType string
	Token       string
}
