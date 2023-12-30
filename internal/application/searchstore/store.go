package searchstore

import (
	"context"
	"sync"

	"github.com/blevesearch/bleve/v2"
	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/listoption"
)

type (
	Store interface {
		List(opts ...listoption.ListOption) ([]*Index, error)
		Get(id uuid.UUID) (*Index, error)
		Insert(index *Index) error
		Remove(id uuid.UUID) error
		RemoveByReference(path string) error

		RegisterListener(ctx context.Context, id string, handler EventHandler) func()
		UnregisterListener(id string)
		ClearListeners()
	}

	store struct {
		logger logger.Logger
		index  bleve.Index

		listeners map[string]EventHandler
		lock      sync.Mutex
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

func New(opts ...Option) (Store, error) {
	options := &Options{
		Logger: logger.NoOp(),
	}

	for _, o := range opts {
		o(options)
	}

	indexMapping := newIndexMapping()
	index, err := bleve.NewMemOnly(indexMapping)
	if err != nil {
		return nil, err
	}

	return &store{
		logger:    options.Logger,
		index:     index,
		listeners: make(map[string]EventHandler),
	}, nil
}

func (s *store) Insert(index *Index) error {
	existing, _ := s.index.Document(index.ID.String())

	err := s.index.Index(index.ID.String(), index)
	if err != nil {
		return err
	}

	if existing != nil {
		s.emitEvent(Event{
			ID:        index.ID,
			IndexType: index.Type,
			Type:      UpdateEvent,
		})
	}

	s.logger.Debug("Indexed successful", map[string]interface{}{
		"id":       index.ID,
		"type":     index.Type,
		"resource": index.Resources,
	})

	return nil
}

func (s *store) Remove(id uuid.UUID) error {
	return s.remove(id)
}

func (s *store) RemoveByReference(path string) error {
	query := bleve.NewMatchQuery(path)
	query.SetField("reference")
	searchResults, err := s.index.Search(bleve.NewSearchRequest(query))
	if err != nil {
		s.logger.Error("Error searching for indexes with reference", map[string]interface{}{
			"err": err,
		})

		return err
	}

	for _, hit := range searchResults.Hits {
		id, err := uuid.Parse(hit.ID)
		if err != nil {
			s.logger.Error("Error parsing id", map[string]interface{}{
				"err":  err,
				"key":  hit.ID,
				"path": path,
			})
		} else {
			if err := s.remove(id); err != nil {
				s.logger.Error("Error deleting index from id", map[string]interface{}{
					"err":  err,
					"key":  hit.ID,
					"path": path,
				})
			} else {
				s.logger.Info("Deleted index from id", map[string]interface{}{
					"key":  hit.ID,
					"path": path,
				})
			}
		}
	}

	return nil
}

func (s *store) remove(id uuid.UUID) error {
	index, err := s.get(id)
	if err != nil {
		return err
	}

	err = s.index.Delete(id.String())
	if err != nil {
		return err
	}

	s.emitEvent(Event{
		ID:        index.ID,
		IndexType: index.Type,
		Type:      RemoveEvent,
	})

	s.logger.Debug("Removed successful", map[string]interface{}{
		"id": id,
	})

	return nil
}
