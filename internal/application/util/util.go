package util

import (
	"fmt"
	"os"
	"time"
)

func GetDirModTime(dirPath string) (time.Time, error) {
	// Get file/folder information
	info, err := os.Stat(dirPath)
	if err != nil {
		return time.Time{}, err
	}

	// Check if it's a directory
	if !info.IsDir() {
		return time.Time{}, fmt.Errorf("%s is not a directory", dirPath)
	}

	// Return the last modification time
	return info.ModTime(), nil
}
