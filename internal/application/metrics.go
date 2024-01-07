package application

import "github.com/rocketblend/rocketblend-desktop/internal/application/metricservice"

const (
	Domain = "development"

	StartupMetric = "application.startup"

	PackageCountMetric = "package.count"

	ProjectCountMetric   = "project.count"
	ProjectRunMetric     = "project.action.run"
	ProjectExploreMetric = "project.action.explore"
	ProjectUpdateMetric  = "project.action.update"

	StoreInsertMetric = "store.insert"
)

func (d *Driver) addApplicationMetrics() error {
	if err := d.ctx.Err(); err != nil {
		return err
	}

	metric, err := d.factory.GetMetricService()
	if err != nil {
		return err
	}

	if err := metric.Add(d.ctx, metricservice.AddOptions{
		Domain: Domain,
		Name:   StartupMetric,
		Value:  1,
	}); err != nil {
		return err
	}

	if err := d.addProjectMetrics(metric); err != nil {
		return err
	}

	if err := d.addPackageMetrics(metric); err != nil {
		return err
	}

	return nil
}

func (d *Driver) addProjectMetrics(metric metricservice.Service) error {
	if err := d.ctx.Err(); err != nil {
		return err
	}

	project, err := d.factory.GetProjectService()
	if err != nil {
		return err
	}

	// TODO: Add count function to project service
	result, err := project.List(d.ctx)
	if err != nil {
		return err
	}

	return metric.Add(d.ctx, metricservice.AddOptions{
		Domain: Domain,
		Name:   ProjectCountMetric,
		Value:  len(result.Projects),
	})
}

func (d *Driver) addPackageMetrics(metric metricservice.Service) error {
	if err := d.ctx.Err(); err != nil {
		return err
	}

	pack, err := d.factory.GetPackageService()
	if err != nil {
		return err
	}

	// TODO: Add count function to package service
	result, err := pack.List(d.ctx)
	if err != nil {
		return err
	}

	return metric.Add(d.ctx, metricservice.AddOptions{
		Domain: Domain,
		Name:   PackageCountMetric,
		Value:  len(result.Packages),
	})
}
