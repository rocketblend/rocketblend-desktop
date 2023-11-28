package listoption

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/indextype"
)

type (
	ListOptions struct {
		Query string
		Type  indextype.IndexType
		Size  int
		From  int
	}

	ListOption func(*ListOptions)
)

func WithQuery(query string) ListOption {
	return func(o *ListOptions) {
		o.Query = query
	}
}

func WithType(indexType indextype.IndexType) ListOption {
	return func(o *ListOptions) {
		o.Type = indexType
	}
}

func WithSize(size int) ListOption {
	return func(o *ListOptions) {
		o.Size = size
	}
}

func WithFrom(from int) ListOption {
	return func(o *ListOptions) {
		o.From = from
	}
}

func (so *ListOptions) SearchRequest() *bleve.SearchRequest {
	// // Create a new query builder
	query := bleve.NewConjunctionQuery()

	// // Build the query based on the provided parameters
	if so.Type != indextype.Unknown {
		query.AddQuery(bleve.NewQueryStringQuery("type:" + so.Type.String()))
	}

	if so.Query != "" {
		fuzzy := bleve.NewFuzzyQuery(so.Query)
		fuzzy.SetFuzziness(2) // Levenshtein distance
		query.AddQuery(fuzzy)
	} else {
		query.AddQuery(bleve.NewMatchAllQuery())
	}

	return bleve.NewSearchRequestOptions(query, so.Size, so.From, false)
}
