package customer

import (
	"database/sql"
	"errors"
	"github.com/nbs-go/errx"
	"github.com/nbs-go/nlogger"
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
		return "", err
	}

	// Get verification
	ver, err := s.repo.FindVerificationByEmailToken(payload.VerificationToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Error("failed to retrieve verification not found", nlogger.Error(err))
			return "", constant.ResourceNotFoundError.Trace()
		}
		s.log.Errorf("failed to retrieve verification. error: %v", nlogger.Error(err))
		return "", errx.Trace(err)
	}

	// Get customer
	customer, err := s.repo.FindCustomerByID(ver.CustomerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("failed to retrieve customer not found", nlogger.Error(err))
			return "", constant.ResourceNotFoundError.Trace()
		}
		s.log.Errorf("failed to retrieve customer. error: %v", nlogger.Error(err))

		return alreadyVerifiedView, errx.Trace(err)
	}

	// If email already verified
	if !ver.EmailVerifiedAt.Time.IsZero() {
		if err != nil {
			return "", err
		}
		return alreadyVerifiedView, nil
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
		s.log.Errorf("Error when update verification. %v", err)
		return alreadyVerifiedView, errx.Trace(err)
	}

	// Load view email verification success
	// TODO: Refactor to Controller
	emailVerifiedView, err := nval.TemplateFile("", "email_verification_success.html")
	if err != nil {
		return alreadyVerifiedView, err
	}

	return emailVerifiedView, nil
}
