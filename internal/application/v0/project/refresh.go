package project

import "context"

func (r *repository) RefreshProjects(ctx context.Context) error {
	if err := r.refresh(ctx); err != nil {
		return err
	}

	return nil
}

func (r *repository) refresh(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	config, err := r.configurator.Get()
	if err != nil {
		return err
	}

	if err := r.watcher.SetPaths(config.Project.Paths...); err != nil {
		return err
	}

	return nil
}
