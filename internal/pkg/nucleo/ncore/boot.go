package ncore

import (
	"path"
)

type BootOptions struct {
	Manifest        Manifest
	NodeNo          int64
	WorkDir         string
	Environment     Environment
	EnvFile         string
	ResponseMapFile string
	LoadEnvFile     bool
}

func Boot(configDest interface{}, args ...BootOptions) *Core {
	// Load Options
	options := getBootOptions(args)

	// Init Core
	core := Core{
		Manifest:    options.Manifest,
		Environment: options.Environment,
		WorkDir:     options.WorkDir,
		NodeNo:      options.NodeNo,
	}

	// Load config
	err := loadConfig(options.LoadEnvFile, configDest, options.EnvFile)
	if err != nil {
		panic(err)
	}

	// Load responses
	responses, err := loadResponseMap(options.ResponseMapFile)
	if err != nil {
		panic(err)
	}
	core.Responses = responses

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

	// If config file is not set, then set default
	if options.EnvFile == "" {
		options.EnvFile = path.Join(options.WorkDir, ".env")
	}

	// If config file is not set, then set default
	if options.ResponseMapFile == "" {
		options.ResponseMapFile = path.Join(options.WorkDir, "responses.yml")
	}

	// If node number is not set, then set to 1
	if options.NodeNo == 0 {
		options.NodeNo = 1
	}

	return options
}
