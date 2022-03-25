package customer

import (
	"github.com/nbs-go/errx"
	logOption "github.com/nbs-go/nlogger/v2/option"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
)

// Endpoint POST /customer/activation
func (s *Service) customerActivation(payload *dto.RestSwitchingOTPPinCreate) (*ResponseSwitchingSuccess, error) {
	// Set payload
	reqBody := map[string]interface{}{
		"cif":      payload.Cif,
		"token":    payload.OTP,
		"username": payload.NoHp,
	}

	sp := PostDataPayload{
		Path: "/customer/activation",
		Data: reqBody,
	}

	data, err := s.RestSwitchingPostData(sp)
	if err != nil {
		s.log.Error("error found when get otp pin activation", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	return data, nil
}

// Endpoint POST /otp/validate
func (s *Service) otpValidate(payload *dto.RestSwitchingOTPForgetPin) (*ResponseSwitchingSuccess, error) {
	// Set payload
	reqBody := map[string]interface{}{
		"cif":         payload.Cif,
		"flag":        payload.Flag,
		"noHp":        payload.NoHp,
		"norek":       payload.NoHp,
		"requestType": payload.RequestType,
		"token":       payload.OTP,
	}

	sp := PostDataPayload{
		Path: "/otp/validate",
		Data: reqBody,
	}

	data, err := s.RestSwitchingPostData(sp)
	if err != nil {
		s.log.Error("error found when get reset pin otp validate", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	return data, nil
}

func (s *Service) ChangePhoneNumber(payload dto.ChangePhoneNumberRequestCore) (*ResponseSwitchingSuccess, error) {
	// Set payload
	reqBody := map[string]interface{}{
		"noHp":         payload.CurrentPhoneNumber,
		"noHpNew":      payload.NewPhoneNumber,
		"clientId":     "9997",
		"channelId":    6017,
		"cif":          payload.Cif,
		"ibuKandung":   payload.MaidenName,
		"namaNasabah":  payload.FullName,
		"tanggalLahir": payload.DateOfBirth,
	}

	sp := PostDataPayload{
		Path: "/customer/phonenumber",
		Data: reqBody,
	}

	// Send OTP Rest Switching
	r, err := s.RestSwitchingPostData(sp)
	if err != nil {
		s.log.Error("Error when change phone number request", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	// Set result
	result := r

	return result, nil
}

// Endpoint POST /customer/inquiry
func (s *Service) CheckCIF(cif string) (*ResponseSwitchingSuccess, error) {
	// Check CIF
	reqBody := map[string]interface{}{
		"cif": cif,
	}

	sp := PostDataPayload{
		Path: "/customer/inquiry",
		Data: reqBody,
	}

	data, err := s.RestSwitchingPostData(sp)
	if err != nil {
		s.log.Error("error found when get gold savings", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	return data, nil
}
