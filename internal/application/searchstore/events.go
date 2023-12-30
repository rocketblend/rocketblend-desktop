package searchstore

import (
	"context"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/indextype"
)

type EventType int

const (
	InsertEvent EventType = iota
	UpdateEvent
	RemoveEvent
)

type (
	Event struct {
		ID        uuid.UUID           `json:"id,omitempty"`
		Type      EventType           `json:"type,omitempty"`
		IndexType indextype.IndexType `json:"indexType,omitempty"`
		//Data      interface{} 	   `json:"data,omitempty"`
	}

	EventHandler func(Event)
)

func (s *store) RegisterListener(ctx context.Context, id string, handler EventHandler) func() {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.listeners == nil {
		s.listeners = make(map[string]EventHandler)
	}
	s.listeners[id] = handler

	go func() {
		<-ctx.Done()
		s.UnregisterListener(id)
	}()

	return func() {
		s.UnregisterListener(id)
	}
}

func (s *store) UnregisterListener(id string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.logger.Debug("Unregistering listener", map[string]interface{}{
		"id": id,
	})

	delete(s.listeners, id)
}

func (s *store) ClearListeners() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.listeners = make(map[string]EventHandler)
}

func (s *store) emitEvent(event Event) {
	s.lock.Lock()
	listeners := make(map[string]EventHandler, len(s.listeners))
	for id, handler := range s.listeners {
		listeners[id] = handler
	}
	s.lock.Unlock()

	for _, handler := range listeners {
		handler(event)
	}
}
