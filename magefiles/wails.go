package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
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

// configureWailsProject configures the Wails project config based on release information.
func configureWailsProject(releaseVersion string) (*WailsConfig, error) {
	nonTaggedReleaseVersion, taggedReleaseVersion, err := compileVersionRegex()
	if err != nil {
		return nil, err
	}

	nsisCompliantVersion, err := getNSISCompliantVersion(releaseVersion, nonTaggedReleaseVersion, taggedReleaseVersion)
	if err != nil {
		return nil, err
	}
	fmt.Printf("NSIS compatible version: [%s]\n", nsisCompliantVersion)

	wailsConfig, err := readAndParseWailsConfig("wails.json")
	if err != nil {
		return nil, err
	}

	wailsConfig.Info.ProductVersion = nsisCompliantVersion
	if err := writeWailsConfig("wails.json", wailsConfig); err != nil {
		return nil, err
	}

	return wailsConfig, nil
}

// compileVersionRegex compiles and returns the version regex patterns.
func compileVersionRegex() (*regexp.Regexp, *regexp.Regexp, error) {
	nonTaggedReleaseVersion, err := regexp.Compile(`^v(\d+\.\d+\.\d+)-(.+)$`)
	if err != nil {
		return nil, nil, fmt.Errorf("error compiling non-tagged release version regex: %v", err)
	}

	taggedReleaseVersion, err := regexp.Compile(`^v(\d+\.\d+\.\d+)$`)
	if err != nil {
		return nil, nil, fmt.Errorf("error compiling tagged release version regex: %v", err)
	}

	return nonTaggedReleaseVersion, taggedReleaseVersion, nil
}

// getNSISCompliantVersion returns the NSIS-compliant version based on the given release version.
func getNSISCompliantVersion(releaseVersion string, nonTaggedReleaseVersion, taggedReleaseVersion *regexp.Regexp) (string, error) {
	if nonTaggedReleaseVersion.MatchString(releaseVersion) {
		return nonTaggedReleaseVersion.ReplaceAllString(releaseVersion, "$1.$2"), nil
	} else if taggedReleaseVersion.MatchString(releaseVersion) {
		return taggedReleaseVersion.ReplaceAllString(releaseVersion, "$1.0"), nil
	} else {
		return "", fmt.Errorf("invalid release version: %s. Expected semantic release in one of the following formats: vX.X.X or vX.X.X-X-XXXXXXX", releaseVersion)
	}
}

// readAndParseWailsConfig reads the wails.json file and parses it into a WailsConfig struct.
func readAndParseWailsConfig(filePath string) (*WailsConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %v", filePath, err)
	}

	var config WailsConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing %s: %v", filePath, err)
	}

	return &config, nil
}

// writeWailsConfig writes the updated WailsConfig back to the wails.json file.
func writeWailsConfig(filePath string, config *WailsConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling %s: %v", filePath, err)
	}

	if err := os.WriteFile(filePath, data, os.ModePerm); err != nil {
		return fmt.Errorf("error writing %s: %v", filePath, err)
	}

	return nil
}
