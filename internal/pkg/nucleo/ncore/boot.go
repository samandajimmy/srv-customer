package ncore

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nvalidate"
)

type BootOptions struct {
	// Values for Core
	Manifest    Manifest    // Application manifest
	Environment Environment // Process environment
	WorkDir     string      // Working directory
	NodeID      string      // Process Node Identifier
	// Boot options
	EnvPrefix string
}

func Boot(destConfig interface{}, destConfigExternal interface{}, args ...BootOptions) *Core {
	// Load Options
	options := getBootOptions(args)

	// Init Core
	core := Core{
		Manifest:    options.Manifest,
		Environment: options.Environment,
		WorkDir:     options.WorkDir,
		NodeID:      options.NodeID,
	}

	// Load config
	err := envconfig.Process(options.EnvPrefix, destConfig)
	if err != nil {
		panic(err)
	}

	// Load config external
	err = envconfig.Process("EXTERNAL", destConfigExternal)
	if err != nil {
		panic(err)
	}

	// Load validation message
	nvalidate.Init()

	return &core
}

func getBootOptions(args []BootOptions) BootOptions {
	// Get options
	var options BootOptions
	if len(args) > 0 {
		options = args[0]
	}

	// If working directory is not set, then set to current directory
	if options.WorkDir == "" {
		options.WorkDir = "."
	}

	// If node number is not set, then set random id
	if options.NodeID == "" {
		nodeID, err := gonanoid.Generate(AlphaNumUpperCharSet, 4)
		if err != nil {
			panic(fmt.Errorf("failed to generate NodeId. Error = %w", err))
		}
		options.NodeID = nodeID
	}

	return options
}
