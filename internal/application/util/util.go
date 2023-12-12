package util

import (
	"crypto/sha1"
	"os"
	"time"

	"github.com/google/uuid"
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
