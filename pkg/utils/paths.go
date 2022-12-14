package utils

import (
	"path/filepath"
	"regexp"
	"strings"
	"utdocs/diagnostics"

	"github.com/iancoleman/strcase"
)

var pathRegex *regexp.Regexp

func init() {
	regex, err := regexp.Compile("^(\\.\\./)*")
	diagnostics.HandleError(err)

	pathRegex = regex
}

func GetFileName(path string) string {
	cleanPath := filepath.Clean(path)
	return strings.TrimSuffix(filepath.Base(cleanPath), filepath.Ext(cleanPath))
}

func PrettifyTitle(path string) string {
	fileName := GetFileName(path)
	return strcase.ToDelimited(fileName, ' ')
}

func StripParentDirectories(path string) string {
	cleanPath := filepath.ToSlash(filepath.Clean(path))
	return pathRegex.ReplaceAllString(cleanPath, "")
}
