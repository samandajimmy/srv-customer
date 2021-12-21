package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
)

var log = nlogger.Get()

func main() {
	// Boot
	bootStartedAt := time.Now()
	app := boot()
	log.Debugf("Boot Time: %s", time.Since(bootStartedAt))

	// Start app
	start(&app)
}

func boot() customer_svc.API {
	// Handle command
	bootOptions := handleCmdFlags()

	// Boot Core
	config := contract.Config{}
	core := ncore.Boot(&config, bootOptions.Core)

	// Boot App
	app := customer_svc.NewAPI(core, config)
	err := app.Boot()
	if err != nil {
		panic(err)
	}

	return app
}

func start(app *customer_svc.API) {
	// Get server config
	config := app.Config.Server

	// Init router
	router := app.InitRouter()

	log.Infof("%s HTTP Server is listening to port %d", AppSlug, config.ListenPort)
	log.Infof("%s HTTP Server Started. Base URL: %s", AppSlug, config.GetHttpBaseUrl())
	err := http.ListenAndServe(config.GetListenPort(), router)
	if err != nil {
		panic(fmt.Errorf("%s: failed to on listen.\n  > %w", AppSlug, err))
	}
}
