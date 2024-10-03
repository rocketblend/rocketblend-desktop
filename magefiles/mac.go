package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/magefile/mage/sh"
)

const BundleID = "io.rocketblend.rocketblend-desktop"

// EnvironmentVariables holds the necessary environment variables for the macOS build process.
type EnvironmentVariables struct {
	Certificate         string
	CertificatePassword string
	DeveloperID         string
	AppleID             string
	Password            string
	TeamID              string
}

// buildMacOSApp builds the macOS universal version of the project and handles signing and notarization.
func buildMacOSApp(ldFlags, appName, appVersion string, releaseBuild bool) error {
	env, err := getEnvVariables()
	if err != nil {
		return err
	}

	if err := buildMacOSWailsApp(ldFlags); err != nil {
		return fmt.Errorf("error building Wails app for macOS: %v", err)
	}

	if releaseBuild {
		appPath := fmt.Sprintf("./build/bin/%s.app", appName)

		if err := importMacOSCertificate(env.Certificate, env.CertificatePassword); err != nil {
			return fmt.Errorf("error importing certificate: %v", err)
		}

		if err := signMacOSFile(appPath, BundleID, env.DeveloperID, "./build/darwin/entitlements.plist"); err != nil {
			return fmt.Errorf("error signing .app: %v", err)
		}

		dmgOutputPath, err := createMacOSDMG(appPath, appName, appVersion)
		if err != nil {
			return err
		}

		if err := signMacOSFile(dmgOutputPath, BundleID, env.DeveloperID, ""); err != nil {
			return fmt.Errorf("error signing DMG: %v", err)
		}

		if err := notarizeMacOSDMG(dmgOutputPath, env.AppleID, env.Password, env.TeamID); err != nil {
			return fmt.Errorf("error submitting DMG for notarization: %v", err)
		}

		if err := stapleMacOSNotarization(dmgOutputPath); err != nil {
			return fmt.Errorf("error stapling notarization ticket: %v", err)
		}
	} else {
		fmt.Println("Skipping signing and notarization for this build.")
	}

	return nil
}

// getEnvVariables gathers and validates the necessary environment variables for the macOS build process.
func getEnvVariables() (EnvironmentVariables, error) {
	env := EnvironmentVariables{
		Certificate:         os.Getenv("AC_CERTIFICATE"),
		CertificatePassword: os.Getenv("AC_CERTIFICATE_PASSWORD"),
		DeveloperID:         os.Getenv("AC_DEVELOPER_ID"),
		AppleID:             os.Getenv("AC_APPLE_ID"),
		Password:            os.Getenv("AC_PASSWORD"),
		TeamID:              os.Getenv("AC_TEAM_ID"),
	}

	if env.Certificate == "" || env.CertificatePassword == "" || env.DeveloperID == "" || env.AppleID == "" || env.Password == "" || env.TeamID == "" {
		return env, fmt.Errorf("missing required environment variables")
	}

	return env, nil
}

// buildMacOSWailsApp compiles the Wails app for macOS with the given ldflags.
func buildMacOSWailsApp(ldFlags string) error {
	fmt.Println("Building Wails app for macOS")
	return sh.RunV("wails", "build", "-m", "-nosyncgomod", "-ldflags", ldFlags, "-platform", "darwin/universal")
}

// importMacOSCertificate imports the certificate into the keychain for code signing.
func importMacOSCertificate(cert, certPassword string) error {
	certBytes, err := base64.StdEncoding.DecodeString(cert)
	if err != nil {
		return fmt.Errorf("error decoding base64 certificate: %v", err)
	}

	certFilePath := "./cert.p12"
	if err := os.WriteFile(certFilePath, certBytes, 0600); err != nil {
		return fmt.Errorf("error writing certificate to file: %v", err)
	}
	defer os.Remove(certFilePath)

	if err := sh.Run("security", "import", certFilePath, "-P", certPassword, "-T", "/usr/bin/codesign"); err != nil {
		return fmt.Errorf("error importing certificate: %v", err)
	}

	return sh.Run("security", "set-key-partition-list", "-S", "apple-tool:,apple:", "-s", "-k", certPassword, certFilePath)
}

// signMacOSFile signs any file (app, DMG, etc.) with the Developer ID Application certificate.
func signMacOSFile(filePath, bundleID, developerID, entitlementsPath string) error {
	fmt.Printf("Signing file: %s\n", filePath)

	args := []string{"--deep", "--force", "--options", "runtime", "--sign", developerID, "--timestamp", "--identifier", bundleID}
	if entitlementsPath != "" {
		args = append(args, "--entitlements", entitlementsPath)
	}

	args = append(args, filePath)
	return sh.RunV("codesign", args...)
}

// createMacOSDMG creates a DMG from the given .app bundle and sets the icon.
func createMacOSDMG(appPath, appName, appVersion string) (string, error) {
	dmgOutputPath := fmt.Sprintf("./build/bin/%s-darwin-universal-%s.dmg", appName, appVersion)

	fmt.Println("Building DMG for macOS")
	if err := sh.RunV("create-dmg", "--window-size", "800", "300", "--no-internet-enable", "--hide-extension", "--app-drop-link", "600", "40", dmgOutputPath, appPath); err != nil {
		return "", fmt.Errorf("error building DMG for macOS: %v", err)
	}

	if err := compileSetIconSwiftScript(); err != nil {
		return "", fmt.Errorf("error compiling seticon Swift script: %v", err)
	}

	iconFilePath := fmt.Sprintf("./build/bin/%s.app/Contents/Resources/iconfile.icns", appName)
	if err := setDMGIcons(iconFilePath, dmgOutputPath); err != nil {
		return "", fmt.Errorf("error setting DMG icons: %v", err)
	}

	return dmgOutputPath, nil
}

// notarizeMacOSDMG submits the DMG for macOS notarization using credentials passed from environment variables.
func notarizeMacOSDMG(dmgPath, appleID, password, teamID string) error {
	fmt.Println("Submitting DMG for macOS notarization")
	return sh.RunV("xcrun", "notarytool", "submit", dmgPath, "--wait", "--apple-id", appleID, "--password", password, "--team-id", teamID)
}

// stapleMacOSNotarization staples the notarization ticket to the macOS DMG.
func stapleMacOSNotarization(dmgPath string) error {
	fmt.Println("Stapling notarization ticket to macOS DMG")
	return sh.RunV("xcrun", "stapler", "staple", dmgPath)
}

// compileSetIconSwiftScript compiles the seticon.swift script for macOS DMG.
func compileSetIconSwiftScript() error {
	fmt.Println("Compiling seticon.swift for macOS DMG")
	if err := sh.Run("swiftc", "./build/darwin/seticon.swift"); err != nil {
		return fmt.Errorf("error compiling seticon for macOS DMG: %v", err)
	}

	if err := sh.Run("chmod", "+x", "./seticon"); err != nil {
		return fmt.Errorf("error setting permissions on seticon: %v", err)
	}

	return nil
}

// setDMGIcons runs the compiled seticon.swift script to set the icon for the DMG.
func setDMGIcons(iconFilePath, dmgOutputPath string) error {
	fmt.Println("Setting icons for macOS DMG")
	if err := sh.RunV("./seticon", iconFilePath, dmgOutputPath); err != nil {
		return fmt.Errorf("error setting icons for macOS DMG: %v", err)
	}

	return nil
}
