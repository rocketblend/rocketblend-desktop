package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/operationservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/util"
	"github.com/rocketblend/rocketblend/pkg/driver/reference"
)

type (
	CreateProjectOperationOpts struct {
		Name string `json:"name"`
	}
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

func (d *Driver) CreateProjectOperation(opts CreateProjectOperationOpts) (uuid.UUID, error) {
	operationservice, err := d.factory.GetOperationService()
	if err != nil {
		d.logger.Error("failed to get operation service", map[string]interface{}{"error": err.Error()})
		return uuid.Nil, err
	}

	projectService, err := d.factory.GetProjectService()
	if err != nil {
		d.logger.Error("failed to get project service", map[string]interface{}{"error": err.Error()})
		return uuid.Nil, err
	}

	projectPath, err := d.getProjectPath()
	if err != nil {
		return uuid.Nil, err
	}

	defaultBuild, err := d.getDefaultBuild()
	if err != nil {
		return uuid.Nil, err
	}

	opid, err := operationservice.Create(d.ctx, func(ctx context.Context, opid uuid.UUID) (interface{}, error) {
		result, err := projectService.Create(ctx, projectservice.CreateProjectOpts{
			DisplayName:   opts.Name,
			BlendFileName: util.DisplayNameToFilename(opts.Name),
			Path:          projectPath,
			Build:         defaultBuild,
		})
		if err != nil {
			d.logger.Error("failed to create project", map[string]interface{}{"error": err.Error(), "opid": opid})
			return nil, err
		}

		d.logger.Debug("project created", map[string]interface{}{"id": result.ID, "opid": opid})
		return result, nil
	})
	if err != nil {
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

func (d *Driver) getDefaultBuild() (reference.Reference, error) {
	_, rbConfig, err := d.getRocketBlendConfig()
	if err != nil {
		return "", err
	}

	return rbConfig.DefaultBuild, nil
}

func (d *Driver) getProjectPath() (string, error) {
	_, aConfig, err := d.getApplicationConfig()
	if err != nil {
		return "", err
	}

	if len(aConfig.Project.Paths) == 0 {
		return "", errors.New("no project path configured")
	}

	// TODO: support multiple project paths
	return aConfig.Project.Paths[0], nil
}
