package packageservice

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend/pkg/driver/reference"
)

func (s *service) Add(ctx context.Context, reference reference.Reference) error {
	return fmt.Errorf("not implemented")
}

func (s *service) Install(ctx context.Context, id uuid.UUID) error {
	return fmt.Errorf("not implemented")
}

func (s *service) Uninstall(ctx context.Context, id uuid.UUID) error {
	return fmt.Errorf("not implemented")
}

func (s *service) Refresh(ctx context.Context) error {
	return fmt.Errorf("not implemented")
}
