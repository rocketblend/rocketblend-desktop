package main

import (
	"fmt"

	"github.com/magefile/mage/sh"
)

// buildLinuxAMD64 builds the Linux AMD64 version of the project.
func buildLinuxAMD64(ldFlags, appVersion string, skipFrontend bool) error {
	outputFileName := fmt.Sprintf("rocketblend-desktop-linux-amd64-%s", appVersion)
	skipBindingsFlag, skipFrontendFlag := "", ""
	if skipFrontend {
		skipBindingsFlag, skipFrontendFlag = "-skipbindings", "-s"
	}

	crossCompileFlags := map[string]string{"GOOS": "linux", "GOARCH": "amd64", "CC": "x86_64-linux-gnu-gcc", "CXX": "x86_64-linux-gnu-g++"}
	return sh.RunWithV(crossCompileFlags, "wails", "build", "-m", "-nosyncgomod", "-ldflags", ldFlags, "-platform", "linux/amd64", "-o", outputFileName, skipBindingsFlag, skipFrontendFlag)
}

// buildLinuxARM64 builds the Linux ARM64 version of the project.
func buildLinuxARM64(ldFlags, appVersion string, skipFrontend bool) error {
	outputFileName := fmt.Sprintf("rocketblend-desktop-linux-arm64-%s", appVersion)
	skipBindingsFlag, skipFrontendFlag := "", ""
	if skipFrontend {
		skipBindingsFlag, skipFrontendFlag = "-skipbindings", "-s"
	}

	crossCompileFlags := map[string]string{"GOOS": "linux", "GOARCH": "arm64", "CC": "aarch64-linux-gnu-gcc", "CXX": "aarch64-linux-gnu-g++"}
	return sh.RunWithV(crossCompileFlags, "wails", "build", "-m", "-nosyncgomod", "-ldflags", ldFlags, "-platform", "linux/arm64", "-o", outputFileName, skipBindingsFlag, skipFrontendFlag)
}
