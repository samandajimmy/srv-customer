package service

import (
	"database/sql"
	"errors"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/convert"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ntime"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
	"time"
)

type Verification struct {
	verificationRepo    contract.VerificationRepository
	customerRepo        contract.CustomerRepository
	cacheService        contract.CacheService
	notificationService contract.NotificationService
	emailConfig         contract.EmailConfig
	clientConfig        contract.ClientConfig
	response            *ncore.ResponseMap
	httpBaseUrl         string
}

func (s *Verification) HasInitialized() bool {
	return true
}

func (s *Verification) Init(app *contract.PdsApp) error {
	s.verificationRepo = app.Repositories.Verification
	s.customerRepo = app.Repositories.Customer
	s.cacheService = app.Services.Cache
	s.clientConfig = app.Config.Client
	s.notificationService = app.Services.Notification
	s.response = app.Responses
	s.httpBaseUrl = app.Config.Server.GetHttpBaseUrl()
	s.emailConfig = app.Config.Email
	return nil
}

func (s *Verification) VerifyEmailCustomer(payload dto.VerificationPayload) (string, error) {
	// Load view already verified
	alreadyVerfiedView, err := nval.TemplateFile("", "already_verified.html")

	// Get verification
	ver, err := s.verificationRepo.FindByEmailToken(payload.VerificationToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("failed to retrieve verfication not found", nlogger.Error(err))
			return alreadyVerfiedView, s.response.GetError("E_RES_1")
		}
		log.Errorf("failed to retrieve verfication. error: %v", err)
		return alreadyVerfiedView, ncore.TraceError(err)
	}

	// Get customer
	customer, err := s.customerRepo.FindById(ver.CustomerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("failed to retrieve customer not found", nlogger.Error(err))
			return alreadyVerfiedView, s.response.GetError("E_RES_1")
		}
		log.Errorf("failed to retrieve customer. error: %v", err)
		return alreadyVerfiedView, ncore.TraceError(err)
	}

	// If email already verified
	if !ver.EmailVerifiedAt.Time.IsZero() {
		if err != nil {
			return "", err
		}
		return alreadyVerfiedView, nil
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
	err = s.verificationRepo.UpdateByCustomerID(ver)
	if err != nil {
		log.Errorf("Error when update verification. %v", err)
		return alreadyVerfiedView, ncore.TraceError(err)
	}

	// Load view email verification success
	emailVerifiedView, err := nval.TemplateFile("", "email_verification_success.html")
	if err != nil {
		return alreadyVerfiedView, err
	}

	return emailVerifiedView, nil
}
