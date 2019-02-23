package app

import (
	"github.com/sirupsen/logrus"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/localFileStructure"
)

func AuthenticateToPiwigo() {
	logrus.Warnln("Authenticating to piwigo server (NotImplemented)")
}

func ScanLocalDirectories(root string) {
	fileNodes := localFileStructure.ScanLocalFileStructure(root)
	logrus.Debugln("filepath.Walk() returned %v\n", fileNodes)
}

func GetAllCategoriesFromServer()  {
	// get all categories from server and flatten structure to match directory names
	// 2018/2018 album blah
	logrus.Warnln("Loading all categories from the server (NotImplemented)")
}

func FindMissingAlbums()  {
	logrus.Warnln("Looking up missing albums (NotImplemented)")
}

func CreateMissingAlbums()  {
	logrus.Warnln("Creating missing albums (NotImplemented)")
}

func FindMissingImages()  {
	logrus.Warnln("Finding missing images (NotImplemented)")
}

func UploadImages()  {
	logrus.Warnln("Uploading missing images (NotImplemented)")
}