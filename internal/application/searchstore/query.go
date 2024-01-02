package searchstore

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/blevesearch/bleve/v2/document"
	index "github.com/blevesearch/bleve_index_api"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/indextype"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/listoption"
)

func (s *store) List(ctx context.Context, opts ...listoption.ListOption) ([]*Index, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	options := &listoption.ListOptions{
		Size: 50,
	}

	for _, o := range opts {
		o(options)
	}

	result, err := s.index.SearchInContext(ctx, options.SearchRequest())
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
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		id, err := uuid.Parse(hit.ID)
		if err != nil {
			return nil, err
		}

		index, err := s.get(ctx, id)
		if err != nil {
			return nil, err
		}

		indexes = append(indexes, index)
	}

	return indexes, nil
}

func (s *store) Get(ctx context.Context, id uuid.UUID) (*Index, error) {
	return s.get(ctx, id)
}

func (s *store) get(ctx context.Context, id uuid.UUID) (*Index, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

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
			case "reference":
				result.Reference = value
			case "name":
				result.Name = value
			case "category":
				result.Category = value
			case "state":
				if state, err := strconv.Atoi(value); err == nil {
					result.State = state
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
