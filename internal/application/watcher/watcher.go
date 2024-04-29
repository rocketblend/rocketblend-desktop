package watcher

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rjeczalik/notify"
)

type (
	Watcher interface {
		Close() error
		SetPaths(paths ...string) error
	}

	UpdateObjectFunc      func(path string) error
	RemoveObjectFunc      func(path string) error
	ResolveObjectPathFunc func(path string) string
	IsWatchableFileFunc   func(path string) bool

	service struct {
		logger logger.Logger
		paths  map[string]struct{}

		updateObjectFunc      UpdateObjectFunc
		removeObjectFunc      RemoveObjectFunc
		resolveObjectPathFunc ResolveObjectPathFunc
		isWatchableFileFunc   IsWatchableFileFunc

		debounceDuration time.Duration

		watchers map[string]*watcher
		events   map[string]*projectEvent

		mu  sync.RWMutex
		emu sync.RWMutex
	}

	Options struct {
		Logger           logger.Logger
		Paths            []string
		DebounceDuration time.Duration

		UpdateObjectFunc      UpdateObjectFunc
		RemoveObjectFunc      RemoveObjectFunc
		ResolveObjectPathFunc ResolveObjectPathFunc
		IsWatchableFileFunc   IsWatchableFileFunc
	}

	Option func(*Options)
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithPaths(paths ...string) Option {
	return func(o *Options) {
		o.Paths = paths
	}
}

func WithEventDebounceDuration(duration time.Duration) Option {
	return func(o *Options) {
		o.DebounceDuration = duration
	}
}

func WithUpdateObjectFunc(f UpdateObjectFunc) Option {
	return func(o *Options) { o.UpdateObjectFunc = f }
}

func WithRemoveObjectFunc(f RemoveObjectFunc) Option {
	return func(o *Options) { o.RemoveObjectFunc = f }
}

func WithResolveObjectPathFunc(f ResolveObjectPathFunc) Option {
	return func(o *Options) { o.ResolveObjectPathFunc = f }
}

func WithIsWatchableFileFunc(f IsWatchableFileFunc) Option {
	return func(o *Options) { o.IsWatchableFileFunc = f }
}

func New(opts ...Option) (Watcher, error) {
	options := &Options{
		Logger:           logger.NoOp(),
		DebounceDuration: 500 * time.Millisecond,
	}

	for _, o := range opts {
		o(options)
	}

	s := &service{
		logger:                options.Logger,
		debounceDuration:      options.DebounceDuration,
		watchers:              make(map[string]*watcher),
		events:                make(map[string]*projectEvent),
		paths:                 make(map[string]struct{}),
		updateObjectFunc:      options.UpdateObjectFunc,
		removeObjectFunc:      options.RemoveObjectFunc,
		resolveObjectPathFunc: options.ResolveObjectPathFunc,
		isWatchableFileFunc:   options.IsWatchableFileFunc,
	}

	if err := s.setPaths(options.Paths...); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *service) SetPaths(newPaths ...string) error {
	return s.setPaths(newPaths...)
}

func (s *service) setPaths(paths ...string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Convert newPaths to a map for efficient lookup
	pathMap := make(map[string]struct{}, len(paths))
	for _, path := range paths {
		if path == "" {
			continue
		}

		pathMap[path] = struct{}{}
	}

	// Unregister paths not in the new set of paths
	for path := range s.paths {
		if _, exists := pathMap[path]; !exists {
			if err := s.unregisterPath(path); err != nil {
				s.logger.Error("failed to unregister path", map[string]interface{}{"path": path, "error": err})
			}
		}
	}

	// Register new paths
	for path := range pathMap {
		if _, alreadyRegistered := s.paths[path]; !alreadyRegistered {
			if err := s.registerPath(path); err != nil {
				s.logger.Error("failed to register path", map[string]interface{}{"path": path, "error": err})
			}
		}
	}

	return nil
}

func (s *service) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var closeErrors []string
	s.logger.Debug("closing watcher")

	for path := range s.paths {
		if err := s.unregisterPath(path); err != nil {
			closeErrors = append(closeErrors, fmt.Sprintf("path %s: %v", path, err))
			s.logger.Error("failed to unregister path", map[string]interface{}{
				"path": path,
				"err":  err,
			})
		}
	}

	s.paths = make(map[string]struct{})

	if len(closeErrors) > 0 {
		return fmt.Errorf("close completed with errors: %s", strings.Join(closeErrors, "; "))
	}

	return nil
}

func (s *service) registerPath(path string) error {
	// Check if the path is already registered or it is a subpath of a registered path
	for registeredPath := range s.paths {
		rel, err := filepath.Rel(registeredPath, path)
		if err != nil || !strings.HasPrefix(rel, "..") {
			return fmt.Errorf("path %s is already registered or is a subpath of a registered path", path)
		}
	}

	// Walk the file tree starting at 'path'
	if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}

		if s.isWatchableFileFunc != nil && s.isWatchableFileFunc(path) {
			objectPath := s.resolveObjectPath(path)

			// Assume handleEventDebounced and updateObject are appropriately adjusted to handle the map-based structure
			s.handleEventDebounced(&objectEventInfo{
				ObjectPath: objectPath,
				EventInfo: eventInfo{
					path:  path,
					event: notify.Write,
				},
			})

			// Trigger initial update
			if err := s.updateObject(objectPath); err != nil {
				s.logger.Error("failed to update watched object", map[string]interface{}{
					"err":  err,
					"path": path,
				})
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error while walking the path %s: %w", path, err)
	}

	// Watch the path
	if err := s.watchPath(path); err != nil {
		return fmt.Errorf("failed to watch path %s: %w", path, err)
	}

	// Add the path to the registered paths map
	s.paths[path] = struct{}{}

	return nil
}

func (s *service) unregisterPath(path string) error {
	if _, exists := s.paths[path]; !exists {
		return fmt.Errorf("path %s not found", path)
	}

	// Remove the path from the registered paths
	delete(s.paths, path)

	// Remove indexed projects within unregistered path.
	if err := s.removeObject(path); err != nil {
		return fmt.Errorf("failed to remove projects in path %s: %w", path, err)
	}

	// Unwatch the path
	if err := s.unwatchPath(path); err != nil {
		return fmt.Errorf("failed to unwatch path %s: %w", path, err)
	}

	return nil
}

func (s *service) updateObject(path string) error {
	if s.updateObjectFunc != nil {
		if err := s.updateObjectFunc(path); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) removeObject(path string) error {
	if s.removeObjectFunc != nil {
		return s.removeObjectFunc(path)
	}

	return nil
}

func (s *service) resolveObjectPath(path string) string {
	if s.resolveObjectPathFunc != nil {
		return s.resolveObjectPathFunc(path)
	}

	return path
}
