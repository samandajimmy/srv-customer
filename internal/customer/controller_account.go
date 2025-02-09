package customer

import (
	"encoding/json"
	"fmt"
	"github.com/nbs-go/errx"
	logOption "github.com/nbs-go/nlogger/v2/option"
	"net/http"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/validate"
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

func (c *AccountController) HandleAuthUser(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get token
	tokenString, err := nhttp.ExtractBearerAuth(rx.Request)
	if err != nil {
		log.Error("error when extract token", logOption.Error(err), logOption.Context(ctx))
		return nil, errx.Trace(err)
	}

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Get UserRefID
	userRefID, err := svc.ValidateTokenAndRetrieveUserRefID(tokenString)
	if err != nil {
		return nil, errx.Trace(err)
	}

	rx.SetContextValue(constant.UserRefIDContextKey, userRefID)

	return nhttp.Continue(), nil
}

func GetUserRefID(rx *nhttp.Request) string {
	v := rx.GetContextValue(constant.UserRefIDContextKey)

	userRefID, ok := v.(string)
	if !ok {
		panic(fmt.Errorf("unexpected userRefID value in context. Type: %T", v))
	}

	return userRefID
}

func (c *AccountController) PostLogin(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get Payload
	var payload dto.LoginPayload
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = validate.PostLogin(&payload)
	if err != nil {
		log.Error("Bad request validate payload")
		return nil, err
	}

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Call service
	resp, err := svc.Login(payload)
	if err != nil {
		log.Error("error found when call service", logOption.Error(err), logOption.Context(ctx))
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
	err := validate.GetVerifyEmail(&payload)
	if err != nil {
		log.Errorf("Invalid payload. err: %v", err)
		c.renderError(w, 400, err)
		return //nolint:gosimple
	}

	// Init service
	svc := c.NewService(r.Context())
	defer svc.Close()

	// Call service
	resp, err := svc.VerifyEmailCustomer(payload)
	if err != nil {
		log.Errorf("Error when processing service. err: %v", err)
		c.renderSuccess(w, resp)
		return //nolint:gosimple
	}

	// Render response
	c.renderSuccess(w, resp)
	return //nolint:gosimple
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
		log.Error("failed to write response", logOption.Error(err))
	}
}

func (c *AccountController) PostSendOTP(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get Payload
	var payload dto.SendOTPPayload
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = validate.PostSendOTP(&payload)
	if err != nil {
		log.Error("Bad request validate payload", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	// Init service
	svc := c.NewService(rx.Context())
	defer svc.Close()

	// Call service
	resp, err := svc.RegisterStepOne(payload)
	if err != nil {
		log.Error("error when processing service", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (c *AccountController) PostResendOTP(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get Payload
	var payload dto.RegisterResendOTPPayload
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = validate.PostResendOTP(&payload)
	if err != nil {
		log.Error("Bad request validate payload", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	// Init service
	svc := c.NewService(rx.Context())
	defer svc.Close()

	// Call service
	resp, err := svc.RegisterResendOTP(payload)
	if err != nil {
		log.Error("error found when call register resend otp service", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (c *AccountController) PostVerifyOTP(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get Payload
	var payload dto.RegisterVerifyOTPPayload
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = validate.PostVerifyOTP(&payload)
	if err != nil {
		log.Error("Bad request validate payload", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Trace(errx.Source(err))
	}

	// Init service
	svc := c.NewService(rx.Context())
	defer svc.Close()

	// Call service
	resp, err := svc.RegisterStepTwo(payload)
	if err != nil {
		log.Error("error when processing service", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (c *AccountController) PostRegister(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get Payload
	var payload dto.RegisterPayload
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = validate.PostRegister(&payload)
	if err != nil {
		log.Error("Bad request validate payload", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Trace(errx.Source(err))
	}

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Check is force update password
	validatePassword := svc.ValidatePassword(payload.Password)
	if !validatePassword.IsValid {
		err = fmt.Errorf("password: %s", validatePassword.Message)
		return nil, nhttp.BadRequestError.Trace(errx.Source(err))
	}

	// Call service
	resp, err := svc.Register(payload)
	if err != nil {
		log.Error("error found when call service", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (c *AccountController) PostUpdatePasswordCheck(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get user UserRefID
	userRefID := GetUserRefID(rx)

	// Get payload
	var payload dto.UpdatePasswordCheckPayload
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = validate.PostUpdatePasswordCheck(&payload)
	if err != nil {
		log.Error("Bad request validate payload", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	// Init service
	svc := c.NewService(rx.Context())
	defer svc.Close()

	// Call service
	valid, err := svc.IsValidPassword(userRefID, payload.CurrentPassword)
	if err != nil {
		log.Error("error when processing service", logOption.Error(err), logOption.Context(ctx))
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
	userRefID := GetUserRefID(rx)

	// Get payload
	var payload dto.UpdatePasswordPayload
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = validate.PutUpdatePassword(&payload)
	if err != nil {
		log.Error("Bad request validate payload", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Trace(errx.Source(err))
	}

	// Init service
	svc := c.NewService(rx.Context())
	defer svc.Close()

	// Call service
	err = svc.UpdatePassword(userRefID, payload)
	if err != nil {
		log.Error("error when call update service", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetMessage("Password diperbarui"), nil
}

func (c *AccountController) PostValidatePin(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get Payload
	var payload dto.ValidatePinPayload
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = validate.PostValidatePin(&payload)
	if err != nil {
		log.Error("Bad request validate payload", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Trace(errx.Source(err))
	}

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Call service
	err = svc.ValidatePin(&payload)
	if err != nil {
		log.Error("error when call update service", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetMessage(constant.IsValidPINMessage), nil
}

func (c *AccountController) PostCheckPin(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get user UserRefID
	userRefID := GetUserRefID(rx)

	// Get Payload
	var payload dto.CheckPinPayload
	payload.UserRefID = userRefID
	payload.CheckPIN = true
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = validate.PostCheckPin(&payload)
	if err != nil {
		log.Error("Bad request validate payload", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Trace(errx.Source(err))
	}

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Call service
	err = svc.CheckPinUser(&payload)
	if err != nil {
		log.Error("error when call update service", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetMessage(constant.IsCheckValidMessage), nil
}

func (c *AccountController) PostUpdatePin(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get user UserRefID
	userRefID := GetUserRefID(rx)

	// Get Payload
	var payload dto.UpdatePinPayload
	payload.UserRefID = userRefID
	payload.CheckPIN = true
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = validate.PostUpdatePin(&payload)
	if err != nil {
		log.Error("Bad request validate payload", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Call service
	resp, err := svc.UpdatePin(&payload)
	if err != nil {
		log.Error("error found when call service", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp).SetMessage(resp.Title), nil
}

func (c *AccountController) PostCheckOTPPinCreate(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get user UserRefID
	userRefID := GetUserRefID(rx)

	// Get Payload
	var payload dto.CheckOTPPinPayload
	payload.UserRefID = userRefID
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = validate.CheckPostOTP(&payload)
	if err != nil {
		log.Error("Bad request validate payload", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Call service
	err = svc.CheckOTPPinCreate(&payload)
	if err != nil {
		log.Error("error found when call service", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetMessage(constant.OTPValidMessage), nil
}

func (c *AccountController) PostCreatePin(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get user UserRefID
	userRefID := GetUserRefID(rx)

	// Get Payload
	var payload dto.PostCreatePinPayload
	payload.UserRefID = userRefID
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = validate.PostCreatePin(&payload)
	if err != nil {
		log.Error("Bad request validate payload", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Call service
	err = svc.CreatePinUser(&payload)
	if err != nil {
		log.Error("error found when call service", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetMessage(constant.IsValidPINMessage), nil
}

func (c *AccountController) PostOTPForgetPin(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get user UserRefID
	userRefID := GetUserRefID(rx)

	// Get Payload
	var payload dto.CheckOTPPinPayload
	payload.UserRefID = userRefID
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = validate.CheckPostOTP(&payload)
	if err != nil {
		log.Error("Bad request validate payload", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Call service
	err = svc.CheckOTPForgetPin(&payload)
	if err != nil {
		log.Error("error found when call service", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetMessage(constant.OTPValidMessage), nil
}

func (c *AccountController) PostForgetPin(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get user UserRefID
	userRefID := GetUserRefID(rx)

	// Get Payload
	var payload dto.ForgetPinPayload
	payload.UserRefID = userRefID
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = validate.PostForgetPin(&payload)
	if err != nil {
		log.Error("Bad request validate payload", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Call service
	err = svc.ForgetPin(&payload)
	if err != nil {
		log.Error("error found when call service", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetMessage(constant.IsUpdatedPINMessage), nil
}

func (c *AccountController) PostSendOTPResetPassword(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get Payload
	var payload dto.OTPResetPasswordPayload
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = validate.PostSendOTPPassword(&payload)
	if err != nil {
		log.Error("Bad request validate payload", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Call service
	err = svc.SendOTPResetPassword(payload)
	if err != nil {
		log.Error("error found when call service", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetMessage(constant.ResendOTPPasswordMessage), nil
}

func (c *AccountController) PostVerifyOTPResetPassword(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get Payload
	var payload dto.VerifyOTPResetPasswordPayload
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = validate.PostVerifyOTPResetPassword(&payload)
	if err != nil {
		log.Error("Bad request validate payload", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Call service
	err = svc.VerifyOTPResetPassword(payload)
	if err != nil {
		log.Error("error found when call service", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetMessage(constant.ValidOTPPasswordMessage), nil
}

func (c *AccountController) PostResetPasswordByOTP(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get Payload
	var payload dto.ResetPasswordByOTPPayload
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = validate.PostResetPasswordByOTP(&payload)
	if err != nil {
		log.Error("Bad request validate payload", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Call service
	err = svc.ResetPasswordByOTP(payload)
	if err != nil {
		log.Error("error found when call service", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetMessage(constant.SuccessfullyResetPassword), nil
}

func (c *AccountController) PostChangeEmail(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get payload
	var payload dto.EmailChangePayload
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = validate.PostChangeEmail(&payload)
	if err != nil {
		log.Error("Bad request validate payload", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Trace(errx.Source(err))
	}

	// Init service
	svc := c.NewService(rx.Context())
	defer svc.Close()

	// Call service
	err = svc.ChangeEmail(payload)
	if err != nil {
		log.Error("error when call update service", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetMessage("Email user berhasil diperbaharui."), nil
}

func (c *AccountController) PostChangePhoneNumber(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get payload
	var payload dto.ChangePhoneNumberPayload
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = validate.PostChangePhoneNumber(&payload)
	if err != nil {
		log.Error("Bad request validate payload")
		return nil, err
	}

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Call service
	resp, err := svc.UpdatePhoneNumber(payload)
	if err != nil {
		log.Error("error when call update phone number service", logOption.Error(err), logOption.Context(ctx))
		return nil, errx.Trace(err)
	}

	return nhttp.Success().SetData(resp), nil
}

func (c *AccountController) PostUpdateSmartAccess(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get user UserRefID
	userRefID := GetUserRefID(rx)

	// Get payload
	var payload dto.UpdateSmartAccessPayload
	payload.UserRefID = userRefID
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = validate.PostUpdateSmartAccess(&payload)
	if err != nil {
		log.Error("Bad request validate payload", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Trace(errx.Source(err))
	}

	// Init service
	svc := c.NewService(rx.Context())
	defer svc.Close()

	// Call service
	err = svc.PostUpdateSmartAccess(payload)
	if err != nil {
		log.Error("error when call update service", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetMessage(constant.ActivateSmartAccessMessage), nil
}

func (c *AccountController) GetSmartAccessStatus(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get user UserRefID
	userRefID := GetUserRefID(rx)

	// Get payload
	var payload dto.GetSmartAccessStatusPayload
	payload.UserRefID = userRefID
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = validate.GetSmartAccessStatus(&payload)
	if err != nil {
		log.Error("Bad request validate payload", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	// Init service
	svc := c.NewService(rx.Context())
	defer svc.Close()

	// Call service
	resp, err := svc.GetSmartAccessStatus(payload)
	if err != nil {
		log.Error("error when call update service", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (c *AccountController) PutSynchronizeCustomer(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get payload
	var payload dto.PutSynchronizeCustomerPayload
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("error when parse json body", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = validate.PutSynchronizeCustomer(&payload)
	if err != nil {
		log.Error("Bad request validate payload", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	// Init service
	svc := c.NewService(rx.Context())
	defer svc.Close()

	// Call service
	resp, err := svc.PutSynchronizeCustomer(payload)
	if err != nil {
		log.Error("error when call update service", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}
