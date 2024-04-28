package types

import "encoding/json"

type IndexType int

const (
	IndexTypeUnknown IndexType = iota
	IndexTypeProject
	IndexTypePackage
	IndexTypeOperation
	IndexTypeMetric
)

func (p IndexType) String() string {
	return [...]string{"unknown", "project", "package", "operation", "metric"}[p]
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
		"unknown":   IndexTypeUnknown,
		"project":   IndexTypeProject,
		"package":   IndexTypePackage,
		"operation": IndexTypeOperation,
		"metric":    IndexTypeMetric,
	}

	packageType, ok := packageTypeMap[str]
	if !ok {
		return IndexTypeUnknown
	}

	return packageType
}
