package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/operationservice"
)

func (d *Driver) GetOperation(opid uuid.UUID) (*operationservice.Operation, error) {
	operationService, err := d.factory.GetOperationService()
	if err != nil {
		d.logger.Error("failed to get operation service", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	operation, err := operationService.Get(d.ctx, opid)
	if err != nil {
		d.logger.Error("failed to get operation", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	return operation, nil
}

func (d *Driver) ListOperations() ([]*operationservice.Operation, error) {
	operationService, err := d.factory.GetOperationService()
	if err != nil {
		d.logger.Error("failed to get operation service", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	operations, err := operationService.List(d.ctx)
	if err != nil {
		d.logger.Error("failed to list operations", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	return operations, nil
}

func (d *Driver) CancelOperation(opid uuid.UUID) error {
	operationService, err := d.factory.GetOperationService()
	if err != nil {
		d.logger.Error("failed to get operation service", map[string]interface{}{"error": err.Error()})
		return err
	}

	if err := operationService.Cancel(opid); err != nil {
		d.logger.Error("failed to cancel operation", map[string]interface{}{"error": err.Error()})
		return err
	}

	return nil
}

func (d *Driver) InstallPackageOperation(id uuid.UUID) (uuid.UUID, error) {
	operationservice, err := d.factory.GetOperationService()
	if err != nil {
		d.logger.Error("failed to get operation service", map[string]interface{}{"error": err.Error()})
		return uuid.Nil, err
	}

	packageService, err := d.factory.GetPackageService()
	if err != nil {
		d.logger.Error("failed to get package service", map[string]interface{}{"error": err.Error()})
		return uuid.Nil, err
	}

	opid, err := operationservice.Create(d.ctx, func(ctx context.Context, opid uuid.UUID) (interface{}, error) {
		if err := packageService.Install(ctx, id); err != nil {
			d.logger.Error("failed to install package", map[string]interface{}{"error": err.Error(), "opid": opid})
			return nil, err
		}

		d.logger.Debug("package installed", map[string]interface{}{"id": id, "opid": opid})
		return nil, nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if err := packageService.AppendOperation(d.ctx, id, opid); err != nil {
		d.logger.Error("failed to append operation to package", map[string]interface{}{"error": err.Error(), "opid": opid})
		err := operationservice.Cancel(opid)
		if err != nil {
			d.logger.Error("failed to cancel operation", map[string]interface{}{"error": err.Error(), "opid": opid})
		}

		return uuid.Nil, err
	}

	return opid, nil
}

func (d *Driver) LongRunningOperation() (uuid.UUID, error) {
	operationservice, err := d.factory.GetOperationService()
	if err != nil {
		d.logger.Error("failed to get operation service", map[string]interface{}{"error": err.Error()})
		return uuid.Nil, err
	}

	opid, err := operationservice.Create(d.ctx, func(ctx context.Context, opid uuid.UUID) (interface{}, error) {
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
