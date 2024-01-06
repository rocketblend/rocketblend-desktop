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
		List(ctx context.Context, opts ...listoption.ListOption) ([]*Index, error)
		Get(ctx context.Context, id uuid.UUID) (*Index, error)
		Insert(ctx context.Context, index *Index) error
		Remove(ctx context.Context, id uuid.UUID) error
		RemoveByReference(ctx context.Context, path string) error

		Close() error
	}

	store struct {
		logger logger.Logger
		index  bleve.Index

		dispatcher eventservice.Service
	}

	Options struct {
		Logger     logger.Logger
		Dispatcher eventservice.Service
	}

	Option func(*Options)
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithDispatcherService(dispatcher eventservice.Service) Option {
	return func(o *Options) {
		o.Dispatcher = dispatcher
	}
}

func New(opts ...Option) (Store, error) {
	options := &Options{
		Logger: logger.NoOp(),
	}

	for _, o := range opts {
		o(options)
	}

	if options.Dispatcher == nil {
		return nil, errors.New("dispatcher service is required")
	}

	indexMapping := newIndexMapping()
	index, err := bleve.NewMemOnly(indexMapping)
	if err != nil {
		return nil, err
	}

	return &store{
		logger:     options.Logger,
		index:      index,
		dispatcher: options.Dispatcher,
	}, nil
}

func (s *store) Insert(ctx context.Context, index *Index) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	existing, _ := s.index.Document(index.ID.String())

	err := s.index.Index(index.ID.String(), index)
	if err != nil {
		return err
	}

	if existing != nil {
		event := NewEvent(index.ID, index.Type)
		if err := s.dispatcher.EmitEvent(ctx, InsertEventChannel, event); err != nil {
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

func (s *store) Remove(ctx context.Context, id uuid.UUID) error {
	return s.remove(ctx, id)
}

func (s *store) RemoveByReference(ctx context.Context, path string) error {
	query := bleve.NewMatchQuery(path)
	query.SetField("reference")
	searchResults, err := s.index.SearchInContext(ctx, bleve.NewSearchRequest(query))
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
			if err := s.remove(ctx, id); err != nil {
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

func (s *store) remove(ctx context.Context, id uuid.UUID) error {
	index, err := s.get(ctx, id)
	if err != nil {
		return err
	}

	err = s.index.Delete(id.String())
	if err != nil {
		return err
	}

	event := NewEvent(index.ID, index.Type)
	if err := s.dispatcher.EmitEvent(ctx, RemoveEventChannel, event); err != nil {
		s.logger.Error("error emitting event", map[string]interface{}{
			"err": err,
		})
	}

	s.logger.Debug("Removed successful", map[string]interface{}{
		"id": id,
	})

	return nil
}
