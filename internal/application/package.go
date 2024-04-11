package application

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	pack "github.com/rocketblend/rocketblend-desktop/internal/application/package"
	"github.com/rocketblend/rocketblend-desktop/internal/application/packageservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/listoption"
	"github.com/rocketblend/rocketblend/pkg/driver/reference"
)

func (d *Driver) GetPackage(id uuid.UUID) (*packageservice.GetPackageResponse, error) {
	ctx := context.Background()

	packageService, err := d.factory.GetPackageService()
	if err != nil {
		d.logger.Error("failed to get package service", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	pack, err := packageService.Get(ctx, id)
	if err != nil {
		d.logger.Error("failed to find package by id", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	return pack, err
}

func (d *Driver) ListPackages(query string, packageType pack.PackageType, installed bool) (*packageservice.ListPackagesResponse, error) {
	ctx := context.Background()

	packageService, err := d.factory.GetPackageService()
	if err != nil {
		d.logger.Error("failed to get package service", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	category := ""
	if packageType != pack.Unknown {
		category = strconv.Itoa(int(packageType))
	}

	var state *int = nil
	if installed {
		stateInt := int(pack.Installed)
		state = &stateInt
	}

	response, err := packageService.List(ctx, []listoption.ListOption{
		listoption.WithQuery(query),
		listoption.WithCategory(category),
		listoption.WithState(state),
	}...)
	if err != nil {
		d.logger.Error("failed to find all packages", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	d.logger.Debug("found packages", map[string]interface{}{"packages": len(response.Packages)})

	return response, err
}

func (d *Driver) AddPackage(referenceStr string) error {
	packageService, err := d.factory.GetPackageService()
	if err != nil {
		d.logger.Error("failed to get package service", map[string]interface{}{"error": err.Error()})
		return err
	}

	ref, err := reference.Parse(referenceStr)
	if err != nil {
		d.logger.Error("failed to parse reference", map[string]interface{}{"error": err.Error()})
		return err
	}

	if err := packageService.Add(d.ctx, ref); err != nil {
		d.logger.Error("failed to add package", map[string]interface{}{"error": err.Error()})
		return err
	}

	return nil
}

func (d *Driver) RefreshPackages() error {
	packageService, err := d.factory.GetPackageService()
	if err != nil {
		d.logger.Error("failed to get package service", map[string]interface{}{"error": err.Error()})
		return err
	}

	if err := packageService.Refresh(d.ctx); err != nil {
		d.logger.Error("failed to refresh packages", map[string]interface{}{"error": err.Error()})
		return err
	}

	return nil
}

func (d *Driver) UninstallPackage(id uuid.UUID) error {
	packageService, err := d.factory.GetPackageService()
	if err != nil {
		d.logger.Error("failed to get package service", map[string]interface{}{"error": err.Error()})
		return err
	}

	if err := packageService.Uninstall(d.ctx, id); err != nil {
		d.logger.Error("failed to uninstall package", map[string]interface{}{"error": err.Error()})
		return err
	}

	d.logger.Debug("package uninstalled", map[string]interface{}{"id": id})
	return nil
}
