package listoption

import (
	"strconv"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/indextype"
)

type (
	ListOptions struct {
		Query     string
		Type      indextype.IndexType
		Reference string
		Name      string
		Category  string
		Resource  string
		Operation string
		State     *int
		Size      int
		From      int
		StartTime time.Time
		EndTime   time.Time
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

func WithReference(reference string) ListOption {
	return func(o *ListOptions) {
		o.Reference = reference
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

func WithState(state *int) ListOption {
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

	if so.Reference != "" {
		referenceQuery := bleve.NewMatchPhraseQuery(so.Reference)
		referenceQuery.SetField("reference")
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

	if so.State != nil {
		stateQuery := bleve.NewQueryStringQuery("state:" + strconv.Itoa(*so.State))
		query.AddQuery(stateQuery)
	}

	if !so.StartTime.IsZero() && !so.EndTime.IsZero() {
		dateRangeQuery := bleve.NewDateRangeQuery(so.StartTime, so.EndTime)
		dateRangeQuery.SetField("date")
		query.AddQuery(dateRangeQuery)
	}

	if so.Query != "" {
		textQuery := bleve.NewMatchPhraseQuery(so.Query)
		// textQuery.Fuzziness = 1

		query.AddQuery(textQuery)
	} else {
		matchAllQuery := bleve.NewMatchAllQuery()
		query.AddQuery(matchAllQuery)
	}

	return bleve.NewSearchRequestOptions(query, so.Size, so.From, false)
}
