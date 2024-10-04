package main

import (
	"fmt"
	"runtime"
)

func buildRelease(version, timestamp, commitSha, link, outputDir string, debug bool) error {
	config, err := configureWailsProject(version)
	if err != nil {
		return err
	}

	switch runtime.GOOS {
	case "linux", "windows":
		if err := buildReleaseLinux(config.Name, version, timestamp, commitSha, link, outputDir, debug); err != nil {
			return err
		}

		if err := buildReleaseWindows(config.Name, version, timestamp, commitSha, link, outputDir, debug); err != nil {
			return err
		}
	case "darwin":
		return buildReleaseMacOS(config.Name, version, timestamp, commitSha, link, outputDir, debug)
	default:
		return fmt.Errorf("unsupported OS/architecture: %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	return nil
}

func buildFlags(version, timestamp, commitSha, link, outputDir string, debug bool) string {
	buildType := "release"
	if debug {
		buildType = "debug"
	}

	return fmt.Sprintf("-X 'main.BuildType=%s' -X 'main.Version=%s' -X 'main.BuildTimestamp=%s' -X 'main.CommitSha=%s' -X 'main.BuildLink=%s'", buildType, version, timestamp, commitSha, link)
}
