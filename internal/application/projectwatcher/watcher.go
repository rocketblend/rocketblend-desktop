package projectwatcher

import (
	"path/filepath"
	"sync"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rjeczalik/notify"
	"github.com/rocketblend/rocketblend-desktop/internal/application/project"
)

var watchFileExtensions = []string{
	".blend",
	".yaml",
}

type (
	ProjectWatcher struct {
		logger   logger.Logger
		projects map[string]*project.Project
		events   chan notify.EventInfo
		mu       sync.RWMutex
	}

	Options struct {
		Logger logger.Logger
	}

	Option func(*Options)
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func New(opts ...Option) (*ProjectWatcher, error) {
	options := &Options{
		Logger: logger.NoOp(),
	}

	for _, o := range opts {
		o(options)
	}

	return &ProjectWatcher{
		logger:   options.Logger,
		projects: make(map[string]*project.Project),
		events:   make(chan notify.EventInfo, 1),
	}, nil
}

func (pw *ProjectWatcher) AddWatchPath(path string) error {
	err := notify.Watch(path+"/...", pw.events, notify.Write|notify.Remove|notify.Rename)
	if err != nil {
		pw.logger.Error("Error while adding path to watcher", map[string]interface{}{
			"path": path,
			"err":  err,
		})

		return err
	}

	pw.logger.Info("Added new watch path", map[string]interface{}{
		"path": path,
	})

	return nil
}

func (pw *ProjectWatcher) Close() error {
	notify.Stop(pw.events)
	pw.logger.Info("Project watcher closed", nil)
	return nil
}

func (pw *ProjectWatcher) GetProjects() map[string]*project.Project {
	pw.mu.RLock()
	defer pw.mu.RUnlock()

	// Return a copy of the map to prevent concurrent modification
	projects := make(map[string]*project.Project)
	for key, value := range pw.projects {
		projects[key] = value
	}

	return projects
}

func (pw *ProjectWatcher) updateProject(projectPath string) {
	pw.mu.Lock()
	defer pw.mu.Unlock()

	// Get or create the project
	project, err := project.Find(projectPath)
	if err != nil {
		pw.logger.Error("Error while getting or creating project", map[string]interface{}{
			"projectPath": projectPath,
			"err":         err,
		})

		return
	}

	pw.projects[projectPath] = project

	pw.logger.Trace("Project updated", map[string]interface{}{
		"projectPath": projectPath,
		"project":     project,
	})
}

func (pw *ProjectWatcher) removeProject(projectPath string) {
	pw.mu.Lock()
	defer pw.mu.Unlock()

	delete(pw.projects, projectPath)

	pw.logger.Trace("Project removed", map[string]interface{}{
		"projectPath": projectPath,
	})
}

func (pw *ProjectWatcher) Run() {
	for event := range pw.events {
		switch event.Event() {
		case notify.Write:
			pw.logger.Debug("Modified file", map[string]interface{}{
				"file": event.Path(),
			})

			if isWatchFile(event.Path()) {
				pw.updateProject(filepath.Dir(event.Path()))
			}
		case notify.Remove, notify.Rename:
			pw.logger.Debug("Removed or renamed file", map[string]interface{}{
				"file": event.Path(),
			})

			if isWatchFile(event.Path()) {
				pw.removeProject(filepath.Dir(event.Path()))
			}
		}

		pw.logger.Info("Filesystem event occurred", map[string]interface{}{
			"event": event,
		})
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
