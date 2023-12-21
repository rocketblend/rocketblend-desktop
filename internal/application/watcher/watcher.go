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
		RegisterPaths(path ...string) error
		UnregisterPaths(path ...string) error
		GetRegisteredPaths() []string
	}

	UpdateObjectFunc      func(path string) error
	RemoveObjectFunc      func(path string) error
	ResolveObjectPathFunc func(path string) string
	IsWatchableFileFunc   func(path string) bool

	service struct {
		logger          logger.Logger
		registeredPaths []string

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
		registeredPaths:       make([]string, 0),
		updateObjectFunc:      options.UpdateObjectFunc,
		removeObjectFunc:      options.RemoveObjectFunc,
		resolveObjectPathFunc: options.ResolveObjectPathFunc,
		isWatchableFileFunc:   options.IsWatchableFileFunc,
	}

	if err := s.RegisterPaths(options.Paths...); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *service) RegisterPaths(paths ...string) error {
	errChan := make(chan error, len(paths))
	var wg sync.WaitGroup
	wg.Add(len(paths))

	for _, path := range paths {
		go func(path string) {
			defer wg.Done()

			if err := s.registerPath(path); err != nil {
				errChan <- fmt.Errorf("failed to register path %s: %w", path, err)
			}
		}(path)
	}

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return err
	}

	return nil
}

func (s *service) UnregisterPaths(paths ...string) error {
	errChan := make(chan error, len(paths))
	var wg sync.WaitGroup
	wg.Add(len(paths))

	for _, path := range paths {
		go func(path string) {
			defer wg.Done()

			if err := s.unregisterPath(path); err != nil {
				errChan <- fmt.Errorf("failed to unregister path %s: %w", path, err)
			}
		}(path)
	}

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return err
	}

	return nil
}

func (s *service) GetRegisteredPaths() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return append([]string(nil), s.registeredPaths...) // return copy of slice
}

func (s *service) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var closeErrors []string
	s.logger.Debug("Closing watcher")

	paths := s.registeredPaths
	for _, path := range paths {
		if err := s.unregisterPathLocked(path); err != nil {
			closeErrors = append(closeErrors, fmt.Sprintf("path %s: %v", path, err))
			s.logger.Error("Failed to unregister path", map[string]interface{}{
				"path": path,
				"err":  err,
			})
		}
	}

	s.registeredPaths = []string{}

	if len(closeErrors) > 0 {
		return fmt.Errorf("close completed with errors: %s", strings.Join(closeErrors, "; "))
	}

	return nil
}

func (s *service) registerPath(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the path is already registered or it is a subpath of a registered path
	for _, registeredPath := range s.registeredPaths {
		rel, err := filepath.Rel(registeredPath, path)
		if err != nil || strings.HasPrefix(rel, "..") {
			continue
		}

		return fmt.Errorf("path %s is already registered or is a subpath of a registered path", path)
	}

	// Walk the file tree starting at 'path'
	if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}

		if s.isWatchableFile(path) {
			s.handleEventDebounced(&objectEventInfo{
				ObjectPath: s.resolveObjectPath(path),
				EventInfo: eventInfo{
					path:  path,
					event: notify.Write,
				},
			})

			objectPath := s.resolveObjectPath(path)

			// Trigger initial update
			if err := s.updateObject(objectPath); err != nil {
				s.logger.Debug("Failed to load object", map[string]interface{}{
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

	// Add the path to the registered paths
	s.registeredPaths = append(s.registeredPaths, path)

	return nil
}

func (s *service) unregisterPath(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.unregisterPathLocked(path)
}

// unregisterPathLocked unregisters a path without locking the mutex.
func (s *service) unregisterPathLocked(path string) error {
	pathIndex := -1
	for i, registeredPath := range s.registeredPaths {
		if registeredPath == path {
			pathIndex = i
			break
		}
	}

	if pathIndex == -1 {
		return fmt.Errorf("path %s not found", path)
	}

	// Remove indexed projects within unregistered path.
	if err := s.removeObject(path); err != nil {
		return fmt.Errorf("failed to remove projects in path %s: %w", path, err)
	}

	// Unwatch the path
	if err := s.unwatchPath(path); err != nil {
		return fmt.Errorf("failed to unwatch path %s: %w", path, err)
	}

	// Remove the path from the registered paths
	s.registeredPaths = append(s.registeredPaths[:pathIndex], s.registeredPaths[pathIndex+1:]...)

	return nil
}

func (s *service) updateObject(path string) error {
	if s.updateObjectFunc != nil {
		return s.updateObjectFunc(path)
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
