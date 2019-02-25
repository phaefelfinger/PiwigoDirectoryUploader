package localFileStructure

import (
	"os"
	"path/filepath"
	"strings"
)

func ScanLocalFileStructure(path string) (map[string]*FilesystemNode, error) {
	fileMap := make(map[string]*FilesystemNode)

	relativeRoot := filepath.Base(path)+"/"

	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if path == p {
			return nil
		}

		//TODO: Only allow jpg and png files here

		key := strings.Replace(p,relativeRoot,"",1)

		fileMap[p] = &FilesystemNode{
			Key:   key,
			Name:  info.Name(),
			IsDir: info.IsDir(),
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return fileMap, nil
}
