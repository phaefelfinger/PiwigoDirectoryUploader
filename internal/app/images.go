/*
 * Copyright (C) 2019 Philipp Haefelfinger (http://www.haefelfinger.ch/). All Rights Reserved.
 * This application is licensed under GPLv2. See the LICENSE file in the root directory of the project.
 */

package app

import (
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/localFileStructure"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/piwigo"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

// to make use of the new local data store, we have to rethink and refactor the whole local detection process
// extend the storage of the images to keep track of upload state
// TBD: How to deal with updates -> delete / upload all based on md5 sums

type fileChecksumCalculator func(filePath string) (string, error)

// Update the local image metadata by walking through all found files and check if the modification date has changed
// or if they are new to the local database. If the files is new or changed, the md5sum will be rebuilt as well.
func synchronizeLocalImageMetadata(metadataStorage ImageMetadataProvider, fileSystemNodes map[string]*localFileStructure.FilesystemNode, categories map[string]*piwigo.PiwigoCategory, checksumCalculator fileChecksumCalculator) error {
	logrus.Debugf("Starting synchronizeLocalImageMetadata")
	logrus.Info("Synchronizing local image metadata database with local available images")

	err := synchronizeLocalImageMetadataScanNewFiles(fileSystemNodes, metadataStorage, categories, checksumCalculator)
	if err != nil {
		return err
	}

	logrus.Debugf("Finished synchronizeLocalImageMetadata")
	return nil
}

func synchronizeLocalImageMetadataScanNewFiles(fileSystemNodes map[string]*localFileStructure.FilesystemNode, metadataStorage ImageMetadataProvider, categories map[string]*piwigo.PiwigoCategory, checksumCalculator fileChecksumCalculator) error {
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

		if metadata.LastChange.Equal(file.ModTime) {
			logrus.Infof("No changed detected on file %s -> keeping current state", file.Key)
			continue
		}

		metadata.LastChange = file.ModTime
		metadata.UploadRequired = true
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

// Uploads the pending images to the piwigo gallery and assign the category of to the image.
// Update local metadata and set upload flag to false. Also updates the piwigo image id if there was a difference.
func uploadImages(piwigoCtx piwigo.PiwigoImageApi, metadataProvider ImageMetadataProvider) error {
	logrus.Debugf("Starting uploadImages")

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

	logrus.Debugf("Finished uploadImages successfully")
	return nil
}

// This method aggregates the check for files with missing piwigoids and if changed files need to be uploaded again.
func synchronizePiwigoMetadata(piwigoCtx piwigo.PiwigoImageApi, metadataProvider ImageMetadataProvider) error {
	// TODO: check if category has to be assigned (image possibly added to two albums -> only uploaded once but assigned multiple times) -> implement later
	logrus.Debugf("Starting synchronizePiwigoMetadata")
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
	logrus.Infof("Checking pending files if they really differ from the version in piwigo...")

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
	logrus.Infof("checking for pending files that are already on piwigo and updating piwigoids...")
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
		err = provider.SavePiwigoIdAndUpdateUploadFlag(md5sum, piwigoId)
		if err != nil {
			logrus.Warnf("Could not save piwigo id %d for file %s", piwigoId, md5sum)
		}
	}

	return nil
}
