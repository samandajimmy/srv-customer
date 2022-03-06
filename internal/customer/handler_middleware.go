package customer

import (
	"fmt"
	"github.com/nbs-go/errx"
	"github.com/nbs-go/nlogger/v2"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
)

type Middlewares struct {
	*Handler
}

func NewMiddlewares(h *Handler) *Middlewares {
	return &Middlewares{h}
}

func (h *Middlewares) AuthUser(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get token
	tokenString, err := nhttp.ExtractBearerAuth(rx.Request)
	if err != nil {
		log.Error("error when extract token", nlogger.Error(err), nlogger.Context(ctx))
		return nil, errx.Trace(err)
	}

	// Init service
	svc := h.NewService(ctx)
	defer svc.Close()

	// Get UserRefID
	userRefID, err := svc.validateTokenAndRetrieveUserRefID(tokenString)
	if err != nil {
		return nil, err
	}

	rx.SetContextValue(constant.UserRefIDContextKey, userRefID)

	return nhttp.Continue(), nil
}

func getUserRefID(rx *nhttp.Request) (string, error) {
	v := rx.GetContextValue(constant.UserRefIDContextKey)

	userRefID, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("unexpected userRefID value in context. Type: %T", v)
	}

	return userRefID, nil
}
