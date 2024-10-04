package main

import (
	"fmt"
	"os"
	"runtime"
)

type (
	NotarizeVariables struct {
		DeveloperID string
		AppleID     string
		Password    string
		TeamID      string
	}
)

func Package(appPath, version, bundleID, outputDir, entitlements string, notorize bool) error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("unsupported OS/architecture: %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	variables, err := getNotarizeVariables()
	if err != nil {
		return err
	}

	return packageMacOS(appPath, version, bundleID, outputDir, variables.DeveloperID, variables.AppleID, variables.Password, variables.TeamID, entitlements, notorize)
}

func getNotarizeVariables() (*NotarizeVariables, error) {
	env := NotarizeVariables{
		DeveloperID: os.Getenv("AC_DEVELOPER_ID"),
		AppleID:     os.Getenv("AC_APPLE_ID"),
		Password:    os.Getenv("AC_PASSWORD"),
		TeamID:      os.Getenv("AC_TEAM_ID"),
	}

	if env.DeveloperID == "" || env.AppleID == "" || env.Password == "" || env.TeamID == "" {
		return nil, fmt.Errorf("missing required environment variables")
	}

	return &env, nil
}