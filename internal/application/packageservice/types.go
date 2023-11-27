package packageservice

import pack "github.com/rocketblend/rocketblend-desktop/internal/application/package"

type (
	GetPackageResponse struct {
		Package *pack.Package `json:"package,omitempty"`
	}

	ListPackagesResponse struct {
		Packages []*pack.Package `json:"packages,omitempty"`
	}
)
