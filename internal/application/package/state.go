package pack

type PackageState int

const (
	Available PackageState = iota
	Downloading
	Stopped
	Installed
	Error
)

var AllPackageStates = []struct {
	Value  PackageState
	TSName string
}{
	{Available, "AVAILABLE"},
	{Downloading, "DOWNLOADING"},
	{Stopped, "STOPPED"},
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
