package projectservice

import (
	"time"

	"github.com/rocketblend/rocketblend-desktop/internal/application/project"
	"github.com/rocketblend/rocketblend/pkg/driver/reference"
)

type (
	Project struct {
		ID        string                `json:"id,omitempty"`
		Name      string                `json:"name,omitempty"`
		Tags      []string              `json:"tags,omitempty"`
		Path      string                `json:"path,omitempty"`
		FileName  string                `json:"fileName,omitempty"`
		Build     reference.Reference   `json:"build,omitempty"`
		Addons    []reference.Reference `json:"addons,omitempty"`
		Version   string                `json:"version,omitempty"`
		UpdatedAt time.Time             `json:"updatedAt,omitempty"`
	}

	GetProjectResponse struct {
		Project *Project `json:"project,omitempty"`
	}

	ListProjectsResponse struct {
		Projects []*Project `json:"projects,omitempty"`
	}

	CreateProjectRequest struct {
		Name     string                `json:"name,omitempty"`
		Tags     []string              `json:"tags,omitempty"`
		Path     string                `json:"path,omitempty"`
		FileName string                `json:"fileName,omitempty"`
		Build    reference.Reference   `json:"build,omitempty"`
		Addons   []reference.Reference `json:"addons,omitempty"`
	}

	UpdateProjectRequest struct {
		ID       string                `json:"id,omitempty"`
		Name     string                `json:"name,omitempty"`
		Tags     []string              `json:"tags,omitempty"`
		Path     string                `json:"path,omitempty"`
		FileName string                `json:"fileName,omitempty"`
		Build    reference.Reference   `json:"build,omitempty"`
		Addons   []reference.Reference `json:"addons,omitempty"`
	}
)

func mapProjects(projects ...*project.Project) []*Project {
	result := make([]*Project, 0, len(projects))

	for _, p := range projects {
		newProject := &Project{
			ID:        p.Settings.ID.String(),
			Name:      p.Settings.Name,
			Tags:      p.Settings.Tags,
			Path:      p.BlendFile.ProjectPath,
			FileName:  p.BlendFile.BlendFileName,
			Build:     p.BlendFile.RocketFile.GetBuild(),
			Addons:    p.BlendFile.RocketFile.GetAddons(),
			Version:   p.BlendFile.RocketFile.GetVersion(),
			UpdatedAt: p.UpdatedAt,
		}

		result = append(result, newProject)
	}

	return result
}
