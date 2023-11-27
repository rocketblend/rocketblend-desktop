package searchstore

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/listoptions"
)

type (
	Store interface {
		List(opts ...listoptions.ListOption) ([]*Index, error)
		Get(id uuid.UUID) (*Index, error)
		Insert(index *Index) error
		Remove(id uuid.UUID) error
		RemoveByPath(path string) error
	}

	store struct {
		logger logger.Logger
		index  bleve.Index
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

	s := &store{
		logger: options.Logger,
		index:  index,
	}

	return s, nil
}

func (s *store) Insert(index *Index) error {
	err := s.index.Index(index.ID.String(), index)
	if err != nil {
		return err
	}

	s.logger.Debug("Indexed succusful", map[string]interface{}{
		"id":   index.ID,
		"type": index.Type,
	})

	return nil
}

func (s *store) Remove(id uuid.UUID) error {
	err := s.index.Delete(id.String())
	if err != nil {
		return err
	}

	s.logger.Debug("Removed succusful", map[string]interface{}{
		"id": id,
	})

	return nil
}

func (s *store) RemoveByPath(path string) error {
	query := bleve.NewPrefixQuery(path)
	search := bleve.NewSearchRequest(query)
	searchResults, err := s.index.Search(search)
	if err != nil {
		s.logger.Error("Error searching for projects in path", map[string]interface{}{
			"err": err,
		})

		return err
	}

	for _, hit := range searchResults.Hits {
		if err := s.index.Delete(hit.ID); err != nil {
			s.logger.Error("Error deleting project from index", map[string]interface{}{
				"err":  err,
				"key":  hit.ID,
				"path": path,
			})
		} else {
			s.logger.Info("Deleted project from index", map[string]interface{}{
				"key":  hit.ID,
				"path": path,
			})
		}
	}

	return nil
}
