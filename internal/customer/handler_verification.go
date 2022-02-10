package customer

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
)

type Verification struct {
	*Handler
}

func NewVerification(h *Handler) *Verification {
	return &Verification{h}
}

func (h *Verification) VerifyEmail(rx *nhttp.Request) (*nhttp.Response, error) {
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

	// Init service
	svc := h.NewService(rx.Context())
	defer svc.Close()

	// Call service
	resp, err := svc.VerifyEmailCustomer(payload)
	if err != nil {
		log.Errorf("Error when processing service. err: %v", err)
		return nhttp.View().SetData(resp), nil
	}

	return nhttp.View().SetData(resp), err
}
