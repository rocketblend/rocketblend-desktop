package projectstore

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend-desktop/internal/application/project"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectstore/listoptions"
)

type (
	Store interface {
		Close() error
		ListProjects(opts ...listoptions.ListOption) ([]*project.Project, error)
		GetProject(key string) (*project.Project, error)
		RegisterPaths(path ...string) error
		UnregisterPaths(path ...string) error
		GetRegisteredPaths() []string
	}

	store struct {
		logger          logger.Logger
		index           bleve.Index
		registeredPaths []string

		watcherEnabled   bool
		debounceDuration time.Duration

		watchers map[string]*watcher
		events   map[string]*projectEvent

		mu  sync.RWMutex
		emu sync.RWMutex
	}

	Options struct {
		Logger           logger.Logger
		Paths            []string
		WatcherEnabled   bool
		DebounceDuration time.Duration
	}

	Option func(*Options)
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithWatcher() Option {
	return func(o *Options) {
		o.WatcherEnabled = true
	}
}

func WithPaths(paths ...string) Option {
	return func(o *Options) {
		o.Paths = paths
	}
}

func WithDebounceDuration(duration time.Duration) Option {
	return func(o *Options) {
		o.DebounceDuration = duration
	}
}

func New(opts ...Option) (Store, error) {
	options := &Options{
		Logger:           logger.NoOp(),
		DebounceDuration: 500 * time.Millisecond,
	}

	for _, o := range opts {
		o(options)
	}

	indexMapping := newIndexMapping()
	index, err := bleve.NewMemOnly(indexMapping)
	if err != nil {
		return nil, err
	}

	s := &store{
		logger:           options.Logger,
		index:            index,
		watcherEnabled:   options.WatcherEnabled,
		debounceDuration: options.DebounceDuration,
		watchers:         make(map[string]*watcher),
		events:           make(map[string]*projectEvent),
		registeredPaths:  make([]string, 0),
	}

	if err := s.RegisterPaths(options.Paths...); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *store) RegisterPaths(paths ...string) error {
	errChan := make(chan error, len(paths))
	var wg sync.WaitGroup
	wg.Add(len(paths))

	for _, path := range paths {
		go func(path string) {
			defer wg.Done()

			if err := s.registerPath(path); err != nil {
				errChan <- fmt.Errorf("failed to register path %s: %w", path, err)
			}
		}(path)
	}

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return err
	}

	return nil
}

func (s *store) UnregisterPaths(paths ...string) error {
	errChan := make(chan error, len(paths))
	var wg sync.WaitGroup
	wg.Add(len(paths))

	for _, path := range paths {
		go func(path string) {
			defer wg.Done()

			if err := s.unregisterPath(path); err != nil {
				errChan <- fmt.Errorf("failed to unregister path %s: %w", path, err)
			}
		}(path)
	}

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return err
	}

	return nil
}

func (s *store) GetRegisteredPaths() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return append([]string(nil), s.registeredPaths...) // return copy of slice
}

func (s *store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	paths := s.GetRegisteredPaths()
	if err := s.UnregisterPaths(paths...); err != nil {
		s.logger.Error("Failed to unregister paths", map[string]interface{}{
			"paths": paths,
			"err":   err,
		})
	}

	s.logger.Info("Project watcher closed successfully")
	return nil
}

func (s *store) registerPath(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the path is already registered or it is a subpath of a registered path
	for _, registeredPath := range s.registeredPaths {
		rel, err := filepath.Rel(registeredPath, path)
		if err != nil || strings.HasPrefix(rel, "..") {
			continue
		}

		return fmt.Errorf("path %s is already registered or is a subpath of a registered path", path)
	}

	// Walk the file tree starting at 'path'
	if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}

		if info.IsDir() {
			s.loadProject(path)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error while walking the path %s: %w", path, err)
	}

	// Watch the path
	if err := s.watchPath(path); err != nil {
		return fmt.Errorf("failed to watch path %s: %w", path, err)
	}

	// Add the path to the registered paths
	s.registeredPaths = append(s.registeredPaths, path)

	return nil
}

func (s *store) unregisterPath(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	pathIndex := -1
	for i, registeredPath := range s.registeredPaths {
		if registeredPath == path {
			pathIndex = i
			break
		}
	}

	if pathIndex == -1 {
		return fmt.Errorf("path %s not found", path)
	}

	// Remove indexed projects within unregistered path.
	s.removeProjectsInPath(path)

	// Unwatch the path
	if err := s.unwatchPath(path); err != nil {
		return fmt.Errorf("failed to unwatch path %s: %w", path, err)
	}

	// Remove the path from the registered paths
	s.registeredPaths = append(s.registeredPaths[:pathIndex], s.registeredPaths[pathIndex+1:]...)

	return nil
}

func (s *store) loadProject(key string) error {
	// Get or create the project
	project, err := project.Load(key)
	if err != nil {
		return fmt.Errorf("failed to load project %s: %w", key, err)
	}

	if err := s.updateIndex(key, project); err != nil {
		return fmt.Errorf("failed to update index for project %s: %w", key, err)
	}

	s.logger.Debug("Project loaded and indexed", map[string]interface{}{
		"key":     key,
		"project": project,
	})

	return nil
}

func (s *store) removeProject(key string) error {
	if err := s.removeIndex(key); err != nil {
		return fmt.Errorf("failed to remove index for project %s: %w", key, err)
	}

	s.logger.Debug("Project removed from index", map[string]interface{}{
		"key": key,
	})

	return nil
}

func (s *store) removeProjectsInPath(path string) {
	query := bleve.NewPrefixQuery(path)
	search := bleve.NewSearchRequest(query)
	searchResults, err := s.index.Search(search)
	if err != nil {
		s.logger.Error("Error searching for projects in path", map[string]interface{}{
			"err": err,
		})

		return
	}

	for _, hit := range searchResults.Hits {
		if err := s.index.Delete(hit.ID); err != nil {
			s.logger.Error("Error deleting project from index", map[string]interface{}{
				"err":  err,
				"key":  hit.ID,
				"path": path,
			})
		} else {
			s.logger.Info("Deleted project from index", map[string]interface{}{
				"key":  hit.ID,
				"path": path,
			})
		}
	}
}

func (s *store) getProjectPath(path string) string {
	projectPath := filepath.Dir(path)

	if projectPath == project.ConfigDir {
		projectPath = filepath.Dir(projectPath)
	}

	return projectPath
}
