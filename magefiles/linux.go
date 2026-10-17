//go:build mage

package main

import (
	"fmt"
	"runtime"

	"github.com/magefile/mage/sh"
)

// buildLinux builds the Linux binary for the host architecture.
func buildLinux(name, version, timestamp, commitSha, link, buildtype string) error {
	arch := runtime.GOARCH
	fmt.Printf("Building Linux %s binary for %s\n", arch, name)

	ldFlags := buildFlags(version, timestamp, commitSha, link, buildtype)
	outputFile := fmt.Sprintf("%s-linux-%s-%s", name, arch, version)

	// webkit2_41 links against webkit2gtk-4.1, as current distros no longer ship webkit2gtk-4.0.
	if err := sh.RunV("wails", "build", "-m", "-nosyncgomod", "-tags", "webkit2_41", "-ldflags", ldFlags, "-platform", fmt.Sprintf("linux/%s", arch), "-o", outputFile); err != nil {
		return fmt.Errorf("error building Linux %s: %w", arch, err)
	}

	return nil
}
