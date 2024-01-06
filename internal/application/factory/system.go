package factory

import (
	"github.com/rocketblend/rocketblend-desktop/internal/application/eventservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/metricservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/operationservice"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore"
)

func (f *factory) GetEventService() (eventservice.Service, error) {
	if err := f.checkClosing(); err != nil {
		return nil, err
	}

	f.eventMutex.RLock()
	if f.eventService != nil {
		defer f.eventMutex.RUnlock()
		return f.eventService, nil
	}
	f.eventMutex.RUnlock()

	f.eventMutex.Lock()
	defer f.eventMutex.Unlock()
	eventService, err := eventservice.New(
		eventservice.WithLogger(f.logger),
	)
	if err != nil {
		return nil, err
	}

	f.eventService = eventService

	return f.eventService, nil
}

func (f *factory) GetSearchStore() (searchstore.Store, error) {
	if err := f.checkClosing(); err != nil {
		return nil, err
	}

	f.searchstoreMutex.RLock()
	if f.searchStore != nil {
		defer f.searchstoreMutex.RUnlock()
		return f.searchStore, nil
	}
	f.searchstoreMutex.RUnlock()

	event, err := f.GetEventService()
	if err != nil {
		return nil, err
	}

	f.searchstoreMutex.Lock()
	defer f.searchstoreMutex.Unlock()
	store, err := searchstore.New(
		searchstore.WithLogger(f.logger),
		searchstore.WithDispatcherService(event),
	)
	if err != nil {
		return nil, err
	}

	f.searchStore = store

	return f.searchStore, nil
}

func (f *factory) GetMetricService() (metricservice.Service, error) {
	if err := f.checkClosing(); err != nil {
		return nil, err
	}

	f.metricMutex.RLock()
	if f.metricService != nil {
		defer f.metricMutex.RUnlock()
		return f.metricService, nil
	}

	f.metricMutex.RUnlock()

	store, err := f.GetSearchStore()
	if err != nil {
		return nil, err
	}

	f.metricMutex.Lock()
	defer f.metricMutex.Unlock()
	metricService, err := metricservice.New(
		metricservice.WithLogger(f.logger),
		metricservice.WithStore(store),
	)
	if err != nil {
		return nil, err
	}

	f.metricService = metricService

	return f.metricService, nil
}

func (f *factory) GetOperationService() (operationservice.Service, error) {
	if err := f.checkClosing(); err != nil {
		return nil, err
	}

	f.operationMutex.RLock()
	if f.operationService != nil {
		defer f.operationMutex.RUnlock()
		return f.operationService, nil
	}
	f.operationMutex.RUnlock()

	dispatcher, err := f.GetEventService()
	if err != nil {
		return nil, err
	}

	store, err := f.GetSearchStore()
	if err != nil {
		return nil, err
	}

	f.operationMutex.Lock()
	defer f.operationMutex.Unlock()
	operationService, err := operationservice.New(
		operationservice.WithLogger(f.logger),
		operationservice.WithStore(store),
		operationservice.WithDispatcher(dispatcher),
	)
	if err != nil {
		return nil, err
	}

	f.operationService = operationService

	return f.operationService, nil
}
