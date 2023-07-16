package projectwatcher

import (
	"context"
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
	Watcher interface {
		AddWatchPath(path string) error
		Close() error
		GetProjects() []*project.Project
		GetProject(projectPath string) (*project.Project, bool)
	}

	watcher struct {
		logger    logger.Logger
		projects  map[string]*project.Project
		events    chan notify.EventInfo
		mu        sync.RWMutex
		ctx       context.Context
		cancel    context.CancelFunc
		isRunning bool
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

func New(opts ...Option) (Watcher, error) {
	options := &Options{
		Logger: logger.NoOp(),
	}

	for _, o := range opts {
		o(options)
	}

	return &watcher{
		logger:   options.Logger,
		projects: make(map[string]*project.Project),
		events:   make(chan notify.EventInfo, 1),
	}, nil
}

func (w *watcher) AddWatchPath(path string) error {
	err := notify.Watch(path+"/...", w.events, notify.Write|notify.Remove|notify.Rename)
	if err != nil {
		w.logger.Error("Error while adding path to watcher", map[string]interface{}{
			"path": path,
			"err":  err,
		})

		return err
	}

	w.mu.Lock()
	if !w.isRunning {
		w.ctx, w.cancel = context.WithCancel(context.Background())
		go w.run(w.ctx)
		w.isRunning = true
	}
	w.mu.Unlock()

	w.logger.Info("Added new watch path", map[string]interface{}{
		"path": path,
	})

	return nil
}

func (w *watcher) Close() error {
	notify.Stop(w.events)
	w.mu.Lock()
	if w.isRunning {
		w.cancel()
		w.isRunning = false
	}
	w.mu.Unlock()

	w.logger.Info("Project watcher closed", nil)
	return nil
}

func (w *watcher) GetProjects() []*project.Project {
	w.mu.RLock()
	defer w.mu.RUnlock()

	projects := make([]*project.Project, len(w.projects))
	for _, value := range w.projects {
		projects = append(projects, value)
	}

	return projects
}

func (w *watcher) GetProject(projectPath string) (*project.Project, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	project, ok := w.projects[projectPath]
	return project, ok
}

func (w *watcher) updateProject(projectPath string) {
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

func (w *watcher) removeProject(projectPath string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	delete(w.projects, projectPath)

	w.logger.Trace("Project removed", map[string]interface{}{
		"projectPath": projectPath,
	})
}
func (w *watcher) run(ctx context.Context) {
	for {
		select {
		case event := <-w.events:
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
