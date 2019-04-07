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

// Uploads the pending images to the piwigo gallery and assign the category of to the image.
// Update local metadata and set upload flag to false. Also updates the piwigo image id if there was a difference.
func UploadImages(piwigoCtx piwigo.PiwigoImageApi, metadataProvider datastore.ImageMetadataProvider) error {
	logrus.Debug("Starting uploadImages")
	defer logrus.Debug("Finished uploadImages successfully")

	images, err := metadataProvider.ImageMetadataToUpload()
	if err != nil {
		return err
	}

	logrus.Infof("Uploading %d images to piwigo", len(images))

	for _, img := range images {

		imgId, err := piwigoCtx.UploadImage(img.PiwigoId, img.FullImagePath, img.Md5Sum, img.CategoryPiwigoId)
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
