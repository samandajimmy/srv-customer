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
	router.Handle(http.MethodPost, "/accounts/login", router.HandleFunc(controllers.Account.PostLogin))

	// Verification
	router.Handle(http.MethodGet, "/accounts/verify-email", http.HandlerFunc(controllers.Account.GetVerifyEmail))

	// Register Step-1
	router.Handle(http.MethodPost, "/accounts/register/send-otp", router.HandleFunc(controllers.Account.PostSendOTP))
	// Register Resend OTP
	router.Handle(http.MethodPost, "/accounts/register/resend-otp", router.HandleFunc(controllers.Account.PostResendOTP))
	// Register Step-2
	router.Handle(http.MethodPost, "/accounts/register/verify-otp", router.HandleFunc(controllers.Account.PostVerifyOTP))
	// Register Step-3
	router.Handle(http.MethodPost, "/accounts/register", router.HandleFunc(controllers.Account.PostRegister))

	// Check password
	router.Handle(http.MethodPost, "/accounts/check-password", router.HandleFunc(controllers.Account.HandleAuthUser),
		router.HandleFunc(controllers.Account.PostUpdatePasswordCheck))
	// Update password
	router.Handle(http.MethodPut, "/accounts/password", router.HandleFunc(controllers.Account.HandleAuthUser),
		router.HandleFunc(controllers.Account.PutUpdatePassword))

	// Profile
	router.Handle(http.MethodGet, "/profiles", router.HandleFunc(controllers.Account.HandleAuthUser),
		router.HandleFunc(controllers.Profile.GetDetail))

	router.Handle(http.MethodPut, "/profiles", router.HandleFunc(controllers.Account.HandleAuthUser),
		router.HandleFunc(controllers.Profile.PutUpdate))

	router.Handle(http.MethodPost, "/profiles/avatar", router.HandleFunc(controllers.Account.HandleAuthUser),
		router.HandleFunc(controllers.Profile.PostUpdateAvatar))

	router.Handle(http.MethodPost, "/profiles/ktp", router.HandleFunc(controllers.Account.HandleAuthUser),
		router.HandleFunc(controllers.Profile.PostUpdateKTP))

	router.Handle(http.MethodPost, "/profiles/npwp", router.HandleFunc(controllers.Account.HandleAuthUser),
		router.HandleFunc(controllers.Profile.PostUpdateNPWP))

	router.Handle(http.MethodPost, "/profiles/sid", router.HandleFunc(controllers.Account.HandleAuthUser),
		router.HandleFunc(controllers.Profile.PostUpdateSID))

	router.Handle(http.MethodGet, "/profiles/status", router.HandleFunc(controllers.Account.HandleAuthUser),
		router.HandleFunc(controllers.Profile.GetStatus))

	// PIN
	router.Handle(http.MethodPost, "/accounts/pin/validation",
		router.HandleFunc(controllers.Account.PostValidatePin))
	router.Handle(http.MethodPost, "/accounts/pin/check",
		router.HandleFunc(controllers.Account.HandleAuthUser),
		router.HandleFunc(controllers.Account.PostCheckPin))
	router.Handle(http.MethodPost, "/accounts/pin/update",
		router.HandleFunc(controllers.Account.HandleAuthUser),
		router.HandleFunc(controllers.Account.PostUpdatePin))
	router.Handle(http.MethodPost, "/accounts/pin/otp-create",
		router.HandleFunc(controllers.Account.HandleAuthUser),
		router.HandleFunc(controllers.Account.PostCheckOTPPinCreate))
	router.Handle(http.MethodPost, "/accounts/pin/create",
		router.HandleFunc(controllers.Account.HandleAuthUser),
		router.HandleFunc(controllers.Account.PostCreatePin))
	router.Handle(http.MethodPost, "/accounts/pin/otp-forget",
		router.HandleFunc(controllers.Account.HandleAuthUser),
		router.HandleFunc(controllers.Account.PostOTPForgetPin))
	router.Handle(http.MethodPost, "/accounts/pin/forget",
		router.HandleFunc(controllers.Account.HandleAuthUser),
		router.HandleFunc(controllers.Account.PostForgetPin))

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
