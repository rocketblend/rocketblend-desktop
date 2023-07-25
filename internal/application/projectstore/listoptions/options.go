package listoptions

import "github.com/blevesearch/bleve/v2"

type (
	ListOptions struct {
		Name string
		Size int
		From int
	}

	ListOption func(*ListOptions)
)

func WithName(name string) ListOption {
	return func(o *ListOptions) {
		o.Name = name
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
	// Create a new query builder
	query := bleve.NewConjunctionQuery()

	// Build the query based on the provided parameters
	if so.Name != "" {
		query.AddQuery(bleve.NewQueryStringQuery("name:" + so.Name))
	}

	return bleve.NewSearchRequestOptions(query, so.Size, so.From, false)
}
