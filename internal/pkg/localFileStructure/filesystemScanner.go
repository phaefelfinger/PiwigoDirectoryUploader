package localFileStructure

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

func ScanLocalFileStructure(path string) (map[string]*FilesystemNode, error) {
	fullPathRoot, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	logrus.Infof("Scanning %s for images...", fullPathRoot)

	fileMap := make(map[string]*FilesystemNode)
	fullPathReplace := fmt.Sprintf("%s%c", fullPathRoot, os.PathSeparator)
	numberOfDirectories := 0
	numberOfImages := 0

	err = filepath.Walk(fullPathRoot, func(path string, info os.FileInfo, err error) error {
		if fullPathRoot == path {
			return nil
		}

		//TODO: Only allow jpg and png files here

		key := strings.Replace(path, fullPathReplace, "", 1)

		fileMap[path] = &FilesystemNode{
			Key:     key,
			Path:    path,
			Name:    info.Name(),
			IsDir:   info.IsDir(),
			ModTime: info.ModTime(),
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
