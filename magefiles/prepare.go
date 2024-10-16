package main

const (
	appIdentifier    = "io.rocketblend.rocketblend-desktop"
	appPath          = "./build/bin/rocketblend-desktop.app"
	buildBinDir      = "./build/bin/"
	entitlementsPath = "./build/darwin/entitlements.plist"
)

func Prepare(version, timestamp, commitSha, link, buildtype string, notarize bool) error {
	cleannedVersion, err := getCleannedVersion(version)
	if err != nil {
		return err
	}

	if err := build(cleannedVersion, timestamp, commitSha, link, buildtype, notarize); err != nil {
		return err
	}

	if err := pack(appPath, cleannedVersion, appIdentifier, buildBinDir, entitlementsPath, notarize); err != nil {
		return err
	}

	return nil
}
