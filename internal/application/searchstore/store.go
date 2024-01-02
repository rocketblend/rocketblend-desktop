package searchstore

import (
	"context"
	"errors"

	"github.com/blevesearch/bleve/v2"
	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/eventservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/listoption"
)

const (
	InsertEventChannel = "searchstore.insert"
	RemoveEventChannel = "searchstore.remove"
)

type (
	Store interface {
		List(opts ...listoption.ListOption) ([]*Index, error)
		Get(id uuid.UUID) (*Index, error)
		Insert(index *Index) error
		Remove(id uuid.UUID) error
		RemoveByReference(path string) error

		Close() error
	}

	store struct {
		logger logger.Logger
		index  bleve.Index

		event eventservice.Service
	}

	Options struct {
		Logger logger.Logger
		Event  eventservice.Service
	}

	Option func(*Options)
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithEventService(event eventservice.Service) Option {
	return func(o *Options) {
		o.Event = event
	}
}

func New(opts ...Option) (Store, error) {
	options := &Options{
		Logger: logger.NoOp(),
	}

	for _, o := range opts {
		o(options)
	}

	if options.Event == nil {
		return nil, errors.New("event service is required")
	}

	indexMapping := newIndexMapping()
	index, err := bleve.NewMemOnly(indexMapping)
	if err != nil {
		return nil, err
	}

	return &store{
		logger: options.Logger,
		index:  index,
		event:  options.Event,
	}, nil
}

func (s *store) Insert(index *Index) error {
	existing, _ := s.index.Document(index.ID.String())

	err := s.index.Index(index.ID.String(), index)
	if err != nil {
		return err
	}

	ctx := context.Background()
	if existing != nil {
		event := NewEvent(index.ID, index.Type)
		if err := s.event.EmitEvent(ctx, InsertEventChannel, event); err != nil {
			s.logger.Error("error emitting event", map[string]interface{}{
				"err": err,
			})
		}
	}

	s.logger.Debug("indexed successful", map[string]interface{}{
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
		s.logger.Error("error searching for indexes with reference", map[string]interface{}{
			"err": err,
		})

		return err
	}

	for _, hit := range searchResults.Hits {
		id, err := uuid.Parse(hit.ID)
		if err != nil {
			s.logger.Error("error parsing id", map[string]interface{}{
				"err":  err,
				"key":  hit.ID,
				"path": path,
			})
		} else {
			if err := s.remove(id); err != nil {
				s.logger.Error("error deleting index from id", map[string]interface{}{
					"err":  err,
					"key":  hit.ID,
					"path": path,
				})
			} else {
				s.logger.Info("deleted index from id", map[string]interface{}{
					"key":  hit.ID,
					"path": path,
				})
			}
		}
	}

	return nil
}

func (s *store) Close() error {
	if s.index == nil {
		return nil
	}

	return s.index.Close()
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

	// TODO: make functions take context
	ctx := context.Background()
	event := NewEvent(index.ID, index.Type)
	if err := s.event.EmitEvent(ctx, RemoveEventChannel, event); err != nil {
		s.logger.Error("error emitting event", map[string]interface{}{
			"err": err,
		})
	}

	s.logger.Debug("Removed successful", map[string]interface{}{
		"id": id,
	})

	return nil
}
