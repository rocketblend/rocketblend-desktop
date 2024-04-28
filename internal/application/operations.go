package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
)

type (
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

func (d *Driver) GetOperation(opts GetOperationOpts) (*GetOperationResult, error) {
	operation, err := d.operator.Get(d.ctx, opts.ID)
	if err != nil {
		d.logger.Error("failed to get operation", map[string]interface{}{
			"error": err.Error(),
			"id":    opts.ID,
		})
		return nil, err
	}

	return &GetOperationResult{
		Operation: operation,
	}, nil
}

func (d *Driver) ListOperations(opts ListOperationsOpts) (*ListOperationsResult, error) {
	operations, err := d.operator.List(d.ctx)
	if err != nil {
		d.logger.Error("failed to list operations", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	return &ListOperationsResult{
		Operations: operations,
	}, nil
}

func (d *Driver) CancelOperation(opts CancelOperationOpts) error {
	if err := d.operator.Cancel(opts.ID); err != nil {
		d.logger.Error("failed to cancel operation", map[string]interface{}{
			"error": err.Error(),
			"id":    opts.ID,
		})
		return err
	}

	return nil
}

func (d *Driver) LongRunningOperation() (uuid.UUID, error) {
	opid, err := d.operator.Create(d.ctx, func(ctx context.Context, opid uuid.UUID) (interface{}, error) {
		// Simulate a long-running operation
		for i := 0; i < 10; i++ {
			select {
			case <-ctx.Done():
				d.logger.Debug("long running operation canceled", map[string]interface{}{"opid": opid})
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
