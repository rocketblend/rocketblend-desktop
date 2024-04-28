package operator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/store/indextype"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/store/listoption"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
)

type (
	operation struct {
		ctx    context.Context
		cancel context.CancelFunc
	}

	Operator struct {
		logger types.Logger

		store      types.Store
		dispatcher types.Dispatcher

		operations    map[uuid.UUID]*operation
		operationsMux sync.RWMutex
	}

	Options struct {
		Logger types.Logger

		Store      types.Store
		Dispatcher types.Dispatcher
	}

	Option func(*Options)
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithStore(store types.Store) Option {
	return func(o *Options) {
		o.Store = store
	}
}

func WithDispatcher(dispatcher types.Dispatcher) Option {
	return func(o *Options) {
		o.Dispatcher = dispatcher
	}
}

func New(opts ...Option) (*Operator, error) {
	options := &Options{
		Logger: logger.NoOp(),
	}

	for _, o := range opts {
		o(options)
	}

	if options.Store == nil {
		return nil, fmt.Errorf("store is required")
	}

	if options.Dispatcher == nil {
		return nil, fmt.Errorf("dispatcher is required")
	}

	return &Operator{
		logger:     options.Logger,
		operations: make(map[uuid.UUID]*operation),
		store:      options.Store,
		dispatcher: options.Dispatcher,
	}, nil
}

func (o *Operator) Create(ctx context.Context, opFunc func(ctx context.Context, opid uuid.UUID) (interface{}, error)) (uuid.UUID, error) {
	opid := uuid.New()
	opctx, cancel := context.WithCancel(ctx)

	o.operationsMux.Lock()
	o.operations[opid] = &operation{ctx: opctx, cancel: cancel}
	o.operationsMux.Unlock()

	operation := types.Operation{ID: opid}
	index, err := convertToSearchIndex(operation)
	if err != nil {
		return uuid.Nil, err
	}

	if err := o.store.Insert(ctx, index); err != nil {
		return uuid.Nil, err
	}

	o.logger.Info("starting operation", map[string]interface{}{"id": opid})

	go func() {
		defer cancel()
		result, err := opFunc(opctx, opid)

		operation := types.Operation{
			ID:        opid,
			Completed: true,
			Result:    result,
		}

		if err != nil {
			operation.ErrorMsg = err.Error()
			o.logger.Error("operation failed", map[string]interface{}{"error": err.Error()})
		}

		// TODO: Handle failure better. We still want to update the operation state.
		index, err := convertToSearchIndex(operation)
		if err != nil {
			o.logger.Error("failed to marshal Operation", map[string]interface{}{"error": err.Error()})
			return
		}

		if err := o.store.Insert(context.Background(), index); err != nil {
			o.logger.Error("failed to insert Operation", map[string]interface{}{"error": err.Error()})
			return
		}

		o.logger.Info("operation ended", map[string]interface{}{"id": opid})

	}()

	return opid, nil
}

func (o *Operator) Get(ctx context.Context, opid uuid.UUID) (*types.Operation, error) {
	index, err := o.store.Get(ctx, opid)
	if err != nil {
		return nil, err
	}

	operation, err := convertIndexToOperation(index)
	if err != nil {
		return nil, err
	}

	return operation, nil
}

func (o *Operator) List(ctx context.Context, opts ...listoption.ListOption) ([]*types.Operation, error) {
	opts = append(opts, listoption.WithType(indextype.Operation))
	indexes, err := o.store.List(ctx, opts...)
	if err != nil {
		return nil, err
	}

	operations := make([]*types.Operation, 0, len(indexes))
	for _, index := range indexes {
		op, err := convertIndexToOperation(index)
		if err != nil {
			return nil, err
		}

		operations = append(operations, op)
	}

	// s.logger.Debug("found operations", map[string]interface{}{
	// 	"operations": len(operations),
	// 	"indexes":    len(indexes),
	// })

	return operations, nil
}

func (o *Operator) Cancel(opid uuid.UUID) error {
	o.operationsMux.RLock()
	op, exists := o.operations[opid]
	o.operationsMux.RUnlock()

	if !exists {
		return errors.New("operation does not exist")
	}

	op.cancel()
	o.logger.Info("cancelled operation", map[string]interface{}{"id": opid})

	operation := types.Operation{
		ID:        opid,
		ErrorMsg:  errors.New("operation cancelled").Error(),
		Completed: true,
	}

	index, err := convertToSearchIndex(operation)
	if err != nil {
		return err
	}

	o.store.Insert(context.Background(), index)
	return nil
}

func convertIndexToOperation(index *types.Index) (*types.Operation, error) {
	status := &types.Operation{}
	if err := json.Unmarshal([]byte(index.Data), status); err != nil {
		return nil, err
	}

	return status, nil
}

func convertToSearchIndex(operation types.Operation) (*types.Index, error) {
	data, err := json.Marshal(operation)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OperationStatus: %w", err)
	}

	state := 0
	if operation.Completed {
		state = 1
	}

	return &types.Index{
		ID:    operation.ID,
		Type:  indextype.Operation,
		State: state,
		Data:  string(data),
	}, nil
}
