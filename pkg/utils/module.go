package utils

import (
	"errors"
	"github.com/go-redis/redis/v7"
	"io/ioutil"
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

// RegisterModule registers a module in Redis
func RegisterModule(r *redis.Client, prefix, suffix, m string) {
	r.Publish(prefix+":"+suffix, WithTimestamp(m))
}

// UnregisterModule unregisters a module from Redis
func UnregisterModule(r *redis.Client, prefix, suffix, m string) {
	r.Publish(prefix+":"+suffix, WithTimestamp(m))
}
