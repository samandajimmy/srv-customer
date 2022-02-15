package customer

import (
	"errors"
	"fmt"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nvalidate"
)

type Customer struct {
	*Handler
}

func NewCustomer(h *Handler) *Customer {
	return &Customer{h}
}

func (h *Customer) PostLogin(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get Payload
	var payload dto.LoginRequest
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
	svc := h.NewService(ctx)
	defer svc.Close()

	// Call service
	resp, err := svc.Login(payload)
	if err != nil {
		log.Error("error found when call service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (h *Customer) PostRegister(rx *nhttp.Request) (*nhttp.Response, error) {
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
		data := nvalidate.Message(err.Error())
		return nhttp.UnprocessableEntity(data), nil
	}

	// Init service
	svc := h.NewService(ctx)
	defer svc.Close()

	// Check is force update password
	validatePassword := svc.validatePassword(payload.Password)
	if validatePassword.IsValid == false {
		err := errors.New(fmt.Sprintf("password: %s.", validatePassword.Message))
		data := nvalidate.Message(err.Error())
		return nhttp.UnprocessableEntity(data), nil
	}

	// Call service
	resp, err := svc.Register(payload)
	if err != nil {
		log.Errorf("error found when call service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

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
		log.Errorf("error when processing service", nlogger.Error(err), nlogger.Context(ctx))
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
		log.Error("error when processing service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (h *Customer) GetProfile(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get token
	tokenString, err := nhttp.ExtractBearerAuth(rx.Request)
	if err != nil {
		log.Error("error when extract token", nlogger.Error(err), nlogger.Context(ctx))
		return nil, ncore.TraceError("failed to extract token", err)
	}

	// Init service
	svc := h.NewService(ctx)
	defer svc.Close()

	// Get UserRefID
	userRefId, err := svc.validateTokenAndRetrieveUserRefID(tokenString)
	if err != nil {
		return nil, err
	}

	// Call service
	resp, err := svc.CustomerProfile(userRefId)
	if err != nil {
		log.Error("error when call customer profile service")
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (h *Customer) UpdateProfile(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get token
	tokenString, err := nhttp.ExtractBearerAuth(rx.Request)
	if err != nil {
		log.Error("error when extract token", nlogger.Error(err), nlogger.Context(ctx))
		return nil, ncore.TraceError("failed to extract token", err)
	}

	// Get payload
	var payload dto.UpdateProfileRequest
	err = rx.ParseJSONBody(&payload)
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

	// Get UserRefID
	userRefId, err := svc.validateTokenAndRetrieveUserRefID(tokenString)
	if err != nil {
		return nil, err
	}

	// Call service
	err = svc.UpdateCustomerProfile(userRefId, payload)
	if err != nil {
		log.Error("error when processing service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetMessage("Update data user berhasil").SetData(false), nil
}

func (h *Customer) UpdatePasswordCheck(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get token
	tokenString, err := nhttp.ExtractBearerAuth(rx.Request)
	if err != nil {
		log.Error("error when extract token", nlogger.Error(err), nlogger.Context(ctx))
		return nil, ncore.TraceError("failed to extract token", err)
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
		data := nvalidate.Message(err.Error())
		return nhttp.UnprocessableEntity(data), nil
	}

	// Init service
	svc := h.NewService(rx.Context())
	defer svc.Close()

	// Call service
	valid, err := svc.isValidPassword(tokenString, payload.CurrentPassword)
	if err != nil {
		log.Error("error when processing service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, err
	}

	if !valid {
		return nil, h.Responses.GetError("E_USR_1")
	}

	return nhttp.Success().SetMessage("Password Sesuai"), nil
}

func (h *Customer) UpdatePassword(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get user UserRefID
	userRefID, err := getUserRefID(rx)
	if err != nil {
		log.Errorf("error: %v", err, nlogger.Error(err), nlogger.Context(ctx))
		return nil, ncore.TraceError("", err)
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
		data := nvalidate.Message(err.Error())
		return nhttp.UnprocessableEntity(data), nil
	}

	// Init service
	svc := h.NewService(rx.Context())
	defer svc.Close()

	// Call service
	err = svc.UpdatePassword(userRefID, payload)
	if err != nil {
		log.Error("error when call update service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetMessage("Password diperbarui"), nil
}

func (s *Service) validateJWT(token string) (jwt.Token, error) {
	// Parsing Token
	t, err := jwt.ParseString(token, jwt.WithVerify(constant.JWTSignature, []byte(s.config.JWTKey)))
	if err != nil {
		s.log.Error("parsing jwt token", nlogger.Error(err), nlogger.Context(s.ctx))
		return nil, err
	}

	if err = jwt.Validate(t); err != nil {
		s.log.Error("error when validate", nlogger.Error(err), nlogger.Context(s.ctx))
		return nil, err
	}

	err = jwt.Validate(t, jwt.WithIssuer(constant.JWTIssuer))
	if err != nil {
		s.log.Error("error found when validate with issuer", nlogger.Error(err), nlogger.Context(s.ctx))
		return nil, err
	}

	return t, nil
}

func (s *Service) validateTokenAndRetrieveUserRefID(tokenString string) (string, error) {
	// Get Context
	ctx := s.ctx

	// validate JWT
	token, err := s.validateJWT(tokenString)
	if err != nil {
		s.log.Error("error when validate JWT", nlogger.Error(err), nlogger.Context(ctx))
		return "", err
	}

	accessToken, _ := token.Get("access_token")

	tokenId, _ := token.Get("id")

	// Session token
	key := fmt.Sprintf("%s:%s:%s", constant.Prefix, constant.CacheTokenJWT, tokenId)

	tokenFromCache, err := s.CacheGet(key)
	if err != nil {
		s.log.Error("error get token from cache", nlogger.Error(err), nlogger.Context(ctx))
		return "", err
	}

	if accessToken != tokenFromCache {
		return "", s.responses.GetError("E_AUTH_11")
	}

	userRefId := nval.ParseStringFallback(tokenId, "")

	return userRefId, nil
}
