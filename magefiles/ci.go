//go:build mage

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"runtime"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

func Build(buildType, appVersion, buildTimestamp, commitSha, buildLink string) error {
	if buildType != "release" && buildType != "debug" {
		return fmt.Errorf("Invalid build type: %s. Expected either \"release\" or \"debug\"", buildType)
	}

	var ldFlags = fmt.Sprintf("-X 'main.BuildType=%s' -X 'main.Version=%s' -X 'main.BuildTimestamp=%s' -X 'main.CommitSha=%s' -X 'main.BuildLink=%s'", buildType, appVersion, buildTimestamp, commitSha, buildLink)

	if runtime.GOOS == "linux" {
		mg.Deps(mg.F(configureWailsProject, appVersion))

		if runtime.GOARCH == "amd64" {
			if err := buildLinuxAMD64(ldFlags, appVersion, false); err != nil {
				return fmt.Errorf("Error building Linux AMD64: %s", err)
			}

			if err := buildLinuxARM64(ldFlags, appVersion, true); err != nil {
				return fmt.Errorf("Error building Linux ARM64: %s", err)
			}
		} else {
			if err := buildLinuxARM64(ldFlags, appVersion, false); err != nil {
				return fmt.Errorf("Error building Linux ARM64: %s", err)
			}

			if err := buildLinuxAMD64(ldFlags, appVersion, true); err != nil {
				return fmt.Errorf("Error building Linux AMD64: %s", err)
			}
		}

		var buildWindowsErr = buildWindowsAMD64(ldFlags, appVersion, true)

		if buildWindowsErr != nil {
			buildWindowsErr = fmt.Errorf("Error building Windows AMD64: %s", buildWindowsErr)
		}

		return buildWindowsErr
	} else if runtime.GOOS == "darwin" {
		mg.Deps(mg.F(configureWailsProject, appVersion))

		return buildDarwinUniversal(ldFlags, appVersion)
	} else {
		return fmt.Errorf("Unsupported OS/architecture: %s/%s", runtime.GOOS, runtime.GOARCH)
	}
}

func configureWailsProject(releaseVersion string) error {
	var nonTaggedReleaseVersion, error = regexp.Compile("^v(\\d+\\.\\d+\\.\\d+)-(.+)$")

	if error != nil {
		fmt.Println("Error compiling regex", error)
		return error
	}

	var taggedReleaseVersion, error2 = regexp.Compile("^v(\\d+\\.\\d+\\.\\d+)$")

	if error2 != nil {
		fmt.Println("Error compiling regex", error2)
		return error2
	}

	var nsisCompliantVersion = ""

	if nonTaggedReleaseVersion.MatchString(releaseVersion) == true {
		nsisCompliantVersion = nonTaggedReleaseVersion.ReplaceAllString(releaseVersion, "$1.$2")
	} else if taggedReleaseVersion.MatchString(releaseVersion) == true {
		nsisCompliantVersion = taggedReleaseVersion.ReplaceAllString(releaseVersion, "$1.0")
	} else {
		return fmt.Errorf("Invalid release version: %s. Expected semantic release in one of the following two formats: vX.X.X or vX.X.X-X-XXXXXXX", releaseVersion)
	}

	fmt.Printf("NSIS compatible version: [%s]\n", nsisCompliantVersion)

	type WailsProjectConfigAuthor struct {
		Name string `json:"name"`
	}

	type FileAssociation struct {
		Ext         string `json:"ext"`
		Name        string `json:"name"`
		Description string `json:"description"`
		IconName    string `json:"iconName"`
		Role        string `json:"role"`
	}

	type WailsProjectConfigInfo struct {
		CompanyName      string            `json:"companyName"`
		ProductVersion   string            `json:"productVersion"`
		Copyright        string            `json:"copyright"`
		Comments         string            `json:"comments"`
		FileAssociations []FileAssociation `json:"fileAssociations"`
	}

	type WailsProjectConfig struct {
		Schema               string                   `json:"$schema"`
		Name                 string                   `json:"name"`
		OutputFilename       string                   `json:"outputfilename"`
		FrontendInstall      string                   `json:"frontend:install"`
		FrontendBuild        string                   `json:"frontend:build"`
		FrontendDevWatcher   string                   `json:"frontend:dev:watcher"`
		FrontendDevServerUrl string                   `json:"frontend:dev:serverUrl"`
		Author               WailsProjectConfigAuthor `json:"author"`
		Info                 WailsProjectConfigInfo   `json:"info"`
	}

	fmt.Println("Reading Wails Config")
	var wailsConfigJson, read_error = os.ReadFile("wails.json")

	if read_error != nil {
		fmt.Println("Error reading wails.json", read_error)
		return read_error
	}

	var wailsConfig WailsProjectConfig

	var parse_error = json.Unmarshal(wailsConfigJson, &wailsConfig)

	if parse_error != nil {
		fmt.Println("Error parsing wails.json", parse_error)
		return parse_error
	}

	fmt.Println("Setting Wails Product Version")
	wailsConfig.Info.ProductVersion = nsisCompliantVersion

	var updatedWailsConfig, marshal_error = json.MarshalIndent(wailsConfig, "", "  ")

	if marshal_error != nil {
		fmt.Println("Error marshalling wails.json", marshal_error)
		return marshal_error
	}

	fmt.Println("Writing Wails Config")
	return os.WriteFile("wails.json", updatedWailsConfig, os.ModePerm)
}

