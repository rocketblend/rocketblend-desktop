package searchstore

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/google/uuid"
)

type (
	Index struct {
		ID   uuid.UUID `json:"id,omitempty"`
		Type IndexType `json:"type,omitempty"`
		Path string    `json:"path,omitempty"`
		Name string    `json:"name,omitempty"`
		Data string    `json:"data,omitempty"`
	}
)

func (v *Index) BleveType() string {
	return "index"
}

func newIndexMapping() mapping.IndexMapping {
	mapping := bleve.NewDocumentMapping()
	// projectMapping.Dynamic = false

	// source data store - this is where original doc will be stored
	dataTextFieldMapping := bleve.NewTextFieldMapping()
	dataTextFieldMapping.Store = true
	dataTextFieldMapping.Index = false
	dataTextFieldMapping.IncludeInAll = false
	dataTextFieldMapping.IncludeTermVectors = false
	mapping.AddFieldMappingsAt("data", dataTextFieldMapping)

	// create
	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("item", mapping)
	indexMapping.TypeField = "type"
	indexMapping.DefaultAnalyzer = "en"

	return indexMapping
}
