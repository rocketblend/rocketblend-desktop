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
	Watcher struct {
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

func New(opts ...Option) (*Watcher, error) {
	options := &Options{
		Logger: logger.NoOp(),
	}

	for _, o := range opts {
		o(options)
	}

	return &Watcher{
		logger:   options.Logger,
		projects: make(map[string]*project.Project),
		events:   make(chan notify.EventInfo, 1),
	}, nil
}

func (w *Watcher) AddWatchPath(path string) error {
	err := notify.Watch(path+"/...", w.events, notify.Write|notify.Remove|notify.Rename)
	if err != nil {
		w.logger.Error("Error while adding path to watcher", map[string]interface{}{
			"path": path,
			"err":  err,
		})

		return err
	}

	w.logger.Info("Added new watch path", map[string]interface{}{
		"path": path,
	})

	return nil
}

func (w *Watcher) Close() error {
	notify.Stop(w.events)
	w.logger.Info("Project watcher closed", nil)
	return nil
}

func (w *Watcher) GetProjects() []*project.Project {
	w.mu.RLock()
	defer w.mu.RUnlock()

	projects := make([]*project.Project, len(w.projects))
	for _, value := range w.projects {
		projects = append(projects, value)
	}

	return projects
}

func (w *Watcher) GetProject(projectPath string) (*project.Project, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	project, ok := w.projects[projectPath]
	return project, ok
}

func (w *Watcher) updateProject(projectPath string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Get or create the project
	project, err := project.Find(projectPath)
	if err != nil {
		w.logger.Error("Error while getting or creating project", map[string]interface{}{
			"projectPath": projectPath,
			"err":         err,
		})

		return
	}

	w.projects[projectPath] = project

	w.logger.Trace("Project updated", map[string]interface{}{
		"projectPath": projectPath,
		"project":     project,
	})
}

func (w *Watcher) removeProject(projectPath string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	delete(w.projects, projectPath)

	w.logger.Trace("Project removed", map[string]interface{}{
		"projectPath": projectPath,
	})
}

func (w *Watcher) Run() {
	for event := range w.events {
		switch event.Event() {
		case notify.Write:
			w.logger.Debug("Modified file", map[string]interface{}{
				"file": event.Path(),
			})

			if isWatchFile(event.Path()) {
				w.updateProject(filepath.Dir(event.Path()))
			}
		case notify.Remove, notify.Rename:
			w.logger.Debug("Removed or renamed file", map[string]interface{}{
				"file": event.Path(),
			})

			if isWatchFile(event.Path()) {
				w.removeProject(filepath.Dir(event.Path()))
			}
		}

		w.logger.Info("Filesystem event occurred", map[string]interface{}{
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
