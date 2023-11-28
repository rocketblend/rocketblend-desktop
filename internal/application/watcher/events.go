package watcher

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rjeczalik/notify"
)

type (
	watcher struct {
		EventChannel chan notify.EventInfo
		Ctx          context.Context
		Cancel       context.CancelFunc
	}

	objectEventInfo struct {
		ObjectPath string
		EventInfo  notify.EventInfo
	}

	projectEvent struct {
		event     *objectEventInfo
		timer     *time.Timer
		eventLock sync.Mutex
	}

	// Create a fake notify.EventInfo to trigger initial load.
	eventInfo struct {
		event notify.Event
		path  string
		sys   interface{}
	}
)

func (m eventInfo) Event() notify.Event {
	return m.event
}

func (m eventInfo) Path() string {
	return m.path
}

func (m eventInfo) Sys() interface{} {
	return m.sys
}

func (s *service) watchPath(path string) error {
	eventChannel := make(chan notify.EventInfo, 1)

	err := notify.Watch(path+"/...", eventChannel, notify.Write|notify.Remove|notify.Rename)
	if err != nil {
		return fmt.Errorf("unable to add path %s to watcher: %w", path, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go s.monitorEvents(eventChannel, ctx)

	s.watchers[path] = &watcher{
		EventChannel: eventChannel,
		Ctx:          ctx,
		Cancel:       cancel,
	}

	s.logger.Debug("Watching path", map[string]interface{}{
		"path": path,
	})

	return nil
}

func (s *service) unwatchPath(path string) error {
	notify.Stop(s.watchers[path].EventChannel)
	s.watchers[path].Cancel()
	delete(s.watchers, path)

	s.logger.Debug("Unwatching path", map[string]interface{}{
		"path": path,
	})

	return nil
}

func (s *service) monitorEvents(events chan notify.EventInfo, ctx context.Context) {
	for {
		select {
		case event := <-events:
			// Only handle events for files we care about. Change to function set in options.
			if s.isWatchableFile(event.Path()) {
				s.handleEventDebounced(&objectEventInfo{
					ObjectPath: s.resolveObjectPath(event.Path()),
					EventInfo:  event,
				})
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *service) handleEventDebounced(event *objectEventInfo) {
	s.logger.Info("Filesystem event occurred", map[string]interface{}{
		"event": event,
	})

	s.emu.Lock()
	pe, ok := s.events[event.ObjectPath]
	if !ok {
		pe = &projectEvent{}
		s.events[event.ObjectPath] = pe
	}
	s.emu.Unlock()

	pe.eventLock.Lock()
	defer pe.eventLock.Unlock()

	pe.event = event

	if pe.timer != nil {
		return
	}

	pe.timer = time.AfterFunc(s.debounceDuration, func() {
		pe.eventLock.Lock()
		defer pe.eventLock.Unlock()

		s.handleEvent(pe.event)
		pe.timer = nil
	})
}

func (s *service) handleEvent(event *objectEventInfo) {
	s.logger.Info("Filesystem event occurred", map[string]interface{}{
		"event":      event.EventInfo.Event(),
		"path":       event.EventInfo.Path(),
		"objectPath": event.ObjectPath,
	})

	switch event.EventInfo.Event() {
	case notify.Create, notify.Write:
		// Export event.
		if err := s.updateObject(event.ObjectPath); err != nil {
			s.logger.Error("Error while loading project", map[string]interface{}{
				"err": err,
			})
		}

	case notify.Remove, notify.Rename:
		// Export event.
		if err := s.removeObject(event.ObjectPath); err != nil {
			s.logger.Error("Error while removing project", map[string]interface{}{
				"err": err,
			})
		}
	}
}

func (s *service) isWatchableFile(filename string) bool {
	if s.isWatchableFileFunc != nil {
		return s.isWatchableFileFunc(filename)
	}

	return true
}
