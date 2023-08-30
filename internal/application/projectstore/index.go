package projectstore

import (
	"encoding/json"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	index "github.com/blevesearch/bleve_index_api"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/project"
)

type (
	ProjectIndexMeta struct {
		ID   uuid.UUID `json:"id,omitempty"`
		Path string    `json:"path,omitempty"`
		Name string    `json:"name,omitempty"`
		Data string    `json:"data,omitempty"`
	}
)

func (v *ProjectIndexMeta) Type() string {
	return "project"
}

func (v *ProjectIndexMeta) BleveType() string {
	return "project"
}

func (s *store) updateIndex(id uuid.UUID, project *project.Project) error {
	data, err := json.Marshal(project)
	if err != nil {
		return err
	}

	return s.index.Index(id.String(), &ProjectIndexMeta{
		ID:   id,
		Path: project.Path,
		Name: project.Name,
		Data: string(data),
	})
}

func (s *store) get(id uuid.UUID) (*project.Project, error) {
	var project project.Project
	doc, err := s.index.Document(id.String())
	if err != nil {
		return nil, err
	}

	doc.VisitFields(func(field index.Field) {
		if field.Name() == "data" {
			err := json.Unmarshal(field.Value(), &project)
			if err != nil {
				s.logger.Error("failed to unmarshal the document field", map[string]interface{}{
					"error": err,
				})
			}
		}
	})

	return &project, nil
}

func newIndexMapping() mapping.IndexMapping {
	projectMapping := bleve.NewDocumentMapping()
	// projectMapping.Dynamic = false

	// source data store - this is where original doc will be stored
	dataTextFieldMapping := bleve.NewTextFieldMapping()
	dataTextFieldMapping.Store = true
	dataTextFieldMapping.Index = false
	dataTextFieldMapping.IncludeInAll = false
	dataTextFieldMapping.IncludeTermVectors = false
	projectMapping.AddFieldMappingsAt("data", dataTextFieldMapping)

	// create
	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("project", projectMapping)
	indexMapping.TypeField = "type"
	indexMapping.DefaultAnalyzer = "en"

	return indexMapping
}
