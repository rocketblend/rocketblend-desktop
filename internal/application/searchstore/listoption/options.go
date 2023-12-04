package listoption

import (
	"strconv"

	"github.com/blevesearch/bleve/v2"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/indextype"
)

type (
	ListOptions struct {
		Query    string
		Type     indextype.IndexType
		Category string
		Resource string
		Ready    bool
		Size     int
		From     int
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

func WithCategory(category string) ListOption {
	return func(o *ListOptions) {
		o.Category = category
	}
}

func WithResource(resource string) ListOption {
	return func(o *ListOptions) {
		o.Resource = resource
	}
}

func WithReady(ready bool) ListOption {
	return func(o *ListOptions) {
		o.Ready = ready
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
	query := bleve.NewConjunctionQuery()

	if so.Type != indextype.Unknown {
		query.AddQuery(bleve.NewQueryStringQuery("type:" + strconv.Itoa(int(so.Type))))
	}

	if so.Category != "" {
		query.AddQuery(bleve.NewQueryStringQuery("category:" + so.Category))
	}

	if so.Resource != "" {
		resourceQuery := bleve.NewMatchPhraseQuery(so.Resource)
		resourceQuery.SetField("resources")
		query.AddQuery(resourceQuery)
	}

	if so.Ready {
		readyQuery := bleve.NewBoolFieldQuery(true)
		readyQuery.SetField("ready")
		query.AddQuery(readyQuery)
	}

	if so.Query != "" {
		fuzzy := bleve.NewFuzzyQuery(so.Query)
		fuzzy.SetFuzziness(2) // Levenshtein distance
		query.AddQuery(fuzzy)
	}

	return bleve.NewSearchRequestOptions(query, so.Size, so.From, false)
}
