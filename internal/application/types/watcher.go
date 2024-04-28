package types

type (
	Watcher interface {
		Close() error
		SetPaths(paths ...string) error
	}
)
