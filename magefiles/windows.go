package main

import (
	"fmt"

	"github.com/magefile/mage/sh"
)

func buildReleaseWindows(name, version, timestamp, commitSha, link, outputDir string, debug bool) error {
	ldFlags := buildFlags(version, timestamp, commitSha, link, outputDir, debug)
	if err := buildWindowsAMD64(name, version, ldFlags, true); err != nil {
		return fmt.Errorf("error building Windows AMD64: %s", err)
	}

	return nil
}

// buildWindowsAMD64 builds the Windows AMD64 version of the project.
func buildWindowsAMD64(name, version, ldFlags string, skipFrontend bool) error {
	outputFileName := fmt.Sprintf("%s-windows-amd64-%s.exe", name, version)
	skipBindingsFlag, skipFrontendFlag := "", ""
	if skipFrontend {
		skipBindingsFlag, skipFrontendFlag = "-skipbindings", "-s"
	}

	crossCompileFlags := map[string]string{
		"GOOS":   "windows",
		"GOARCH": "amd64",
		"CC":     "x86_64-w64-mingw32-gcc",
		"CXX":    "x86_64-w64-mingw32-g++",
	}

	err := sh.RunWithV(crossCompileFlags, "wails", "build", "-m", "-nosyncgomod", "-ldflags", ldFlags, "-nsis", "-platform", "windows/amd64", "-o", outputFileName, skipBindingsFlag, skipFrontendFlag)
	if err != nil {
		return fmt.Errorf("error building Windows AMD64: %v", err)
	}

	return nil
}
