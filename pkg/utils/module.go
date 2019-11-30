package utils

import (
	"path/filepath"
	"strings"
)

// ParseModuleFromMessage gets the module name and event timestamp from a message
func ParseModuleFromMessage(m string) (name, timestamp string) {
	moduleParts := strings.Split(m, "@")
	return moduleParts[0], moduleParts[1]
}

// GetPathForModule builds the path for a module
func GetPathForModule(baseDir, m string) string {
	pathParts := append([]string{baseDir}, strings.Split(m, "/")...)

	fullModulePath := filepath.Join(pathParts...)

	return fullModulePath
}
