package projectsearcher

import (
	"fmt"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend-desktop/internal/application/project"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectwatcher"
)

type (
	FindOptions struct {
	}

	FindOption func(*FindOptions)

	Searcher interface {
		FindAll(opts ...FindOption) ([]*project.Project, error)
		FindByPath(projectPath string) (*project.Project, error)
	}

	searcher struct {
		logger  logger.Logger
		watcher projectwatcher.Watcher
	}

	Options struct {
		Logger  logger.Logger
		Watcher projectwatcher.Watcher
	}

	Option func(*Options)
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithWatcher(watcher projectwatcher.Watcher) Option {
	return func(o *Options) {
		o.Watcher = watcher
	}
}

func New(opts ...Option) (Searcher, error) {
	options := &Options{
		Logger: logger.NoOp(),
	}

	for _, o := range opts {
		o(options)
	}

	if options.Watcher == nil {
		return nil, fmt.Errorf("watcher is required")
	}

	return &searcher{
		logger:  options.Logger,
		watcher: options.Watcher,
	}, nil
}

func (s *searcher) FindAll(opts ...FindOption) ([]*project.Project, error) {
	return s.watcher.GetProjects(), nil
}

func (s *searcher) FindByPath(projectPath string) (*project.Project, error) {
	project, ok := s.watcher.GetProject(projectPath)
	if !ok {
		return nil, fmt.Errorf("project not found")
	}

	return project, nil
}
