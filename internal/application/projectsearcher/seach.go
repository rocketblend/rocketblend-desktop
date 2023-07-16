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
		watcher *projectwatcher.Watcher
	}

	Options struct {
		Logger    logger.Logger
		WatchPath string
	}

	Option func(*Options)
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithWatchPath(watchPath string) Option {
	return func(o *Options) {
		o.WatchPath = watchPath
	}
}

func New(opts ...Option) (Searcher, error) {
	options := &Options{
		Logger: logger.NoOp(),
	}

	for _, o := range opts {
		o(options)
	}

	if options.WatchPath == "" {
		return nil, fmt.Errorf("root path is required")
	}

	watcher, err := projectwatcher.New(projectwatcher.WithLogger(options.Logger))
	if err != nil {
		return nil, err
	}

	if err := watcher.AddWatchPath(options.WatchPath); err != nil {
		return nil, err
	}

	go watcher.Run()
	//defer watcher.Close()

	return &searcher{
		logger:  options.Logger,
		watcher: watcher,
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
