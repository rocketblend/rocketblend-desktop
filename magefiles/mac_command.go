package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/sh"
)

func buildReleaseMacOS(name, version, timestamp, commitSha, link, outputDir string, debug bool) error {
	ldFlags := buildFlags(version, timestamp, commitSha, link, outputDir, debug)
	return buildMacOSWailsApp(ldFlags, name, outputDir)
}

func packageAppMacOS(appPath, version, outputDir, developerID, appleID, password, teamID, entitlementsPath string) error {
	name := strings.TrimSuffix(filepath.Base(appPath), ".app")
	dmgOutputPath := filepath.Join(outputDir, fmt.Sprintf("%s-darwin-universal-%s.dmg", name, version))

	if err := signMacOSFile(appPath, developerID, BundleID, entitlementsPath); err != nil {
		return err
	}

	if err := createMacOSDMG(appPath, dmgOutputPath); err != nil {
		return err
	}

	if err := signMacOSFile(dmgOutputPath, developerID, BundleID, ""); err != nil {
		return err
	}

	if err := notarizeMacOSFile(dmgOutputPath, appleID, password, teamID); err != nil {
		return err
	}

	if err := stapleNotarization(dmgOutputPath); err != nil {
		return err
	}

	return nil
}

func signMacOSFile(filePath, developerID, bundleID, entitlementsPath string) error {
	fmt.Printf("Signing file: %s with Developer ID: %s\n", filePath, developerID)

	args := []string{"--verbose", "--force", "--options", "runtime", "--sign", developerID, "--timestamp", "--identifier", bundleID}
	if entitlementsPath != "" {
		args = append(args, "--entitlements", entitlementsPath)
	}

	args = append(args, filePath)
	return sh.RunV("codesign", args...)
}

func buildMacOSWailsApp(ldFlags, name, version, outputDir string) error {
	outputFileName := fmt.Sprintf("%s-%s.app", name, version)
	fmt.Printf("Building Wails app for macOS: %s to %s\n", outputFileName, outputDir)
	return sh.RunV("wails", "build", "-m", "-nosyncgomod", "-ldflags", ldFlags, "-platform", "darwin/universal", "-o", filepath.Join(outputDir, outputFileName))
}

func createMacOSDMG(appPath, dmgOutputPath string) error {
	fmt.Printf("Building DMG for app at %s\n", dmgOutputPath)
	if err := sh.RunV("create-dmg", "--window-size", "800", "300", "--no-internet-enable", "--hide-extension", filepath.Base(appPath), "--app-drop-link", "600", "40", dmgOutputPath, appPath); err != nil {
		return fmt.Errorf("error building DMG for macOS: %v", err)
	}

	if err := compileSetIconSwiftScript(); err != nil {
		return fmt.Errorf("error compiling seticon Swift script: %v", err)
	}

	iconFilePath := fmt.Sprintf("%s/Contents/Resources/iconfile.icns", appPath)
	if err := setDMGIcons(iconFilePath, dmgOutputPath); err != nil {
		return fmt.Errorf("error setting DMG icons: %v", err)
	}

	return nil
}

func notarizeMacOSFile(filePath, appleID, password, teamID string) error {
	fmt.Printf("Submitting file for notarization: %s\n", filePath)
	return sh.RunV("xcrun", "notarytool", "submit", filePath, "--wait", "--apple-id", appleID, "--password", password, "--team-id", teamID)
}

func stapleNotarization(filePath string) error {
	fmt.Printf("Stapling notarization ticket to file: %s\n", filePath)
	return sh.RunV("xcrun", "stapler", "staple", filePath)
}

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

func setDMGIcons(iconFilePath, dmgOutputPath string) error {
	fmt.Println("Setting icons for macOS DMG")
	if err := sh.RunV("./seticon", iconFilePath, dmgOutputPath); err != nil {
		return fmt.Errorf("error setting icons for macOS DMG: %v", err)
	}

	return nil
}

func getMacOSVariables() (*EnvironmentVariables, error) {
	env := EnvironmentVariables{
		DeveloperID: os.Getenv("AC_DEVELOPER_ID"),
		AppleID:     os.Getenv("AC_APPLE_ID"),
		Password:    os.Getenv("AC_PASSWORD"),
		TeamID:      os.Getenv("AC_TEAM_ID"),
	}

	if env.DeveloperID == "" || env.AppleID == "" || env.Password == "" || env.TeamID == "" {
		return nil, fmt.Errorf("missing required environment variables")
	}

	return &env, nil
}
