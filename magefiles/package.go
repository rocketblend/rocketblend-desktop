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

func Package(path, version, bundleID, outputDir, entitlements string, notarize bool) error {
	if runtime.GOOS != "darwin" {
		return nil
	}

	variables, err := getNotarizeVariables()
	if err != nil {
		return err
	}

	return packageMacOS(path, version, bundleID, outputDir, variables.DeveloperID, variables.AppleID, variables.Password, variables.TeamID, entitlements, notarize)
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
