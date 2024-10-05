package main

import (
	"fmt"
	"runtime"
)

func Build(version, timestamp, commitSha, link, buildtype string) error {
	config, err := configureWailsProject(version)
	if err != nil {
		return err
	}

	switch runtime.GOOS {
	case "linux", "windows":
		if err := buildLinux(config.Name, version, timestamp, commitSha, link, buildtype); err != nil {
			return err
		}

		if err := buildWindows(config.Name, version, timestamp, commitSha, link, buildtype); err != nil {
			return err
		}
	case "darwin":
		return buildMacOS(config.Name, version, timestamp, commitSha, link, buildtype)
	default:
		return fmt.Errorf("unsupported OS/architecture: %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	return nil
}

func buildFlags(version, timestamp, commitSha, link, buildType string) string {
	return fmt.Sprintf("-X 'main.BuildType=%s' -X 'main.Version=%s' -X 'main.BuildTimestamp=%s' -X 'main.CommitSha=%s' -X 'main.BuildLink=%s'", buildType, version, timestamp, commitSha, link)
}
