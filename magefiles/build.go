package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"runtime"

	"github.com/magefile/mage/mg"
)

type (
	// Author represents the author of the Wails project.
	Author struct {
		Name string `json:"name"`
	}

	// FileAssociation represents a file association in the Wails project configuration.
	FileAssociation struct {
		Ext         string `json:"ext"`
		Name        string `json:"name"`
		Description string `json:"description"`
		IconName    string `json:"iconName"`
		Role        string `json:"role"`
	}

	// Info contains various information about the Wails project.
	Info struct {
		CompanyName      string            `json:"companyName"`
		ProductVersion   string            `json:"productVersion"`
		Copyright        string            `json:"copyright"`
		Comments         string            `json:"comments"`
		FileAssociations []FileAssociation `json:"fileAssociations"`
	}

	// WailsConfig represents the overall configuration structure for a Wails project.
	WailsConfig struct {
		Schema               string `json:"$schema"`
		Name                 string `json:"name"`
		OutputFilename       string `json:"outputfilename"`
		FrontendInstall      string `json:"frontend:install"`
		FrontendBuild        string `json:"frontend:build"`
		FrontendDevWatcher   string `json:"frontend:dev:watcher"`
		FrontendDevServerUrl string `json:"frontend:dev:serverUrl"`
		WailsJSDir           string `json:"wailsjsdir"`
		Author               Author `json:"author"`
		Info                 Info   `json:"info"`
	}
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

// configureWailsProject configures the Wails project based on the given version.
func configureWailsProject(releaseVersion string) error {
	nonTaggedReleaseVersion, err := regexp.Compile(`^v(\d+\.\d+\.\d+)-(.+)$`)
	if err != nil {
		fmt.Printf("Error compiling regex: %v\n", err)
		return err
	}

	taggedReleaseVersion, err := regexp.Compile(`^v(\d+\.\d+\.\d+)$`)
	if err != nil {
		fmt.Printf("Error compiling regex: %v\n", err)
		return err
	}

	nsisCompliantVersion := ""
	if nonTaggedReleaseVersion.MatchString(releaseVersion) {
		nsisCompliantVersion = nonTaggedReleaseVersion.ReplaceAllString(releaseVersion, "$1.$2")
	} else if taggedReleaseVersion.MatchString(releaseVersion) {
		nsisCompliantVersion = taggedReleaseVersion.ReplaceAllString(releaseVersion, "$1.0")
	} else {
		return fmt.Errorf("invalid release version: %s. Expected semantic release in one of the following two formats: vX.X.X or vX.X.X-X-XXXXXXX", releaseVersion)
	}

	fmt.Printf("NSIS compatible version: [%s]\n", nsisCompliantVersion)

	wailsConfigJSON, err := os.ReadFile("wails.json")
	if err != nil {
		fmt.Printf("Error reading wails.json: %v\n", err)
		return err
	}

	var wailsConfig WailsConfig
	if err := json.Unmarshal(wailsConfigJSON, &wailsConfig); err != nil {
		fmt.Printf("Error parsing wails.json: %v\n", err)
		return err
	}

	wailsConfig.Info.ProductVersion = nsisCompliantVersion

	updatedWailsConfigJSON, err := json.MarshalIndent(wailsConfig, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling wails.json: %v\n", err)
		return err
	}

	return os.WriteFile("wails.json", updatedWailsConfigJSON, os.ModePerm)
}
