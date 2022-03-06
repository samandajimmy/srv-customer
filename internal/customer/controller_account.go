package customer

import (
	"encoding/json"
	"fmt"
	"github.com/nbs-go/errx"
	"github.com/nbs-go/nlogger/v2"
	"net/http"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
)

type AccountController struct {
	*Handler
}

func NewAccountController(h *Handler) *AccountController {
	return &AccountController{
		h,
	}
}

func (c *AccountController) PostLogin(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get Payload
	var payload dto.LoginRequest
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", nlogger.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = payload.Validate()
	if err != nil {
		log.Error("unprocessable entity")
		return nil, nhttp.BadRequestError.Trace(errx.Source(err))
	}

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Call service
	resp, err := svc.Login(payload)
	if err != nil {
		log.Error("error found when call service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (c *AccountController) GetVerifyEmail(w http.ResponseWriter, r *http.Request) {
	// Get query param
	var payload dto.VerificationPayload
	q := r.URL.Query()
	payload.VerificationToken = q.Get("t")

	// Validate payload
	err := payload.Validate()
	if err != nil {
		log.Errorf("Invalid payload. err: %v", err)
		c.renderError(w, 400, err)
		return
	}

	// Init service
	svc := c.NewService(r.Context())
	defer svc.Close()

	// Call service
	resp, err := svc.VerifyEmailCustomer(payload)
	if err != nil {
		log.Errorf("Error when processing service. err: %v", err)
		return
	}

	// Render response
	c.renderSuccess(w, resp)
	return
}

func (c *AccountController) renderError(w http.ResponseWriter, statusCode int, err error) {
	// Write header in
	w.Header().Add(nhttp.ContentTypeHeader, nhttp.ContentTypeJSON)

	// Write header
	w.WriteHeader(statusCode)

	// Write error in JSON
	err = json.NewEncoder(w).Encode(err)
	if err != nil {
		log.Errorf("failed to write response to json ( payload = %+v )", err)
	}
}

func (c *AccountController) renderSuccess(w http.ResponseWriter, htmlBody string) {
	w.Header().Add(nhttp.ContentTypeHeader, "text/html")
	_, err := w.Write([]byte(htmlBody))
	if err != nil {
		log.Errorf("failed to write response", nlogger.Error(err))
	}
}

func (c *AccountController) PostSendOTP(rx *nhttp.Request) (*nhttp.Response, error) {
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
		return nil, nhttp.BadRequestError.Trace(errx.Source(err))
	}

	// Init service
	svc := c.NewService(rx.Context())
	defer svc.Close()

	// Call service
	resp, err := svc.RegisterStepOne(payload)
	if err != nil {
		log.Error("error when processing service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (c *AccountController) PostResendOTP(rx *nhttp.Request) (*nhttp.Response, error) {
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
		return nil, nhttp.BadRequestError.Trace(errx.Source(err))
	}

	// Init service
	svc := c.NewService(rx.Context())
	defer svc.Close()

	// Call service
	resp, err := svc.RegisterResendOTP(payload)
	if err != nil {
		log.Error("error found when call register resend otp service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (c *AccountController) PostVerifyOTP(rx *nhttp.Request) (*nhttp.Response, error) {
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
		return nil, nhttp.BadRequestError.Trace(errx.Source(err))
	}

	// Init service
	svc := c.NewService(rx.Context())
	defer svc.Close()

	// Call service
	resp, err := svc.RegisterStepTwo(payload)
	if err != nil {
		log.Error("error when processing service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (c *AccountController) PostRegister(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get Payload
	var payload dto.RegisterNewCustomer
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", nlogger.Error(err), nlogger.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = payload.Validate()
	if err != nil {
		log.Error("unprocessable entity", nlogger.Error(err), nlogger.Context(ctx))
		return nil, nhttp.BadRequestError.Trace(errx.Source(err))
	}

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Check is force update password
	validatePassword := svc.validatePassword(payload.Password)
	if !validatePassword.IsValid {
		err = fmt.Errorf("password: %s", validatePassword.Message)
		return nil, nhttp.BadRequestError.Trace(errx.Source(err))
	}

	// Call service
	resp, err := svc.Register(payload)
	if err != nil {
		log.Error("error found when call service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (c *AccountController) PostUpdatePasswordCheck(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get user UserRefID
	userRefID, err := getUserRefID(rx)
	if err != nil {
		log.Errorf("error: %v", err, nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	// Get payload
	var payload dto.UpdatePasswordCheckRequest
	err = rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", nlogger.Error(err), nlogger.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = payload.Validate()
	if err != nil {
		log.Error("unprocessable entity", nlogger.Error(err), nlogger.Context(ctx))
		return nil, nhttp.BadRequestError.Trace(errx.Source(err))
	}

	// Init service
	svc := c.NewService(rx.Context())
	defer svc.Close()

	// Call service
	valid, err := svc.isValidPassword(userRefID, payload.CurrentPassword)
	if err != nil {
		log.Error("error when processing service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, err
	}

	if !valid {
		return nil, constant.InvalidPasswordError
	}

	return nhttp.Success().SetMessage("Password Sesuai"), nil
}

func (c *AccountController) PutUpdatePassword(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get user UserRefID
	userRefID, err := getUserRefID(rx)
	if err != nil {
		log.Errorf("error: %v", err, nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	// Get payload
	var payload dto.UpdatePasswordRequest
	err = rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", nlogger.Error(err), nlogger.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = payload.Validate()
	if err != nil {
		log.Error("unprocessable entity", nlogger.Error(err), nlogger.Context(ctx))
		return nil, nhttp.BadRequestError.Trace(errx.Source(err))
	}

	// Init service
	svc := c.NewService(rx.Context())
	defer svc.Close()

	// Call service
	err = svc.UpdatePassword(userRefID, payload)
	if err != nil {
		log.Error("error when call update service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetMessage("Password diperbarui"), nil
}
