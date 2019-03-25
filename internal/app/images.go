/*
 * Copyright (C) 2019 Philipp Haefelfinger (http://www.haefelfinger.ch/). All Rights Reserved.
 * This application is licensed under GPLv2. See the LICENSE file in the root directory of the project.
 */

package app

import (
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/localFileStructure"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/piwigo"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

// to make use of the new local data store, we have to rethink and refactor the whole local detection process
// extend the storage of the images to keep track of upload state
// TBD: How to deal with updates -> delete / upload all based on md5 sums

type fileChecksumCalculator func(filePath string) (string, error)

// Update the local image metadata by walking through all found files and check if the modification date has changed
// or if they are new to the local database. If the files is new or changed, the md5sum will be rebuilt as well.
func synchronizeLocalImageMetadata(metadataStorage ImageMetadataProvider, fileSystemNodes map[string]*localFileStructure.FilesystemNode, categories map[string]*piwigo.PiwigoCategory, checksumCalculator fileChecksumCalculator) error {
	logrus.Debug("Starting synchronizeLocalImageMetadata")
	defer logrus.Debug("Leaving synchronizeLocalImageMetadata")

	logrus.Info("Synchronizing local image metadata database with local available images")

	err := synchronizeLocalImageMetadataScanNewFiles(fileSystemNodes, metadataStorage, categories, checksumCalculator)
	if err != nil {
		return err
	}

	err = synchronizeLocalImageMetadataFindFilesToDelete(metadataStorage)
	if err != nil {
		return err
	}
	return nil
}

func synchronizeLocalImageMetadataScanNewFiles(fileSystemNodes map[string]*localFileStructure.FilesystemNode, metadataStorage ImageMetadataProvider, categories map[string]*piwigo.PiwigoCategory, checksumCalculator fileChecksumCalculator) error {
	logrus.Debug("Entering synchronizeLocalImageMetadataScanNewFiles")
	defer logrus.Debug("Leaving synchronizeLocalImageMetadataScanNewFiles")

	for _, file := range fileSystemNodes {
		if file.IsDir {
			// we are only interested in files not directories
			continue
		}

		metadata, err := metadataStorage.ImageMetadata(file.Path)
		if err == ErrorRecordNotFound {
			logrus.Debugf("No metadata for %s found. Creating new entry.", file.Key)
			metadata = ImageMetaData{}
			metadata.Filename = file.Name
			metadata.FullImagePath = file.Path
			metadata.CategoryPath = filepath.Dir(file.Key)

			category, exist := categories[metadata.CategoryPath]
			if exist {
				metadata.CategoryId = category.Id
			} else {
				logrus.Warnf("No category found for image %s", file.Path)
			}

		} else if err != nil {
			logrus.Errorf("Could not get metadata due to trouble. Cancelling - %s", err)
			return err
		}

		metadata.UploadRequired = !metadata.LastChange.Equal(file.ModTime) || metadata.PiwigoId == 0
		metadata.DeleteRequired = false
		metadata.LastChange = file.ModTime
		metadata.Md5Sum, err = checksumCalculator(file.Path)
		if err != nil {
			logrus.Warnf("Could not calculate checksum for file %s. Skipping...", file.Path)
			continue
		}

		err = metadataStorage.SaveImageMetadata(metadata)
		if err != nil {
			return err
		}
	}
	return nil
}

func synchronizeLocalImageMetadataFindFilesToDelete(provider ImageMetadataProvider) error {
	logrus.Debug("Entering synchronizeLocalImageMetadataFindFilesToDelete")
	defer logrus.Debug("Leaving synchronizeLocalImageMetadataFindFilesToDelete")

	images, err := provider.ImageMetadataAll()
	if err != nil {
		return err
	}

	for _, img := range images {
		if _, err := os.Stat(img.FullImagePath); os.IsNotExist(err) {
			img.UploadRequired = false
			img.DeleteRequired = true
			err := provider.SaveImageMetadata(img)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Uploads the pending images to the piwigo gallery and assign the category of to the image.
// Update local metadata and set upload flag to false. Also updates the piwigo image id if there was a difference.
func uploadImages(piwigoCtx piwigo.PiwigoImageApi, metadataProvider ImageMetadataProvider) error {
	logrus.Debug("Starting uploadImages")
	defer logrus.Debug("Finished uploadImages successfully")

	images, err := metadataProvider.ImageMetadataToUpload()
	if err != nil {
		return err
	}

	logrus.Infof("Uploading %d images to piwigo", len(images))

	for _, img := range images {

		imgId, err := piwigoCtx.UploadImage(img.PiwigoId, img.FullImagePath, img.Md5Sum, img.CategoryId)
		if err != nil {
			logrus.Warnf("could not upload image %s. Continuing with the next image.", img.FullImagePath)
			continue
		}

		if imgId > 0 && imgId != img.PiwigoId {
			img.PiwigoId = imgId
			logrus.Debugf("Updating image %d with piwigo id %d", img.ImageId, img.PiwigoId)
		}

		logrus.Infof("Successfully uploaded %s", img.FullImagePath)

		img.UploadRequired = false
		err = metadataProvider.SaveImageMetadata(img)
		if err != nil {
			logrus.Warnf("could not save uploaded image %s. Continuing with the next image.", img.FullImagePath)
			continue
		}
	}

	return nil
}

func deleteImages(piwigoCtx piwigo.PiwigoImageApi, metadataProvider ImageMetadataProvider) error {
	logrus.Debug("Starting deleteImages")
	defer logrus.Debug("Finished deleteImages successfully")

	images, err := metadataProvider.ImageMetadataToDelete()
	if err != nil {
		return err
	}

	if len(images) == 0 {
		logrus.Info("There are not images scheduled for deletion.")
		return nil
	}

	logrus.Infof("Deleting %d images from piwigo", len(images))

	piwigoIds := []int{}
	for _, img := range images {
		if img.PiwigoId > 0 {
			logrus.Tracef("Adding %d to deletable list", img.PiwigoId)
			piwigoIds = append(piwigoIds, img.PiwigoId)
		}
	}

	if len(piwigoIds) > 0 {
		err = piwigoCtx.DeleteImages(piwigoIds)
		if err != nil {
			return err
		}
	} else {
		logrus.Info("Only local images found to delete. No call to piwigo is made.")
	}

	return metadataProvider.DeleteMarkedImages()
}

// This method aggregates the check for files with missing piwigoids and if changed files need to be uploaded again.
func synchronizePiwigoMetadata(piwigoCtx piwigo.PiwigoImageApi, metadataProvider ImageMetadataProvider) error {
	logrus.Debug("Entering synchronizePiwigoMetadata")
	defer logrus.Debug("Leaving synchronizePiwigoMetadata")

	// TODO: check if category has to be assigned (image possibly added to two albums -> only uploaded once but assigned multiple times) -> implement later
	err := updatePiwigoIdIfAlreadyUploaded(metadataProvider, piwigoCtx)
	if err != nil {
		return err
	}

	err = checkPiwigoForChangedImages(metadataProvider, piwigoCtx)
	if err != nil {
		return err
	}

	return nil
}

// Check all images with upload required if they are really changed and need to be uploaded to the server.
func checkPiwigoForChangedImages(provider ImageMetadataProvider, piwigoCtx piwigo.PiwigoImageApi) error {
	logrus.Info("Checking pending files if they really differ from the version in piwigo...")
	defer logrus.Info("Finished checking pending files if they really differ from the version in piwigo...")

	images, err := provider.ImageMetadataToUpload()
	if err != nil {
		return err
	}

	if len(images) == 0 {
		logrus.Info("There are no existing images to check for modification on the server.")
		return nil
	}

	for _, img := range images {
		if img.PiwigoId == 0 {
			continue
		}
		state, err := piwigoCtx.ImageCheckFile(img.PiwigoId, img.Md5Sum)
		if err != nil {
			logrus.Warnf("Error during file change check of file %s", img.FullImagePath)
			continue
		}

		if state == piwigo.ImageStateUptodate {
			logrus.Debugf("File %s - %d has not changed", img.FullImagePath, img.PiwigoId)
			img.UploadRequired = false
			err = provider.SaveImageMetadata(img)
			if err != nil {
				logrus.Warnf("Could not save image data of image %s", img.FullImagePath)
			}
		}
	}

	return nil
}

// This function calls piwigo and checks if the given md5sum is already present.
// Only files without a piwigo id are used to query the server.
func updatePiwigoIdIfAlreadyUploaded(provider ImageMetadataProvider, piwigoCtx piwigo.PiwigoImageApi) error {
	logrus.Info("checking for pending files that are already on piwigo and updating piwigoids...")
	defer logrus.Info("finshed checking for pending files that are already on piwigo and updating piwigoids...")

	images, err := provider.ImageMetadataToUpload()
	if err != nil {
		return err
	}

	if len(images) == 0 {
		logrus.Info("There are no existing images to check for modification on the server.")
		return nil
	}

	logrus.Debugln("Preparing lookuplist for missing piwigo ids...")
	files := make([]string, 0, len(images))
	for _, img := range images {
		if img.PiwigoId == 0 {
			files = append(files, img.Md5Sum)
		}
	}

	if len(files) == 0 {
		logrus.Info("There are no images without piwigo id to check for modification on the server.")
		return nil
	}

	missingResults, err := piwigoCtx.ImagesExistOnPiwigo(files)
	if err != nil {
		return err
	}
	for md5sum, piwigoId := range missingResults {
		logrus.Debugf("Setting piwigo id of %s to %d", md5sum, piwigoId)
		if piwigoId > 0 {
			err = provider.SavePiwigoIdAndUpdateUploadFlag(md5sum, piwigoId)
			if err != nil {
				logrus.Warnf("Could not save piwigo id %d for file %s", piwigoId, md5sum)
			}
		} else {
			logrus.Tracef("Image %s not found on server", md5sum)
		}
	}

	return nil
}
