package main

import (
	"fmt"
	"runtime"

	"github.com/magefile/mage/mg"
)

// Build compiles the project based on the given parameters.
func Build(buildType, appVersion, buildTimestamp, commitSha, buildLink string) error {
	if buildType != "release" && buildType != "debug" {
		return fmt.Errorf("invalid build type: %s. Expected either \"release\" or \"debug\"", buildType)
	}

	ldFlags := fmt.Sprintf("-X 'main.BuildType=%s' -X 'main.Version=%s' -X 'main.BuildTimestamp=%s' -X 'main.CommitSha=%s' -X 'main.BuildLink=%s'", buildType, appVersion, buildTimestamp, commitSha, buildLink)

	switch runtime.GOOS {
	case "linux":
		mg.Deps(mg.F(configureWailsProject, appVersion))

		if runtime.GOARCH == "amd64" {
			if err := buildLinuxAMD64(ldFlags, appVersion, false); err != nil {
				return fmt.Errorf("error building Linux AMD64: %s", err)
			}

			if err := buildLinuxARM64(ldFlags, appVersion, true); err != nil {
				return fmt.Errorf("error building Linux ARM64: %s", err)
			}
		} else {
			if err := buildLinuxARM64(ldFlags, appVersion, false); err != nil {
				return fmt.Errorf("error building Linux ARM64: %s", err)
			}

			if err := buildLinuxAMD64(ldFlags, appVersion, true); err != nil {
				return fmt.Errorf("error building Linux AMD64: %s", err)
			}
		}

		if err := buildWindowsAMD64(ldFlags, appVersion, true, false); err != nil {
			return fmt.Errorf("error building Windows AMD64: %s", err)
		}

		return nil
	case "darwin":
		mg.Deps(mg.F(configureWailsProject, appVersion))
		// TODO: Grab name from config.
		releaseBuild := buildType == "release"
		return buildMacOSApp(ldFlags, "rocketblend-desktop", appVersion, releaseBuild)
	default:
		return fmt.Errorf("unsupported OS/architecture: %s/%s", runtime.GOOS, runtime.GOARCH)
	}
}
