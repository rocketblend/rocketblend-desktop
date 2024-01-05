package statisticservice

import (
	"context"
	"errors"
	"time"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore"
)

type (
	Statistic struct {
		ID         uuid.UUID   `json:"id"`
		Domain     string      `json:"domain"`
		Metric     string      `json:"metric"`
		Period     Period      `json:"period"`
		DataType   DataType    `json:"dataType"`
		Value      interface{} `json:"value"`
		RecordedAt time.Time   `json:"recordedAt"`
	}

	Service interface {
		Get(ctx context.Context, id uuid.UUID) (*Statistic, error)
		List(ctx context.Context) ([]*Statistic, error)

		Record(ctx context.Context, request *RecordStatisticRequest) error
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

func (s *service) Get(ctx context.Context, id uuid.UUID) (*Statistic, error) {
	return nil, errors.New("not implemented")
}

func (s *service) List(ctx context.Context) ([]*Statistic, error) {
	return nil, errors.New("not implemented")
}

func (s *service) Record(ctx context.Context, request *RecordStatisticRequest) error {
	return errors.New("not implemented")
}
