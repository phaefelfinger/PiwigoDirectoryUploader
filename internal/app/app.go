package app

import (
	"flag"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/localFileStructure"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/piwigo"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	imagesRootPath            = flag.String("imagesRootPath", "", "This is the images root path that should be mirrored to piwigo.")
	sqliteDb                  = flag.String("sqliteDb", "", "The connection string to the sql lite database file.")
	piwigoUrl                 = flag.String("piwigoUrl", "", "The root url without tailing slash to your piwigo installation.")
	piwigoUser                = flag.String("piwigoUser", "", "The username to use during sync.")
	piwigoPassword            = flag.String("piwigoPassword", "", "This is password to the given username.")
	piwigoUploadChunkSizeInKB = flag.Int("piwigoUploadChunkSizeInKB", 512, "The chunksize used to upload an image to piwigo.")
)

func Run() {
	context, err := createAppContext()
	if err != nil {
		logErrorAndExit(err, 1)
	}

	err = context.piwigo.LoginToPiwigoAndConfigureContext()
	if err != nil {
		logErrorAndExit(err, 2)
	}

	filesystemNodes, err := localFileStructure.ScanLocalFileStructure(context.localRootPath)
	if err != nil {
		logErrorAndExit(err, 3)
	}

	categories, err := getAllCategoriesFromServer(context)
	if err != nil {
		logErrorAndExit(err, 4)
	}

	err = synchronizeCategories(context, filesystemNodes, categories)
	if err != nil {
		logErrorAndExit(err, 5)
	}

	err = synchronizeLocalImageMetadata(context.dataStore, filesystemNodes, localFileStructure.CalculateFileCheckSums)
	if err != nil {
		logErrorAndExit(err, 6)
	}

	err = synchronizeImages(context.piwigo, context.dataStore, categories)
	if err != nil {
		logErrorAndExit(err, 7)
	}

	_ = piwigo.Logout(context.piwigo)
}

func logErrorAndExit(err error, exitCode int) {
	logrus.Errorln(err)
	os.Exit(exitCode)
}
