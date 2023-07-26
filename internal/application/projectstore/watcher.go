package projectstore

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/rjeczalik/notify"
)

var watchFileExtensions = []string{
	//".blend",
	".yaml",
}

type (
	Watcher struct {
		EventChannel chan notify.EventInfo
		Ctx          context.Context
		Cancel       context.CancelFunc
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

	s.mu.Lock()
	s.watchers[path] = &Watcher{
		EventChannel: eventChannel,
		Ctx:          ctx,
		Cancel:       cancel,
	}
	s.mu.Unlock()

	return nil
}

func (s *store) unwatchPath(path string) error {
	if !s.watcherEnabled {
		return nil
	}

	s.mu.Lock()
	notify.Stop(s.watchers[path].EventChannel)
	s.watchers[path].Cancel()
	delete(s.watchers, path)
	s.mu.Unlock()

	return nil
}

func (s *store) monitorEvents(events chan notify.EventInfo, ctx context.Context) {
	for {
		select {
		case event := <-events:
			s.handleEvent(event)
		case <-ctx.Done():
			return
		}
	}
}

func (s *store) handleEvent(event notify.EventInfo) {
	s.logger.Info("Filesystem event occurred", map[string]interface{}{
		"event": event,
	})

	if !isWatchFile(event.Path()) {
		return
	}

	projectPath := s.getProjectPath(event.Path())

	switch event.Event() {
	case notify.Write:
		s.logger.Debug("Modified file", map[string]interface{}{
			"file": event.Path(),
		})

		if err := s.loadProject(projectPath); err != nil {
			s.logger.Error("Error while loading project", map[string]interface{}{
				"err": err,
			})
		}

	case notify.Remove, notify.Rename:
		s.logger.Debug("Removed or renamed file", map[string]interface{}{
			"file": event.Path(),
		})

		if err := s.removeProject(projectPath); err != nil {
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
