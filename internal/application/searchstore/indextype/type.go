package indextype

import "encoding/json"

type IndexType int

const (
	Unknown IndexType = iota
	Project
	Package
	Operation
	Statistic
)

func (p IndexType) String() string {
	return [...]string{"unknown", "project", "package", "operation", "statistic"}[p]
}

func (p IndexType) Int() int {
	return int(p)
}

func (p IndexType) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p *IndexType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*p = PackageTypeFromString(s)
	return nil
}

func PackageTypeFromString(str string) IndexType {
	packageTypeMap := map[string]IndexType{
		"unknown":   Unknown,
		"project":   Project,
		"package":   Package,
		"operation": Operation,
		"statistic": Statistic,
	}

	packageType, ok := packageTypeMap[str]
	if !ok {
		return Unknown
	}

	return packageType
}
