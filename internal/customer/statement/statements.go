package statement

import "repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"

type Statements struct {
	// Customer

	AccessSession   *AccessSession
	Address         *Address
	AuditLogin      *AuditLogin
	Credential      *Credential
	Customer        *Customer
	FinancialData   *FinancialData
	OTP             *OTP
	Verification    *Verification
	VerificationOTP *VerificationOTP

	// External

	User    *User
	UserPin *UserPin
}

func New(db *nsql.DatabaseContext) *Statements {
	return &Statements{
		AccessSession:   NewAccessSession(db),
		Address:         NewAddress(db),
		AuditLogin:      NewAuditLogin(db),
		Credential:      NewCredential(db),
		Customer:        NewCustomer(db),
		FinancialData:   NewFinancialData(db),
		OTP:             NewOTP(db),
		Verification:    NewVerification(db),
		VerificationOTP: NewVerificationOTP(db),
	}
}

func NewExternal(db *nsql.DatabaseContext) *Statements {
	return &Statements{
		User:    NewUser(db),
		UserPin: NewUserPin(db),
	}
}
