package utils

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// ParseModuleFromMessage gets the module name and event timestamp from a message
func ParseModuleFromMessage(message string) (name, timestamp string) {
	moduleParts := strings.Split(message, "@")

	return moduleParts[0], moduleParts[1]
}

// GetPathForModule builds the path for a module
func GetPathForModule(baseDir, module string) string {
	pathParts := append([]string{baseDir}, strings.Split(module, "/")...)

	fullModulePath := filepath.Join(pathParts...)

	return fullModulePath
}

// GetModuleName returns the module name from `go.mod`
func GetModuleName(goModFilePath string) (error, string) {
	f, err := ioutil.ReadFile(goModFilePath)
	if err != nil {
		return errors.New("Could not read module file"), ""
	}

	for _, line := range strings.Split(string(f), "\n") {
		if strings.Contains(line, "module") {
			return nil, strings.Split(line, "module ")[1]
		}
	}

	return errors.New("Could find module declaration"), ""
}
