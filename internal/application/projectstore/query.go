package projectstore

import (
	"fmt"

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

	var matchingProjects []*project.Project
	for _, hit := range result.Hits {
		s.mu.RLock()
		proj, ok := s.projects[hit.ID]
		s.mu.RUnlock()

		if !ok {
			continue
		}

		matchingProjects = append(matchingProjects, proj.Copy())
	}

	return matchingProjects, nil
}

func (s *store) Get(key string) (*project.Project, error) {
	s.mu.RLock()
	proj, ok := s.projects[key]
	s.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("project not found")
	}

	return proj.Copy(), nil
}
