package project

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/types"
	"github.com/rocketblend/rocketblend-desktop/internal/application/watcher"
	"github.com/rocketblend/rocketblend-desktop/internal/helpers"
	rbhelpers "github.com/rocketblend/rocketblend/pkg/helpers"
	"github.com/rocketblend/rocketblend/pkg/reference"
	rbtypes "github.com/rocketblend/rocketblend/pkg/types"
)

type (
	Repository struct {
		logger    types.Logger
		validator types.Validator

		rbRepository   types.RBRepository
		rbDriver       types.RBDriver
		rbConfigurator types.RBConfigurator

		blender      types.Blender
		configurator types.Configurator

		store      types.Store
		watcher    types.Watcher
		dispatcher types.Dispatcher
	}

	Options struct {
		Logger    types.Logger
		Validator types.Validator

		RBRepository   types.RBRepository
		RBDriver       types.RBDriver
		RBConfigurator types.RBConfigurator

		Blender      types.Blender
		Configurator types.Configurator

		Store      types.Store
		Dispatcher types.Dispatcher

		WatcherDebounceDuration time.Duration
	}

	Option func(*Options)
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithValidator(validator types.Validator) Option {
	return func(o *Options) {
		o.Validator = validator
	}
}

func WithStore(store types.Store) Option {
	return func(o *Options) {
		o.Store = store
	}
}

func WithDispatcher(dispatcher types.Dispatcher) Option {
	return func(o *Options) {
		o.Dispatcher = dispatcher
	}
}

func WithConfigurator(configurator types.Configurator) Option {
	return func(o *Options) {
		o.Configurator = configurator
	}
}

func WithRocketBlendConfigurator(configurator types.RBConfigurator) Option {
	return func(o *Options) {
		o.RBConfigurator = configurator
	}
}

func WithRocketBlendRepository(repository types.RBRepository) Option {
	return func(o *Options) {
		o.RBRepository = repository
	}
}

func WithRocketBlendDriver(driver types.RBDriver) Option {
	return func(o *Options) {
		o.RBDriver = driver
	}
}

func WithBlender(blender types.Blender) Option {
	return func(o *Options) {
		o.Blender = blender
	}
}

func WithWatcherDebounceDuration(duration time.Duration) Option {
	return func(o *Options) {
		o.WatcherDebounceDuration = duration
	}
}

func New(opts ...Option) (*Repository, error) {
	options := &Options{
		Logger:                  logger.NoOp(),
		WatcherDebounceDuration: 2 * time.Second,
	}

	for _, o := range opts {
		o(options)
	}

	if options.Validator == nil {
		return nil, errors.New("validator is required")
	}

	if options.Configurator == nil {
		return nil, errors.New("configurator is required")
	}

	if options.Store == nil {
		return nil, errors.New("store is required")
	}

	if options.Dispatcher == nil {
		return nil, errors.New("dispatcher is required")
	}

	if options.RBConfigurator == nil {
		return nil, errors.New("rocketblend configurator is required")
	}

	if options.RBRepository == nil {
		return nil, errors.New("rocketblend repository is required")
	}

	if options.RBDriver == nil {
		return nil, errors.New("rocketblend driver is required")
	}

	if options.Blender == nil {
		return nil, errors.New("blender is required")
	}

	config, err := options.Configurator.Get()
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
			project, err := load(options.Validator, options.RBConfigurator, path)
			if err != nil {
				return err
			}

			index, err := convertToIndex(project)
			if err != nil {
				return err
			}

			options.Logger.Debug("updating project index", map[string]interface{}{
				"id":        index.ID,
				"reference": index.Reference,
			})

			return options.Store.Insert(context.Background(), index)
		}),
		watcher.WithRemoveObjectFunc(func(removePath string) error {
			return options.Store.RemoveByReference(context.Background(), path.Clean(removePath))
		}),
	)
	if err != nil {
		return nil, err
	}

	return &Repository{
		logger:         options.Logger,
		configurator:   options.Configurator,
		validator:      options.Validator,
		rbConfigurator: options.RBConfigurator,
		rbRepository:   options.RBRepository,
		rbDriver:       options.RBDriver,
		blender:        options.Blender,
		store:          options.Store,
		dispatcher:     options.Dispatcher,
		watcher:        watcher,
	}, nil
}

