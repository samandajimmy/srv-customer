package customer

import (
	logOption "github.com/nbs-go/nlogger/v2/option"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
)

func (s *Service) SynchronizeCustomer(reqBody map[string]interface{}) (*ResponsePdsAPI, error) {
	// Set header
	reqHeader := map[string]string{
		"Accept":       "*/*",
		"Content-Type": "application/x-www-form-urlencoded",
	}

	// Set payload
	postDataPayload := PostDataPayload{
		Path:   "/synchronize/customer",
		Data:   reqBody,
		Header: reqHeader,
	}

	resp, err := s.PdsPostData(postDataPayload)
	if err != nil {
		s.log.Error("error found when sync customer to PDS API", logOption.Error(err))
		return nil, err
	}

	// Set result
	result := resp

	return result, nil
}

func (s *Service) SynchronizePassword(payload dto.RegisterPayload) (*ResponsePdsAPI, error) {
	// Set payload
	reqBody := map[string]interface{}{
		"no_hp":    payload.PhoneNumber,
		"password": payload.Password,
	}

	// Set header
	reqHeader := map[string]string{
		"Accept":       "*/*",
		"Content-Type": "application/x-www-form-urlencoded",
	}

	// Set payload
	postDataPayload := PostDataPayload{
		Path:   "/synchronize/customer",
		Data:   reqBody,
		Header: reqHeader,
	}

	resp, err := s.PdsPostData(postDataPayload)
	if err != nil {
		s.log.Error("error found when sync customer password to PDS API", logOption.Error(err))
		return nil, err
	}

	// Set result
	result := resp

	return result, nil
}
