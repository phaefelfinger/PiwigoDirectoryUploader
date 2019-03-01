package main

import (
	"git.haefelfinger.net/piwigo/DirectoriesToAlbums/internal/app"
	"github.com/sirupsen/logrus"
	"github.com/vharitonsky/iniflags"
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
