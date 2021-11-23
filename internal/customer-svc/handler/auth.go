package handler

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
)

func NewAuth(authService contract.AuthService) *Auth {
	a := Auth{
		authService: authService,
	}
	return &a
}

type Auth struct {
	router      *nhttp.Router
	authService contract.AuthService
}

func (h *Auth) ValidateClient(r *nhttp.Request) (*nhttp.Response, error) {
	// Extract Basic Auth
	clientID, clientSecret, err := nhttp.ExtractBasicAuth(r.Request)
	if err != nil {
		return nil, err
	}

	// Authentication app
	err = h.authService.ValidateClient(dto.ClientCredential{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})
	if err != nil {
		log.Errorf("an error occurred on authentication app. Error => %s", err)
		return nil, err
	}

	return nhttp.Continue(), nil
}
