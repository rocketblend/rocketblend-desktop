package searchstore

import (
	"encoding/json"

	"github.com/blevesearch/bleve/v2/document"
	index "github.com/blevesearch/bleve_index_api"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/listoptions"
)

func (s *store) List(opts ...listoptions.ListOption) ([]*Index, error) {
	options := &listoptions.ListOptions{
		Size: 100,
	}

	for _, o := range opts {
		o(options)
	}

	result, err := s.index.Search(options.SearchRequest())
	if err != nil {
		return nil, err
	}

	s.logger.Debug("search result", map[string]interface{}{
		"total":    result.Total,
		"took":     result.Took,
		"maxScore": result.MaxScore,
	})

	var indexes []*Index
	for _, hit := range result.Hits {
		id, err := uuid.Parse(hit.ID)
		if err != nil {
			return nil, err
		}

		s.logger.Debug("search hit", map[string]interface{}{
			"id":     id,
			"score":  hit.Score,
			"fields": hit.Fields,
		})

		index, err := s.get(id)
		if err != nil {
			return nil, err
		}

		indexes = append(indexes, index)
	}

	return indexes, nil
}

func (s *store) Get(id uuid.UUID) (*Index, error) {
	return s.get(id)
}

func (s *store) get(id uuid.UUID) (*Index, error) {
	doc, err := s.index.Document(id.String())
	if err != nil {
		return nil, err
	}

	// Convert the entire document to a map
	docMap := make(map[string]interface{})
	doc.VisitFields(func(field index.Field) {
		switch field := field.(type) {
		case *document.TextField:
			docMap[field.Name()] = string(field.Value())
		}
	})

	// Marshal the map into JSON
	docJson, err := json.Marshal(docMap)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON into an Index struct
	var result Index
	if err := json.Unmarshal(docJson, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
