//go:build mage

package main

import (
	"fmt"
	"runtime"

	"github.com/magefile/mage/sh"
)

// buildWindows builds the Windows binary and NSIS installer for the host architecture.
func buildWindows(name, version, timestamp, commitSha, link, buildType string) error {
	arch := runtime.GOARCH
	fmt.Printf("Building Windows %s binary for %s\n", arch, name)

	ldFlags := buildFlags(version, timestamp, commitSha, link, buildType)
	outputFile := fmt.Sprintf("%s-windows-%s-%s.exe", name, arch, version)

	if err := sh.RunV("wails", "build", "-m", "-nosyncgomod", "-ldflags", ldFlags, "-nsis", "-platform", fmt.Sprintf("windows/%s", arch), "-o", outputFile); err != nil {
		return fmt.Errorf("error building Windows %s: %w", arch, err)
	}

	return nil
}
