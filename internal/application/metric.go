package application

import (
	"github.com/rocketblend/rocketblend-desktop/internal/application/types"
)

const (
	Domain = "development"

	StartupMetric = "application.startup"

	PackageCountMetric = "package.count"

	ProjectCountMetric = "project.count"
	ProjectRunMetric   = "project.action.run"
	// ProjectExploreMetric = "project.action.explore"
	ProjectUpdateMetric = "project.action.update"

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

func (s *System) ListMetrics(opts ListMetricsOpts) (*ListMetricsResult, error) {
	result, err := s.tracker.ListMetrics(s.ctx, &types.FilterMetricsOpts{
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

func (s *System) AggregateMetrics(opts AggregateMetricsOpts) (*AggregateMetricsResult, error) {
	result, err := s.tracker.AggregateMetrics(s.ctx, &types.FilterMetricsOpts{
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

func (l *lifecycle) addApplicationMetrics() error {
	if err := l.ctx.Err(); err != nil {
		return err
	}

	if err := l.tracker.CreateMetric(l.ctx, &types.CreateMetricOpts{
		Domain: Domain,
		Name:   StartupMetric,
		Value:  1,
	}); err != nil {
		return err
	}

	if err := l.addProjectMetrics(); err != nil {
		return err
	}

	if err := l.addPackageMetrics(); err != nil {
		return err
	}

	return nil
}

// TODO: Create a track function in portfolio and catalog and call it from startup instead.
func (l *lifecycle) addProjectMetrics() error {
	if err := l.ctx.Err(); err != nil {
		return err
	}

	result, err := l.portfolio.ListProjects(l.ctx)
	if err != nil {
		return err
	}

	return l.tracker.CreateMetric(l.ctx, &types.CreateMetricOpts{
		Domain: Domain,
		Name:   ProjectCountMetric,
		Value:  len(result.Projects),
	})
}

func (l *lifecycle) addPackageMetrics() error {
	if err := l.ctx.Err(); err != nil {
		return err
	}

	result, err := l.catalog.ListPackages(l.ctx)
	if err != nil {
		return err
	}

	if err := l.tracker.CreateMetric(l.ctx, &types.CreateMetricOpts{
		Domain: Domain,
		Name:   PackageCountMetric,
		Value:  len(result.Packages),
	}); err != nil {
		return err
	}

	return nil
}
