package customer

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/nbs-go/errx"
	"github.com/nbs-go/nlogger/v2"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/validate"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
)

type BankAccountController struct {
	*Handler
}

func NewBankAccountController(h *Handler) *BankAccountController {
	return &BankAccountController{
		h,
	}
}

func (c *BankAccountController) GetListBankAccount(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	payload, err := getListPayload(rx)
	if err != nil {
		return nil, errx.Trace(err)
	}

	// Get user UserRefID
	userRefID, err := getUserRefID(rx)
	if err != nil {
		log.Errorf("error: %v", err, nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Call service
	resp, err := svc.ListBankAccount(userRefID, payload)
	if err != nil {
		log.Error("error when call list bank account service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, err
	}

	return nhttp.OK().SetData(resp), nil
}

func (c *BankAccountController) PostCreateBankAccount(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get user UserRefID
	userRefID, err := getUserRefID(rx)
	if err != nil {
		log.Errorf("error: %v", err, nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	// Get Payload
	var payload dto.CreateBankAccountPayload
	err = rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("Error when parse json body from request", nlogger.Error(err))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate
	err = validate.PostCreateBankAccount(&payload)
	if err != nil {
		log.Error("Bad request validate payload", nlogger.Error(err), nlogger.Context(ctx))
		return nil, nhttp.BadRequestError.Trace(errx.Source(err))
	}

	// Set subject and requestID
	payload.RequestID = GetRequestID(rx)
	payload.Subject = GetSubject(rx)

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Call service
	resp, err := svc.CreateBankAccount(userRefID, &payload)
	if err != nil {
		log.Error("error when call list bank account service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (c *BankAccountController) GetDetailBankAccount(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get xid
	xid := mux.Vars(rx.Request)["xid"]
	if xid == "" {
		err := errors.New("xid is not found on params")
		log.Errorf("xid is not found on params. err: %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Get user UserRefID
	userRefID, err := getUserRefID(rx)
	if err != nil {
		log.Error("error when get userRefID", nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	// Set payload
	var payload dto.GetDetailBankAccountPayload
	payload.RequestID = GetRequestID(rx)
	payload.XID = xid

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Call service
	resp, err := svc.GetDetailBankAccount(userRefID, &payload)
	if err != nil {
		log.Error("error when call detail bank account service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, err
	}

	return nhttp.OK().SetData(resp), nil
}

func (c *BankAccountController) PutUpdateBankAccount(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get user UserRefID
	userRefID, err := getUserRefID(rx)
	if err != nil {
		log.Errorf("error: %v", err, nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	// Get payload
	var payload dto.UpdateBankAccountPayload
	err = rx.ParseJSONBody(&payload)
	if err != nil {
		log.Errorf("Error when parse json body from request.", nlogger.Error(err))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Set modifier and id
	payload.RequestID = GetRequestID(rx)
	payload.XID = mux.Vars(rx.Request)["xid"]
	payload.Subject = GetSubject(rx)

	err = validate.PutUpdateBankAccount(&payload)
	if err != nil {
		log.Error("Bad request validate payload", nlogger.Error(err), nlogger.Context(ctx))
		return nil, nhttp.BadRequestError.Trace(errx.Source(err))
	}

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Call service
	resp, err := svc.UpdateBankAccount(userRefID, &payload)
	if err != nil {
		log.Error("error when call list bank account service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}
