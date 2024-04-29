package enums

type PackageState string

const (
	PackageStateAvailable   PackageState = "Available"
	PackageStateDownloading PackageState = "Downloading"
	PackageStateCancelled   PackageState = "Cancelled"
	PackageStateInstalled   PackageState = "Installed"
	PackageStateError       PackageState = "Error"
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
