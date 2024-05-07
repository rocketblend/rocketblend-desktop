package types

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/store/indextype"
	"github.com/rocketblend/rocketblend-desktop/internal/application/store/listoption"
)

type (
	Index struct {
		ID         uuid.UUID           `json:"id,omitempty"`
		Type       indextype.IndexType `json:"type,omitempty"`
		Reference  string              `json:"reference,omitempty"`
		Name       string              `json:"name,omitempty"`
		Category   string              `json:"category,omitempty"`
		State      string              `json:"state"`
		Resources  []string            `json:"resources,omitempty"`
		Operations []string            `json:"operations,omitempty"`
		Date       time.Time           `json:"date,omitempty"`
		Data       string              `json:"data,omitempty"`
	}

	Store interface {
		List(ctx context.Context, opts ...listoption.ListOption) ([]*Index, error)
		Get(ctx context.Context, id uuid.UUID) (*Index, error)
		Insert(ctx context.Context, index *Index) error
		Remove(ctx context.Context, id uuid.UUID) error
		RemoveByReference(ctx context.Context, path string) error

		Close() error
	}
)

func (v *Index) BleveType() string {
	return "index"
}
