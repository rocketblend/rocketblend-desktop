package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/types"
)

type (
	Operation struct {
		logger types.Logger

		operator types.Operator

		ctx context.Context
	}

	GetOperationOpts struct {
		ID uuid.UUID `json:"id"`
	}

	GetOperationResult struct {
		Operation *types.Operation `json:"operation"`
	}

	ListOperationsOpts struct {
	}

	ListOperationsResult struct {
		Operations []*types.Operation `json:"operations"`
	}

	CancelOperationOpts struct {
		ID uuid.UUID `json:"id"`
	}
)

func newOperation(ctx context.Context, deps *dependencies) *Operation {
	return &Operation{
		logger:   deps.logger,
		operator: deps.operator,
		ctx:      ctx,
	}
}

func (s *Operation) GetOperation(opts GetOperationOpts) (*GetOperationResult, error) {
	operation, err := s.operator.Get(s.ctx, opts.ID)
	if err != nil {
		s.logger.Error("failed to get operation", map[string]interface{}{
			"error": err.Error(),
			"id":    opts.ID,
		})
		return nil, err
	}

	return &GetOperationResult{
		Operation: operation,
	}, nil
}

func (s *Operation) ListOperations() (*ListOperationsResult, error) {
	operations, err := s.operator.List(s.ctx)
	if err != nil {
		s.logger.Error("failed to list operations", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	return &ListOperationsResult{
		Operations: operations,
	}, nil
}

func (s *Operation) CancelOperation(opts CancelOperationOpts) error {
	if err := s.operator.Cancel(opts.ID); err != nil {
		s.logger.Error("failed to cancel operation", map[string]interface{}{
			"error": err.Error(),
			"id":    opts.ID,
		})
		return err
	}

	return nil
}

func (s *Operation) LongRunningOperation() (uuid.UUID, error) {
	opid, err := s.operator.Create(s.ctx, func(ctx context.Context, opid uuid.UUID) (interface{}, error) {
		// Simulate a long-running operation
		for i := 0; i < 10; i++ {
			select {
			case <-ctx.Done():
				s.logger.Debug("long running operation canceled", map[string]interface{}{"opid": opid})
				return nil, ctx.Err()
			default:
				time.Sleep(2 * time.Second)
			}
		}

		return nil, nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	return opid, nil
}
