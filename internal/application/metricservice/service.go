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

	AggregateOptions struct {
		Domain    string
		Name      string
		StartTime time.Time
		EndTime   time.Time
		MinCount  int
	}

	Aggregate struct {
		Domain string
		Name   string
		Sum    int
		Avg    float64
		Count  int
		Min    int
		Max    int
	}

	Service interface {
		Get(ctx context.Context, id uuid.UUID) (*Metric, error)
		List(ctx context.Context, opts ...listoption.ListOption) ([]*Metric, error)

		Add(ctx context.Context, options AddOptions) error
		Aggregate(ctx context.Context, options AggregateOptions) (*Aggregate, error)
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

func New(options ...Option) Service {
	opts := &Options{
		Logger: logger.NoOp(),
	}

	for _, option := range options {
		option(opts)
	}

	return &service{
		logger: opts.Logger,
		store:  opts.Store,
	}
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*Metric, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

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

func (s *service) List(ctx context.Context, opts ...listoption.ListOption) ([]*Metric, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	opts = append(opts, listoption.WithType(indextype.Metric))
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

func (s *service) Aggregate(ctx context.Context, options AggregateOptions) (*Aggregate, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return nil, errors.New("not implemented")
}

func (s *service) Remove(ctx context.Context, id uuid.UUID) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	return s.store.Remove(ctx, id)
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
