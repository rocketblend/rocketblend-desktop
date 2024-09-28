//go:build mage

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"runtime"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
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

		if err := buildWindowsAMD64(ldFlags, appVersion, true, false); err != nil {
			return fmt.Errorf("Error building Windows AMD64: %s", err)
		}

		return nil
	case "darwin":
		mg.Deps(mg.F(configureWailsProject, appVersion))
		return buildDarwinUniversal(ldFlags, appVersion, false)
	default:
		return fmt.Errorf("unsupported OS/architecture: %s/%s", runtime.GOOS, runtime.GOARCH)
	}
}

// configureWailsProject configures the Wails project based on the given version.
func configureWailsProject(releaseVersion string) error {
	nonTaggedReleaseVersion, err := regexp.Compile("^v(\\d+\\.\\d+\\.\\d+)-(.+)$")
	if err != nil {
		fmt.Printf("Error compiling regex: %v\n", err)
		return err
	}

	taggedReleaseVersion, err := regexp.Compile("^v(\\d+\\.\\d+\\.\\d+)$")
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

// buildLinuxAMD64 builds the Linux AMD64 version of the project.
func buildLinuxAMD64(ldFlags, appVersion string, skipFrontend bool) error {
	outputFileName := fmt.Sprintf("rocketblend-desktop-linux-amd64-%s", appVersion)
	skipBindingsFlag, skipFrontendFlag := "", ""
	if skipFrontend {
		skipBindingsFlag, skipFrontendFlag = "-skipbindings", "-s"
	}

	crossCompileFlags := map[string]string{"GOOS": "linux", "GOARCH": "amd64", "CC": "x86_64-linux-gnu-gcc", "CXX": "x86_64-linux-gnu-g++"}
	return sh.RunWithV(crossCompileFlags, "wails", "build", "-m", "-nosyncgomod", "-ldflags", ldFlags, "-platform", "linux/amd64", "-o", outputFileName, skipBindingsFlag, skipFrontendFlag)
}

// buildLinuxARM64 builds the Linux ARM64 version of the project.
func buildLinuxARM64(ldFlags, appVersion string, skipFrontend bool) error {
	outputFileName := fmt.Sprintf("rocketblend-desktop-linux-arm64-%s", appVersion)
	skipBindingsFlag, skipFrontendFlag := "", ""
	if skipFrontend {
		skipBindingsFlag, skipFrontendFlag = "-skipbindings", "-s"
	}

	crossCompileFlags := map[string]string{"GOOS": "linux", "GOARCH": "arm64", "CC": "aarch64-linux-gnu-gcc", "CXX": "aarch64-linux-gnu-g++"}
	return sh.RunWithV(crossCompileFlags, "wails", "build", "-m", "-nosyncgomod", "-ldflags", ldFlags, "-platform", "linux/arm64", "-o", outputFileName, skipBindingsFlag, skipFrontendFlag)
}

// buildWindowsAMD64 builds the Windows AMD64 version of the project.
func buildWindowsAMD64(ldFlags, appVersion string, skipFrontend bool, sign bool) error {
	outputFileName := fmt.Sprintf("rocketblend-desktop-windows-amd64-%s.exe", appVersion)
	skipBindingsFlag, skipFrontendFlag := "", ""
	if skipFrontend {
		skipBindingsFlag, skipFrontendFlag = "-skipbindings", "-s"
	}

	crossCompileFlags := map[string]string{
		"GOOS":   "windows",
		"GOARCH": "amd64",
		"CC":     "x86_64-w64-mingw32-gcc",
		"CXX":    "x86_64-w64-mingw32-g++",
	}

	err := sh.RunWithV(crossCompileFlags, "wails", "build", "-m", "-nosyncgomod", "-ldflags", ldFlags, "-nsis", "-platform", "windows/amd64", "-o", outputFileName, skipBindingsFlag, skipFrontendFlag)
	if err != nil {
		return fmt.Errorf("error building Windows AMD64: %v", err)
	}

	if sign {
		// Recommended to just use Azure Trusted Certificate for signing instead.
		// Keeping for reference.
		fmt.Println("Importing Code Signing Certificates")
		certFilePath := "certificate/certificate.pfx"
		pemFilePath := "certificate/certificate.pem"
		signCert := os.Getenv("WIN_CERTIFICATE")
		signCertPassword := os.Getenv("WIN_CERTIFICATE_PASSWORD")

		if signCert == "" || signCertPassword == "" {
			return fmt.Errorf("missing required environment variables for code-signing")
		}

		if _, err := os.Stat("certificate"); os.IsNotExist(err) {
			if err := os.Mkdir("certificate", os.ModePerm); err != nil {
				return fmt.Errorf("error creating certificate directory: %v", err)
			}
		}

		if err := os.WriteFile("certificate/certificate.txt", []byte(signCert), 0600); err != nil {
			return fmt.Errorf("error writing base64 certificate to file: %v", err)
		}

		if err := sh.Run("certutil", "-decode", "certificate/certificate.txt", certFilePath); err != nil {
			return fmt.Errorf("error decoding certificate: %v", err)
		}

		if err := sh.Run("openssl", "pkcs12", "-in", certFilePath, "-out", pemFilePath, "-nodes", "-passin", fmt.Sprintf("pass:%s", signCertPassword)); err != nil {
			return fmt.Errorf("error converting PFX to PEM: %v", err)
		}

		fmt.Println("Signing Build")
		if err := sh.Run("osslsigncode", "sign", "-certs", pemFilePath, "-key", pemFilePath, "-pass", signCertPassword, "-in", outputFileName, "-out", outputFileName, "-t", "http://timestamp.digicert.com", "-h", "sha256"); err != nil {
			return fmt.Errorf("error signing executable: %v", err)
		}
	}

	return nil
}

// buildDarwinUniversal builds the Darwin universal version of the project.
func buildDarwinUniversal(ldFlags, appVersion string, sign bool) error {
	if err := sh.RunV("wails", "build", "-m", "-nosyncgomod", "-ldflags", ldFlags, "-platform", "darwin/universal"); err != nil {
		return fmt.Errorf("error building Darwin Wails App: %v", err)
	}

	if sign {
		fmt.Println("Importing Code Signing Certificates")
		if err := importDarwinCodeSigningCertificates(os.Getenv("AC_CERTIFICATE"), os.Getenv("AC_CERTIFICATE_PASSWORD")); err != nil {
			return fmt.Errorf("error importing code signing certificates: %v", err)
		}

		fmt.Println("Signing Build")
		if err := sh.RunV("gon", "-log-level=info", "./build/darwin/gon-sign.json"); err != nil {
			return fmt.Errorf("error signing build: %v", err)
		}
	}

	fmt.Println("Building DMG")
	dmgOutputPath := fmt.Sprintf("./build/bin/rocketblend-desktop-darwin-universal-%s.dmg", appVersion)
	if err := sh.RunV("create-dmg", "--window-size", "800", "300", "--no-internet-enable", "--hide-extension", "rocketblend-desktop.app", "--app-drop-link", "600", "40", dmgOutputPath, "./build/bin/rocketblend-desktop.app"); err != nil {
		return fmt.Errorf("error building DMG: %v", err)
	}

	fmt.Println("Compiling seticon.swift")
	if err := sh.Run("swiftc", "./build/darwin/seticon.swift"); err != nil {
		return fmt.Errorf("error compiling seticon with Swift: %v", err)
	}

	if err := sh.Run("chmod", "+x", "./seticon"); err != nil {
		return fmt.Errorf("error setting permissions on seticon: %v", err)
	}

	fmt.Println("Setting DMG icons")
	return sh.RunV("./seticon", "./build/bin/rocketblend-desktop.app/Contents/Resources/iconfile.icns", dmgOutputPath)
}

// importDarwinCodeSigningCertificates imports the code signing certificates into the keychain.
func importDarwinCodeSigningCertificates(certBase64, password string) error {
	if certBase64 == "" || password == "" {
		return fmt.Errorf("missing required environment variables for code-signing")
	}

	certBytes, err := base64.StdEncoding.DecodeString(certBase64)
	if err != nil {
		return fmt.Errorf("error decoding base64 certificate: %v", err)
	}

	certFilePath := "./cert.p12"
	if err := os.WriteFile(certFilePath, certBytes, 0600); err != nil {
		return fmt.Errorf("error writing certificate to file: %v", err)
	}
	defer os.Remove(certFilePath)

	if err := sh.Run("security", "import", certFilePath, "-P", password, "-T", "/usr/bin/codesign"); err != nil {
		return fmt.Errorf("error importing certificate: %v", err)
	}

	return nil
}
