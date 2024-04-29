package container

import (
	"fmt"

	"github.com/rocketblend/rocketblend-desktop/internal/application/configurator"
	"github.com/rocketblend/rocketblend-desktop/internal/application/operator"
	pack "github.com/rocketblend/rocketblend-desktop/internal/application/package"
	"github.com/rocketblend/rocketblend-desktop/internal/application/project"
	"github.com/rocketblend/rocketblend-desktop/internal/application/store"
	"github.com/rocketblend/rocketblend-desktop/internal/application/tracker"
	"github.com/rocketblend/rocketblend-desktop/internal/application/types"
)

func (c *Container) GetConfigurator() (types.Configurator, error) {
	var err error
	c.configuratorHolder.once.Do(func() {
		c.configuratorHolder.instance, err = configurator.New(
			configurator.WithLogger(c.logger),
			configurator.WithValidator(c.validator),
			configurator.WithLocation(c.applicationDir),
		)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get/create configurator: %w", err)
	}

	return c.configuratorHolder.instance, nil
}

func (c *Container) GetStore() (types.Store, error) {
	var err error
	c.storeHolder.once.Do(func() {
		dispatcher, errDispatcher := c.GetDispatcher()
		if err != nil {
			err = errDispatcher
			return
		}

		c.storeHolder.instance, err = store.New(
			store.WithLogger(c.logger),
			store.WithDispatcher(dispatcher),
		)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get/create store: %w", err)
	}

	return c.storeHolder.instance, nil
}

func (c *Container) GetTracker() (types.Tracker, error) {
	var err error
	c.trackerHolder.once.Do(func() {
		store, errStore := c.GetStore()
		if err != nil {
			err = errStore
			return
		}

		c.trackerHolder.instance, err = tracker.New(
			tracker.WithLogger(c.logger),
			tracker.WithStore(store),
		)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get/create tracker: %w", err)
	}

	return c.trackerHolder.instance, nil
}

func (c *Container) GetOperator() (types.Operator, error) {
	var err error
	c.operatorHolder.once.Do(func() {
		store, errStore := c.GetStore()
		if err != nil {
			err = errStore
			return
		}

		dispatcher, errDispatcher := c.GetDispatcher()
		if err != nil {
			err = errDispatcher
			return
		}

		c.operatorHolder.instance, err = operator.New(
			operator.WithLogger(c.logger),
			operator.WithStore(store),
			operator.WithDispatcher(dispatcher),
		)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get/create operator: %w", err)
	}

	return c.operatorHolder.instance, nil
}

func (c *Container) GetPortfolio() (types.Portfolio, error) {
	var err error
	c.portfolioHolder.once.Do(func() {
		store, errStore := c.GetStore()
		if err != nil {
			err = errStore
			return
		}

		dispatcher, errDispatcher := c.GetDispatcher()
		if err != nil {
			err = errDispatcher
			return
		}

		configurator, errConfigurator := c.GetConfigurator()
		if err != nil {
			err = errConfigurator
			return
		}

		rbConfigurator, errRBConfigurator := c.rbContainer.GetConfigurator()
		if err != nil {
			err = errRBConfigurator
			return
		}

		rbRepository, errRBRepository := c.rbContainer.GetRepository()
		if err != nil {
			err = errRBRepository
			return
		}

		rbDriver, errRBDriver := c.rbContainer.GetDriver()
		if err != nil {
			err = errRBDriver
			return
		}

		blender, errBlender := c.rbContainer.GetBlender()
		if err != nil {
			err = errBlender
			return
		}

		c.portfolioHolder.instance, err = project.New(
			project.WithLogger(c.logger),
			project.WithValidator(c.validator),
			project.WithStore(store),
			project.WithDispatcher(dispatcher),
			project.WithConfigurator(configurator),
			project.WithRocketBlendConfigurator(rbConfigurator),
			project.WithRocketBlendRepository(rbRepository),
			project.WithRocketBlendDriver(rbDriver),
			project.WithBlender(blender),
			project.WithWatcherDebounceDuration(c.watcherDebounce),
		)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get/create portfolio: %w", err)
	}

	return c.portfolioHolder.instance, nil
}

func (c *Container) GetCatalog() (types.Catalog, error) {
	var err error
	c.catalogHolder.once.Do(func() {
		store, errStore := c.GetStore()
		if err != nil {
			err = errStore
			return
		}

		dispatcher, errDispatcher := c.GetDispatcher()
		if err != nil {
			err = errDispatcher
			return
		}

		rbConfigurator, errRBConfigurator := c.rbContainer.GetConfigurator()
		if err != nil {
			err = errRBConfigurator
			return
		}

		rbRepository, errRBRepository := c.rbContainer.GetRepository()
		if err != nil {
			err = errRBRepository
			return
		}

		c.catalogHolder.instance, err = pack.New(
			pack.WithLogger(c.logger),
			pack.WithValidator(c.validator),
			pack.WithStore(store),
			pack.WithDispatcher(dispatcher),
			pack.WithRocketBlendRepository(rbRepository),
			pack.WithRocketBlendConfigurator(rbConfigurator),
			pack.WithWatcherDebounceDuration(c.watcherDebounce),
		)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get/create catalog: %w", err)
	}

	return c.catalogHolder.instance, nil
}
