package projectservice

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/project"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/indextype"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/listoption"
	"github.com/rocketblend/rocketblend-desktop/internal/application/watcher"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/factory"
)

type (
	Service interface {
		Get(ctx context.Context, id uuid.UUID) (*GetProjectResponse, error)
		List(ctx context.Context, opts ...listoption.ListOption) (*ListProjectsResponse, error)

		Create(ctx context.Context, request *CreateProjectRequest) error
		Update(ctx context.Context, request *UpdateProjectRequest) error

		Render(ctx context.Context, id uuid.UUID) error
		Run(ctx context.Context, id uuid.UUID) error

		Explore(ctx context.Context, id uuid.UUID) error

		Close() error
	}

	service struct {
		logger logger.Logger

		factory factory.Factory

		store   searchstore.Store
		watcher watcher.Watcher
	}

	Options struct {
		Logger                  logger.Logger
		Factory                 factory.Factory
		Store                   searchstore.Store
		WatcherDebounceDuration time.Duration
	}

	Option func(*Options)
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithFactory(factory factory.Factory) Option {
	return func(o *Options) {
		o.Factory = factory
	}
}

func WithStore(store searchstore.Store) Option {
	return func(o *Options) {
		o.Store = store
	}
}

func WithWatcherDebounceDuration(duration time.Duration) Option {
	return func(o *Options) {
		o.WatcherDebounceDuration = duration
	}
}

func New(opts ...Option) (Service, error) {
	options := &Options{
		Logger:                  logger.NoOp(),
		WatcherDebounceDuration: 2 * time.Second,
	}

	for _, o := range opts {
		o(options)
	}

	if options.Store == nil {
		return nil, fmt.Errorf("store is required")
	}

	if options.Factory == nil {
		return nil, fmt.Errorf("factory is required")
	}

	// TODO: Move to config
	watchPaths := "E:\\Blender\\Projects\\Testing\\RocketBlend"
	watchFiles := []string{".blend", ".yaml"}

	watcher, err := watcher.New(
		watcher.WithLogger(options.Logger),
		watcher.WithEventDebounceDuration(options.WatcherDebounceDuration),
		watcher.WithPaths(watchPaths),
		watcher.WithIsWatchableFileFunc(func(path string) bool {
			for _, ext := range watchFiles {
				if filepath.Ext(path) == ext {
					return true
				}
			}

			return false
		}),
		watcher.WithResolveObjectPathFunc(func(path string) string {
			projectPath := filepath.Dir(path)
			if strings.HasSuffix(projectPath, project.ConfigDir) {
				projectPath = filepath.Dir(projectPath)
			}

			return projectPath
		}),
		watcher.WithUpdateObjectFunc(func(path string) error {
			project, err := project.Load(path)
			if err != nil {
				return fmt.Errorf("failed to load project %s: %w", path, err)
			}

			data, err := json.Marshal(project)
			if err != nil {
				return err
			}

			return options.Store.Insert(&searchstore.Index{
				ID:   project.ID,
				Name: project.Name,
				Type: indextype.Project,
				Path: path,
				Data: string(data),
			})
		}),
		watcher.WithRemoveObjectFunc(func(path string) error {
			return options.Store.RemoveByPath(path)
		}),
	)
	if err != nil {
		return nil, err
	}

	return &service{
		logger:  options.Logger,
		factory: options.Factory,
		store:   options.Store,
		watcher: watcher,
	}, nil
}

func (s *service) Close() error {
	return s.watcher.Close()
}

func (s *service) Create(ctx context.Context, request *CreateProjectRequest) error {
	return nil
}

func (s *service) Update(ctx context.Context, request *UpdateProjectRequest) error {
	return nil
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*GetProjectResponse, error) {
	project, err := s.get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &GetProjectResponse{
		Project: project,
	}, nil
}

func (s *service) List(ctx context.Context, opts ...listoption.ListOption) (*ListProjectsResponse, error) {
	opts = append(opts, listoption.WithType(indextype.Project))
	indexes, err := s.store.List(opts...)
	if err != nil {
		return nil, err
	}

	response := &ListProjectsResponse{
		Projects: make([]*project.Project, 0),
	}

	for _, index := range indexes {
		pack, err := s.convert(index)
		if err != nil {
			return nil, err
		}

		response.Projects = append(response.Projects, pack)
	}

	return response, nil
}

func (s *service) get(ctx context.Context, id uuid.UUID) (*project.Project, error) {
	index, err := s.store.Get(id)
	if err != nil {
		return nil, err
	}

	project, err := s.convert(index)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (s *service) convert(index *searchstore.Index) (*project.Project, error) {
	var result project.Project
	if err := json.Unmarshal([]byte(index.Data), &result); err == nil {
		return nil, err
	}

	return &result, nil
}
