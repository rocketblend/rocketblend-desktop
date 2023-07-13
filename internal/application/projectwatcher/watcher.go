package projectwatcher

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/fsnotify/fsnotify"
)

type Project struct {
	BlendFile string
	YamlFile  string
}

type ProjectWatcher struct {
	logger   logger.Logger
	projects map[string]*Project
	watcher  *fsnotify.Watcher
	updates  chan struct{}
	mu       sync.RWMutex
}

type (
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

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &ProjectWatcher{
		logger:   options.Logger,
		projects: make(map[string]*Project),
		updates:  make(chan struct{}),
		watcher:  watcher,
	}, nil
}

func (pw *ProjectWatcher) Close() error {
	return pw.watcher.Close()
}

func (pw *ProjectWatcher) Watch(paths ...string) error {
	for _, path := range paths {
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			return pw.watcher.Add(path)
		})
		if err != nil {
			pw.logger.Error("Error while adding path to watcher", map[string]interface{}{
				"path": path,
				"err":  err,
			})
			return err
		}
	}

	return nil
}

func (pw *ProjectWatcher) Run() {
	for {
		select {
		case event := <-pw.watcher.Events:
			pw.logger.Info("Filesystem event occurred", map[string]interface{}{
				"event": event,
			})

			if event.Op&fsnotify.Write == fsnotify.Write {
				pw.logger.Debug("Modified file", map[string]interface{}{
					"file": event.Name,
				})
				pw.updateProjects(event.Name)
				pw.updates <- struct{}{}
			}
		case err := <-pw.watcher.Errors:
			pw.logger.Error("Error from filesystem watcher", map[string]interface{}{
				"err": err,
			})
		}
	}
}

func (pw *ProjectWatcher) GetProjects() map[string]*Project {
	pw.mu.RLock()
	defer pw.mu.RUnlock()

	copiedProjects := make(map[string]*Project)
	for k, v := range pw.projects {
		copiedProjects[k] = v
	}

	return copiedProjects
}

func (pw *ProjectWatcher) GetUpdates() <-chan struct{} {
	return pw.updates
}

func (pw *ProjectWatcher) updateProjects(filePath string) {
	pw.mu.Lock()
	defer pw.mu.Unlock()

	switch filepath.Ext(filePath) {
	case ".blend":
		project := pw.getOrCreateProject(filePath)
		project.BlendFile = filePath
	case ".yaml":
		project := pw.getOrCreateProject(filePath)
		project.YamlFile = filePath
	}

	pw.logger.Trace("Projects updated", map[string]interface{}{
		"projects": pw.projects,
	})
}

func (pw *ProjectWatcher) getOrCreateProject(filePath string) *Project {
	dir := filepath.Dir(filePath)
	project, exists := pw.projects[dir]
	if !exists {
		project = &Project{}
		pw.projects[dir] = project
	}

	return project
}
