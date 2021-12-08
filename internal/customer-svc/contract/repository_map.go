package contract

type RepositoryMap struct {
	Customer        CustomerRepository
	VerificationOTP VerificationOTPRepository
	OTP             OTPRepository
	Credential      CredentialRepository
	AccessSession   AccessSessionRepository
	AuditLogin      AuditLoginRepository
	Verification    VerificationRepository
}