func (r *Repository) Close() error {
	return r.watcher.Close()
}

func (r *Repository) saveDetail(path string, detail *types.Detail) error {
	if err := rbhelpers.Save(r.validator, filepath.Join(path, types.DetailFileName), detail); err != nil {
		return err
	}

	return nil
}

func load(validator rbtypes.Validator, configurator rbtypes.Configurator, path string) (*types.Project, error) {
	if ignoreProject(path) {
		return nil, errors.New("project is ignored")
	}

	blendFilePaths, err := findFilePathForExtension(path, rbtypes.BlendFileExtension)
	if err != nil {
		return nil, err
	}

	config, err := configurator.Get()
	if err != nil {
		return nil, err
	}

	profile, err := loadOrCreateProfile(validator, path, config.DefaultBuild)
	if err != nil {
		return nil, err
	}

	blendFilePath := blendFilePaths[0]
	detail, err := loadOrCreateDetail(validator, path, blendFilePath)
	if err != nil {
		return nil, err
	}

	modTime, err := helpers.GetModTime(path)
	if err != nil {
		return nil, err
	}

	builds := profile.FindAll(rbtypes.PackageBuild)
	if len(builds) == 0 {
		return nil, errors.New("no build found in profile")
	}

	addons := make([]reference.Reference, 0, len(profile.Dependencies)-1)
	for _, dep := range profile.Dependencies {
		if dep.Type == rbtypes.PackageAddon {
			addons = append(addons, dep.Reference)
		}
	}

	// TODO: Switch to using a list/map of image paths.

	return &types.Project{
		ID:            detail.ID,
		Name:          detail.Name,
		Tags:          detail.Tags,
		SplashPath:    imagePath(path, detail.SplashPath),
		ThumbnailPath: imagePath(path, detail.ThumbnailPath),
		Path:          path,
		FileName:      filepath.Base(blendFilePath),
		Build:         builds[0].Reference,
		Addons:        addons,
		UpdatedAt:     modTime,
	}, nil
}

func loadOrCreateProfile(validator rbtypes.Validator, path string, defaultBuild reference.Reference) (*rbtypes.Profile, error) {
	profileFilePath := filepath.Join(path, rbtypes.ProfileFileName)
	_, err := os.Stat(profileFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			profile := &rbtypes.Profile{
				Dependencies: []*rbtypes.Dependency{{Reference: defaultBuild, Type: rbtypes.PackageBuild}},
			}

			if err := rbhelpers.Save(validator, profileFilePath, profile); err != nil {
				return nil, err
			}

			return profile, nil
		}

		return nil, err
	}

	profile, err := rbhelpers.Load[rbtypes.Profile](validator, profileFilePath)
	if err == nil {
		return profile, nil
	}

	return nil, err
}

func loadOrCreateDetail(validator rbtypes.Validator, path string, blendFilePath string) (*types.Detail, error) {
	detailFilePath := filepath.Join(path, types.DetailFileName)
	_, err := os.Stat(detailFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			detail := &types.Detail{
				ID:   uuid.New(),
				Name: helpers.FilenameToDisplayName(blendFilePath),
			}

			if err := rbhelpers.Save(validator, detailFilePath, detail); err != nil {
				return nil, err
			}

			return detail, nil
		}

		return nil, err
	}

	detail, err := rbhelpers.Load[types.Detail](validator, detailFilePath)
	if err == nil {
		return detail, nil
	}

	return nil, err
}

func imagePath(rootPath string, imagePath string) string {
	if imagePath == "" {
		return ""
	}

	return filepath.ToSlash(filepath.Join(rootPath, imagePath))
}

func ignoreProject(projectPath string) bool {
	_, err := os.Stat(filepath.Join(projectPath, types.IgnoreFileName))
	return !os.IsNotExist(err)
}

func findFilePathForExtension(dir string, ext string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*"+ext))
	if err != nil {
		return nil, fmt.Errorf("failed to list files in current directory: %w", err)
	}

	if len(files) == 0 {
		return nil, errors.New("no files found in current directory")
	}

	return files, nil
}
