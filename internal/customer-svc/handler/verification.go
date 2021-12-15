package handler

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
)

func NewVerification(verificationService contract.VerificationService) *Verification {
	return &Verification{
		verificationService: verificationService,
	}
}

type Verification struct {
	router              *nhttp.Router
	verificationService contract.VerificationService
}

func (h *Verification) VerfiyEmail(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get query param
	var payload dto.VerificationPayload
	q := rx.URL.Query()
	payload.VerificationToken = q.Get("t")

	// Validate payload
	err := payload.Validate()
	if err != nil {
		log.Errorf("Unprocessable Entity. err: %v", err)
		return nhttp.View(), nil
	}

	// Call service
	resp, err := h.verificationService.VerifyEmailCustomer(payload)
	if err != nil {
		log.Errorf("Error when processing service. err: %v", err)
		return nhttp.View().SetData(resp), nil
	}

	return nhttp.View().SetData(resp), err
}
