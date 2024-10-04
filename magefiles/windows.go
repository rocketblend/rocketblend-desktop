package main

import (
	"fmt"
	"path/filepath"

	"github.com/magefile/mage/sh"
)

func buildWindows(name, version, timestamp, commitSha, link, outputDir string, debug bool) error {
	ldFlags := buildFlags(version, timestamp, commitSha, link, debug)
	if err := buildWindowsAMD64(name, version, ldFlags, outputDir, true); err != nil {
		return fmt.Errorf("error building Windows AMD64: %s", err)
	}

	return nil
}

// buildWindowsAMD64 builds the Windows AMD64 version of the project.
func buildWindowsAMD64(name, version, ldFlags, outputDir string, skipFrontend bool) error {
	outputFilePath := filepath.Join(outputDir, fmt.Sprintf("%s-windows-amd64-%s.exe", name, version))
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

	err := sh.RunWithV(crossCompileFlags, "wails", "build", "-m", "-nosyncgomod", "-ldflags", ldFlags, "-nsis", "-platform", "windows/amd64", "-o", outputFilePath, skipBindingsFlag, skipFrontendFlag)
	if err != nil {
		return fmt.Errorf("error building Windows AMD64: %v", err)
	}

	return nil
}
