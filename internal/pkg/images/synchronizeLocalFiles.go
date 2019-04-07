/*
 * Copyright (C) 2019 Philipp Haefelfinger (http://www.haefelfinger.ch/). All Rights Reserved.
 * This application is licensed under GPLv2. See the LICENSE file in the root directory of the project.
 */

package images

import (
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/datastore"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/localFileStructure"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type fileChecksumCalculator func(filePath string) (string, error)

// Update the local image metadata by walking through all found files and check if the modification date has changed
// or if they are new to the local database. If the files is new or changed, the md5sum will be rebuilt as well.
func SynchronizeLocalImageMetadata(imageDb datastore.ImageMetadataProvider, categoryDb datastore.CategoryProvider, fileSystemNodes map[string]*localFileStructure.FilesystemNode, checksumCalculator fileChecksumCalculator) error {
	logrus.Debug("Starting SynchronizeLocalImageMetadata")
	defer logrus.Debug("Leaving SynchronizeLocalImageMetadata")

	logrus.Info("Synchronizing local image metadata database with local available images")

	err := synchronizeLocalImageMetadataScanNewFiles(fileSystemNodes, imageDb, categoryDb, checksumCalculator)
	if err != nil {
		return err
	}

	err = synchronizeLocalImageMetadataFindFilesToDelete(imageDb)
	if err != nil {
		return err
	}
	return nil
}

func synchronizeLocalImageMetadataScanNewFiles(fileSystemNodes map[string]*localFileStructure.FilesystemNode, imageDb datastore.ImageMetadataProvider, categoryDb datastore.CategoryProvider, checksumCalculator fileChecksumCalculator) error {
	logrus.Debug("Entering synchronizeLocalImageMetadataScanNewFiles")
	defer logrus.Debug("Leaving synchronizeLocalImageMetadataScanNewFiles")

	for _, file := range fileSystemNodes {
		if file.IsDir {
			// we are only interested in files not directories
			continue
		}

		metadata, err := imageDb.ImageMetadata(file.Path)
		if err == datastore.ErrorRecordNotFound {
			logrus.Debugf("Creating new metadata entry for %s.", file.Path)
			metadata = datastore.ImageMetaData{}
			metadata.Filename = file.Name
			metadata.FullImagePath = file.Path
			metadata.CategoryPath = filepath.Dir(file.Key)

			category, err := categoryDb.GetCategoryByKey(metadata.CategoryPath)
			if err == nil {
				metadata.CategoryPiwigoId = category.PiwigoId
			} else {
				logrus.Warnf("No category found for image %s - %s", file.Path, err)
			}

		} else if err != nil {
			logrus.Errorf("Could not get metadata due to trouble. Cancelling - %s", err)
			return err
		}

		if fileDidNotChange(&metadata, file) {
			logrus.Debugf("No changes found for file %s", file.Path)
			continue
		}

		metadata.UploadRequired = !metadata.LastChange.Equal(file.ModTime) || metadata.PiwigoId == 0
		metadata.DeleteRequired = false
		metadata.LastChange = file.ModTime
		metadata.Md5Sum, err = checksumCalculator(file.Path)
		if err != nil {
			logrus.Warnf("Could not calculate checksum for file %s. Skipping...", file.Path)
			continue
		}

		err = imageDb.SaveImageMetadata(metadata)
		if err != nil {
			return err
		}
	}
	return nil
}

func synchronizeLocalImageMetadataFindFilesToDelete(imageDb datastore.ImageMetadataProvider) error {
	logrus.Debug("Entering SynchronizeLocalImageMetadataFindFilesToDelete")
	defer logrus.Debug("Leaving SynchronizeLocalImageMetadataFindFilesToDelete")

	images, err := imageDb.ImageMetadataAll()
	if err != nil {
		return err
	}

	for _, img := range images {
		if _, err := os.Stat(img.FullImagePath); os.IsNotExist(err) {
			img.UploadRequired = false
			img.DeleteRequired = true
			err := imageDb.SaveImageMetadata(img)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func fileDidNotChange(metadata *datastore.ImageMetaData, file *localFileStructure.FilesystemNode) bool {
	return metadata.LastChange.Equal(file.ModTime) && !metadata.DeleteRequired
}
