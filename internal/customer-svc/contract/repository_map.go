package contract

type RepositoryMap struct {
	Customer        CustomerRepository
	VerificationOTP VerificationOTPRepository
	OTP             OTPRepository
}
