package factory

import (
	rocketblendBlendFile "github.com/rocketblend/rocketblend/pkg/driver/blendfile"
	rocketblendInstallation "github.com/rocketblend/rocketblend/pkg/driver/installation"
	rocketblendRocketPack "github.com/rocketblend/rocketblend/pkg/driver/rocketpack"
	rocketblendConfig "github.com/rocketblend/rocketblend/pkg/rocketblend/config"
)

func (f *factory) GetConfigService() (rocketblendConfig.Service, error) {
	if err := f.checkClosing(); err != nil {
		return nil, err
	}

	return f.rocketblendFactory.GetConfigService()
}

func (f *factory) GetRocketPackService() (rocketblendRocketPack.Service, error) {
	if err := f.checkClosing(); err != nil {
		return nil, err
	}

	return f.rocketblendFactory.GetRocketPackService()
}

func (f *factory) GetInstallationService() (rocketblendInstallation.Service, error) {
	if err := f.checkClosing(); err != nil {
		return nil, err
	}

	return f.rocketblendFactory.GetInstallationService()
}

func (f *factory) GetBlendFileService() (rocketblendBlendFile.Service, error) {
	if err := f.checkClosing(); err != nil {
		return nil, err
	}

	return f.rocketblendFactory.GetBlendFileService()
}
