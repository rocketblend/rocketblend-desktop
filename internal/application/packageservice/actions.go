package packageservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	pack "github.com/rocketblend/rocketblend-desktop/internal/application/package"
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

func (s *service) Install(ctx context.Context, id uuid.UUID) (err error) {
	item, err := s.get(ctx, id)
	if err != nil {
		return err
	}

	if item.State != pack.Available && item.State != pack.Cancelled {
		return fmt.Errorf("package not in state for installation (%s)", item.State)
	}

	rocketpacks, err := s.rocketblendPackageService.GetPackages(ctx, false, item.Reference)
	if err != nil {
		return err
	}

	// Don't use updateWithContext on these as we don't want to cancel state updates.
	if err := s.update(id, pack.Downloading); err != nil {
		return err
	}

	if _, err := s.rocketblendInstallationService.GetInstallations(ctx, rocketpacks, false); err != nil {
		newState := pack.Error
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			newState = pack.Cancelled
		}

		if uErr := s.update(id, newState); uErr != nil {
			return fmt.Errorf("error updating to %s state after GetInstallations error: %w", newState, uErr)
		}

		return fmt.Errorf("error in GetInstallations: %w", err)
	}

	if err := s.update(id, pack.Installed); err != nil {
		return fmt.Errorf("error updating to installed state: %w", err)
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

func (s *service) update(id uuid.UUID, state pack.PackageState) error {
	ctx := context.Background()

	item, err := s.get(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting item in update: %w", err)
	}

	item.State = state

	if err = s.insert(ctx, item); err != nil {
		return fmt.Errorf("error inserting item in update: %w", err)
	}

	return nil
}

func (s *service) updateWithContext(ctx context.Context, id uuid.UUID, state pack.PackageState) error {
	item, err := s.get(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting item in update: %w", err)
	}

	item.State = state

	if err = s.insert(ctx, item); err != nil {
		return fmt.Errorf("error inserting item in update: %w", err)
	}

	return nil
}
