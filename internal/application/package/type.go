package pack

import "encoding/json"

type Type int

const (
	Unknown Type = iota
	Addon
	Build
)

func (p Type) String() string {
	return [...]string{"unknown", "addon", "build"}[p]
}

func (p Type) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p *Type) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*p = PackageTypeFromString(s)
	return nil
}

func PackageTypeFromString(str string) Type {
	packageTypeMap := map[string]Type{
		"unknown": Unknown,
		"addon":   Addon,
		"build":   Build,
	}

	packageType, ok := packageTypeMap[str]
	if !ok {
		return Unknown
	}

	return packageType
}
