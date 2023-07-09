package driver

type (
	Driver interface {
	}

	driver struct {
	}
)

// New creates a new Driver struct
func New() Driver {
	return &driver{}
}
