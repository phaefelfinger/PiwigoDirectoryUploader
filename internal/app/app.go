package app

import (
	"flag"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/localFileStructure"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	imagesRootPath            = flag.String("imagesRootPath", "", "This is the images root path that should be mirrored to piwigo.")
	sqliteDb                  = flag.String("sqliteDb", "./localstate.db", "The connection string to the sql lite database file.")
	noUpload                  = flag.Bool("noUpload", false, "If set to true, the metadata gets prepared but the upload is not called and the application is exited with code 90")
	piwigoUrl                 = flag.String("piwigoUrl", "", "The root url without tailing slash to your piwigo installation.")
	piwigoUser                = flag.String("piwigoUser", "", "The username to use during sync.")
	piwigoPassword            = flag.String("piwigoPassword", "", "This is password to the given username.")
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

	categories, err := getAllCategoriesFromServer(context.piwigo)
	if err != nil {
		logErrorAndExit(err, 4)
	}

	err = synchronizeCategories(context.piwigo, filesystemNodes, categories)
	if err != nil {
		logErrorAndExit(err, 5)
	}

	err = synchronizeLocalImageMetadata(context.dataStore, filesystemNodes, categories, localFileStructure.CalculateFileCheckSums)
	if err != nil {
		logErrorAndExit(err, 6)
	}

	err = synchronizePiwigoMetadata(context.piwigo, context.dataStore)
	if err != nil {
		logErrorAndExit(err, 7)
	}

	if *noUpload {
		logrus.Warnln("Skipping upload of images as flag noUpload is set to true!")
		_ = context.piwigo.Logout()
		os.Exit(90)
	}

	err = uploadImages(context.piwigo, context.dataStore)
	if err != nil {
		logErrorAndExit(err, 8)
	}

	_ = context.piwigo.Logout()
}

func logErrorAndExit(err error, exitCode int) {
	logrus.Errorln(err)
	os.Exit(exitCode)
}
