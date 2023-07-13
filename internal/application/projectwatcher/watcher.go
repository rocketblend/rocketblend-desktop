package projectwatcher

import (
	"path/filepath"
	"sync"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rjeczalik/notify"
)

type (
	Project struct {
		BlendFile string
		YamlFile  string
	}

	ProjectWatcher struct {
		logger   logger.Logger
		projects map[string]*Project
		events   chan notify.EventInfo
		updates  chan *Project
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
		projects: make(map[string]*Project),
		events:   make(chan notify.EventInfo, 1),
		updates:  make(chan *Project),
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

func (pw *ProjectWatcher) GetProjects() map[string]*Project {
	pw.mu.RLock()
	defer pw.mu.RUnlock()

	// Return a copy of the map to prevent concurrent modification
	projects := make(map[string]*Project)
	for key, value := range pw.projects {
		projects[key] = value
	}

	return projects
}

func (pw *ProjectWatcher) GetUpdates() <-chan *Project {
	return pw.updates
}

func (pw *ProjectWatcher) getOrCreateProject(filename string) *Project {
	dir := filepath.Dir(filename)

	if project, ok := pw.projects[dir]; ok {
		return project
	}

	project := &Project{}
	pw.projects[dir] = project

	return project
}

func (pw *ProjectWatcher) updateProjects(filename string) {
	pw.mu.Lock()
	defer pw.mu.Unlock()

	switch filepath.Ext(filename) {
	case ".blend":
		project := pw.getOrCreateProject(filename)
		project.BlendFile = filename
		pw.updates <- project
	case ".yaml":
		project := pw.getOrCreateProject(filename)
		project.YamlFile = filename
		pw.updates <- project
	}

	pw.logger.Trace("Projects updated", map[string]interface{}{
		"projects": pw.projects,
	})
}

func (pw *ProjectWatcher) removeProject(filename string) {
	pw.mu.Lock()
	defer pw.mu.Unlock()

	dir := filepath.Dir(filename)
	delete(pw.projects, dir)

	pw.logger.Trace("Project removed", map[string]interface{}{
		"project_dir": dir,
	})
}

func (pw *ProjectWatcher) Run() {
	for event := range pw.events {
		switch event.Event() {
		case notify.Write:
			pw.logger.Debug("Modified file", map[string]interface{}{
				"file": event.Path(),
			})

			pw.updateProjects(event.Path())
		case notify.Remove, notify.Rename:
			pw.logger.Debug("Removed or renamed file", map[string]interface{}{
				"file": event.Path(),
			})

			pw.removeProject(event.Path())
		}

		pw.logger.Info("Filesystem event occurred", map[string]interface{}{
			"event": event,
		})
	}
}
