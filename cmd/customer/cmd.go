package main

import (
	"flag"
	"fmt"
	"os"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
)

type CmdFlags struct {
	CmdShowHelp    *bool
	CmdShowVersion *bool
	OptEnvironment *string
	OptWorkDir     *string
	OptNodeId      *string
}

type BootOptions struct {
	Core              ncore.BootOptions
	CmdSeedSuperAdmin bool
}

/// initCmdFlags initiate available command line interface commands and options for parsing
func initCmdFlags() CmdFlags {
	return CmdFlags{
		// Command
		CmdShowHelp:    flag.Bool("help", false, "Command: Show available commands and options"),
		CmdShowVersion: flag.Bool("version", false, "Command: Show version"),
		// Options
		OptEnvironment: flag.String("env", "", "Option: Set app environment"),
		OptWorkDir:     flag.String("dir", ".", "Option: Set working directory"),
		OptNodeId:      flag.String("node-id", "", "Option: App instance number"),
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
			Manifest:    ncore.NewManifest(AppName, AppVersion, BuildSignature),
			Environment: ncore.ParseEnvironment(*cmdFlags.OptEnvironment),
			WorkDir:     *cmdFlags.OptWorkDir,
			NodeId:      *cmdFlags.OptNodeId,
			EnvPrefix:   EnvPrefix,
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
		"  Version         : %s\n"+
		"  Build Signature : %s\n",
		AppName, AppVersion, BuildSignature)
}
