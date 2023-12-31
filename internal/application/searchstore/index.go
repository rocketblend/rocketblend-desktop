package searchstore

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/indextype"
)

type (
	Index struct {
		ID         uuid.UUID           `json:"id,omitempty"`
		Type       indextype.IndexType `json:"type,omitempty"`
		Reference  string              `json:"reference,omitempty"`
		Name       string              `json:"name,omitempty"`
		Category   string              `json:"category,omitempty"`
		State      int                 `json:"state"`
		Resources  []string            `json:"resources,omitempty"`
		Operations []string            `json:"operations,omitempty"`
		Error      string              `json:"error,omitempty"`
		Data       string              `json:"data,omitempty"`
	}
)

func (v *Index) BleveType() string {
	return "index"
}

func newIndexMapping() mapping.IndexMapping {
	mapping := bleve.NewDocumentMapping()
	// mapping.Dynamic = false

	// source data store - this is where original doc will be stored
	dataTextFieldMapping := bleve.NewTextFieldMapping()
	dataTextFieldMapping.Store = true
	dataTextFieldMapping.Index = false
	dataTextFieldMapping.IncludeInAll = false
	dataTextFieldMapping.IncludeTermVectors = false
	mapping.AddFieldMappingsAt("data", dataTextFieldMapping)

	// create
	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("index", mapping)
	indexMapping.TypeField = "type"
	indexMapping.DefaultAnalyzer = "en"

	return indexMapping
}
