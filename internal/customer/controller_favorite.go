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

type FavoriteController struct {
	*Handler
}

func NewFavoriteController(h *Handler) *FavoriteController {
	return &FavoriteController{
		h,
	}
}

func (c *FavoriteController) PostCreate(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get Payload
	var payload dto.CreateFavoritePayload
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("Error when parse json body from request", logOption.Error(err))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate
	err = validate.PostCreateFavorite(&payload)
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
	resp, err := svc.CreateFavorite(&payload)
	if err != nil {
		log.Error("error when call list favorite service", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (c *FavoriteController) GetList(rx *nhttp.Request) (*nhttp.Response, error) {
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
	resp, err := svc.ListFavorite(userRefID, payload)
	if err != nil {
		log.Error("error when call list transaction favorite service", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (c *FavoriteController) Delete(rx *nhttp.Request) (*nhttp.Response, error) {
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
	var payload dto.GetDetailFavoritePayload
	payload.RequestID = rx.GetRequestId()
	payload.UserRefID = GetUserRefID(rx)
	payload.XID = xid

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Call service
	err := svc.DeleteFavorite(&payload)
	if err != nil {
		return nil, errx.Trace(err)
	}

	return nhttp.OK(), nil
}
