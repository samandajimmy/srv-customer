package main

import (
	"flag"
	"fmt"
	"os"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
)

type CmdFlags struct {
	CmdShowHelp        *bool
	CmdShowVersion     *bool
	OptEnvironment     *ncore.Environment
	OptWorkDir         *string
	OptResponseMapFile *string
	OptNodeNo          *int64
	OptEnvFile         *string
	OptLoadEnvFile     *bool
}

type BootOptions struct {
	Core              ncore.BootOptions
	CmdSeedSuperAdmin bool
}

/// initCmdFlags initiate available command line interface commands and options for parsing
func initCmdFlags() CmdFlags {
	return CmdFlags{
		CmdShowHelp:        flag.Bool("help", false, "Command: Show available commands and options"),
		CmdShowVersion:     flag.Bool("version", false, "Command: Show version"),
		OptEnvironment:     flag.Int("env", ncore.DevelopmentEnvironment, "Option: Set app environment"),
		OptEnvFile:         flag.String("env-file", "", "Option: Set config file"),
		OptResponseMapFile: flag.String("response-map", "", "Option: Set error codes file"),
		OptWorkDir:         flag.String("dir", ".", "Option: Set working directory"),
		OptNodeNo:          flag.Int64("node-no", 1, "Option: App instance number"),
		OptLoadEnvFile:     flag.Bool("load-env-file", false, "Option: Load environment from file that is set in -config or .env as default"),
	}
}

func handleCmdFlags() BootOptions {
	// Parse CLI commands and options
	cmdFlags := initCmdFlags()
	flag.Parse()

	// Intercept help command
	if *cmdFlags.CmdShowHelp {
		printHelp()
		os.Exit(0)
	}

	// Intercept version command
	if *cmdFlags.CmdShowVersion {
		printVersion()
		os.Exit(0)
	}

	return BootOptions{
		Core: ncore.BootOptions{
			Manifest: ncore.Manifest{
				AppName:    AppName,
				AppVersion: AppVersion,
				Metadata: map[string]interface{}{
					"build_hash": BuildHash,
				},
			},
			NodeNo:          *cmdFlags.OptNodeNo,
			WorkDir:         *cmdFlags.OptWorkDir,
			Environment:     *cmdFlags.OptEnvironment,
			EnvFile:         *cmdFlags.OptEnvFile,
			ResponseMapFile: *cmdFlags.OptResponseMapFile,
			LoadEnvFile:     *cmdFlags.OptLoadEnvFile,
		},
	}
}

func printHelp() {
	fmt.Printf("%s. Available Commands and Options:\n\n", AppName)
	flag.PrintDefaults()
}

// printVersion print app version and integrity
func printVersion() {
	fmt.Printf("%s\n"+
		"  Version    : %s\n"+
		"  Build Hash : %s\n",
		AppName, AppVersion, BuildHash)
}
