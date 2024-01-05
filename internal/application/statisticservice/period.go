package statisticservice

type Period int

const (
	PeriodUnknown Period = iota
	PeriodLifetime
	PeriodDay
	PeriodWeek
)

var AllPeriodTypes = []struct {
	Value  Period
	TSName string
}{
	{PeriodUnknown, "UNKNOWN"},
	{PeriodLifetime, "LIFETIME"},
	{PeriodDay, "DAY"},
	{PeriodWeek, "WEEK"},
}

func (p Period) String() string {
	for _, t := range AllPeriodTypes {
		if t.Value == p {
			return t.TSName
		}
	}

	return "UNKNOWN"
}
