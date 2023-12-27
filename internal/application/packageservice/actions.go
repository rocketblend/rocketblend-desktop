package packageservice

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend/pkg/driver/reference"
)

func (s *service) Add(ctx context.Context, reference reference.Reference) error {
	rocketpack, err := s.rocketblendPackageService.GetPackages(ctx, true, reference)
	if err != nil {
		return err
	}

	if len(rocketpack) == 0 {
		return fmt.Errorf("package not found")
	}

	return nil
}

func (s *service) Install(ctx context.Context, id uuid.UUID) error {
	pack, err := s.get(ctx, id)
	if err != nil {
		return err
	}

	rocketpacks, err := s.rocketblendPackageService.GetPackages(ctx, false, pack.Reference)
	if err != nil {
		return err
	}

	installs, err := s.rocketblendInstallationService.GetInstallations(ctx, rocketpacks, false)
	if err != nil {
		return err
	}

	_, ok := installs[pack.Reference]
	if !ok {
		return fmt.Errorf("installation not found")
	}

	return nil
}

func (s *service) Uninstall(ctx context.Context, id uuid.UUID) error {
	pack, err := s.get(ctx, id)
	if err != nil {
		return err
	}

	rocketpacks, err := s.rocketblendPackageService.GetPackages(ctx, false, pack.Reference)
	if err != nil {
		return err
	}

	if err := s.rocketblendInstallationService.RemoveInstallations(ctx, rocketpacks); err != nil {
		return err
	}

	return nil
}

func (s *service) Refresh(ctx context.Context) error {
	return fmt.Errorf("not implemented")
}
