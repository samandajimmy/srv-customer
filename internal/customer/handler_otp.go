package customer

import (
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nvalidate"
)

func (h *Customer) SendOTP(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get Payload
	var payload dto.RegisterStepOne
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", nlogger.Error(err), nlogger.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = payload.Validate()
	if err != nil {
		log.Error("unprocessable Entity", nlogger.Error(err), nlogger.Context(ctx))
		data := nvalidate.Message(err.Error())
		return nhttp.UnprocessableEntity(data), nil
	}

	// Init service
	svc := h.NewService(rx.Context())
	defer svc.Close()

	// Call service
	resp, err := svc.RegisterStepOne(payload)
	if err != nil {
		log.Error("error when processing service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (h *Customer) VerifyOTP(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get Payload
	var payload dto.RegisterStepTwo
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", nlogger.Error(err), nlogger.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = payload.Validate()
	if err != nil {
		log.Error("unprocessable entity", nlogger.Error(err), nlogger.Context(ctx))
		data := nvalidate.Message(err.Error())
		return nhttp.UnprocessableEntity(data), nil
	}

	// Init service
	svc := h.NewService(rx.Context())
	defer svc.Close()

	// Call service
	resp, err := svc.RegisterStepTwo(payload)
	if err != nil {
		log.Error("error when processing service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (h *Customer) ResendOTP(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get Payload
	var payload dto.RegisterResendOTP
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", nlogger.Error(err), nlogger.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = payload.Validate()
	if err != nil {
		log.Error("unprocessable entity", nlogger.Error(err), nlogger.Context(ctx))
		data := nvalidate.Message(err.Error())
		return nhttp.UnprocessableEntity(data), nil
	}

	// Init service
	svc := h.NewService(rx.Context())
	defer svc.Close()

	// Call service
	resp, err := svc.RegisterResendOTP(payload)
	if err != nil {
		log.Error("error found when call register resend otp service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}
