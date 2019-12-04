package utils

import (
	"errors"
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
func GetModuleName(content string) (error, string) {
	for _, line := range strings.Split(string(content), "\n") {
		if strings.Contains(line, "module") {
			return nil, strings.Split(line, "module ")[1]
		}
	}

	return errors.New("Could find module declaration"), ""
}

// GetModuleWithReplaces returns the module file content with additional replacement declarations
func GetModuleWithReplaces(content string, modulesToReplace []string, hostReplacement string) string {
	var requires []string
	var isInRequiresBlock bool

	for _, line := range strings.Split(content, "\n") {
		if strings.Contains(line, "require (") {
			isInRequiresBlock = true
			continue
		}

		if strings.Contains(line, ")") {
			isInRequiresBlock = false
			continue
		}

		if isInRequiresBlock {
			lineParts := strings.Split(line, " ")
			requireWithoutVersion := strings.TrimSpace(strings.Join(lineParts[:len(lineParts)-1], ""))
			requires = append(requires, requireWithoutVersion)
		}
	}

	var replaces []string

	for _, require := range requires {
		for _, moduleToReplace := range modulesToReplace {
			if require == moduleToReplace {
				moduleSuffixes := strings.Split(require, "/")[1:]
				modulePartsWithHostReplace := append([]string{hostReplacement}, moduleSuffixes...)
				moduleWithReplacedHost := strings.Join(modulePartsWithHostReplace, "/")

				replaces = append(replaces, moduleWithReplacedHost)
			}
		}
	}

	var replaceBlock string

	replaceBlock = replaceBlock + "\n// GODIBS:TEMPREPLACE:START"

	for index, replace := range replaces {
		moduleWithReplacePrefix := "replace " + requires[index] + " => " + replace + " master"

		replaceBlock = replaceBlock + "\n" + moduleWithReplacePrefix
	}

	replaceBlock = replaceBlock + "\n// GODIBS:TEMPREPLACE:END"

	return content + "\n" + replaceBlock
}

// GetModuleWithoutReplaces returns the content without the replaces
func GetModuleWithoutReplaces(content string) string {
	var contentWithoutReplaces string
	var isInReplacesBlock bool

	for _, line := range strings.Split(content, "\n") {
		if strings.Contains(line, "// GODIBS:TEMPREPLACE:START") {
			isInReplacesBlock = true
			continue
		}

		if strings.Contains(line, "// GODIBS:TEMPREPLACE:END") {
			isInReplacesBlock = false
			continue
		}

		if !isInReplacesBlock {
			contentWithoutReplaces += line + "\n"
		}
	}

	return contentWithoutReplaces
}
