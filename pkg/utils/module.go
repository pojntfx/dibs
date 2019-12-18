package utils

import (
	"errors"
	"path/filepath"
	"strings"
)

var (
	HeaderReplaceStart = "// DIBS:TEMPREPLACE:START" // The comment that marks the start of a module replace directives block
	HeaderReplaceEnd   = "// DIBS:TEMPREPLACE:END"   // The comment that marks the end of a module replace directives block
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
	for _, line := range strings.Split(content, "\n") {
		if strings.Contains(line, "module") {
			return nil, strings.Split(line, "module ")[1]
		}
	}

	return errors.New("could find module declaration"), ""
}

// GetModuleWithReplaces returns the module file content with additional replacement declarations
func GetModuleWithReplaces(content string, modulesToReplace []string, dirToReplaceHost string) (string, error) {
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
				moduleSuffixes := strings.Split(require, "/")
				modulePartsWithHostReplace := append([]string{dirToReplaceHost}, moduleSuffixes...)
				moduleWithReplacedHost := strings.Join(modulePartsWithHostReplace, "/")

				replaces = append(replaces, moduleWithReplacedHost)
			}
		}
	}

	replaceBlock := HeaderReplaceStart

	for index, replace := range replaces {
		if requires != nil {
			moduleWithReplacePrefix := "replace " + requires[index] + " => " + replace

			replaceBlock = replaceBlock + "\n" + moduleWithReplacePrefix
		} else {
			return "", errors.New("failure in parsing requires")
		}
	}

	replaceBlock = replaceBlock + "\n" + HeaderReplaceEnd

	return content + replaceBlock, nil
}

// GetModuleWithoutReplaces returns the content without the replaces
func GetModuleWithoutReplaces(content string) string {
	var contentWithoutReplaces string
	var isInReplacesBlock bool

	for _, line := range strings.Split(content, "\n") {
		if strings.Contains(line, HeaderReplaceStart) {
			isInReplacesBlock = true
			continue
		}

		if strings.Contains(line, HeaderReplaceEnd) {
			isInReplacesBlock = false
			continue
		}

		if !isInReplacesBlock {
			contentWithoutReplaces += line + "\n"
		}
	}

	return contentWithoutReplaces
}

// GetModulesFromRawInputString returns the modules for a comma-separated list of modules
func GetModulesFromRawInputString(rawInput string) []string {
	modules := strings.Split(rawInput, ",")

	return modules
}
