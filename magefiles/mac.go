package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/sh"
)

func buildMacOS(name, version, timestamp, commitSha, link, buildType string) error {
	fmt.Printf("Building macOS universal binary for %s\n", name)
	ldFlags := buildFlags(version, timestamp, commitSha, link, buildType)
	return sh.RunV("wails", "build", "-m", "-nosyncgomod", "-ldflags", ldFlags, "-platform", "darwin/universal", "-obfuscated", "-garbleargs", garbleFlags())
}

func packageMacOS(appPath, version, bundleID, outputDir, developerID, appleID, password, teamID, entitlementsPath string, notorize bool) error {
	if err := signMacOSFile(appPath, developerID, bundleID, entitlementsPath); err != nil {
		return err
	}

	name := strings.TrimSuffix(filepath.Base(appPath), ".app")
	dmgOutputPath := filepath.Join(outputDir, fmt.Sprintf("%s-darwin-universal%s.dmg", name, formatVersion(version)))

	if err := createMacOSDMG(appPath, dmgOutputPath); err != nil {
		return err
	}

	if err := signMacOSFile(dmgOutputPath, developerID, bundleID, ""); err != nil {
		return err
	}

	if notorize {
		if err := notarizeMacOSFile(dmgOutputPath, appleID, password, teamID); err != nil {
			return err
		}

		if err := stapleNotarization(dmgOutputPath); err != nil {
			return err
		}
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

func formatVersion(version string) string {
	if version == "" {
		return ""
	}

	return "-" + version
}
