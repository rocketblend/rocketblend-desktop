package metricservice

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/indextype"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/listoption"
)

type (
	Metric struct {
		ID         uuid.UUID `json:"id"`
		Domain     string    `json:"domain"`
		Name       string    `json:"name"`
		Value      int       `json:"value"`
		RecordedAt time.Time `json:"recordedAt"`
	}

	AddOptions struct {
		Domain string `json:"domain"`
		Name   string `json:"Name"`
		Value  int    `json:"value"`
	}

	FilterOptions struct {
		Domain    string    `json:"domain"`
		Name      string    `json:"name"`
		StartTime time.Time `json:"startTime"`
		EndTime   time.Time `json:"endTime"`
	}

	Aggregate struct {
		Domain string  `json:"domain"`
		Name   string  `json:"name"`
		Sum    int     `json:"sum"`
		Avg    float64 `json:"avg"`
		Count  int     `json:"count"`
		Min    int     `json:"min"`
		Max    int     `json:"max"`
	}

	Service interface {
		Get(ctx context.Context, id uuid.UUID) (*Metric, error)
		List(ctx context.Context, options FilterOptions) ([]*Metric, error)
		Aggregate(ctx context.Context, options FilterOptions) (*Aggregate, error)

		Add(ctx context.Context, options AddOptions) error
		Remove(ctx context.Context, id uuid.UUID) error
	}

	service struct {
		logger logger.Logger
		store  searchstore.Store
	}

	Options struct {
		Logger logger.Logger

		Store searchstore.Store
	}

	Option func(*Options)
)

func WithLogger(logger logger.Logger) Option {
	return func(opts *Options) {
		opts.Logger = logger
	}
}

func WithStore(store searchstore.Store) Option {
	return func(opts *Options) {
		opts.Store = store
	}
}

func New(options ...Option) (Service, error) {
	opts := &Options{
		Logger: logger.NoOp(),
	}

	for _, option := range options {
		option(opts)
	}

	if opts.Store == nil {
		return nil, errors.New("store is required")
	}

	return &service{
		logger: opts.Logger,
		store:  opts.Store,
	}, nil
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*Metric, error) {
	index, err := s.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	metric, err := convertIndexToMetric(index)
	if err != nil {
		return nil, err
	}

	return metric, nil
}

func (s *service) List(ctx context.Context, options FilterOptions) ([]*Metric, error) {
	return s.list(ctx, options)
}

func (s *service) Add(ctx context.Context, options AddOptions) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if options.Domain == "" {
		return errors.New("domain is required")
	}

	if options.Name == "" {
		return errors.New("name is required")
	}

	metric := &Metric{
		ID:         uuid.New(),
		Domain:     options.Domain,
		Name:       options.Name,
		Value:      options.Value,
		RecordedAt: time.Now(),
	}

	index, err := convertMetricToIndex(metric)
	if err != nil {
		return err
	}

	if err := s.store.Insert(ctx, index); err != nil {
		return err
	}

	return nil
}

func (s *service) Aggregate(ctx context.Context, options FilterOptions) (*Aggregate, error) {
	metrics, err := s.list(ctx, options)
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

	return &Aggregate{
		Domain: options.Domain,
		Name:   options.Name,
		Sum:    sum,
		Avg:    avg,
		Count:  count,
		Min:    min,
		Max:    max,
	}, nil
}

func (s *service) Remove(ctx context.Context, id uuid.UUID) error {
	return s.store.Remove(ctx, id)
}

func (s *service) list(ctx context.Context, options FilterOptions) ([]*Metric, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	opts := []listoption.ListOption{
		listoption.WithType(indextype.Metric),
		listoption.WithReference(options.Domain),
		listoption.WithName(options.Name),
		listoption.WithDateRange(options.StartTime, options.EndTime),
	}

	indexes, err := s.store.List(ctx, opts...)
	if err != nil {
		return nil, err
	}

	metrics := make([]*Metric, 0, len(indexes))
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

func convertMetricToIndex(metric *Metric) (*searchstore.Index, error) {
	data, err := json.Marshal(metric)
	if err != nil {
		return nil, err
	}

	return &searchstore.Index{
		ID:        metric.ID,
		Name:      metric.Name,
		Reference: metric.Domain,
		Type:      indextype.Metric,
		Date:      metric.RecordedAt,
		Data:      string(data),
	}, nil
}

func convertIndexToMetric(index *searchstore.Index) (*Metric, error) {
	metric := &Metric{}
	if err := json.Unmarshal([]byte(index.Data), metric); err != nil {
		return nil, err
	}

	return metric, nil
}