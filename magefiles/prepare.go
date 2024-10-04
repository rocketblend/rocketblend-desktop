package main

func Prepare(version, timestamp, commitSha, link, outputDir, buildtype string, notorize bool) error {
	cleannedVersion, err := getCleannedVersion(version)
	if err != nil {
		return err
	}

	if err := Build(cleannedVersion, timestamp, commitSha, link, outputDir, buildtype); err != nil {
		return err
	}

	if err := Package(outputDir, cleannedVersion, "io.rocketblend.rocketblend-desktop", ".", "./build/darwin/entitlements.plist", notorize); err != nil {
		return err
	}

	return nil
}
