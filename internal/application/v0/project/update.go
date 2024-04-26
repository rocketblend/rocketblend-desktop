package project

import (
	"context"

	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
)

func (r *repository) UpdateProject(ctx context.Context, opts *types.UpdateProjectOpts) error {
	project, err := r.get(ctx, opts.ID)
	if err != nil {
		return err
	}

	// settings := project.Settings()
	// if opts.Name != nil {
	// 	settings.Name = *opts.Name
	// }

	// if opts.Tags != nil {
	// 	settings.Tags = *opts.Tags
	// }

	// if opts.ThumbnailPath != nil {
	// 	settings.ThumbnailPath = *opts.ThumbnailPath
	// }

	// if opts.SplashPath != nil {
	// 	settings.SplashPath = *opts.SplashPath
	// }

	// filePath := filepath.Join(project.Path, projectsettings.FileName)
	// if err := projectsettings.Save(settings, filePath); err != nil {
	// 	return err
	// }

	r.emitEvent(ctx, project.ID, UpdateEventChannel)

	return nil
}
