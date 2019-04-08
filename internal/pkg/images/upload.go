/*
 * Copyright (C) 2019 Philipp Haefelfinger (http://www.haefelfinger.ch/). All Rights Reserved.
 * This application is licensed under GPLv2. See the LICENSE file in the root directory of the project.
 */

package images

import (
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/datastore"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/piwigo"
	"github.com/sirupsen/logrus"
	"sync"
)

// Uploads the pending images to the piwigo gallery and assign the category of to the image.
// Update local metadata and set upload flag to false. Also updates the piwigo image id if there was a difference.
func UploadImages(piwigoCtx piwigo.PiwigoImageApi, metadataProvider datastore.ImageMetadataProvider, numberOfWorkers int) error {
	logrus.Debug("Starting uploadImages")
	defer logrus.Debug("Finished uploadImages successfully")

	images, err := metadataProvider.ImageMetadataToUpload()
	if err != nil {
		return err
	}

	if len(images) == 0 {
		logrus.Info("No images to upload.")
		return nil
	}

	if numberOfWorkers <= 0 {
		logrus.Warnf("Invalid numbers of worker set: %d falling back to default of 4", numberOfWorkers)
		numberOfWorkers = 4
	}

	logrus.Infof("Uploading %d images to piwigo using %d workers", len(images), numberOfWorkers)
	workQueue := make(chan datastore.ImageMetaData, numberOfWorkers)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go uploadQueueProducer(images, workQueue, &wg)

	for i := 0; i < numberOfWorkers; i++ {
		logrus.Debugf("Starting image upload worker %d", i)
		wg.Add(1)
		go uploadQueueWorker(workQueue, piwigoCtx, metadataProvider, &wg)
	}

	wg.Wait()
	return nil
}

func uploadQueueWorker(workQueue <-chan datastore.ImageMetaData, piwigoCtx piwigo.PiwigoImageApi, metadataProvider datastore.ImageMetadataProvider, waitGroup *sync.WaitGroup) {
	for img := range workQueue {
		logrus.Debugf("%s: uploading image to piwigo", img.FullImagePath)

		imgId, err := piwigoCtx.UploadImage(img.PiwigoId, img.FullImagePath, img.Md5Sum, img.CategoryPiwigoId)
		if err != nil {
			logrus.Warnf("%s: could not upload image. Continuing with the next image.", img.FullImagePath)
			continue
		}

		if imgId > 0 && imgId != img.PiwigoId {
			img.PiwigoId = imgId
			logrus.Debugf("%s: Updating image %d with piwigo id %d", img.FullImagePath, img.ImageId, img.PiwigoId)
		}
		logrus.Infof("%s: Successfully uploaded", img.FullImagePath)

		img.UploadRequired = false
		err = metadataProvider.SaveImageMetadata(img)
		if err != nil {
			logrus.Warnf("%s: could not save uploaded image. Continuing with the next image.", img.FullImagePath)
			continue
		}
	}
	waitGroup.Done()
}

func uploadQueueProducer(imagesToUpload []datastore.ImageMetaData, workQueue chan<- datastore.ImageMetaData, waitGroup *sync.WaitGroup) {
	for _, img := range imagesToUpload {
		logrus.Debugf("%s: Adding image to queue", img.FullImagePath)
		workQueue <- img
	}
	waitGroup.Done()
	close(workQueue)
}
