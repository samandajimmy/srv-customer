package dto

type EmailVerification struct {
	FullName        string
	VerificationURL string
	Email           string
}

type EmailBlock struct {
	Title        string
	Text         string
	Message      string
	LastTryLogin string
	BaseURL      string
}

type EmailChangePayload struct {
	Email     string `json:"email"`
	UserRefID string `json:"userRefID"`
}

type EmailWelcome struct {
	Title    string
	FullName string
	Email    string
	BaseURL  string
}
