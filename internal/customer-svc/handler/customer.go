package handler

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
)

func NewCustomer(customerService contract.CustomerService) *Customer {
	return &Customer{
		customerService: customerService,
	}
}

type Customer struct {
	router          *nhttp.Router
	customerService contract.CustomerService
}

func (h *Customer) PostLogin(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get Payload
	var payload dto.LoginRequest
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Errorf("Error when parse json body. err: %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = payload.Validate()
	if err != nil {
		log.Errorf("Bad request. err: %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Call service
	resp, err := h.customerService.Login(payload)
	if err != nil {
		log.Errorf("Error when processing service. err: %v", err)
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (h *Customer) PostRegister(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get Payload
	var payload dto.RegisterNewCustomer
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Errorf("Error when parse json body. err: %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = payload.Validate()
	if err != nil {
		log.Errorf("Bad request. err: %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Call service
	resp, err := h.customerService.Register(payload)
	if err != nil {
		log.Errorf("Error when processing service. err: %v", err)
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (h *Customer) SendOTP(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get Payload
	var payload dto.RegisterStepOne
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Errorf("Error when parse json body. err: %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = payload.Validate()
	if err != nil {
		log.Errorf("Bad request. err: %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Call service
	resp, err := h.customerService.RegisterStepOne(payload)
	if err != nil {
		log.Errorf("Error when processing service. err: %v", err)
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (h *Customer) VerifyOTP(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get Payload
	var payload dto.RegisterStepTwo
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Errorf("Error when parse json body. err: %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = payload.Validate()
	if err != nil {
		log.Errorf("Bad request. err: %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Call service
	resp, err := h.customerService.RegisterStepTwo(payload)
	if err != nil {
		log.Errorf("Error when processing service. err: %v", err)
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (h *Customer) ResendOTP(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get Payload
	var payload dto.RegisterResendOTP
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Errorf("Error when parse json body. err: %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = payload.Validate()
	if err != nil {
		log.Errorf("Bad request. err: %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Call service
	resp, err := h.customerService.RegisterResendOTP(payload)
	if err != nil {
		log.Errorf("Error when processing service. err: %v", err)
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}
