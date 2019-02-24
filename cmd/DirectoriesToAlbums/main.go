package main

import (
	"github.com/sirupsen/logrus"
	"github.com/vharitonsky/iniflags"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/app"
	"os"
)

func main() {
	iniflags.Parse()

	InitializeLog()

	app.Run()
}

func InitializeLog() {
	//TODO: make log configurable to file instead of console
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)
	logrus.Infoln("Starting Piwigo directories to albums...")
}
