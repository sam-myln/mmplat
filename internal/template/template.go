package template

import (
	fs "mmplat/internal/filesystem"
	"mmplat/internal/util"
	"strconv"
)

type FileMetadata struct {
	ID   string
	Name string
	Size string
	Type string
	Path string
}

func NewFileMetadata(id, name, size, fileType, path string) FileMetadata {
	return FileMetadata{
		ID:   id,
		Name: name,
		Size: size,
		Type: fileType,
		Path: path,
	}
}

// PrepFileMetadata templ[name,size,type]
func PrepFileMetadata(id fs.Id, item fs.IItem) FileMetadata {
	return NewFileMetadata(
		strconv.Itoa(int(id)),
		item.Name(),
		strconv.Itoa(int(item.Size())),
		util.ExtToMetadata(item),
		// TODO obfuscate actual file behind exposed API interface
		item.Path(),
	)
}
