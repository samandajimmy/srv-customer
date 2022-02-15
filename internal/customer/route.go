package customer

import (
	"net/http"
	"path"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
)

func setUpRoute(router *nhttp.Router, handlers *HandlerMap) {

	// Common
	router.Handle(http.MethodGet, "/", router.HandleFunc(handlers.Common.GetAPIStatus))

	// Login
	router.Handle(http.MethodPost, "/auth/login", router.HandleFunc(handlers.Customer.PostLogin))

	// Verification
	router.Handle(http.MethodGet, "/auth/verify_email", router.HandleFunc(handlers.Verification.VerifyEmail))

	// Register Step-1
	router.Handle(http.MethodPost, "/register/step-1", router.HandleFunc(handlers.Customer.SendOTP))
	// Register Resend OTP
	router.Handle(http.MethodPost, "/register/resend-otp", router.HandleFunc(handlers.Customer.ResendOTP))
	// Register Step-2
	router.Handle(http.MethodPost, "/register/step-2", router.HandleFunc(handlers.Customer.VerifyOTP))
	// Register Step-3
	router.Handle(http.MethodPost, "/register/step-3", router.HandleFunc(handlers.Customer.PostRegister))

	// Customer
	router.Handle(http.MethodGet, "/profile", router.HandleFunc(handlers.Customer.GetProfile))

	router.Handle(http.MethodPut, "/profile", router.HandleFunc(handlers.Customer.UpdateProfile))

	// Static asset
	staticDir := "/web/assets/"
	router.PathPrefix(staticDir).Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))
}

func InitRouter(workDir string, config *Config, handlers *HandlerMap) http.Handler {
	var debug bool
	if config.Debug != "" {
		debug = true
	} else {
		debug = false
	}

	// Init router
	router := nhttp.NewRouter(nhttp.RouterOptions{
		LogRequest: true,
		Debug:      debug,
		TrustProxy: nval.ParseBooleanFallback(config.TrustProxy, false),
	})

	// Enable cors
	if config.CORS.Enabled {
		log.Debug("CORS Enabled")
		router.Use(config.CORS.NewMiddleware())
	}

	// Set-up Routes
	setUpRoute(router, handlers)

	// Set-up Static
	staticPath := path.Join(workDir, "/web/static")
	staticDir := http.Dir(staticPath)
	staticServer := http.FileServer(staticDir)
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", staticServer))

	return router
}
