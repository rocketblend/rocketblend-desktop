package helpers

import (
	"context"
	"crypto/sha1"
	"os"
	"time"

	"github.com/google/uuid"
)

type HeartbeatFuncType func()

// startHeartbeat starts a heartbeat goroutine that will execute the provided heartbeat function at each interval
func StartHeartbeat(ctx context.Context, interval time.Duration, heartbeatFunc HeartbeatFuncType) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			heartbeatFunc()
		case <-ctx.Done():
			return
		}
	}
}

// TouchFile updates the timestamps of an existing file. It returns an error if the file does not exist.
func TouchFile(filename string) error {
	currentTime := time.Now()

	// Open the file with READ and WRITE permissions without creating it if it doesn't exist
	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Update the access and modification times of the file
	return os.Chtimes(filename, currentTime, currentTime)
}

func GetModTime(path string) (time.Time, error) {
	// Get file/folder information
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}

	// Return the last modification time
	return info.ModTime(), nil
}

func StringToUUID(seed string) (uuid.UUID, error) {
	hasher := sha1.New()
	hasher.Write([]byte(seed))
	hashBytes := hasher.Sum(nil)

	// Truncate or pad the hash to 16 bytes
	uuidBytes := make([]byte, 16)
	copy(uuidBytes, hashBytes)

	// Generate UUID from the hash
	u, err := uuid.FromBytes(uuidBytes)
	if err != nil {
		return uuid.UUID{}, err
	}

	return u, nil
}
