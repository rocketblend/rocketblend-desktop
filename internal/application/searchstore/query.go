package searchstore

import (
	"encoding/json"
	"strconv"

	"github.com/blevesearch/bleve/v2/document"
	index "github.com/blevesearch/bleve_index_api"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/indextype"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/listoption"
)

func (s *store) List(opts ...listoption.ListOption) ([]*Index, error) {
	options := &listoption.ListOptions{
		Size: 500,
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

	var result Index
	doc.VisitFields(func(field index.Field) {
		switch field := field.(type) {
		case *document.TextField:
			value := string(field.Value())

			switch field.Name() {
			case "id":
				if uid, err := uuid.Parse(value); err == nil {
					result.ID = uid
				}
			case "type":
				if typeInt, err := strconv.Atoi(value); err == nil {
					result.Type = indextype.IndexType(typeInt)
				}
			case "path":
				result.Path = value
			case "name":
				result.Name = value
			case "category":
				result.Category = value
			case "ready":
				if readyBool, err := strconv.ParseBool(value); err == nil {
					result.Ready = readyBool
				}
			case "data":
				result.Data = value
			case "resources":
				var resources []string
				if err := json.Unmarshal([]byte(value), &resources); err == nil {
					result.Resources = resources
				}
			}
		}
	})

	return &result, nil
}
