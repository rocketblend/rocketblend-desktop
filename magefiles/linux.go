package main

import (
	"fmt"
	"runtime"

	"github.com/magefile/mage/sh"
)

func buildReleaseLinux(name, version, timestamp, commitSha, link, outputDir string, debug bool) error {
	ldFlags := buildFlags(version, timestamp, commitSha, link, outputDir, debug)
	if runtime.GOARCH == "arm64" {
		if err := buildLinuxARM64(name, version, ldFlags, false); err != nil {
			return fmt.Errorf("error building Linux ARM64: %s", err)
		}

		if err := buildLinuxAMD64(name, version, ldFlags, true); err != nil {
			return fmt.Errorf("error building Linux AMD64: %s", err)
		}

		return nil
	}

	if err := buildLinuxAMD64(name, version, ldFlags, false); err != nil {
		return fmt.Errorf("error building Linux AMD64: %s", err)
	}

	if err := buildLinuxARM64(name, version, ldFlags, true); err != nil {
		return fmt.Errorf("error building Linux ARM64: %s", err)
	}

	return nil
}

// buildLinuxAMD64 builds the Linux AMD64 version of the project.
func buildLinuxAMD64(name, version, ldFlags string, skipFrontend bool) error {
	outputFileName := fmt.Sprintf("%s-linux-amd64-%s", name, version)
	skipBindingsFlag, skipFrontendFlag := "", ""
	if skipFrontend {
		skipBindingsFlag, skipFrontendFlag = "-skipbindings", "-s"
	}

	crossCompileFlags := map[string]string{"GOOS": "linux", "GOARCH": "amd64", "CC": "x86_64-linux-gnu-gcc", "CXX": "x86_64-linux-gnu-g++"}
	return sh.RunWithV(crossCompileFlags, "wails", "build", "-m", "-nosyncgomod", "-ldflags", ldFlags, "-platform", "linux/amd64", "-o", outputFileName, skipBindingsFlag, skipFrontendFlag)
}

// buildLinuxARM64 builds the Linux ARM64 version of the project.
func buildLinuxARM64(name, version, ldFlags string, skipFrontend bool) error {
	outputFileName := fmt.Sprintf("%s-linux-arm64-%s", name, version)
	skipBindingsFlag, skipFrontendFlag := "", ""
	if skipFrontend {
		skipBindingsFlag, skipFrontendFlag = "-skipbindings", "-s"
	}

	crossCompileFlags := map[string]string{"GOOS": "linux", "GOARCH": "arm64", "CC": "aarch64-linux-gnu-gcc", "CXX": "aarch64-linux-gnu-g++"}
	return sh.RunWithV(crossCompileFlags, "wails", "build", "-m", "-nosyncgomod", "-ldflags", ldFlags, "-platform", "linux/arm64", "-o", outputFileName, skipBindingsFlag, skipFrontendFlag)
}
