package pack

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/util"
	"github.com/rocketblend/rocketblend/pkg/downloader"
	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/driver/rocketpack"
	"github.com/rocketblend/rocketblend/pkg/semver"
)

type (
	Addon struct {
		ID        uuid.UUID           `json:"id,omitempty"`
		Reference reference.Reference `json:"reference,omitempty"`
		Name      string              `json:"name,omitempty"`
		Path      string              `json:"path,omitempty"`
		Version   *semver.Version     `json:"version,omitempty"`
		Resource  string              `json:"resource,omitempty"`
		URI       *downloader.URI     `json:"uri,omitempty"`
		UpdatedAt time.Time           `json:"updatedAt,omitempty"`
	}
)

func NewAddon(path string, rocketpack *rocketpack.RocketPack) (*Addon, error) {
	if rocketpack == nil || !rocketpack.IsAddon() {
		return nil, fmt.Errorf("rocketpack is not an addon")
	}

	reference, err := pathToReference(path)
	if err != nil {
		return nil, err
	}

	name, err := reference.GetRepo()
	if err != nil {
		return nil, err
	}

	modTime, err := util.GetDirModTime(path)
	if err != nil {
		return nil, err
	}

	return &Addon{
		ID:        uuid.New(),
		Name:      name,
		Reference: reference,
		Path:      path,
		Version:   rocketpack.Addon.Version,
		Resource:  rocketpack.Addon.Source.Resource,
		URI:       rocketpack.Addon.Source.URI,
		UpdatedAt: modTime,
	}, nil
}
