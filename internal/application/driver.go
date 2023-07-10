package application

import (
	"context"
	"fmt"

	"github.com/rocketblend/rocketblend-desktop/internal/application/services/ipcService"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Driver struct
type Driver struct {
	ctx context.Context

	argumentChannel chan ipcService.Args
}

// NewApp creates a new App application struct
func NewDriver() *Driver {
	return &Driver{
		argumentChannel: make(chan ipcService.Args),
	}
}

// Greet returns a greeting for the given name
func (a *Driver) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// Quit quits the application
func (a *Driver) Quit() {
	runtime.Quit(a.ctx)
}

// eventEmitter emits an event to the frontend
// func (a *Driver) eventEmitter(eventName string, optionalData ...interface{}) error {
// 	if a.ctx == nil {
// 		return fmt.Errorf("context is nil")
// 	}

// 	runtime.EventsEmit(a.ctx, eventName, optionalData...)

// 	return nil
// }

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *Driver) startup(ctx context.Context) {
	a.ctx = ctx
}

// shutdown is called when the app is shutting down
func (b *Driver) shutdown(ctx context.Context) {}

// menu returns the application menu
func (a *Driver) menu() *menu.Menu {
	AppMenu := menu.NewMenu()
	FileMenu := AppMenu.AddSubmenu("File")
	// FileMenu.AddText("&Open", keys.CmdOrCtrl("o"), openFile)
	// FileMenu.AddSeparator()
	FileMenu.AddText("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		a.Quit()
	})

	return AppMenu
}
