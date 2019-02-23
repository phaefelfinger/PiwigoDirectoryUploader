package localFileStructure

import (
	"os"
	"path/filepath"
)

func ScanLocalFileStructure(path string) map[string]FilesystemNode {
	fileMap := make(map[string]FilesystemNode)

	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if path == p {
			return nil
		}

		//TODO: Only allow jpg and png files here

		fileMap[p] = FilesystemNode{
			Key:   p,
			Name:  info.Name(),
			IsDir: info.IsDir(),
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return fileMap
}
