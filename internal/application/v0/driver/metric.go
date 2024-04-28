package driver

import (
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
)

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

type (
	ListMetricsOpts types.FilterMetricsOpts

	ListMetricsResult struct {
		Metrics []*types.Metric
	}

	AggregateMetricsOpts types.FilterMetricsOpts

	AggregateMetricsResult struct {
		Aggregate *types.Aggregate
	}
)

func (d *Driver) ListMetrics(opts ListMetricsOpts) (*ListMetricsResult, error) {
	result, err := d.tracker.ListMetrics(d.ctx, &types.FilterMetricsOpts{
		Domain:    opts.Domain,
		Name:      opts.Name,
		StartTime: opts.StartTime,
		EndTime:   opts.EndTime,
	})
	if err != nil {
		return nil, err
	}

	return &ListMetricsResult{
		Metrics: result.Metrics,
	}, nil
}

func (d *Driver) AggregateMetrics(opts AggregateMetricsOpts) (*AggregateMetricsResult, error) {
	result, err := d.tracker.AggregateMetrics(d.ctx, &types.FilterMetricsOpts{
		Domain:    opts.Domain,
		Name:      opts.Name,
		StartTime: opts.StartTime,
		EndTime:   opts.EndTime,
	})
	if err != nil {
		return nil, err
	}

	return &AggregateMetricsResult{
		Aggregate: result.Aggregate,
	}, nil
}

func (d *Driver) addApplicationMetrics() error {
	if err := d.ctx.Err(); err != nil {
		return err
	}

	if err := d.tracker.CreateMetric(d.ctx, &types.CreateMetricOpts{
		Domain: Domain,
		Name:   StartupMetric,
		Value:  1,
	}); err != nil {
		return err
	}

	if err := d.addProjectMetrics(); err != nil {
		return err
	}

	if err := d.addPackageMetrics(); err != nil {
		return err
	}

	return nil
}

func (d *Driver) addProjectMetrics() error {
	if err := d.ctx.Err(); err != nil {
		return err
	}

	result, err := d.portfolio.ListProjects(d.ctx)
	if err != nil {
		return err
	}

	return d.tracker.CreateMetric(d.ctx, &types.CreateMetricOpts{
		Domain: Domain,
		Name:   ProjectCountMetric,
		Value:  len(result.Projects),
	})
}

func (d *Driver) addPackageMetrics() error {
	if err := d.ctx.Err(); err != nil {
		return err
	}

	result, err := d.catalog.ListPackages(d.ctx)
	if err != nil {
		return err
	}

	if err := d.tracker.CreateMetric(d.ctx, &types.CreateMetricOpts{
		Domain: Domain,
		Name:   PackageCountMetric,
		Value:  len(result.Packages),
	}); err != nil {
		return err
	}

	return nil
}
