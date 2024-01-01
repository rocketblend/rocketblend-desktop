package pack

type PackageState int

const (
	Available PackageState = iota
	Downloading
	Cancelled
	Installed
	Error
)

var AllPackageStates = []struct {
	Value  PackageState
	TSName string
}{
	{Available, "AVAILABLE"},
	{Downloading, "DOWNLOADING"},
	{Cancelled, "CANCELLED"},
	{Installed, "INSTALLED"},
	{Error, "ERROR"},
}

func (ps PackageState) String() string {
	for _, state := range AllPackageStates {
		if state.Value == ps {
			return state.TSName
		}
	}
	return "ERROR"
}
