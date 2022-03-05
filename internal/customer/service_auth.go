package customer

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
)

func (s *Service) ValidateClient(payload dto.ClientCredential) error {
	if payload.ClientID != s.config.ClientID || payload.ClientSecret != s.config.ClientSecret {
		return nhttp.UnauthorizedError.Trace()
	}

	return nil
}
