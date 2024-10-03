package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/magefile/mage/sh"
)

// buildMacOSApp builds the macOS universal version of the project.
func buildMacOSApp(ldFlags, appName, appVersion string, releaseBuild bool) error {
	if err := buildMacOSWailsApp(ldFlags); err != nil {
		return fmt.Errorf("error building Wails app for macOS: %v", err)
	}

	if releaseBuild {
		if err := importAndSignMacOSCode(); err != nil {
			return err
		}

		dmgOutputPath, err := buildAndSignMacOSDMG(appVersion, appName)
		if err != nil {
			return err
		}

		if err := notarizeMacOSDMG(dmgOutputPath); err != nil {
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

// buildMacOSWailsApp compiles the Wails app for macOS with the given ldflags.
func buildMacOSWailsApp(ldFlags string) error {
	fmt.Println("Building Wails app for macOS")
	return sh.RunV("wails", "build", "-m", "-nosyncgomod", "-ldflags", ldFlags, "-platform", "darwin/universal")
}

// importAndSignMacOSCode handles code signing and certificate import for macOS.
func importAndSignMacOSCode() error {
	cert := os.Getenv("AC_CERTIFICATE")
	certPassword := os.Getenv("AC_CERTIFICATE_PASSWORD")

	if cert == "" || certPassword == "" {
		return fmt.Errorf("missing required environment variables: AC_CERTIFICATE or AC_CERTIFICATE_PASSWORD")
	}

	if err := importMacOSCodeSigningCertificates(cert, certPassword); err != nil {
		return fmt.Errorf("error importing code signing certificates: %v", err)
	}

	fmt.Println("Signing macOS Build")
	if err := sh.RunV("gon", "-log-level=info", "./build/darwin/gon-sign.json"); err != nil {
		return fmt.Errorf("error signing macOS build: %v", err)
	}

	return nil
}

// buildAndSignMacOSDMG creates and signs the DMG for macOS, returning its output path.
func buildAndSignMacOSDMG(appName, appVersion string) (string, error) {
	dmgOutputPath := fmt.Sprintf("./build/bin/%s-darwin-universal-%s.dmg", appName, appVersion)

	fmt.Println("Building DMG for macOS")
	if err := sh.RunV("create-dmg", "--window-size", "800", "300", "--no-internet-enable", "--hide-extension", fmt.Sprintf("%s.app", appName), "--app-drop-link", "600", "40", dmgOutputPath, fmt.Sprintf("./build/bin/%s.app", appName)); err != nil {
		return "", fmt.Errorf("error building DMG for macOS: %v", err)
	}

	developerID := os.Getenv("AC_DEVELOPER_ID")
	if developerID == "" {
		return "", fmt.Errorf("missing required environment variable: AC_DEVELOPER_ID")
	}

	fmt.Println("Signing macOS DMG")
	if err := sh.RunV("codesign", "--sign", developerID, dmgOutputPath); err != nil {
		return "", fmt.Errorf("error signing macOS DMG: %v", err)
	}

	if err := compileAndSetMacOSDMGIcons(dmgOutputPath, appName); err != nil {
		return "", err
	}

	return dmgOutputPath, nil
}

// compileAndSetMacOSDMGIcons compiles and sets the DMG icons for macOS.
func compileAndSetMacOSDMGIcons(dmgOutputPath, appName string) error {
	fmt.Println("Compiling seticon.swift for macOS DMG")
	if err := sh.Run("swiftc", "./build/darwin/seticon.swift"); err != nil {
		return fmt.Errorf("error compiling seticon for macOS DMG: %v", err)
	}

	if err := sh.Run("chmod", "+x", "./seticon"); err != nil {
		return fmt.Errorf("error setting permissions on seticon for macOS DMG: %v", err)
	}

	fmt.Println("Setting icons for macOS DMG")
	if err := sh.RunV("./seticon", fmt.Sprintf("./build/bin/%s.app/Contents/Resources/iconfile.icns", appName), dmgOutputPath); err != nil {
		return fmt.Errorf("error setting icons for macOS DMG: %v", err)
	}

	return nil
}

// notarizeMacOSDMG submits the DMG for macOS notarization using credentials passed from environment variables.
func notarizeMacOSDMG(dmgPath string) error {
	appleID := os.Getenv("AC_APPLE_ID")
	password := os.Getenv("AC_PASSWORD")
	teamID := os.Getenv("AC_TEAM_ID")

	// Error if notarization credentials are missing
	if appleID == "" || password == "" || teamID == "" {
		return fmt.Errorf("missing required environment variables: AC_APPLE_ID, AC_PASSWORD, or AC_TEAM_ID")
	}

	// Submit the DMG for notarization and wait for completion
	fmt.Println("Submitting DMG for macOS notarization")
	return sh.RunV("xcrun", "notarytool", "submit", dmgPath,
		"--wait",
		"--apple-id", appleID,
		"--password", password,
		"--team-id", teamID)
}

// stapleMacOSNotarization staples the notarization ticket to the macOS DMG.
func stapleMacOSNotarization(dmgPath string) error {
	fmt.Println("Stapling notarization ticket to macOS DMG")
	return sh.RunV("xcrun", "stapler", "staple", dmgPath)
}

// importMacOSCodeSigningCertificates imports the code signing certificates into the keychain.
func importMacOSCodeSigningCertificates(certBase64, password string) error {
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
