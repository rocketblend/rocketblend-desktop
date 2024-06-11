package pack

import (
	"context"

	"github.com/rocketblend/rocketblend/pkg/reference"
	rbtypes "github.com/rocketblend/rocketblend/pkg/types"
)

// RefreshPackages refreshes the packages.
// Currently, it only refreshes the default build repo.
func (r *Repository) RefreshPackages(ctx context.Context) error {
	config, err := r.rbConfigurator.Get()
	if err != nil {
		return err
	}

	if _, err := r.rbRepository.GetPackages(ctx, &rbtypes.GetPackagesOpts{
		References: []reference.Reference{config.DefaultBuild},
		Update:     true,
	}); err != nil {
		return err
	}

	return nil
}
