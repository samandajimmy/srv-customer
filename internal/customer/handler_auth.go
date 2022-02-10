package customer

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
)

type Auth struct {
	*Handler
}

func NewAuth(h *Handler) *Auth {
	return &Auth{h}
}

func (h *Auth) ValidateClient(rx *nhttp.Request) (*nhttp.Response, error) {
	// Extract Basic Auth
	clientID, clientSecret, err := nhttp.ExtractBasicAuth(rx.Request)
	if err != nil {
		return nil, err
	}

	// Payload
	payload := dto.ClientCredential{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	// Init service
	svc := h.NewService(rx.Context())

	// Authentication app
	err = svc.ValidateClient(payload)
	if err != nil {
		log.Errorf("an error occurred on authentication app. Error => %s", err)
		return nil, err
	}

	return nhttp.Continue(), nil
}
