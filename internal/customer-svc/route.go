package customer_svc

import (
	"net/http"

	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
)

func setUpRoute(router *nhttp.Router, handlers *HandlerMap) {
	// Common
	router.Handle(http.MethodGet, "/", router.HandleFunc(handlers.Common.GetAPIStatus))

	// Customer
	router.Handle(http.MethodPost, "/auth/register", router.HandleFunc(handlers.Customer.PostCreate))
	// Register Step-1
	router.Handle(http.MethodPost, "/register/step-1", router.HandleFunc(handlers.Customer.SendOTP))
	// Register Step-2
	router.Handle(http.MethodPost, "/register/step-2", router.HandleFunc(handlers.Customer.VerifyOTP))
}
