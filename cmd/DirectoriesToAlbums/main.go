package main

import (
	"flag"
	"git.haefelfinger.net/piwigo/DirectoriesToAlbums/internal/app"
	"github.com/sirupsen/logrus"
	"github.com/vharitonsky/iniflags"
	"os"
)

var (
	logLevel = flag.String("logLevel", "info", "The minimum log level required to write out a log message. (panic,fatal,error,warn,info,debug,trace)")
)

func main() {
	iniflags.Parse()
	initializeLog()
	app.Run()
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
