package customer

import (
	"net/http"
	"path"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
)

func setUpRoute(router *nhttp.Router, controllers *Controllers) {
	// Common
	router.Handle(http.MethodGet, "/", router.HandleFunc(controllers.Common.GetAPIStatus))

	// Login
	router.Handle(http.MethodPost, "/auth/login", router.HandleFunc(controllers.Customer.PostLogin))

	// Verification
	router.Handle(http.MethodGet, "/auth/verify_email", http.HandlerFunc(controllers.Verification.VerifyEmail))

	// Register Step-1
	router.Handle(http.MethodPost, "/register/step-1", router.HandleFunc(controllers.Customer.SendOTP))
	// Register Resend OTP
	router.Handle(http.MethodPost, "/register/resend-otp", router.HandleFunc(controllers.Customer.ResendOTP))
	// Register Step-2
	router.Handle(http.MethodPost, "/register/step-2", router.HandleFunc(controllers.Customer.VerifyOTP))
	// Register Step-3
	router.Handle(http.MethodPost, "/register/step-3", router.HandleFunc(controllers.Customer.PostRegister))

	// Profile
	router.Handle(http.MethodGet, "/profile", router.HandleFunc(controllers.Middlewares.AuthUser),
		router.HandleFunc(controllers.Profile.GetDetail))

	router.Handle(http.MethodPut, "/profile", router.HandleFunc(controllers.Middlewares.AuthUser),
		router.HandleFunc(controllers.Profile.PutUpdate))

	router.Handle(http.MethodPost, "/profile/check_password", router.HandleFunc(controllers.Middlewares.AuthUser),
		router.HandleFunc(controllers.Customer.UpdatePasswordCheck))

	router.Handle(http.MethodPut, "/profile/password", router.HandleFunc(controllers.Middlewares.AuthUser),
		router.HandleFunc(controllers.Customer.UpdatePassword))

	router.Handle(http.MethodPost, "/profile/avatar", router.HandleFunc(controllers.Middlewares.AuthUser),
		router.HandleFunc(controllers.Profile.PostUpdateAvatar))

	router.Handle(http.MethodPost, "/profile/ktp", router.HandleFunc(controllers.Middlewares.AuthUser),
		router.HandleFunc(controllers.Profile.PostUpdateKTP))

	router.Handle(http.MethodPost, "/profile/npwp", router.HandleFunc(controllers.Middlewares.AuthUser),
		router.HandleFunc(controllers.Profile.PostUpdateNPWP))

	router.Handle(http.MethodPost, "/profile/sid", router.HandleFunc(controllers.Middlewares.AuthUser),
		router.HandleFunc(controllers.Profile.PostUpdateSID))

	router.Handle(http.MethodGet, "/profile/status", router.HandleFunc(controllers.Middlewares.AuthUser),
		router.HandleFunc(controllers.Profile.GetStatus))

	// Static asset
	staticDir := "/web/assets/"
	router.PathPrefix(staticDir).Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))
}

func InitRouter(workDir string, config *Config, controllers *Controllers) http.Handler {
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
	setUpRoute(router, controllers)

	// Set-up Static
	staticPath := path.Join(workDir, "/web/static")
	staticDir := http.Dir(staticPath)
	staticServer := http.FileServer(staticDir)
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", staticServer))

	return router
}
