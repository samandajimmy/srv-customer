package dto

type EmailVerification struct {
	FullName        string
	VerificationUrl string
	Email           string
}

type EmailBlock struct {
	Title        string
	Text         string
	Message      string
	LastTryLogin string
	BaseUrl      string
}
