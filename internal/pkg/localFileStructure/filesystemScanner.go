package localFileStructure

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

func ScanLocalFileStructure(path string) (map[string]*FilesystemNode, error) {
	fileMap := make(map[string]*FilesystemNode)

	relativeRoot := filepath.Base(path) + "/"

	numberOfDirectories := 0
	numberOfImages := 0

	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if path == p {
			return nil
		}

		//TODO: Only allow jpg and png files here

		key := strings.Replace(p, relativeRoot, "", 1)

		fileMap[p] = &FilesystemNode{
			Key:   key,
			Name:  info.Name(),
			IsDir: info.IsDir(),
		}

		if info.IsDir() {
			numberOfDirectories += 1
		} else {
			numberOfImages += 1
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	logrus.Infof("Found %d directories and %d images on the local filesystem", numberOfDirectories, numberOfImages)

	return fileMap, nil
}
