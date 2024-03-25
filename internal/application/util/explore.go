package util

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
)

// Explore opens the file explorer at the specified path
func Explore(ctx context.Context, path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.CommandContext(ctx, "explorer", path)
	case "darwin":
		cmd = exec.CommandContext(ctx, "open", path)
	case "linux":
		cmd = exec.CommandContext(ctx, "xdg-open", path)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Start()
}
