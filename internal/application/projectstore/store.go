package projectstore

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/blevesearch/bleve/v2"
	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rjeczalik/notify"
	"github.com/rocketblend/rocketblend-desktop/internal/application/project"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectstore/listoptions"
)

var watchFileExtensions = []string{
	//".blend",
	".yaml",
}

type (
	Store interface {
		Close() error
		ListProjects(opts ...listoptions.ListOption) ([]*project.Project, error)
		GetProject(key string) (*project.Project, error)
		RegisterPaths(path ...string) error
		//UnregisterPaths(path ...string) error
		GetRegisteredPaths() []string
	}

	store struct {
		logger          logger.Logger
		index           bleve.Index
		registeredPaths []string

		watcherEnabled bool
		events         chan notify.EventInfo
		mu             sync.RWMutex
		ctx            context.Context
		cancel         context.CancelFunc
		isRunning      bool
	}

	Options struct {
		Logger         logger.Logger
		Paths          []string
		WatcherEnabled bool
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

func New(opts ...Option) (Store, error) {
	options := &Options{
		Logger: logger.NoOp(),
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
		logger:          options.Logger,
		index:           index,
		watcherEnabled:  options.WatcherEnabled,
		events:          make(chan notify.EventInfo, 1),
		registeredPaths: make([]string, 0),
	}

	if err := s.RegisterPaths(options.Paths...); err != nil {
		return nil, err
	}

	if options.WatcherEnabled {
		s.ctx, s.cancel = context.WithCancel(context.Background())
		go s.run(s.ctx)
		s.isRunning = true
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

func (s *store) GetRegisteredPaths() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return append([]string(nil), s.registeredPaths...) // return copy of slice
}

func (s *store) Close() error {
	notify.Stop(s.events)
	s.logger.Info("Stopped notification events")

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		s.cancel()
		s.isRunning = false
		s.logger.Debug("Stopped the watcher and updated the state")
	} else {
		s.logger.Debug("The watcher was already stopped")
	}

	s.logger.Info("Project watcher closed successfully")
	return nil
}

func (s *store) registerPath(path string) error {
	// TODO: Check if the path is already registered

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

	if s.watcherEnabled {
		// Start watching the path
		err := notify.Watch(path+"/...", s.events, notify.Write|notify.Remove|notify.Rename)
		if err != nil {
			return fmt.Errorf("unable to add path %s to watcher: %w", path, err)
		}
	}

	s.mu.Lock()
	s.registeredPaths = append(s.registeredPaths, path)
	s.mu.Unlock()

	s.logger.Info("Path registered", map[string]interface{}{
		"path": path,
	})

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

func (s *store) run(ctx context.Context) {
	for {
		select {
		case event := <-s.events:
			s.handleEvent(event)
		case <-ctx.Done():
			return
		}
	}
}

func (s *store) handleEvent(event notify.EventInfo) {
	// Log the event
	s.logger.Info("Filesystem event occurred", map[string]interface{}{
		"event": event,
	})

	if !isWatchFile(event.Path()) {
		return
	}

	projectPath := s.getProjectPath(event.Path())

	switch event.Event() {
	case notify.Write:
		s.logger.Debug("Modified file", map[string]interface{}{
			"file": event.Path(),
		})

		if err := s.loadProject(projectPath); err != nil {
			s.logger.Error("Error while loading project", map[string]interface{}{
				"err": err,
			})
		}

	case notify.Remove, notify.Rename:
		s.logger.Debug("Removed or renamed file", map[string]interface{}{
			"file": event.Path(),
		})

		if err := s.removeProject(projectPath); err != nil {
			s.logger.Error("Error while removing project", map[string]interface{}{
				"err": err,
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

func isWatchFile(filename string) bool {
	for _, ext := range watchFileExtensions {
		if filepath.Ext(filename) == ext {
			return true
		}
	}

	return false
}
