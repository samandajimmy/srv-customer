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
