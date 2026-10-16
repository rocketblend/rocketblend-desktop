package application

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/enums"
	"github.com/rocketblend/rocketblend-desktop/internal/application/store/listoption"
	"github.com/rocketblend/rocketblend-desktop/internal/application/types"
	"github.com/rocketblend/rocketblend/pkg/downloader"
	rbtypes "github.com/rocketblend/rocketblend/pkg/types"
)

type (
	GetPackageOpts struct {
		ID uuid.UUID `json:"id"`
	}

	GetPackageResult struct {
		Package *types.Package `json:"package,omitempty"`
	}

	ListPackagesOpts struct {
		Query      string             `json:"query"`
		References []string           `json:"references,omitempty"`
		Type       enums.PackageType  `json:"type"`
		State      enums.PackageState `json:"state"`
	}

	ListPackagesResult struct {
		Packages []*types.Package `json:"packages,omitempty"`
	}

	InstallPackageOpts struct {
		ID                   uuid.UUID `json:"id"`
		AllowUntrustedDomain bool      `json:"allowUntrustedDomain"`
	}

	InstallPackageResult struct {
		OperationID          uuid.UUID `json:"operationID"`
		RequiresConfirmation bool      `json:"requiresConfirmation"`
		Domain               string    `json:"domain,omitempty"`
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

	d.logger.Debug("finding all packages", map[string]interface{}{
		"query": opts.Query,
		"type":  opts.Type,
		"state": opts.State,
		"refs":  opts.References,
	})

	response, err := d.catalog.ListPackages(ctx, []listoption.ListOption{
		listoption.WithQuery(opts.Query),
		listoption.WithReferences(opts.References...),
		listoption.WithCategory(string(opts.Type)),
		listoption.WithState(string(opts.State)),
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
	domain, err := d.checkPackageDomain(opts.ID, opts.AllowUntrustedDomain)
	if err != nil {
		d.logger.Error("package download domain not allowed", map[string]interface{}{
			"error": err.Error(),
			"id":    opts.ID,
		})
		return nil, err
	}

	if domain != "" {
		d.logger.Info("package download requires domain confirmation", map[string]interface{}{
			"domain": domain,
			"id":     opts.ID,
		})

		return &InstallPackageResult{
			RequiresConfirmation: true,
			Domain:               domain,
		}, nil
	}

	opid, err := d.operator.Create(d.ctx, func(ctx context.Context, opid uuid.UUID) (interface{}, error) {
		if err := d.catalog.AddPackageOperation(d.ctx, &types.AddPackageOperationOpts{
			ID:          opts.ID,
			OperationID: opid,
		}); err != nil {
			d.logger.Error("failed to append operation to package", map[string]interface{}{
				"error": err.Error(),
				"opid":  opid,
			})

			return nil, err
		}

		defer func() {
			if err := d.catalog.RemovePackageOperation(d.ctx, &types.RemovePackageOperationOpts{
				ID:          opts.ID,
				OperationID: opid,
			}); err != nil {
				d.logger.Error("failed to remove operation from package", map[string]interface{}{
					"error": err.Error(),
					"opid":  opid,
				})
			}
		}()

		if err := d.catalog.InstallPackage(ctx, &types.InstallPackageOpts{
			ID:                         opts.ID,
			BypassUntrustedDomainCheck: opts.AllowUntrustedDomain,
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

	return &InstallPackageResult{
		OperationID: opid,
	}, nil
}

func (d *Driver) checkPackageDomain(id uuid.UUID, allowUntrusted bool) (string, error) {
	result, err := d.catalog.GetPackage(d.ctx, &types.GetPackageOpts{
		ID: id,
	})
	if err != nil {
		return "", err
	}

	config, err := d.rbConfigurator.Get()
	if err != nil {
		return "", err
	}

	uri := result.Package.URI
	err = downloader.CheckDomain(uri, config.TrustedDomains, config.UntrustedDomainPolicy, allowUntrusted)
	if errors.Is(err, rbtypes.ErrUntrustedDomain) {
		return uri.Host, nil
	}

	return "", err
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

// autoRefreshPackages refreshes the packages, if auto pull is enabled.
func (d *Driver) autoRefreshPackages(ctx context.Context) error {
	config, err := d.configurator.Get()
	if err != nil {
		return err
	}

	if !config.Package.AutoPull {
		return nil
	}

	if err := d.catalog.RefreshPackages(ctx); err != nil {
		return err
	}

	return nil
}
