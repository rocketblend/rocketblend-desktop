package driver

import (
	"context"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/store/listoption"
	"github.com/rocketblend/rocketblend-desktop/internal/application/v0/types"
)

type (
	GetPackageOpts struct {
		ID uuid.UUID `json:"id"`
	}

	GetPackageResult struct {
		Package *types.Package `json:"package,omitempty"`
	}

	ListPackagesOpts struct {
		Query       string `json:"query"`
		PackageType string `json:"packageType"`
		Installed   bool   `json:"installed"`
	}

	ListPackagesResult struct {
		Packages []*types.Package `json:"packages,omitempty"`
	}

	InstallPackageOpts struct {
		ID uuid.UUID `json:"id"`
	}

	InstallPackageResult struct {
		OperationID uuid.UUID `json:"operationID"`
	}

	AddPackageOpts struct {
		Reference string `json:"reference"`
	}

	UninstallPackageOpts struct {
		ID uuid.UUID `json:"id"`
	}
)

func (d *Driver) GetPackage(opts GetPackageOpts) (*GetPackageResult, error) {
	ctx := context.Background()

	result, err := d.catalog.GetPackage(ctx, &types.GetPackageOpts{
		ID: opts.ID,
	})
	if err != nil {
		d.logger.Error("failed to find package by id", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	return &GetPackageResult{
		Package: result.Package,
	}, nil
}

func (d *Driver) ListPackages(opts ListPackagesOpts) (*ListPackagesResult, error) {
	ctx := context.Background()

	category := ""
	if opts.PackageType != "" {
		category = opts.PackageType
	}

	// var state *int = nil
	// if opts.Installed {
	// 	stateInt := int(pack.Installed)
	// 	state = &stateInt
	// }

	response, err := d.catalog.ListPackages(ctx, []listoption.ListOption{
		listoption.WithQuery(opts.Query),
		listoption.WithCategory(category),
		//listoption.WithState(state),
	}...)
	if err != nil {
		d.logger.Error("failed to find all packages", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	d.logger.Debug("found packages", map[string]interface{}{"packages": len(response.Packages)})

	return &ListPackagesResult{
		Packages: response.Packages,
	}, nil
}

func (d *Driver) AddPackage(opts AddPackageOpts) error {
	// ref, err := reference.Parse(opts.Reference)
	// if err != nil {
	// 	d.logger.Error("failed to parse reference", map[string]interface{}{"error": err.Error()})
	// 	return err
	// }

	// if err := d.catalog.AddPackage(d.ctx, ref); err != nil {
	// 	d.logger.Error("failed to add package", map[string]interface{}{"error": err.Error()})
	// 	return err
	// }

	return nil
}

func (d *Driver) InstallPackage(opts InstallPackageOpts) (*InstallPackageResult, error) {
	opid, err := d.operator.Create(d.ctx, func(ctx context.Context, opid uuid.UUID) (interface{}, error) {
		if err := d.catalog.InstallPackage(ctx, &types.InstallPackageOpts{
			ID: opts.ID,
		}); err != nil {
			d.logger.Error("failed to install package", map[string]interface{}{
				"error": err.Error(),
				"opid":  opid,
			})
			return nil, err
		}

		d.logger.Debug("package installed", map[string]interface{}{
			"id":   opts.ID,
			"opid": opid,
		})

		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	if err := d.catalog.AddPackageOperation(d.ctx, &types.AddPackageOperationOpts{
		ID:          opts.ID,
		OperationID: opid,
	}); err != nil {
		d.logger.Error("failed to append operation to package", map[string]interface{}{
			"error": err.Error(),
			"opid":  opid,
		})

		if err := d.operator.Cancel(opid); err != nil {
			d.logger.Error("failed to cancel operation", map[string]interface{}{
				"error": err.Error(),
				"opid":  opid,
			})
		}

		return nil, err
	}

	return &InstallPackageResult{
		OperationID: opid,
	}, nil
}

func (d *Driver) RefreshPackages() error {
	if err := d.catalog.RefreshPackages(d.ctx); err != nil {
		d.logger.Error("failed to refresh packages", map[string]interface{}{"error": err.Error()})
		return err
	}

	return nil
}

func (d *Driver) UninstallPackage(opts UninstallPackageOpts) error {
	if err := d.catalog.UninstallPackage(d.ctx, &types.UninstallPackageOpts{
		ID: opts.ID,
	}); err != nil {
		d.logger.Error("failed to uninstall package", map[string]interface{}{"error": err.Error()})
		return err
	}

	d.logger.Debug("package uninstalled", map[string]interface{}{"id": opts.ID})
	return nil
}
