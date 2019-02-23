package localFileStructure

import "fmt"

type FilesystemNode struct {
	Key   string
	Name  string
	IsDir bool
}

func (n *FilesystemNode) String() string {
	return fmt.Sprintf("FilesystemNode: %s", n.Key)
}
