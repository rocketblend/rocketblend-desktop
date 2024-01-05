package statisticservice

type DataType int

const (
	UnknownType DataType = iota
	IntType
	FloatType
	StringType
)

var AllDataTypes = []struct {
	Value  DataType
	TSName string
}{
	{UnknownType, "UNKNOWN"},
	{IntType, "INT"},
	{FloatType, "FLOAT"},
	{StringType, "STRING"},
}

func (p DataType) String() string {
	for _, t := range AllDataTypes {
		if t.Value == p {
			return t.TSName
		}
	}

	return "UNKNOWN"
}
