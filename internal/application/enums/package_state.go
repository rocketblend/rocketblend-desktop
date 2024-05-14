package enums

type PackageState string

const (
	PackageStateAvailable   PackageState = "available"
	PackageStateDownloading PackageState = "downloading"
	PackageStateIncomplete  PackageState = "incomplete"
	PackageStateInstalled   PackageState = "installed"
	PackageStateError       PackageState = "error"
)

var PackageStates = []struct {
	Value  PackageState
	TSName string
}{
	{PackageStateAvailable, "AVAILABLE"},
	{PackageStateDownloading, "DOWNLOADING"},
	{PackageStateIncomplete, "INCOMPLETE"},
	{PackageStateInstalled, "INSTALLED"},
	{PackageStateError, "ERROR"},
}
