package util

import (
	"os"
	"time"
)

func GetModTime(path string) (time.Time, error) {
	// Get file/folder information
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}

	// Return the last modification time
	return info.ModTime(), nil
}
