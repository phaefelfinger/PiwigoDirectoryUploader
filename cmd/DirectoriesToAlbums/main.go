package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/app"
	"os"
)

func main() {
	flag.Parse()
	rootPath := flag.Arg(0)

	InitializeLog()

	app.Run(rootPath)
}

func InitializeLog() {
	//TODO: make log configurable to file instead of console
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)
	logrus.Infoln("Starting Piwigo directories to albums...")
}
