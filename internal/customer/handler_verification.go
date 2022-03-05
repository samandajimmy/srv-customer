package customer

import (
	"encoding/json"
	"github.com/nbs-go/nlogger"
	"net/http"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
)

type Verification struct {
	*Handler
}

func NewVerification(h *Handler) *Verification {
	return &Verification{h}
}

// VerifyEmail TODO: Refactor using http.Handler
func (h *Verification) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	// Get query param
	var payload dto.VerificationPayload
	q := r.URL.Query()
	payload.VerificationToken = q.Get("t")

	// Validate payload
	err := payload.Validate()
	if err != nil {
		log.Errorf("Invalid payload. err: %v", err)
		h.renderError(w, 400, err)
		return
	}

	// Init service
	svc := h.NewService(r.Context())
	defer svc.Close()

	// Call service
	resp, err := svc.VerifyEmailCustomer(payload)
	if err != nil {
		log.Errorf("Error when processing service. err: %v", err)
		return
	}

	// Render response
	h.renderSuccess(w, resp)
	return
}

// renderError TODO: Render default error view for HTML
func (h *Verification) renderError(w http.ResponseWriter, statusCode int, err error) {
	// Write header in
	w.Header().Add(nhttp.ContentTypeHeader, nhttp.ContentTypeJSON)

	// Write header
	w.WriteHeader(statusCode)

	// Write error in JSON
	err = json.NewEncoder(w).Encode(err)
	if err != nil {
		log.Errorf("failed to write response to json ( payload = %+v )", err)
	}
}

func (h *Verification) renderSuccess(w http.ResponseWriter, htmlBody string) {
	w.Header().Add(nhttp.ContentTypeHeader, "text/html")
	_, err := w.Write([]byte(htmlBody))
	if err != nil {
		log.Errorf("failed to write response", nlogger.Error(err))
	}
}
