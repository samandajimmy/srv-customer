package customer

import "repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"

func (s *Service) ValidateClient(payload dto.ClientCredential) error {
	if payload.ClientID != s.config.ClientID || payload.ClientSecret != s.config.ClientSecret {
		return s.responses.GetError("E_AUTH_1")
	}

	return nil
}
