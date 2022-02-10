package customer

import (
	"database/sql"
	"errors"
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/convert"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ntime"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
	"time"
)

func (s *Service) VerifyEmailCustomer(payload dto.VerificationPayload) (string, error) {
	// Load view already verified
	alreadyVerifiedView, err := nval.TemplateFile("", "already_verified.html")

	// Get verification
	ver, err := s.repo.FindVerificationByEmailToken(payload.VerificationToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Error("failed to retrieve verification not found", nlogger.Error(err))
			return alreadyVerifiedView, s.responses.GetError("E_RES_1")
		}
		s.log.Errorf("failed to retrieve verification. error: %v", nlogger.Error(err))
		return alreadyVerifiedView, ncore.TraceError("error", err)
	}

	// Get customer
	customer, err := s.repo.FindCustomerByID(ver.CustomerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("failed to retrieve customer not found", nlogger.Error(err))
			return alreadyVerifiedView, s.responses.GetError("E_RES_1")
		}
		s.log.Errorf("failed to retrieve customer. error: %v", nlogger.Error(err))

		return alreadyVerifiedView, ncore.TraceError("error", err)
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
	modifiedBy := convert.ModifierDTOToModel(dto.Modifier{
		ID:       nval.ParseStringFallback(ver.CustomerId, "0"),
		Role:     constant.UserModifierRole,
		FullName: customer.FullName,
	})
	ver.UpdatedAt = time.Now()
	ver.ModifiedBy = &modifiedBy
	ver.Version += 1

	// Update verification
	err = s.repo.UpdateVerificationByCustomerID(ver)
	if err != nil {
		s.log.Errorf("Error when update verification. %v", err)
		return alreadyVerifiedView, ncore.TraceError("error when update verification", err)
	}

	// Load view email verification success
	emailVerifiedView, err := nval.TemplateFile("", "email_verification_success.html")
	if err != nil {
		return alreadyVerifiedView, err
	}

	return emailVerifiedView, nil
}
