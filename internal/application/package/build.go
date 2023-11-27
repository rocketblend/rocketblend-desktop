package pack

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/util"
	"github.com/rocketblend/rocketblend/pkg/downloader"
	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/driver/rocketpack"
	"github.com/rocketblend/rocketblend/pkg/driver/runtime"
	"github.com/rocketblend/rocketblend/pkg/semver"
)

type (
	Build struct {
		ID        uuid.UUID             `json:"id,omitempty"`
		Reference reference.Reference   `json:"reference,omitempty"`
		Name      string                `json:"name,omitempty"`
		Path      string                `json:"path,omitempty"`
		Args      string                `json:"args,omitempty"`
		Version   *semver.Version       `json:"version,omitempty"`
		Addons    []reference.Reference `json:"addons,omitempty"`
		Resource  string                `json:"resource,omitempty"`
		URI       *downloader.URI       `json:"uri,omitempty"`
		UpdatedAt time.Time             `json:"updatedAt,omitempty"`
	}
)

func NewBuild(path string, platform runtime.Platform, rocketpack *rocketpack.RocketPack) (*Build, error) {
	if rocketpack == nil || !rocketpack.IsBuild() {
		return nil, fmt.Errorf("rocketpack is not a build")
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

	return &Build{
		ID:        uuid.New(),
		Name:      name,
		Reference: reference,
		Path:      path,
		Args:      rocketpack.Build.Args,
		Version:   rocketpack.Build.Version,
		Addons:    rocketpack.Build.Addons,
		Resource:  rocketpack.Build.Sources[platform].Resource,
		URI:       rocketpack.Build.Sources[platform].URI,
		UpdatedAt: modTime,
	}, nil
}
