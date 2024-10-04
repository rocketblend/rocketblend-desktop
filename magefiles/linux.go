package main

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/magefile/mage/sh"
)

func buildLinux(name, version, timestamp, commitSha, link, outputDir string, debug bool) error {
	ldFlags := buildFlags(version, timestamp, commitSha, link, debug)
	nativeArch := runtime.GOARCH
	if err := buildLinuxTarget(name, version, ldFlags, outputDir, nativeArch, false); err != nil {
		return fmt.Errorf("error building Linux %s: %s", nativeArch, err)
	}

	otherArch := getOtherArch(nativeArch)
	if err := buildLinuxTarget(name, version, ldFlags, outputDir, otherArch, true); err != nil {
		return fmt.Errorf("error building Linux %s: %s", otherArch, err)
	}

	return nil
}

func buildLinuxTarget(name, version, ldFlags, outputDir, arch string, skipFrontend bool) error {
	outputFilePath := filepath.Join(outputDir, fmt.Sprintf("%s-linux-%s-%s", name, arch, version))
	skipBindingsFlag, skipFrontendFlag := "", ""
	if skipFrontend {
		skipBindingsFlag, skipFrontendFlag = "-skipbindings", "-s"
	}

	cc, cxx := getCompiler(arch)
	crossCompileFlags := map[string]string{
		"GOOS":   "linux",
		"GOARCH": arch,
		"CC":     cc,
		"CXX":    cxx,
	}

	return sh.RunWithV(crossCompileFlags, "wails", "build", "-m", "-nosyncgomod", "-ldflags", ldFlags, "-platform", fmt.Sprintf("linux/%s", arch), "-o", outputFilePath, skipBindingsFlag, skipFrontendFlag)
}

func getOtherArch(currentArch string) string {
	if currentArch == "arm64" {
		return "amd64"
	}

	return "arm64"
}

// getCompiler returns the appropriate compiler based on the architecture.
func getCompiler(arch string) (cc, cxx string) {
	if arch == "arm64" {
		return "aarch64-linux-gnu-gcc", "aarch64-linux-gnu-g++"
	}

	return "x86_64-linux-gnu-gcc", "x86_64-linux-gnu-g++"
}
