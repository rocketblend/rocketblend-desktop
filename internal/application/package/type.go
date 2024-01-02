package pack

type PackageType int

const (
	Unknown PackageType = iota
	Addon
	Build
)

func (p PackageType) String() string {
	for _, t := range AllPackageTypes {
		if t.Value == p {
			return t.TSName
		}
	}
	return "UNKNOWN"
}

var AllPackageTypes = []struct {
	Value  PackageType
	TSName string
}{
	{Unknown, "UNKNOWN"},
	{Addon, "ADDON"},
	{Build, "BUILD"},
}
