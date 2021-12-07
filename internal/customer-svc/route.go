package customer_svc

import (
	"net/http"

	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
)

func setUpRoute(router *nhttp.Router, handlers *HandlerMap) {
	// Common
	router.Handle(http.MethodGet, "/", router.HandleFunc(handlers.Common.GetAPIStatus))

	// Login
	router.Handle(http.MethodPost, "/auth/login", router.HandleFunc(handlers.Customer.PostLogin))

	// Register Step-1
	router.Handle(http.MethodPost, "/register/step-1", router.HandleFunc(handlers.Customer.SendOTP))
	// Register Resend OTP
	router.Handle(http.MethodPost, "/register/resend-otp", router.HandleFunc(handlers.Customer.ResendOTP))
	// Register Step-2
	router.Handle(http.MethodPost, "/register/step-2", router.HandleFunc(handlers.Customer.VerifyOTP))
	// Register Step-3
	router.Handle(http.MethodPost, "/register/step-3", router.HandleFunc(handlers.Customer.PostRegister))
}
