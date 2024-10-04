package main

func Prepare(version, timestamp, commitSha, link, buildtype string, notorize bool) error {
	cleannedVersion, err := getCleannedVersion(version)
	if err != nil {
		return err
	}

	if err := Build(cleannedVersion, timestamp, commitSha, link, buildtype); err != nil {
		return err
	}

	if err := Package("./build/bin/rocketblend-desktop.app", cleannedVersion, "io.rocketblend.rocketblend-desktop", "./build/bin/", "./build/darwin/entitlements.plist", notorize); err != nil {
		return err
	}

	return nil
}
