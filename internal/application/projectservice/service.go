package projectservice

import (
	"fmt"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectstore"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectstore/listoptions"
)

type (
	Service interface {
		FindAll(opts ...listoptions.ListOption) ([]*Project, error)
		FindByKey(key string) (*Project, error)
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

func (s *service) FindAll(opts ...listoptions.ListOption) ([]*Project, error) {
	projects, err := s.store.ListProjects(opts...)
	if err != nil {
		return nil, err
	}

	return mapProjects(projects...), nil
}

func (s *service) FindByKey(key string) (*Project, error) {
	project, err := s.store.GetProject(key)
	if err != nil {
		return nil, err
	}

	return mapProjects(project)[0], nil
}
