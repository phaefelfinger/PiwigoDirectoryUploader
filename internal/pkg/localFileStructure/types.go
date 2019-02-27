package localFileStructure

import (
	"fmt"
	"time"
)

type FilesystemNode struct {
	Key     string
	Path    string
	Name    string
	IsDir   bool
	ModTime time.Time
}

func (n *FilesystemNode) String() string {
	return fmt.Sprintf("FilesystemNode: %s", n.Path)
}

type ImageNode struct {
	Path    string
	ModTime time.Time
	Md5Sum  string
}
