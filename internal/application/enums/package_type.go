package enums

type PackageType string

const (
	PackageTypeBuild PackageType = "build"
	PackageTypeAddon PackageType = "addon"
)

var PackageTypes = []struct {
	Value  PackageType
	TSName string
}{
	{PackageTypeBuild, "BUILD"},
	{PackageTypeAddon, "ADDON"},
}
