package listoption

import (
	"strconv"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/rocketblend/rocketblend-desktop/internal/application/store/indextype"
)

// TODO: This package should be moved/changed.
type (
	ListOptions struct {
		Query      string
		Type       indextype.IndexType
		References []string
		Name       string
		Category   string
		Resource   string
		Operation  string
		State      string
		Size       int
		From       int
		StartTime  time.Time
		EndTime    time.Time
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

func WithReferences(references ...string) ListOption {
	return func(o *ListOptions) {
		o.References = references
	}
}

func WithName(name string) ListOption {
	return func(o *ListOptions) {
		o.Name = name
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

func WithOperation(operation string) ListOption {
	return func(o *ListOptions) {
		o.Operation = operation
	}
}

func WithState(state string) ListOption {
	return func(o *ListOptions) {
		o.State = state
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

func WithDateRange(startTime, endTime time.Time) ListOption {
	return func(o *ListOptions) {
		o.StartTime = startTime
		o.EndTime = endTime
	}
}

func (so *ListOptions) SearchRequest() *bleve.SearchRequest {
	query := bleve.NewConjunctionQuery()

	if so.Type != indextype.Unknown {
		typeQuery := bleve.NewQueryStringQuery("type:" + strconv.Itoa(int(so.Type)))
		query.AddQuery(typeQuery)
	}

	if len(so.References) > 0 {
		referenceQuery := bleve.NewDisjunctionQuery()
		for _, reference := range so.References {
			phraseQuery := bleve.NewMatchPhraseQuery(reference)
			phraseQuery.SetField("reference")
			referenceQuery.AddQuery(phraseQuery)
		}

		query.AddQuery(referenceQuery)
	}

	if so.Name != "" {
		nameQuery := bleve.NewMatchPhraseQuery(so.Name)
		nameQuery.SetField("name")
		query.AddQuery(nameQuery)
	}

	if so.Category != "" {
		categoryQuery := bleve.NewMatchPhraseQuery(so.Category)
		categoryQuery.SetField("category")
		query.AddQuery(categoryQuery)
	}

	if so.Resource != "" {
		resourceQuery := bleve.NewMatchPhraseQuery(so.Resource)
		resourceQuery.SetField("resources")
		query.AddQuery(resourceQuery)
	}

	if so.Operation != "" {
		operationQuery := bleve.NewMatchPhraseQuery(so.Operation)
		operationQuery.SetField("operations")
		query.AddQuery(operationQuery)
	}

	if so.State != "" {
		stateQuery := bleve.NewMatchPhraseQuery(so.State)
		stateQuery.SetField("state")
		query.AddQuery(stateQuery)
	}

	if !so.StartTime.IsZero() && !so.EndTime.IsZero() {
		dateRangeQuery := bleve.NewDateRangeQuery(so.StartTime, so.EndTime)
		dateRangeQuery.SetField("date")
		query.AddQuery(dateRangeQuery)
	}

	if so.Query != "" {
		searchQuery := bleve.NewBooleanQuery()
		fuzzy_query := bleve.NewFuzzyQuery(so.Query)
		fuzzy_query.SetFuzziness(2)
		searchQuery.AddShould(fuzzy_query)
		searchQuery.AddShould(bleve.NewQueryStringQuery("*" + so.Query + "*"))
		query.AddQuery(searchQuery)
	} else {
		matchAllQuery := bleve.NewMatchAllQuery()
		query.AddQuery(matchAllQuery)
	}

	return bleve.NewSearchRequestOptions(query, so.Size, so.From, false)
}
