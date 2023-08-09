package projectstore

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/rjeczalik/notify"
)

var watchFileExtensions = []string{
	//".blend",
	".yaml",
}

type (
	watcher struct {
		EventChannel chan notify.EventInfo
		Ctx          context.Context
		Cancel       context.CancelFunc
	}

	projectEventInfo struct {
		ProjectPath string
		EventInfo   notify.EventInfo
	}

	projectEvent struct {
		event     *projectEventInfo
		timer     *time.Timer
		eventLock sync.Mutex
	}
)

func (s *store) watchPath(path string) error {
	if !s.watcherEnabled {
		return nil
	}

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

func (s *store) unwatchPath(path string) error {
	if !s.watcherEnabled {
		return nil
	}

	notify.Stop(s.watchers[path].EventChannel)
	s.watchers[path].Cancel()
	delete(s.watchers, path)

	s.logger.Debug("Unwatching path", map[string]interface{}{
		"path": path,
	})

	return nil
}

func (s *store) monitorEvents(events chan notify.EventInfo, ctx context.Context) {
	for {
		select {
		case event := <-events:
			if isWatchFile(event.Path()) {
				s.handleEventDebounced(&projectEventInfo{
					ProjectPath: s.getProjectPath(event.Path()),
					EventInfo:   event,
				})
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *store) handleEventDebounced(event *projectEventInfo) {
	s.logger.Info("Filesystem event occurred", map[string]interface{}{
		"event": event,
	})

	s.emu.Lock()
	pe, ok := s.events[event.ProjectPath]
	if !ok {
		pe = &projectEvent{}
		s.events[event.ProjectPath] = pe
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

func (s *store) handleEvent(event *projectEventInfo) {
	s.logger.Info("Filesystem event occurred", map[string]interface{}{
		"event":   event.EventInfo.Event(),
		"path":    event.EventInfo.Path(),
		"project": event.ProjectPath,
	})

	switch event.EventInfo.Event() {
	case notify.Write:
		if err := s.loadProject(event.ProjectPath); err != nil {
			s.logger.Error("Error while loading project", map[string]interface{}{
				"err": err,
			})
		}

	case notify.Remove, notify.Rename:
		if err := s.removeProjectsInPath(event.ProjectPath); err != nil {
			s.logger.Error("Error while removing project", map[string]interface{}{
				"err": err,
			})
		}
	}
}

func isWatchFile(filename string) bool {
	for _, ext := range watchFileExtensions {
		if filepath.Ext(filename) == ext {
			return true
		}
	}

	return false
}
