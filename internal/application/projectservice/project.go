package projectservice

import (
	"time"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/project"
	"github.com/rocketblend/rocketblend/pkg/driver/reference"
)

type (
	Project struct {
		ID        uuid.UUID             `json:"id,omitempty"`
		Name      string                `json:"name,omitempty"`
		Tags      []string              `json:"tags,omitempty"`
		Path      string                `json:"path,omitempty"`
		FileName  string                `json:"fileName,omitempty"`
		Build     reference.Reference   `json:"build,omitempty"`
		Addons    []reference.Reference `json:"addons,omitempty"`
		Version   string                `json:"version,omitempty"`
		ARGS      string                `json:"args,omitempty"`
		UpdatedAt time.Time             `json:"updatedAt,omitempty"`
	}
)

func mapProjects(projects ...*project.Project) []*Project {
	result := make([]*Project, 0, len(projects))

	for _, p := range projects {
		newProject := &Project{
			ID:        p.Settings.ID,
			Name:      p.Settings.Name,
			Tags:      p.Settings.Tags,
			Path:      p.BlendFile.ProjectPath,
			FileName:  p.BlendFile.BlendFileName,
			Build:     p.BlendFile.RocketFile.GetBuild(),
			Addons:    p.BlendFile.RocketFile.GetAddons(),
			Version:   p.BlendFile.RocketFile.GetVersion(),
			ARGS:      p.BlendFile.RocketFile.GetArgs(),
			UpdatedAt: p.UpdatedAt,
		}

		result = append(result, newProject)
	}

	return result
}
