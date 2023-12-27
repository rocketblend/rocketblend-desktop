package hook

import (
	"github.com/rs/zerolog"
)

type (
	Hook interface {
		Run(e *zerolog.Event, level zerolog.Level, msg string)
	}

	hook struct {
		onLogFunc func(level string, msg string, fields map[string]interface{})
	}

	LogFunc func(level string, msg string, fields map[string]interface{})

	Options struct {
		OnLogFunc LogFunc
	}

	Option func(*Options)
)

// WithOnLogFunc is an option to set the onLogFunc in the service
func WithOnLogFunc(f LogFunc) Option {
	return func(o *Options) {
		o.OnLogFunc = f
	}
}

func New(opts ...Option) (Hook, error) {
	options := &Options{}

	for _, o := range opts {
		o(options)
	}

	return &hook{onLogFunc: options.OnLogFunc}, nil
}

func (srv *hook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	// Extract fields from zerolog.Event
	fields := make(map[string]interface{})
	e.Fields(fields)

	// Call the onLogFunc if it's configured
	if srv.onLogFunc != nil {
		srv.onLogFunc(level.String(), msg, fields)
	}
}
