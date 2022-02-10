package customer

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nvalidate"
)

type Customer struct {
	*Handler
}

func NewCustomer(h *Handler) *Customer {
	return &Customer{h}
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
		log.Errorf("Unprocessable Entity. err: %v", err)
		data := nvalidate.Message(err.Error())
		return nhttp.UnprocessableEntity(data), nil
	}

	// Init service
	svc := h.NewService(rx.Context())
	defer svc.Close()

	// Call service
	resp, err := svc.Login(payload)
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
		log.Error("Error when parse json body.", nlogger.Error(err))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = payload.Validate()
	if err != nil {
		log.Error("Unprocessable Entity.", nlogger.Error(err))
		data := nvalidate.Message(err.Error())
		return nhttp.UnprocessableEntity(data), nil
	}

	// Init service
	svc := h.NewService(rx.Context())
	defer svc.Close()

	// Check is force update password
	validatePassword := svc.ValidatePassword(payload.Password)
	if validatePassword.IsValid == false {
		err := errors.New(fmt.Sprintf("password: %s.", validatePassword.Message))
		data := nvalidate.Message(err.Error())
		return nhttp.UnprocessableEntity(data), nil
	}

	// Call service
	resp, err := svc.Register(payload)
	if err != nil {
		log.Errorf("Error when processing service.", nlogger.Error(err))
		return nil, ncore.TraceError("error", err)
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
		log.Errorf("Unprocessable Entity. err: %v", err)
		data := nvalidate.Message(err.Error())
		return nhttp.UnprocessableEntity(data), nil
	}

	// Init service
	svc := h.NewService(rx.Context())
	defer svc.Close()

	// Call service
	resp, err := svc.RegisterStepOne(payload)
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
		log.Errorf("Unprocessable Entity. err: %v", err)
		data := nvalidate.Message(err.Error())
		return nhttp.UnprocessableEntity(data), nil
	}

	// Init service
	svc := h.NewService(rx.Context())
	defer svc.Close()

	// Call service
	resp, err := svc.RegisterStepTwo(payload)
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
		log.Errorf("Unprocessable Entity. err: %v", err)
		data := nvalidate.Message(err.Error())
		return nhttp.UnprocessableEntity(data), nil
	}

	// Init service
	svc := h.NewService(rx.Context())
	defer svc.Close()

	// Call service
	resp, err := svc.RegisterResendOTP(payload)
	if err != nil {
		log.Errorf("Error when processing service. err: %v", err)
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (h *Customer) GetProfile(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get xid
	id := mux.Vars(rx.Request)["id"]
	if id == "" {
		err := errors.New("id is not found on params")
		log.Errorf("id is not found on params. err: %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	return nil, nil
}
