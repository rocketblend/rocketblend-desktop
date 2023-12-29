package operationservice

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/indextype"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/listoption"
)

type (
	Service interface {
		Start(ctx context.Context, opFunc func(ctx context.Context) (interface{}, error)) (uuid.UUID, error)
		Get(ctx context.Context, opid uuid.UUID) (*Operation, error)
		List(ctx context.Context, opts ...listoption.ListOption) ([]*Operation, error)
		Cancel(opid uuid.UUID) error
	}

	operation struct {
		ctx    context.Context
		cancel context.CancelFunc
	}

	service struct {
		logger logger.Logger
		store  searchstore.Store

		operations    map[uuid.UUID]*operation
		operationsMux sync.RWMutex
	}

	Options struct {
		Logger logger.Logger
		Store  searchstore.Store
	}

	Option func(*Options)
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithStore(store searchstore.Store) Option {
	return func(o *Options) {
		o.Store = store
	}
}

func New(opts ...Option) (Service, error) {
	options := &Options{
		Logger: logger.NoOp(),
	}

	for _, o := range opts {
		o(options)
	}

	if options.Store == nil {
		return nil, fmt.Errorf("store is required")
	}

	return &service{
		logger:     options.Logger,
		operations: make(map[uuid.UUID]*operation),
		store:      options.Store,
	}, nil
}

func (s *service) Start(ctx context.Context, opFunc func(ctx context.Context) (interface{}, error)) (uuid.UUID, error) {
	opid := uuid.New()
	opctx, cancel := context.WithCancel(ctx)

	s.operationsMux.Lock()
	s.operations[opid] = &operation{ctx: opctx, cancel: cancel}
	s.operationsMux.Unlock()

	operation := Operation{ID: opid}
	index, err := operation.ToSearchIndex()
	if err != nil {
		return uuid.Nil, err
	}

	if err := s.store.Insert(index); err != nil {
		return uuid.Nil, err
	}

	go func() {
		defer cancel()
		result, err := opFunc(opctx)

		if opctx.Err() == context.Canceled {
			return
		}

		operation := Operation{
			ID:        opid,
			Completed: true,
			Result:    result,
		}

		if err != nil {
			operation.ErrorMsg = err.Error()
		}

		index, err := operation.ToSearchIndex()
		if err != nil {
			s.logger.Error("failed to marshal Operation", map[string]interface{}{"error": err.Error()})
			return
		}

		s.store.Insert(index)
	}()

	return opid, nil
}

func (s *service) Get(ctx context.Context, opid uuid.UUID) (*Operation, error) {
	index, err := s.store.Get(opid)
	if err != nil {
		return nil, err
	}

	operation, err := convertIndexToOperation(index)
	if err != nil {
		return nil, err
	}

	return operation, nil
}

func (s *service) List(ctx context.Context, opts ...listoption.ListOption) ([]*Operation, error) {
	opts = append(opts, listoption.WithType(indextype.Operation))
	indexes, err := s.store.List(opts...)
	if err != nil {
		return nil, err
	}

	operations := make([]*Operation, 0, len(indexes))
	for _, index := range indexes {
		op, err := convertIndexToOperation(index)
		if err != nil {
			return nil, err
		}

		operations = append(operations, op)
	}

	s.logger.Debug("Found operations", map[string]interface{}{
		"operations": len(operations),
		"indexes":    len(indexes),
	})

	return operations, nil
}

func (s *service) Cancel(opid uuid.UUID) error {
	s.operationsMux.RLock()
	op, exists := s.operations[opid]
	s.operationsMux.RUnlock()

	if !exists {
		return fmt.Errorf("operation not found")
	}

	op.cancel()

	operation := Operation{
		ID:        opid,
		ErrorMsg:  fmt.Errorf("operation cancelled").Error(),
		Completed: true,
	}
	index, err := operation.ToSearchIndex()
	if err != nil {
		return err
	}

	s.store.Insert(index)
	return nil
}

func convertIndexToOperation(index *searchstore.Index) (*Operation, error) {
	status := &Operation{}
	if err := json.Unmarshal([]byte(index.Data), status); err != nil {
		return nil, err
	}

	return status, nil
}
