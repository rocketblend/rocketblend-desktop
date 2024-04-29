package enums

type PackageState string

const (
	PackageStateAvailable   PackageState = "available"
	PackageStateDownloading PackageState = "downloading"
	PackageStateCancelled   PackageState = "cancelled"
	PackageStateInstalled   PackageState = "installed"
	PackageStateError       PackageState = "error"
)

var PackageStates = []struct {
	Value  PackageState
	TSName string
}{
	{PackageStateAvailable, "AVAILABLE"},
	{PackageStateDownloading, "DOWNLOADING"},
	{PackageStateCancelled, "CANCELLED"},
	{PackageStateInstalled, "INSTALLED"},
	{PackageStateError, "ERROR"},
}
