package customer

import (
	"github.com/nbs-go/errx"
	"github.com/nbs-go/nlogger/v2"
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
		log.Error("Error when parse json body from request", nlogger.Error(err))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate
	err = validate.PostCreateFavorite(&payload)
	if err != nil {
		log.Error("Bad request validate payload", nlogger.Error(err), nlogger.Context(ctx))
		return nil, nhttp.BadRequestError.Trace(errx.Source(err))
	}

	// Set subject and requestID
	payload.RequestID = GetRequestID(rx)
	payload.Subject = GetSubject(rx)
	payload.UserRefID = GetUserRefID(rx)

	// Init service
	svc := c.NewService(ctx)
	defer svc.Close()

	// Call service
	resp, err := svc.CreateFavorite(&payload)
	if err != nil {
		log.Error("error when call list favorite service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}
