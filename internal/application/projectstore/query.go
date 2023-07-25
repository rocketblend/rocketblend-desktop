package projectstore

import (
	"github.com/rocketblend/rocketblend-desktop/internal/application/project"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectstore/listoptions"
)

func (s *store) List(opts ...listoptions.ListOption) ([]*project.Project, error) {
	options := &listoptions.ListOptions{
		Size: 100,
	}

	for _, o := range opts {
		o(options)
	}

	result, err := s.index.Search(options.SearchRequest())
	if err != nil {
		return nil, err
	}

	s.logger.Debug("search result", map[string]interface{}{
		"total":    result.Total,
		"took":     result.Took,
		"maxScore": result.MaxScore,
	})

	var matchingProjects []*project.Project
	for _, hit := range result.Hits {
		project, err := s.get(hit.ID)
		if err != nil {
			return nil, err
		}

		matchingProjects = append(matchingProjects, project)
	}

	return matchingProjects, nil
}

func (s *store) Get(key string) (*project.Project, error) {
	return s.get(key)
}
