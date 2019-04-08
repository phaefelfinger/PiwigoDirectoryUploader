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

func DeleteImages(piwigoCtx piwigo.ImageApi, metadataProvider datastore.ImageMetadataProvider) error {
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

	var piwigoIds []int = nil
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
