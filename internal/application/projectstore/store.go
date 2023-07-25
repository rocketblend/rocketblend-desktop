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
		RegisterPath(path string) error
		Close() error
		List(opts ...listoptions.ListOption) ([]*project.Project, error)
		Get(key string) (*project.Project, error)
	}

	store struct {
		logger logger.Logger
		index  bleve.Index

		watcherEnabled bool
		events         chan notify.EventInfo
		mu             sync.RWMutex
		ctx            context.Context
		cancel         context.CancelFunc
		isRunning      bool
	}

	Options struct {
		Logger         logger.Logger
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
		logger:         options.Logger,
		index:          index,
		watcherEnabled: options.WatcherEnabled,
		events:         make(chan notify.EventInfo, 1),
	}

	if options.WatcherEnabled {
		s.ctx, s.cancel = context.WithCancel(context.Background())
		go s.run(s.ctx)
		s.isRunning = true
	}

	return s, nil
}

func (s *store) RegisterPath(path string) error {
	// Load all projects initially
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

	s.logger.Info("Added new watch path", map[string]interface{}{
		"path": path,
	})

	return nil
}

func (s *store) Close() error {
	notify.Stop(s.events)
	s.mu.Lock()
	if s.isRunning {
		s.cancel()
		s.isRunning = false
	}
	s.mu.Unlock()

	s.logger.Info("Project watcher closed", nil)
	return nil
}

func (s *store) loadProject(key string) {
	// Get or create the project
	project, err := project.Load(key)
	if err != nil {
		s.logger.Error("Error while getting or creating project", map[string]interface{}{
			"key": key,
			"err": err,
		})

		return
	}

	if err := s.updateIndex(key, project); err != nil {
		s.logger.Error("Error while indexing project", map[string]interface{}{
			"key": key,
			"err": err,
		})
	}

	s.logger.Debug("Project updated", map[string]interface{}{
		"key":     key,
		"project": project,
	})
}

func (s *store) removeProject(key string) {
	if err := s.removeIndex(key); err != nil {
		s.logger.Error("Error while unindexing project", map[string]interface{}{
			"key": key,
			"err": err,
		})
	}

	s.logger.Debug("Project removed", map[string]interface{}{
		"key": key,
	})
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
		s.loadProject(projectPath)

	case notify.Remove, notify.Rename:
		s.logger.Debug("Removed or renamed file", map[string]interface{}{
			"file": event.Path(),
		})
		s.removeProject(projectPath)
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
