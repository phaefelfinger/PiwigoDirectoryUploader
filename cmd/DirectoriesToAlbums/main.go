package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/app/DirectoriesToAlbums"
	"os"
)

func main() {
	flag.Parse()
	root := flag.Arg(0)

	InitializeLog()

	app.AuthenticateToPiwigo()
	app.ScanLocalDirectories(root)
	app.GetAllCategoriesFromServer()
	app.FindMissingAlbums()
	app.CreateMissingAlbums()
	app.FindMissingImages()
	app.UploadImages()

}

func InitializeLog() {
	//TODO: make log configurable to file instead of console
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)
	logrus.Infoln("Starting Piwigo directories to albums...")
}