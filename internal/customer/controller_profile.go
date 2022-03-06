package customer

import (
	"github.com/nbs-go/errx"
	"github.com/nbs-go/nlogger/v2"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
)

type ProfileController struct {
	*Handler
}

func NewProfileController(h *Handler) *ProfileController {
	return &ProfileController{
		h,
	}
}

func (h *ProfileController) Get(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get user UserRefID
	userRefID, err := getUserRefID(rx)
	if err != nil {
		log.Errorf("error: %v", err, nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	// Init service
	svc := h.NewService(ctx)
	defer svc.Close()

	// Call service
	resp, err := svc.CustomerProfile(userRefID)
	if err != nil {
		log.Error("error when call customer profile service", nlogger.Error(err), nlogger.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}
