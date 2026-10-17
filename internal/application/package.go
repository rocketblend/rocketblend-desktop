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
	Package struct {
		logger types.Logger

		operator types.Operator
		catalog  types.Catalog

		rbConfigurator rbtypes.Configurator

		ctx context.Context
	}

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

func newPackage(ctx context.Context, deps *dependencies) *Package {
	return &Package{
		logger:         deps.logger,
		operator:       deps.operator,
		catalog:        deps.catalog,
		rbConfigurator: deps.rbConfigurator,
		ctx:            ctx,
	}
}

func (s *Package) GetPackage(opts GetPackageOpts) (*GetPackageResult, error) {
	result, err := s.catalog.GetPackage(s.ctx, &types.GetPackageOpts{
		ID: opts.ID,
	})
	if err != nil {
		s.logger.Error("failed to find package by id", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	return &GetPackageResult{
		Package: result.Package,
	}, nil
}

func (s *Package) ListPackages(opts ListPackagesOpts) (*ListPackagesResult, error) {
	s.logger.Debug("finding all packages", map[string]interface{}{
		"query": opts.Query,
		"type":  opts.Type,
		"state": opts.State,
		"refs":  opts.References,
	})

	response, err := s.catalog.ListPackages(s.ctx, []listoption.ListOption{
		listoption.WithQuery(opts.Query),
		listoption.WithReferences(opts.References...),
		listoption.WithCategory(string(opts.Type)),
		listoption.WithState(string(opts.State)),
	}...)
	if err != nil {
		s.logger.Error("failed to find all packages", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	s.logger.Debug("found packages", map[string]interface{}{"packages": len(response.Packages)})

	return &ListPackagesResult{
		Packages: response.Packages,
	}, nil
}

func (s *Package) AddPackage(opts AddPackageOpts) error {
	// ref, err := reference.Parse(opts.Reference)
	// if err != nil {
	// 	s.logger.Error("failed to parse reference", map[string]interface{}{"error": err.Error()})
	// 	return err
	// }

	// if err := s.catalog.AddPackage(s.ctx, ref); err != nil {
	// 	s.logger.Error("failed to add package", map[string]interface{}{"error": err.Error()})
	// 	return err
	// }

	return nil
}

func (s *Package) InstallPackage(opts InstallPackageOpts) (*InstallPackageResult, error) {
	domain, err := s.checkPackageDomain(opts.ID, opts.AllowUntrustedDomain)
	if err != nil {
		s.logger.Error("package download domain not allowed", map[string]interface{}{
			"error": err.Error(),
			"id":    opts.ID,
		})
		return nil, err
	}

	if domain != "" {
		s.logger.Info("package download requires domain confirmation", map[string]interface{}{
			"domain": domain,
			"id":     opts.ID,
		})

		return &InstallPackageResult{
			RequiresConfirmation: true,
			Domain:               domain,
		}, nil
	}

	opid, err := s.operator.Create(s.ctx, func(ctx context.Context, opid uuid.UUID) (interface{}, error) {
		if err := s.catalog.AddPackageOperation(s.ctx, &types.AddPackageOperationOpts{
			ID:          opts.ID,
			OperationID: opid,
		}); err != nil {
			s.logger.Error("failed to append operation to package", map[string]interface{}{
				"error": err.Error(),
				"opid":  opid,
			})

			return nil, err
		}

		defer func() {
			if err := s.catalog.RemovePackageOperation(s.ctx, &types.RemovePackageOperationOpts{
				ID:          opts.ID,
				OperationID: opid,
			}); err != nil {
				s.logger.Error("failed to remove operation from package", map[string]interface{}{
					"error": err.Error(),
					"opid":  opid,
				})
			}
		}()

		if err := s.catalog.InstallPackage(ctx, &types.InstallPackageOpts{
			ID:                         opts.ID,
			BypassUntrustedDomainCheck: opts.AllowUntrustedDomain,
		}); err != nil {
			s.logger.Error("failed to install package", map[string]interface{}{
				"error": err.Error(),
				"opid":  opid,
			})
			return nil, err
		}

		s.logger.Debug("package installed", map[string]interface{}{
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

func (s *Package) checkPackageDomain(id uuid.UUID, allowUntrusted bool) (string, error) {
	result, err := s.catalog.GetPackage(s.ctx, &types.GetPackageOpts{
		ID: id,
	})
	if err != nil {
		return "", err
	}

	config, err := s.rbConfigurator.Get()
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

func (s *Package) RefreshPackages() error {
	if err := s.catalog.RefreshPackages(s.ctx); err != nil {
		s.logger.Error("failed to refresh packages", map[string]interface{}{"error": err.Error()})
		return err
	}

	return nil
}

func (s *Package) UninstallPackage(opts UninstallPackageOpts) error {
	if err := s.catalog.UninstallPackage(s.ctx, &types.UninstallPackageOpts{
		ID: opts.ID,
	}); err != nil {
		s.logger.Error("failed to uninstall package", map[string]interface{}{"error": err.Error()})
		return err
	}

	s.logger.Debug("package uninstalled", map[string]interface{}{"id": opts.ID})
	return nil
}
