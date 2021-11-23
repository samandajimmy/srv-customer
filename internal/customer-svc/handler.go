package customer_svc

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/handler"
	"time"
)

type HandlerMap struct {
	Common *handler.Common
}

func initHandler(app *API) *HandlerMap {

	return &HandlerMap{
		Common: handler.NewCommon(time.Now(), app.Manifest.AppVersion, app.Manifest.GetStringMetadata(constant.BuildHashKey)),
	}
}
