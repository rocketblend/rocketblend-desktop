package assetstore

import (
	"fmt"
	"sync"

	"github.com/flowshot-io/x/pkg/logger"
)

type (
	Store interface {
		Exists(path string) (bool, error)
		Insert(path string) error
		Remove(path string) error
	}

	store struct {
		logger   logger.Logger
		mutex    sync.RWMutex
		register map[string]int
	}

	Options struct {
		Logger         logger.Logger
		InitialCapcity int
	}

	Option func(*Options)
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithInitialCapacity(capacity int) Option {
	return func(o *Options) {
		o.InitialCapcity = capacity
	}
}

func New(opts ...Option) (Store, error) {
	options := &Options{
		Logger:         logger.NoOp(),
		InitialCapcity: 50,
	}

	for _, o := range opts {
		o(options)
	}

	s := &store{
		logger:   options.Logger,
		register: make(map[string]int, options.InitialCapcity),
	}

	return s, nil
}

func (s *store) Exists(path string) (bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	count, exists := s.register[path]
	if !exists {
		return false, fmt.Errorf("path %s does not exist", path)
	}

	s.logger.Info("Path exists", map[string]interface{}{"path": path, "count": count})
	return true, nil
}

func (s *store) Insert(path string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.register[path]++
	s.logger.Info("Path inserted/incremented", map[string]interface{}{"path": path, "count": s.register[path]})
	return nil
}

func (s *store) Remove(path string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if count, ok := s.register[path]; ok {
		if count > 1 {
			s.register[path]--
			s.logger.Info("Path count decremented", map[string]interface{}{"path": path, "count": s.register[path]})
		} else {
			delete(s.register, path)
			s.logger.Info("Path removed", map[string]interface{}{"path": path})
		}
		return nil
	}

	s.logger.Warn("Attempt to remove non-existent path", map[string]interface{}{"path": path})
	return fmt.Errorf("path %s does not exist", path)
}
