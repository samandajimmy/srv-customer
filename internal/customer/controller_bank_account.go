package customer

import (
	"github.com/nbs-go/errx"
	"github.com/nbs-go/nlogger/v2"
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
