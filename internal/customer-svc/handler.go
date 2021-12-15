package customer_svc

import (
	"time"

	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/handler"
)

type HandlerMap struct {
	Common       *handler.Common
	Auth         *handler.Auth
	Customer     *handler.Customer
	Verification *handler.Verification
}

func initHandler(app *API) *HandlerMap {

	return &HandlerMap{
		Common:       handler.NewCommon(time.Now(), app.Manifest.AppVersion, app.Manifest.GetStringMetadata(constant.BuildHashKey)),
		Auth:         handler.NewAuth(app.PdsApp.Services.Auth),
		Customer:     handler.NewCustomer(app.PdsApp.Services.Customer),
		Verification: handler.NewVerification(app.PdsApp.Services.Verification),
	}
}
