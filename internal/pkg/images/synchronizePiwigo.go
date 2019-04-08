/*
 * Copyright (C) 2019 Philipp Haefelfinger (http://www.haefelfinger.ch/). All Rights Reserved.
 * This application is licensed under GPLv2. See the LICENSE file in the root directory of the project.
 */

package images

import (
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/datastore"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/piwigo"
	"github.com/sirupsen/logrus"
)

// This method aggregates the check for files with missing piwigoids and if changed files need to be uploaded again.
func SynchronizePiwigoMetadata(piwigoCtx piwigo.ImageApi, metadataProvider datastore.ImageMetadataProvider) error {
	logrus.Debug("Entering SynchronizePiwigoMetadata")
	defer logrus.Debug("Leaving SynchronizePiwigoMetadata")

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

// This function calls piwigo and checks if the given md5sum is already present.
// Only files without a piwigo id are used to query the server.
func updatePiwigoIdIfAlreadyUploaded(provider datastore.ImageMetadataProvider, piwigoCtx piwigo.ImageApi) error {
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
		if piwigoId > 0 {
			logrus.Debugf("Setting piwigo id of %s to %d", md5sum, piwigoId)
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

// Check all images with upload required if they are really changed and need to be uploaded to the server.
func checkPiwigoForChangedImages(provider datastore.ImageMetadataProvider, piwigoCtx piwigo.ImageApi) error {
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
