package projectservice

import (
	"fmt"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend-desktop/internal/application/project"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectstore"
)

type (
	FindOptions struct {
	}

	FindOption func(*FindOptions)

	Service interface {
		FindAll(opts ...FindOption) ([]*project.Project, error)
		FindByKey(key string) (*project.Project, error)
	}

	service struct {
		logger logger.Logger
		store  projectstore.Store
	}

	Options struct {
		Logger logger.Logger
		Store  projectstore.Store
	}

	Option func(*Options)
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithStore(store projectstore.Store) Option {
	return func(o *Options) {
		o.Store = store
	}
}

func New(opts ...Option) (Service, error) {
	options := &Options{
		Logger: logger.NoOp(),
	}

	for _, o := range opts {
		o(options)
	}

	if options.Store == nil {
		return nil, fmt.Errorf("watcher is required")
	}

	return &service{
		logger: options.Logger,
		store:  options.Store,
	}, nil
}

func (s *service) FindAll(opts ...FindOption) ([]*project.Project, error) {
	return s.store.ListProjects(), nil
}

func (s *service) FindByKey(key string) (*project.Project, error) {
	project, ok := s.store.GetProject(key)
	if !ok {
		return nil, fmt.Errorf("project not found")
	}

	return project, nil
}
