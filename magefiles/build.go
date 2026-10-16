package main

import (
	"fmt"
	"runtime"
)

func Build(version, timestamp, commitSha, link, buildtype string) error {
	return build(version, timestamp, commitSha, link, buildtype, false)
}

func build(version, timestamp, commitSha, link, buildtype string, obfuscate bool) error {
	config, err := configureWailsProject(version, obfuscate)
	if err != nil {
		return err
	}

	// Each platform is built natively on its own runner, so only the host target is built.
	switch runtime.GOOS {
	case "linux":
		return buildLinux(config.Name, version, timestamp, commitSha, link, buildtype)
	case "windows":
		return buildWindows(config.Name, version, timestamp, commitSha, link, buildtype)
	case "darwin":
		return buildMacOS(config.Name, version, timestamp, commitSha, link, buildtype)
	default:
		return fmt.Errorf("unsupported OS/architecture: %s/%s", runtime.GOOS, runtime.GOARCH)
	}
}

func buildFlags(version, timestamp, commitSha, link, buildType string) string {
	return fmt.Sprintf("-X 'main.BuildType=%s' -X 'main.Version=%s' -X 'main.BuildTimestamp=%s' -X 'main.CommitSha=%s' -X 'main.BuildLink=%s'", buildType, version, timestamp, commitSha, link)
}
