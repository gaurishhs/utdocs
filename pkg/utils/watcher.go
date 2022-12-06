package utils

import (
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func ShouldRebuild(path string, op fsnotify.Op) bool {
	base := filepath.Base(path)

	if base == ".DS_Store" {
		return false
	}

	if base == "4913" {
		return false
	}

	if strings.HasSuffix(base, "~") {
		return false
	}

	return true
}
