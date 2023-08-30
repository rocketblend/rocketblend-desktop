package projectservice

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
)

func openInFileExplorer(ctx context.Context, path string) error {
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
