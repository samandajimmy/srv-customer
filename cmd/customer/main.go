package main

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/nbs-go/nlogger"
	"net/http"
	"os"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"time"
)

var log = nlogger.Get()

func main() {
	// Capture started at
	startedAt := time.Now()

	// Handle command flags
	bootOptions := handleCmdFlags()

	// Boot core
	config := customer.Config{}
	configExternal := customer.DatabaseConfig{}

	c := ncore.Boot(&config, &configExternal, bootOptions.Core)
	config.DatabaseExternal = configExternal

	// Init handler
	h, err := customer.NewHandler(c, &config)
	if err != nil {
		log.Debugf("error while init handler", err)
		panic(err)
	}

	// Register handlers
	handlers := customer.RegisterHandler(c.Manifest, h)

	// Set router
	router := customer.InitRouter(c.WorkDir, &config, handlers)
	serverConfig := config.ServerConfig

	// Start server
	log.Infof("Starting %s...", c.Manifest.AppName)
	log.Infof("NodeId = %s, Environment = %s", c.NodeID, c.GetEnvironmentString())
	log.Debugf("Boot Time: %s", time.Since(startedAt))

	err = http.ListenAndServe(serverConfig.GetListenPort(), router)
	if err != nil {
		log.Fatal("failed to serve", nlogger.Error(err))
		os.Exit(2)
	}
}
