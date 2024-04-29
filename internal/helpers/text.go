package helpers

import (
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// FilenameToDisplayName converts a filename into a more human readable display name
func FilenameToDisplayName(filename string) string {
	name := filepath.Base(filename)
	ext := filepath.Ext(name)

	name = name[0 : len(name)-len(ext)]

	name = strings.ReplaceAll(name, "_", " ")
	name = strings.ReplaceAll(name, "-", " ")

	title := cases.Title(language.English)

	return title.String(name)
}

// DisplayNameToFilename converts a display name into a filename
func DisplayNameToFilename(displayName string) string {
	name := strings.ReplaceAll(displayName, " ", "-")
	name = strings.ToLower(name)
	return cleanFileName(name)
}

// cleanFileName removes characters that are generally unsafe or not allowed in filenames.
func cleanFileName(name string) string {
	return strings.Map(func(r rune) rune {
		// List of characters that are not allowed in filenames.
		if strings.ContainsRune(`/?<>\\:*|"`, r) || unicode.IsControl(r) {
			return -1 // Remove the character
		}
		return r // Keep the character
	}, name)
}
