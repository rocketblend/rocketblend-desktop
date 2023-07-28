package listoptions

import "github.com/blevesearch/bleve/v2"

type (
	ListOptions struct {
		Query string
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
	// query := bleve.NewConjunctionQuery()

	// // Build the query based on the provided parameters
	// if so.Name != "" {
	// 	query.AddQuery(bleve.NewQueryStringQuery("name:" + so.Name))
	// }

	if so.Query != "" {
		query := bleve.NewFuzzyQuery(so.Query)
		query.SetFuzziness(2) // Levenshtein distance
		search := bleve.NewSearchRequestOptions(query, so.Size, so.From, false)
		return search
	}

	query := bleve.NewMatchAllQuery()
	return bleve.NewSearchRequestOptions(query, so.Size, so.From, false)
}
