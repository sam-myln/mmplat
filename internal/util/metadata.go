package util

import (
	fs "mmplat/internal/filesystem"
	"path/filepath"
)

func ExtToMetadata(item fs.IItem) string {
	// ExtToMetadata F-mat w/o dot
	// append f-mats here
	overwriteList := map[string]string{
		"epub": "application/epub+zip",
	}
	ext := filepath.Ext(item.Name())
	if ext != "" {
		ext = ext[1:]
	}
	if overwriteList[ext] != "" {
		return overwriteList[ext]
	}
	return item.Metadata()
}
