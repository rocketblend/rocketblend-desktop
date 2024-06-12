package store

import (
	"context"
	"errors"

	"github.com/blevesearch/bleve/v2"
	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/events"
	"github.com/rocketblend/rocketblend-desktop/internal/application/store/listoption"
	"github.com/rocketblend/rocketblend-desktop/internal/application/types"
)

type (
	Store struct {
		logger types.Logger
		index  bleve.Index

		dispatcher types.Dispatcher
	}

	Options struct {
		Logger     types.Logger
		Dispatcher types.Dispatcher
	}

	Option func(*Options)
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithDispatcher(dispatcher types.Dispatcher) Option {
	return func(o *Options) {
		o.Dispatcher = dispatcher
	}
}

func New(opts ...Option) (*Store, error) {
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

	return &Store{
		logger:     options.Logger,
		index:      index,
		dispatcher: options.Dispatcher,
	}, nil
}

func (s *Store) Insert(ctx context.Context, index *types.Index) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	existing, _ := s.index.Document(index.ID.String())

	err := s.index.Index(index.ID.String(), index)
	if err != nil {
		return err
	}

	if existing != nil {
		event := newEvent(index.ID, index.Type)
		if err := s.dispatcher.EmitEvent(ctx, events.StoreInsertChannel, event); err != nil {
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

func (s *Store) Remove(ctx context.Context, id uuid.UUID) error {
	return s.remove(ctx, id)
}

// TODO: remove this and just filter then call normal remove.
func (s *Store) RemoveByReference(ctx context.Context, path string) error {
	listOpts := listoption.ListOptions{
		Reference: path,
		Size:      10000, // TODO: Have the search request function ignore the size if it is 0
	}

	searchResults, err := s.index.SearchInContext(ctx, listOpts.SearchRequest())
	if err != nil {
		s.logger.Error("error searching for indexes with reference", map[string]interface{}{
			"err": err,
		})

		return err
	}

	s.logger.Debug("searched indexes with reference", map[string]interface{}{
		"path":  path,
		"total": searchResults.Total,
		"took":  searchResults.Took,
		"hits":  len(searchResults.Hits),
	})

	if searchResults.Total == 0 {
		return ErrNotFound
	}

	for _, hit := range searchResults.Hits {
		s.logger.Debug("found index for deletion", map[string]interface{}{
			"key":  hit.ID,
			"path": path,
		})

		id, err := uuid.Parse(hit.ID)
		if err != nil {
			s.logger.Error("error parsing id", map[string]interface{}{
				"err":  err,
				"key":  hit.ID,
				"path": path,
			})

			continue
		}

		if err := s.remove(ctx, id); err != nil {
			s.logger.Error("error deleting index from id", map[string]interface{}{
				"err":  err,
				"key":  hit.ID,
				"path": path,
			})
		}
	}

	return nil
}

func (s *Store) Close() error {
	if s.index == nil {
		return nil
	}

	return s.index.Close()
}

func (s *Store) remove(ctx context.Context, id uuid.UUID) error {
	index, err := s.get(ctx, id)
	if err != nil {
		return err
	}

	s.logger.Debug("removing index", map[string]interface{}{
		"id":        id,
		"reference": index.Reference,
		"type":      index.Type,
		"resources": index.Resources,
	})

	err = s.index.Delete(id.String())
	if err != nil {
		return err
	}

	event := newEvent(index.ID, index.Type)
	if err := s.dispatcher.EmitEvent(ctx, events.StoreRemoveChannel, event); err != nil {
		s.logger.Error("error emitting event", map[string]interface{}{
			"err": err,
		})
	}

	s.logger.Debug("removed successful", map[string]interface{}{
		"id": id,
	})

	return nil
}
