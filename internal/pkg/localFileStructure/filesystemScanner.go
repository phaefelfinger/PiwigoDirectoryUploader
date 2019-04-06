/*
 * Copyright (C) 2019 Philipp Haefelfinger (http://www.haefelfinger.ch/). All Rights Reserved.
 * This application is licensed under GPLv2. See the LICENSE file in the root directory of the project.
 */

package localFileStructure

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
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

		if strings.HasPrefix(info.Name(), ".") {
			logrus.Tracef("Skipping hidden file or directory %s", path)
			return nil
		}

		extension := strings.ToLower(filepath.Ext(path))
		if extension != ".jpg" && extension != ".png" && !info.IsDir() {
			return nil
		}

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
