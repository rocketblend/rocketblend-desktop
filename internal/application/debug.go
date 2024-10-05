package application

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type (
	Debug struct {
		ctx context.Context
	}
)

func (d *Debug) Quit() {
	fmt.Println("quitting application")
	runtime.Quit(d.ctx)
}

func (d *Debug) startup(ctx context.Context) {
	fmt.Println("starting application")
	d.ctx = ctx
}

func (d *Debug) shutdown(ctx context.Context) {
	fmt.Println("shutting down application")
}
