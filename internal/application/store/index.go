package store

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
)

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
