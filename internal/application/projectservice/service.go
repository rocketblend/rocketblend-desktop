package projectservice

import (
	"context"
	"fmt"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectstore"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectstore/listoptions"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/factory"
)

type (
	Service interface {
		Get(ctx context.Context, id uuid.UUID) (*GetProjectResponse, error)
		List(ctx context.Context, opts ...listoptions.ListOption) (*ListProjectsResponse, error)

		Create(ctx context.Context, request *CreateProjectRequest) error
		Update(ctx context.Context, request *UpdateProjectRequest) error

		Render(ctx context.Context, id uuid.UUID) error
		Run(ctx context.Context, id uuid.UUID) error
	}

	service struct {
		logger logger.Logger

		store   projectstore.Store
		factory factory.Factory
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
		return nil, fmt.Errorf("store is required")
	}

	factory, err := factory.New()
	if err != nil {
		return nil, err
	}

	if err := factory.SetLogger(options.Logger); err != nil {
		return nil, err
	}

	return &service{
		logger:  options.Logger,
		store:   options.Store,
		factory: factory,
	}, nil
}

func (s *service) Create(ctx context.Context, request *CreateProjectRequest) error {
	return nil
}

func (s *service) Update(ctx context.Context, request *UpdateProjectRequest) error {
	return nil
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*GetProjectResponse, error) {
	project, err := s.store.GetProject(id)
	if err != nil {
		return nil, err
	}

	return &GetProjectResponse{
		Project: project,
	}, nil
}

func (s *service) List(ctx context.Context, opts ...listoptions.ListOption) (*ListProjectsResponse, error) {
	projects, err := s.store.ListProjects(opts...)
	if err != nil {
		return nil, err
	}

	return &ListProjectsResponse{
		projects,
	}, nil
}
