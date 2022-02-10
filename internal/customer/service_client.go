package customer

import (
	"net/http"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nclient"
)

func (s *Service) ClientPostData(endpoint string, body map[string]interface{}, header map[string]string) (*http.Response, error) {

	client := nclient.NewNucleoClient(
		s.config.CoreOauthUsername,
		s.config.CoreClientId,
		s.config.CoreApiUrl,
	)

	data, err := client.PostData(endpoint, body, header)
	if err != nil {
		return nil, err
	}
	return data, nil
}
