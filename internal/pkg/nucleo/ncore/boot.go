package ncore

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"path"
)

type BootOptions struct {
	// Values for Core
	Manifest    Manifest    // Application manifest
	Environment Environment // Process environment
	WorkDir     string      // Working directory
	NodeId      string      // Process Node Identifier
	// Boot options
	EnvPrefix string
	// ResponseMapFile
	ResponseMapFile string
}

func Boot(destConfig interface{}, destConfigExternal interface{}, args ...BootOptions) *Core {
	// Load Options
	options := getBootOptions(args)

	// Init Core
	core := Core{
		Manifest:    options.Manifest,
		Environment: options.Environment,
		WorkDir:     options.WorkDir,
		NodeId:      options.NodeId,
	}

	// Load responses
	responses, err := loadResponseMap(options.ResponseMapFile)
	if err != nil {
		panic(err)
	}
	core.Responses = responses

	// Load config
	err = envconfig.Process(options.EnvPrefix, destConfig)
	if err != nil {
		panic(err)
	}

	// Load config external
	err = envconfig.Process("EXTERNAL", destConfigExternal)
	if err != nil {
		panic(err)
	}

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
	if options.NodeId == "" {
		nodeId, err := gonanoid.Generate(AlphaNumUpperCharSet, 4)
		if err != nil {
			panic(fmt.Errorf("failed to generate NodeId. Error = %w", err))
		}
		options.NodeId = nodeId
	}

	// If config file is not set, then set default
	if options.ResponseMapFile == "" {
		options.ResponseMapFile = path.Join(options.WorkDir, "responses.yml")
	}

	return options
}
