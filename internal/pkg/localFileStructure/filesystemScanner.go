/*
 * Copyright (C) 2020 Philipp Haefelfinger (http://www.haefelfinger.ch/). All Rights Reserved.
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

func ScanLocalFileStructure(path string, extensions []string, ignoreDirs []string, dirSuffixToSkip int) (map[string]*FilesystemNode, error) {
	fullPathRoot, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	ignoreDirsMap := make(map[string]struct{}, len(ignoreDirs))
	for _, ignoredFolder := range ignoreDirs {
		ignoreDirsMap[strings.ToLower(ignoredFolder)] = struct{}{}
	}

	extensionsMap := make(map[string]struct{}, len(extensions))
	for _, extension := range extensions {
		extensionsMap["."+strings.ToLower(extension)] = struct{}{}
	}

	if len(extensionsMap) == 0 {
		logrus.Debug("No extensions specified, adding jpg and png")
		extensionsMap[".jpg"] = struct{}{}
		extensionsMap[".png"] = struct{}{}
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
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		_, dirIgnored := ignoreDirsMap[strings.ToLower(info.Name())]
		if dirIgnored && info.IsDir() {
			logrus.Tracef("Skipping ignored directory %s", path)
			return filepath.SkipDir
		}

		extension := strings.ToLower(filepath.Ext(path))
		_, extensionSupported := extensionsMap[extension]
		if !extensionSupported && !info.IsDir() {
			return nil
		}

		key := buildKey(path, info, fullPathReplace, dirSuffixToSkip)

		fileMap[path] = &FilesystemNode{
			Key:     key,
			Path:    path,
			Name:    filepath.Base(key),
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

func buildKey(path string, info os.FileInfo, fullPathReplace string, dirSuffixToSkip int) string {
	if info.IsDir() {
		return trimPathForKey(path, fullPathReplace, dirSuffixToSkip)
	}
	fileName := filepath.Base(path)
	directoryName := filepath.Dir(path)
	cleanDir := trimPathForKey(directoryName, fullPathReplace, dirSuffixToSkip)
	return filepath.Join(cleanDir, fileName)
}

func trimPathForKey(path string, fullPathReplace string, dirSuffixToSkip int) string {
	trimmedPath := strings.Replace(path, fullPathReplace, "", 1)
	for i := 0; i < dirSuffixToSkip; i++ {
		trimmedPath = filepath.Clean(strings.TrimSuffix(trimmedPath, filepath.Base(trimmedPath)))
	}
	if trimmedPath == "." {
		return "root"
	}
	return trimmedPath
}
