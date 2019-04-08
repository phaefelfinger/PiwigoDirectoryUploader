/*
 * Copyright (C) 2019 Philipp Haefelfinger (http://www.haefelfinger.ch/). All Rights Reserved.
 * This application is licensed under GPLv2. See the LICENSE file in the root directory of the project.
 */

package app

import (
	"flag"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/category"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/images"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/localFileStructure"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	imagesRootPath  = flag.String("imagesRootPath", "", "This is the images root path that should be mirrored to piwigo.")
	sqliteDb        = flag.String("sqliteDb", "./localstate.db", "The connection string to the sql lite database file.")
	noUpload        = flag.Bool("noUpload", false, "If set to true, the metadata gets prepared but the upload is not called and the application is exited with code 90")
	piwigoUrl       = flag.String("piwigoUrl", "", "The root url without tailing slash to your piwigo installation.")
	piwigoUser      = flag.String("piwigoUser", "", "The username to use during sync.")
	piwigoPassword  = flag.String("piwigoPassword", "", "This is password to the given username.")
	removeImages    = flag.Bool("removeImages", false, "If set to true, images scheduled to delete will be removed from the piwigo server. Be sure you want to delete images before enabling this flag.")
	parallelUploads = flag.Int("parallelUploads", 4, "Set the number of images that get uploaded in parallel.")
)

func Run() {
	context, err := newAppContext()
	if err != nil {
		logErrorAndExit(err, 1)
	}

	err = context.piwigo.Login()
	if err != nil {
		logErrorAndExit(err, 2)
	}

	filesystemNodes, err := localFileStructure.ScanLocalFileStructure(context.localRootPath)
	if err != nil {
		logErrorAndExit(err, 3)
	}

	err = category.SynchronizeCategories(filesystemNodes, context.piwigo, context.dataStore)
	if err != nil {
		logErrorAndExit(err, 4)
	}

	err = images.SynchronizeLocalImageMetadata(context.dataStore, context.dataStore, filesystemNodes, localFileStructure.CalculateFileCheckSums)
	if err != nil {
		logErrorAndExit(err, 5)
	}

	err = images.SynchronizePiwigoMetadata(context.piwigo, context.dataStore)
	if err != nil {
		logErrorAndExit(err, 6)
	}

	if *removeImages {
		err = images.DeleteImages(context.piwigo, context.dataStore)
		if err != nil {
			logErrorAndExit(err, 7)
		}
	} else {
		logrus.Info("The flag removeImages is disabled. Skipping...")
	}

	if !(*noUpload) {
		err = images.UploadImages(context.piwigo, context.dataStore, *parallelUploads)
		if err != nil {
			logErrorAndExit(err, 8)
		}
	} else {
		logrus.Warnln("Skipping upload of images as flag noUpload is set to true!")
	}

	_ = context.piwigo.Logout()
}

func logErrorAndExit(err error, exitCode int) {
	logrus.Errorln(err)
	os.Exit(exitCode)
}
