package customer

import (
	"database/sql"
	"github.com/nbs-go/errx"
	logOption "github.com/nbs-go/nlogger/v2/option"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ntime"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
	"time"
)

func (s *Service) VerifyEmailCustomer(payload dto.VerificationPayload) (string, error) {
	// Load view already verified
	alreadyVerifiedView, err := nval.TemplateFile("", "already_verified.html")
	if err != nil {
		return "Failed to load view", err
	}

	// Get verification
	ver, err := s.repo.FindVerificationByEmailToken(payload.VerificationToken)
	if err != nil {
		s.log.Error("failed to retrieve verification")
		err = handleErrorRepository(err, constant.ResourceNotFoundError.Trace())
		return alreadyVerifiedView, errx.Trace(err)
	}

	// If email already verified
	if !ver.EmailVerifiedAt.Time.IsZero() {
		return alreadyVerifiedView, nil
	}

	// Get customer
	customer, err := s.repo.FindCustomerByID(ver.CustomerID)
	if err != nil {
		log.Error("failed to retrieve customer not found", logOption.Error(err))
		err = handleErrorRepository(err, constant.ResourceNotFoundError.Trace())
		return alreadyVerifiedView, errx.Trace(err)
	}

	// Set new value verification
	ver.EmailVerifiedStatus = 1
	ver.EmailVerifiedAt = sql.NullTime{
		Time:  ntime.NewTimeWIB(time.Now()),
		Valid: true,
	}
	ver.EmailVerificationToken = ""
	ver.UpdatedAt = time.Now()
	ver.ModifiedBy = &model.Modifier{
		ID:       nval.ParseStringFallback(ver.CustomerID, "0"),
		Role:     constant.UserModifierRole,
		FullName: customer.FullName,
	}
	ver.Version++

	// Update verification
	err = s.repo.UpdateVerificationByCustomerID(ver)
	if err != nil {
		s.log.Error("Error when update verification", logOption.Error(err))
		return alreadyVerifiedView, errx.Trace(err)
	}

	// Synchronize Verification PDS
	err = s.HandleSynchronizeVerification(customer, ver)
	if err != nil {
		s.log.Error("Error when synchronize data verification")
		return alreadyVerifiedView,errx.Trace(err)
	}

	// Load view email verification success
	emailVerifiedView, err := nval.TemplateFile("", "email_verification_success.html")
	if err != nil {
		return alreadyVerifiedView, err
	}

	return emailVerifiedView, nil
}
