package main

import (
	"fmt"

	"github.com/magefile/mage/sh"
)

func buildLinux(name, version, timestamp, commitSha, link, buildtype string) error {
	ldFlags := buildFlags(version, timestamp, commitSha, link, buildtype)
	nativeArch := "amd64"
	if err := buildLinuxTarget(name, version, ldFlags, nativeArch, false); err != nil {
		return fmt.Errorf("error building Linux %s: %s", nativeArch, err)
	}

	return nil
}

func buildLinuxTarget(name, version, ldFlags, arch string, skipFrontend bool) error {
	fmt.Printf("Building Linux %s binary for %s\n", arch, name)
	outputFile := fmt.Sprintf("%s-linux-%s-%s", name, arch, version)
	skipBindingsFlag, skipFrontendFlag := "", ""
	if skipFrontend {
		skipBindingsFlag, skipFrontendFlag = "-skipbindings", "-s"
	}

	cc, cxx := getCompilerFlags(arch)
	crossCompileFlags := map[string]string{
		"GOOS":   "linux",
		"GOARCH": arch,
		"CC":     cc,
		"CXX":    cxx,
	}

	return sh.RunWithV(crossCompileFlags, "wails", "build", "-m", "-nosyncgomod", "-ldflags", ldFlags, "-platform", fmt.Sprintf("linux/%s", arch), "-o", outputFile, skipBindingsFlag, skipFrontendFlag)
}

func getCompilerFlags(arch string) (cc, cxx string) {
	if arch == "arm64" {
		return "aarch64-linux-gnu-gcc", "aarch64-linux-gnu-g++"
	}

	return "x86_64-linux-gnu-gcc", "x86_64-linux-gnu-g++"
}
