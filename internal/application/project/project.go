package project

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/enums"
	"github.com/rocketblend/rocketblend-desktop/internal/application/fileserver"
	"github.com/rocketblend/rocketblend-desktop/internal/application/store"
	"github.com/rocketblend/rocketblend-desktop/internal/application/types"
	"github.com/rocketblend/rocketblend-desktop/internal/application/watcher"
	"github.com/rocketblend/rocketblend-desktop/internal/helpers"
	rbhelpers "github.com/rocketblend/rocketblend/pkg/helpers"
	"github.com/rocketblend/rocketblend/pkg/reference"
	rbtypes "github.com/rocketblend/rocketblend/pkg/types"
)

const DefaultMediaPath = "/media"

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

	// TODO: This whole watcher thing is a bit of a mess.
	watcher, err := watcher.New(
		watcher.WithLogger(options.Logger),
		watcher.WithEventDebounceDuration(options.WatcherDebounceDuration),
		watcher.WithPaths(config.Project.Paths...),
		watcher.WithIsWatchableFileFunc(func(path string) bool {
			for _, ext := range ValidExtensions() {
				if filepath.Ext(path) == ext {
					return true
				}
			}

			return false
		}),
		watcher.WithResolveObjectPathFunc(func(filePath string) string {
			// TODO: This is a bit of a mess. We should probably just use the store to resolve the project path.
			rootPath := resolveRootPath(filePath, config.Project.Paths)
			if rootPath == "" {
				return ""
			}

			return findProjectRoot(filePath, rootPath)
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
			if err := options.Store.RemoveByReference(context.Background(), path.Clean(removePath)); err != nil && !errors.Is(err, store.ErrNotFound) {
				return err
			}

			return nil
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

func ValidExtensions() []string {
	extensions := make([]string, 0, len(fileserver.ValidExtensions)+2)
	for ext := range fileserver.ValidExtensions {
		extensions = append(extensions, ext)
	}

	extensions = append(extensions, ".blend", ".json", types.IgnoreFileName)
	return extensions
}

func (r *Repository) saveDetail(path string, detail *types.Detail, ensurePath bool, override bool) error {
	if err := rbhelpers.Save(r.validator, detailFilePath(path), detail, ensurePath, override); err != nil {
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

	if filepath.IsLocal(detail.MediaPath) {
		return nil, errors.New("media path must be relative")
	}

	modTime, err := helpers.GetModTime(path)
	if err != nil {
		return nil, err
	}

	media, err := findMediaFiles(filepath.Join(path, detail.MediaPath))
	if err != nil {
		return nil, err
	}

	return &types.Project{
		ID:           detail.ID,
		Name:         detail.Name,
		Tags:         detail.Tags,
		Path:         path,
		MediaPath:    detail.MediaPath,
		FileName:     filepath.Base(blendFilePath),
		Dependencies: convertDependencies(profile.Dependencies),
		Strict:       profile.Strict,
		Media:        media,
		UpdatedAt:    modTime,
	}, nil
}

func resolveRootPath(filePath string, rootPaths []string) string {
	for _, rootPath := range rootPaths {
		if strings.HasPrefix(filePath, rootPath) {
			return rootPath
		}
	}

	return ""
}

func findProjectRoot(filePath, rootPath string) string {
	if !strings.HasPrefix(filePath, rootPath) {
		return ""
	}

	currentPath := filepath.Dir(filePath)

	for {
		if _, err := os.Stat(filepath.Join(currentPath, rbtypes.ProfileDirName)); err == nil {
			return currentPath
		}

		if currentPath == rootPath {
			break
		}

		currentPath = filepath.Dir(currentPath)
	}

	return ""
}

// func findProjectViaStore(filePath string, store types.Store, logger types.Logger) string {
// 	reference := path.Clean(filepath.Dir(filePath))
// 	indexes, _ := store.List(
// 		context.Background(),
// 		listoption.WithReferences(reference),
// 		// listoption.WithType(indextype.Project),
// 		// listoption.WithSize(1),
// 	)

// 	if len(indexes) > 0 {
// 		index := indexes[0]
// 		logger.Debug("found project index", map[string]interface{}{
// 			"id":        index.ID,
// 			"reference": index.Reference,
// 		})

// 		if strings.HasPrefix(filePath, index.Reference) {
// 			return index.Reference
// 		}
// 	}

// 	logger.Debug("unable to resolve project path via store", map[string]interface{}{
// 		"filePath": filePath,
// 	})

// 	return ""
// }

func loadOrCreateProfile(validator rbtypes.Validator, path string, defaultBuild reference.Reference) (*rbtypes.Profile, error) {
	profileFilePath := profileFilePath(path)
	_, err := os.Stat(profileFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			profile := &rbtypes.Profile{
				Dependencies: []*rbtypes.Dependency{{Reference: defaultBuild, Type: rbtypes.PackageBuild}},
			}

			if err := rbhelpers.Save(validator, profileFilePath, profile, true, false); err != nil {
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
	detailFilePath := detailFilePath(path)
	_, err := os.Stat(detailFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			detail := &types.Detail{
				ID:        uuid.New(),
				Name:      helpers.FilenameToDisplayName(blendFilePath),
				MediaPath: DefaultMediaPath,
			}

			if err := rbhelpers.Save(validator, detailFilePath, detail, true, false); err != nil {
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

func convertDependencies(dependencies []*rbtypes.Dependency) []*types.Dependency {
	var deps []*types.Dependency
	for _, d := range dependencies {
		deps = append(deps, &types.Dependency{
			Reference: d.Reference,
			Type:      enums.PackageType(d.Type),
		})
	}

	return deps
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

func detailFilePath(path string) string {
	return filepath.Join(path, rbtypes.ProfileDirName, types.DetailFileName)
}

func profileFilePath(path string) string {
	return filepath.Join(path, rbtypes.ProfileDirName, rbtypes.ProfileFileName)
}
