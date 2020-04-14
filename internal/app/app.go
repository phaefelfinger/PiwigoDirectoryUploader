/*
 * Copyright (C) 2020 Philipp Haefelfinger (http://www.haefelfinger.ch/). All Rights Reserved.
 * This application is licensed under GPLv2. See the LICENSE file in the root directory of the project.
 */

package app

import (
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/category"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/images"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/localFileStructure"
	"github.com/sirupsen/logrus"
	"os"
)

func Run() {
	initializeFlags()
	initializeLog()

	context, err := newAppContext()
	if err != nil {
		logErrorAndExit(err, 1)
	}

	err = context.piwigo.Login()
	if err != nil {
		logErrorAndExit(err, 2)
	}

	filesystemNodes, err := localFileStructure.ScanLocalFileStructure(context.localRootPath, extensions, ignoreDirs, *dirSuffixToSkip)
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

func initializeLog() {
	level, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		level = logrus.DebugLevel
	}
	logrus.SetLevel(level)

	logrus.SetOutput(os.Stdout)

	logrus.Infoln("Starting Piwigo directories to albums...")
}

func logErrorAndExit(err error, exitCode int) {
	logrus.Errorln(err)
	os.Exit(exitCode)
}
