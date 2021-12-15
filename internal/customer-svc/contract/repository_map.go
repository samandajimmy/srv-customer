package contract

type RepositoryMap struct {
	Customer        CustomerRepository
	VerificationOTP VerificationOTPRepository
	OTP             OTPRepository
	Credential      CredentialRepository
	FinancialData   FinancialDataRepository
	AccessSession   AccessSessionRepository
	AuditLogin      AuditLoginRepository
	Verification    VerificationRepository
	Address         AddressRepository
	UserExternal    UserExternalRepository
	UserPinExternal UserPinExternalRepository
}