func buildLinuxAMD64(ldFlags, appVersion string, skipFrontend bool) error {
	var outputFileName = fmt.Sprintf("rocketblend-desktop-linux-amd64-%s", appVersion)

	skipBindingsFlag := ""
	skipFrontendFlag := ""

	if skipFrontend == true {
		skipBindingsFlag = "-skipbindings"
		skipFrontendFlag = "-s"
	}

	crossCompileFlags := map[string]string{"GOOS": "linux", "GOARCH": "amd64", "CC": "x86_64-linux-gnu-gcc", "CXX": "x86_64-linux-gnu-g++"}

	return sh.RunWithV(crossCompileFlags, "wails", "build", "-m", "-nosyncgomod", "-ldflags", ldFlags, "-platform", "linux/amd64", "-o", outputFileName, skipBindingsFlag, skipFrontendFlag)
}

func buildLinuxARM64(ldFlags, appVersion string, skipFrontend bool) error {
	var outputFileName = fmt.Sprintf("rocketblend-desktop-linux-arm64-%s", appVersion)

	skipBindingsFlag := ""
	skipFrontendFlag := ""

	if skipFrontend == true {
		skipBindingsFlag = "-skipbindings"
		skipFrontendFlag = "-s"
	}

	crossCompileFlags := map[string]string{"GOOS": "linux", "GOARCH": "arm64", "CC": "aarch64-linux-gnu-gcc", "CXX": "aarch64-linux-gnu-g++"}

	return sh.RunWithV(crossCompileFlags, "wails", "build", "-m", "-nosyncgomod", "-ldflags", ldFlags, "-platform", "linux/arm64", "-o", outputFileName, skipBindingsFlag, skipFrontendFlag)
}

func buildWindowsAMD64(ldFlags, appVersion string, skipFrontend bool) error {
	var outputFileName = fmt.Sprintf("rocketblend-desktop-windows-amd64-%s.exe", appVersion)

	skipBindingsFlag := ""
	skipFrontendFlag := ""

	if skipFrontend == true {
		skipBindingsFlag = "-skipbindings"
		skipFrontendFlag = "-s"
	}

	crossCompileFlags := map[string]string{"GOOS": "windows", "GOARCH": "amd64", "CC": "x86_64-w64-mingw32-gcc", "CXX": "x86_64-w64-mingw32-g++"}

	return sh.RunWithV(crossCompileFlags, "wails", "build", "-m", "-nosyncgomod", "-ldflags", ldFlags, "-nsis", "-platform", "windows/amd64", "-o", outputFileName, skipBindingsFlag, skipFrontendFlag)
}

func buildDarwinUniversal(ldFlags, appVersion string) error {
	if err := sh.RunV("wails", "build", "-m", "-nosyncgomod", "-ldflags", ldFlags, "-platform", "darwin/universal"); err != nil {
		return fmt.Errorf("Error building Darwin Wails App: %s", err)
	}

	fmt.Println("Building DMG")
	var dmgOutputPath = fmt.Sprintf("./build/bin/rocketblend-desktop-darwin-universal-%s.dmg", appVersion)
	if err := sh.RunV("create-dmg", "--window-size", "800", "300", "--no-internet-enable", "--hide-extension", "RocketBlend-Desktop.app", "--app-drop-link", "600", "40", dmgOutputPath, "./build/bin/RocketBlend-Desktop.app"); err != nil {
		return fmt.Errorf("Error building DMG: %s", err)
	}

	fmt.Println("Compiling seticon.swift")
	if err := sh.Run("swiftc", "./build/darwin/seticon.swift"); err != nil {
		return fmt.Errorf("Error compiling seticon with Swift: %s", err)
	}

	if err := sh.Run("chmod", "+x", "./seticon"); err != nil {
		return fmt.Errorf("Error setting permissions on seticon: %s", err)
	}

	fmt.Println("Setting DMG icons")
	return sh.RunV("./seticon", "./build/bin/RocketBlend-Desktop.app/Contents/Resources/iconfile.icns", dmgOutputPath)
}
