package customer

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/nbs-go/errx"
	logOption "github.com/nbs-go/nlogger/v2/option"
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

	// Get List Payload
	payload, err := getListPayload(rx)
	if err != nil {
		return nil, errx.Trace(err)
	}

	// Get UserRefID
	userRefID := GetUserRefID(rx)

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Call service
	resp, err := svc.ListBankAccount(userRefID, payload)
	if err != nil {
		log.Error("error when call list bank account service", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (c *BankAccountController) PostCreateBankAccount(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get Payload
	var payload dto.CreateBankAccountPayload
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("Error when parse json body from request", logOption.Error(err))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate
	err = validate.PostCreateBankAccount(&payload)
	if err != nil {
		log.Error("Bad request validate payload", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Trace(errx.Source(err))
	}

	// Set subject and requestID
	payload.RequestID = rx.GetRequestId()
	payload.Subject = GetSubject(rx)
	payload.UserRefID = GetUserRefID(rx)

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Call service
	resp, err := svc.CreateBankAccount(&payload)
	if err != nil {
		log.Error("error when call list bank account service", logOption.Error(err), logOption.Context(ctx))
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
		log.Error("xid is not found on params", logOption.Error(err))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Set payload
	var payload dto.GetDetailBankAccountPayload
	payload.UserRefID = GetUserRefID(rx)
	payload.RequestID = rx.GetRequestId()
	payload.XID = xid

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Call service
	resp, err := svc.GetDetailBankAccount(&payload)
	if err != nil {
		log.Error("error when call detail bank account service", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	return nhttp.OK().SetData(resp), nil
}

func (c *BankAccountController) PutUpdateBankAccount(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get payload
	var payload dto.UpdateBankAccountPayload
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("Error when parse json body from request", logOption.Error(err))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Set modifier and id
	payload.RequestID = rx.GetRequestId()
	payload.XID = mux.Vars(rx.Request)["xid"]
	payload.Subject = GetSubject(rx)
	payload.UserRefID = GetUserRefID(rx)

	err = validate.PutUpdateBankAccount(&payload)
	if err != nil {
		log.Error("Bad request validate payload", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Trace(errx.Source(err))
	}

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Call service
	resp, err := svc.UpdateBankAccount(&payload)
	if err != nil {
		log.Error("error when call list bank account service", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (c *BankAccountController) DeleteBankAccount(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get xid
	xid := mux.Vars(rx.Request)["xid"]
	if xid == "" {
		err := errors.New("xid is not found on params")
		log.Errorf("xid is not found on params. err: %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Set payload
	var payload dto.GetDetailBankAccountPayload
	payload.RequestID = rx.GetRequestId()
	payload.UserRefID = GetUserRefID(rx)
	payload.XID = xid

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Call service
	err := svc.DeleteBankAccount(&payload)
	if err != nil {
		return nil, errx.Trace(err)
	}

	return nhttp.OK(), nil
}
