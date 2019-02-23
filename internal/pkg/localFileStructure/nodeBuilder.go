package localFileStructure

import (
	"os"
	"path/filepath"
)

func ScanLocalFileStructure(path string) map[string]FileNode {
	fileMap := make(map[string]FileNode)

	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if path == p {
			return nil
		}

		//TODO: Only allow jpg and png files here

		fileMap[p] = FileNode{
			key:p,
			name:info.Name(),
			isDir:info.IsDir(),
		}
		return nil;
	})

	if err != nil {
		panic(err)
	}

	return fileMap
}
