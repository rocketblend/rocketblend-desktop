package tracker

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/store/indextype"
	"github.com/rocketblend/rocketblend-desktop/internal/application/store/listoption"
	"github.com/rocketblend/rocketblend-desktop/internal/application/types"
)

type (
	Tracker struct {
		logger types.Logger
		store  types.Store
	}

	Options struct {
		Logger types.Logger

		Store types.Store
	}

	Option func(*Options)
)

func WithLogger(logger types.Logger) Option {
	return func(opts *Options) {
		opts.Logger = logger
	}
}

func WithStore(store types.Store) Option {
	return func(opts *Options) {
		opts.Store = store
	}
}

func New(options ...Option) (*Tracker, error) {
	opts := &Options{
		Logger: logger.NoOp(),
	}

	for _, option := range options {
		option(opts)
	}

	if opts.Store == nil {
		return nil, errors.New("store is required")
	}

	return &Tracker{
		logger: opts.Logger,
		store:  opts.Store,
	}, nil
}

func (t *Tracker) GetMetric(ctx context.Context, opts *types.GetMetricOpts) (*types.GetMetricResult, error) {
	index, err := t.store.Get(ctx, opts.ID)
	if err != nil {
		return nil, err
	}

	metric, err := convertIndexToMetric(index)
	if err != nil {
		return nil, err
	}

	return &types.GetMetricResult{
		Metric: metric,
	}, nil
}

func (t *Tracker) ListMetrics(ctx context.Context, opts *types.FilterMetricsOpts) (*types.ListMetricsResult, error) {
	metrics, err := t.list(ctx, opts.Domain, opts.Name, opts.StartTime, opts.EndTime)
	if err != nil {
		return nil, err
	}

	return &types.ListMetricsResult{
		Metrics: metrics,
	}, nil
}

func (t *Tracker) CreateMetric(ctx context.Context, opts *types.CreateMetricOpts) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if opts.Domain == "" {
		return errors.New("domain is required")
	}

	if opts.Name == "" {
		return errors.New("name is required")
	}

	metric := &types.Metric{
		ID:         uuid.New(),
		Domain:     opts.Domain,
		Name:       opts.Name,
		Value:      opts.Value,
		RecordedAt: time.Now(),
	}

	index, err := convertMetricToIndex(metric)
	if err != nil {
		return err
	}

	if err := t.store.Insert(ctx, index); err != nil {
		return err
	}

	return nil
}

func (t *Tracker) AggregateMetrics(ctx context.Context, opts *types.FilterMetricsOpts) (*types.AggregateMetricsResult, error) {
	metrics, err := t.list(ctx, opts.Domain, opts.Name, opts.StartTime, opts.EndTime)
	if err != nil {
		return nil, err
	}

	var sum, count, min, max int
	var avg float64

	for _, metric := range metrics {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		sum += metric.Value
		count++
		if min == 0 || metric.Value < min {
			min = metric.Value
		}
		if metric.Value > max {
			max = metric.Value
		}
	}

	if count > 0 {
		avg = float64(sum) / float64(count)
	}

	return &types.AggregateMetricsResult{
		Aggregate: &types.Aggregate{
			Domain: opts.Domain,
			Name:   opts.Name,
			Sum:    sum,
			Avg:    avg,
			Count:  count,
			Min:    min,
			Max:    max,
		},
	}, nil
}

func (t *Tracker) DeleteMetric(ctx context.Context, opts *types.DeleteMetricOpts) error {
	return t.store.Remove(ctx, opts.ID)
}

func (t *Tracker) list(ctx context.Context, domain string, name string, start time.Time, end time.Time) ([]*types.Metric, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	opts := []listoption.ListOption{
		listoption.WithType(indextype.Metric),
		listoption.WithReference(domain),
		listoption.WithName(name),
		listoption.WithDateRange(start, end),
		listoption.WithSize(10000),
	}

	indexes, err := t.store.List(ctx, opts...)
	if err != nil {
		return nil, err
	}

	metrics := make([]*types.Metric, 0, len(indexes))
	for _, index := range indexes {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		met, err := convertIndexToMetric(index)
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, met)
	}

	return metrics, nil
}

func convertMetricToIndex(metric *types.Metric) (*types.Index, error) {
	data, err := json.Marshal(metric)
	if err != nil {
		return nil, err
	}

	return &types.Index{
		ID:        metric.ID,
		Name:      metric.Name,
		Reference: metric.Domain,
		Type:      indextype.Metric,
		Date:      metric.RecordedAt,
		Data:      string(data),
	}, nil
}

func convertIndexToMetric(index *types.Index) (*types.Metric, error) {
	metric := &types.Metric{}
	if err := json.Unmarshal([]byte(index.Data), metric); err != nil {
		return nil, err
	}

	return metric, nil
}
