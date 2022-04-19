package utils

import (
	"path/filepath"
	"strings"
)

// GetFileName without extension
func GetFileName(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}
