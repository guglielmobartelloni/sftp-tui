package tui

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/knipferrc/teacup/icons"
)

// ConvertBytesToSizeString converts a byte count to a human readable string.
func ConvertBytesToSizeString(size int64) string {
	const (
		thousand    = 1000
		ten         = 10
		fivePercent = 0.0499
	)

	if size < thousand {
		return fmt.Sprintf("%dB", size)
	}

	suffix := []string{
		"K", // kilo
		"M", // mega
		"G", // giga
		"T", // tera
		"P", // peta
		"E", // exa
		"Z", // zeta
		"Y", // yotta
	}

	curr := float64(size) / thousand
	for _, s := range suffix {
		if curr < ten {
			return fmt.Sprintf("%.1f%s", curr-fivePercent, s)
		} else if curr < thousand {
			return fmt.Sprintf("%d%s", int(curr), s)
		}
		curr /= thousand
	}

	return ""
}

// Get the fancy file description with file permission, file size, and mod timestamp
func getFileDescription(value fs.FileInfo) string {
	status := fmt.Sprintf("%s %s %s",
		value.ModTime().Format("2006-01-02 15:04:05"),
		value.Mode().String(),
		ConvertBytesToSizeString(value.Size()))
	return status
}

// Get the file icons based on its properties
func getFileIcon(value fs.FileInfo) string {
	icon, _ := icons.GetIcon(
		value.Name(),
		filepath.Ext(value.Name()),
		icons.GetIndicator(value.Mode()),
	)
	return icon
}

// Utility function to handle errors
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
