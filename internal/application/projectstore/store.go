package projectstore

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rjeczalik/notify"
	"github.com/rocketblend/rocketblend-desktop/internal/application/project"
)

var watchFileExtensions = []string{
	//".blend",
	".yaml",
}

type (
	Store interface {
		AddWatchPath(path string) error
		Close() error
		ListProjects() []*project.Project
		GetProject(key string) (*project.Project, bool)
	}

	store struct {
		logger   logger.Logger
		projects map[string]*project.Project

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

	return &store{
		logger:         options.Logger,
		watcherEnabled: options.WatcherEnabled,
		projects:       make(map[string]*project.Project),
		events:         make(chan notify.EventInfo, 1),
	}, nil
}

func (s *store) AddWatchPath(path string) error {
	// Load all projects initially
	if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			s.loadProject(path)
		}

		return nil
	}); err != nil {
		s.logger.Error("Error while walking the path", map[string]interface{}{
			"path": path,
			"err":  err,
		})

		return err
	}

	// TODO: Enabled/disable file watching
	// Start watching the path
	err := notify.Watch(path+"/...", s.events, notify.Write|notify.Remove|notify.Rename)
	if err != nil {
		s.logger.Error("Error while adding path to watcher", map[string]interface{}{
			"path": path,
			"err":  err,
		})

		return err
	}

	// Run event handler if not already running
	s.mu.Lock()
	if !s.isRunning {
		s.ctx, s.cancel = context.WithCancel(context.Background())
		go s.run(s.ctx)
		s.isRunning = true
	}
	s.mu.Unlock()

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

func (s *store) ListProjects() []*project.Project {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// TODO: Indexing, filtering, sorting, etc.
	projects := make([]*project.Project, 0, len(s.projects))
	for _, value := range s.projects {
		projects = append(projects, value)
	}

	return projects
}

func (s *store) GetProject(key string) (*project.Project, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	project, ok := s.projects[key]
	return project, ok
}

func (s *store) loadProject(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get or create the project
	project, err := project.Load(key)
	if err != nil {
		s.logger.Error("Error while getting or creating project", map[string]interface{}{
			"key": key,
			"err": err,
		})

		return
	}

	s.projects[key] = project

	s.logger.Debug("Project updated", map[string]interface{}{
		"key":     key,
		"project": project,
	})
}

func (s *store) removeProject(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.projects, key)

	s.logger.Debug("Project removed", map[string]interface{}{
		"key": key,
	})
}

func (s *store) run(ctx context.Context) {
	for {
		select {
		case event := <-s.events:
			switch event.Event() {
			case notify.Write:
				s.logger.Debug("Modified file", map[string]interface{}{
					"file": event.Path(),
				})

				//TODO: If file is in a subdirectory of the project config, then we should reload the project
				if isWatchFile(event.Path()) {
					projectPath := filepath.Dir(event.Path())

					// If the project path is the config dir, then we need to go up one more level
					if projectPath == project.ConfigDir {
						projectPath = filepath.Dir(projectPath)
					}

					// TODO: Add a debounce here. If the file is modified multiple times in a short period of time, then we should only load the project once.
					s.loadProject(projectPath)
				}
			case notify.Remove, notify.Rename:
				s.logger.Debug("Removed or renamed file", map[string]interface{}{
					"file": event.Path(),
				})

				if isWatchFile(event.Path()) {
					projectPath := filepath.Dir(event.Path())

					if projectPath == project.ConfigDir {
						projectPath = filepath.Dir(projectPath)
					}

					s.removeProject(projectPath)
				}
			}

			s.logger.Info("Filesystem event occurred", map[string]interface{}{
				"event": event,
			})
		case <-ctx.Done():
			return
		}
	}
}

func isWatchFile(filename string) bool {
	for _, ext := range watchFileExtensions {
		if filepath.Ext(filename) == ext {
			return true
		}
	}

	return false
}
