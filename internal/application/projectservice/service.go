package projectservice

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"time"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/config"
	"github.com/rocketblend/rocketblend-desktop/internal/application/eventservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/project"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectsettings"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/indextype"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/listoption"
	"github.com/rocketblend/rocketblend-desktop/internal/application/watcher"

	rocketblendBlendFile "github.com/rocketblend/rocketblend/pkg/driver/blendfile"
	rocketblendInstallation "github.com/rocketblend/rocketblend/pkg/driver/installation"
	rocketblendPackage "github.com/rocketblend/rocketblend/pkg/driver/rocketpack"
)

type (
	Service interface {
		Get(ctx context.Context, id uuid.UUID) (*GetProjectResponse, error)
		List(ctx context.Context, opts ...listoption.ListOption) (*ListProjectsResponse, error)

		Create(ctx context.Context, opts *CreateProjectOpts) (*CreateProjectResult, error)
		Update(ctx context.Context, opts *UpdateProjectOpts) error
		AddPackage(ctx context.Context, opts *AddProjectPackageOpts) error
		RemovePackage(ctx context.Context, opts *RemoveProjectPackageOpts) error

		Render(ctx context.Context, id uuid.UUID) error
		Run(ctx context.Context, id uuid.UUID) error

		Explore(ctx context.Context, id uuid.UUID) error

		Refresh(ctx context.Context) error

		Close() error
	}

	service struct {
		logger logger.Logger

		applicationConfigService config.Service

		rocketblendPackageService      rocketblendPackage.Service
		rocketblendInstallationService rocketblendInstallation.Service
		rocketblendBlendFileService    rocketblendBlendFile.Service

		store      searchstore.Store
		watcher    watcher.Watcher
		dispatcher eventservice.Service
	}

	Options struct {
		Logger logger.Logger

		ApplicationConfigService config.Service

		RocketblendPackageService      rocketblendPackage.Service
		RocketblendInstallationService rocketblendInstallation.Service
		RocketblendBlendFileService    rocketblendBlendFile.Service

		Store      searchstore.Store
		Dispatcher eventservice.Service

		WatcherDebounceDuration time.Duration
	}

	Option func(*Options)
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithWatcherDebounceDuration(duration time.Duration) Option {
	return func(o *Options) {
		o.WatcherDebounceDuration = duration
	}
}

func WithApplicationConfigService(srv config.Service) Option {
	return func(o *Options) {
		o.ApplicationConfigService = srv
	}
}

func WithRocketBlendPackageService(srv rocketblendPackage.Service) Option {
	return func(o *Options) {
		o.RocketblendPackageService = srv
	}
}

func WithRocketBlendInstallationService(srv rocketblendInstallation.Service) Option {
	return func(o *Options) {
		o.RocketblendInstallationService = srv
	}
}

func WithRocketBlendBlendFileService(srv rocketblendBlendFile.Service) Option {
	return func(o *Options) {
		o.RocketblendBlendFileService = srv
	}
}

func WithStore(store searchstore.Store) Option {
	return func(o *Options) {
		o.Store = store
	}
}

func WithDispatcher(dispatcher eventservice.Service) Option {
	return func(o *Options) {
		o.Dispatcher = dispatcher
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

	if options.ApplicationConfigService == nil {
		return nil, fmt.Errorf("application config service is required")
	}

	if options.RocketblendPackageService == nil {
		return nil, fmt.Errorf("rocketblend package service is required")
	}

	if options.RocketblendInstallationService == nil {
		return nil, fmt.Errorf("rocketblend installation service is required")
	}

	if options.RocketblendBlendFileService == nil {
		return nil, fmt.Errorf("rocketblend blend file service is required")
	}

	if options.Store == nil {
		return nil, fmt.Errorf("store is required")
	}

	if options.Dispatcher == nil {
		return nil, fmt.Errorf("dispatcher is required")
	}

	config, err := options.ApplicationConfigService.Get()
	if err != nil {
		return nil, err
	}

	watcher, err := watcher.New(
		watcher.WithLogger(options.Logger),
		watcher.WithEventDebounceDuration(options.WatcherDebounceDuration),
		watcher.WithPaths(config.Project.Paths...),
		watcher.WithIsWatchableFileFunc(func(path string) bool {
			for _, ext := range config.Project.Watcher.FileExtensions {
				if filepath.Ext(path) == ext {
					return true
				}
			}

			return false
		}),
		watcher.WithResolveObjectPathFunc(func(path string) string {
			return filepath.Dir(path)
		}),
		watcher.WithUpdateObjectFunc(func(path string) error {
			project, err := project.Load(path)
			if err != nil {
				return fmt.Errorf("failed to load project %s: %w", path, err)
			}

			index, err := convertToIndex(project)
			if err != nil {
				return err
			}

			options.Logger.Debug("updating project index", map[string]interface{}{
				"id":        index.ID,
				"reference": index.Reference,
			})

			// TODO: Pass context from watcher
			return options.Store.Insert(context.Background(), index)
		}),
		watcher.WithRemoveObjectFunc(func(removePath string) error {
			return options.Store.RemoveByReference(context.Background(), path.Clean(removePath))
		}),
	)
	if err != nil {
		return nil, err
	}

	return &service{
		logger:                         options.Logger,
		applicationConfigService:       options.ApplicationConfigService,
		rocketblendPackageService:      options.RocketblendPackageService,
		rocketblendInstallationService: options.RocketblendInstallationService,
		rocketblendBlendFileService:    options.RocketblendBlendFileService,
		store:                          options.Store,
		dispatcher:                     options.Dispatcher,
		watcher:                        watcher,
	}, nil
}

func (s *service) Close() error {
	return s.watcher.Close()
}

func (s *service) Update(ctx context.Context, opts *UpdateProjectOpts) error {
	project, err := s.get(ctx, opts.ID)
	if err != nil {
		return err
	}

	settings := project.GetSettings()
	if opts.Name != nil {
		settings.Name = *opts.Name
	}

	if opts.Tags != nil {
		settings.Tags = *opts.Tags
	}

	if opts.ThumbnailPath != nil {
		settings.ThumbnailPath = *opts.ThumbnailPath
	}

	if opts.SplashPath != nil {
		settings.SplashPath = *opts.SplashPath
	}

	filePath := filepath.Join(project.Path, projectsettings.FileName)
	if err := projectsettings.Save(settings, filePath); err != nil {
		return err
	}

	s.EmitEvent(ctx, project.ID, UpdateEventChannel)

	return nil
}

func (s *service) AddPackage(ctx context.Context, opts *AddProjectPackageOpts) error {
	return nil
}

func (s *service) RemovePackage(ctx context.Context, opts *RemoveProjectPackageOpts) error {
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
	indexes, err := s.store.List(ctx, opts...)
	if err != nil {
		return nil, err
	}

	projects := make([]*project.Project, 0, len(indexes))
	for _, index := range indexes {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		pack, err := convertFromIndex(index)
		if err != nil {
			return nil, err
		}

		projects = append(projects, pack)
	}

	s.logger.Debug("Found projects", map[string]interface{}{
		"projects": len(projects),
		"indexes":  len(indexes),
	})

	return &ListProjectsResponse{
		Projects: projects,
	}, nil
}

func (s *service) Refresh(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	config, err := s.applicationConfigService.Get()
	if err != nil {
		return err
	}

	if err := s.watcher.SetPaths(config.Project.Paths...); err != nil {
		return err
	}

	return nil
}

func (s *service) get(ctx context.Context, id uuid.UUID) (*project.Project, error) {
	index, err := s.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	project, err := convertFromIndex(index)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func convertFromIndex(index *searchstore.Index) (*project.Project, error) {
	var result project.Project
	if err := json.Unmarshal([]byte(index.Data), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func convertToIndex(project *project.Project) (*searchstore.Index, error) {
	data, err := json.Marshal(project)
	if err != nil {
		return nil, err
	}

	resources := []string{}
	if project.ThumbnailPath != "" {
		resources = append(resources, filepath.ToSlash(project.ThumbnailPath))
		resources = append(resources, filepath.ToSlash(project.SplashPath))
	}

	return &searchstore.Index{
		ID:        project.ID,
		Name:      project.Name,
		Type:      indextype.Project,
		Reference: path.Clean(project.Path),
		Resources: resources,
		Data:      string(data),
	}, nil
}
