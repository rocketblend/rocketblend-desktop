package enums

type PackageType string

const (
	PackageTypeBuild PackageType = "Build"
	PackageTypeAddon PackageType = "Addon"
)

var PackageTypes = []struct {
	Value  PackageType
	TSName string
}{
	{PackageTypeBuild, "BUILD"},
	{PackageTypeAddon, "ADDON"},
}
