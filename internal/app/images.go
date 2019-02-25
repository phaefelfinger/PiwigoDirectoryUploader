package app

import (
	"errors"
	"github.com/sirupsen/logrus"
)

func synchronizeImages() error {
	findMissingImages()
	uploadImages()
	return errors.New("NOT IMPLEMENTED")
}

func findMissingImages() {
	logrus.Warnln("Finding missing images (NotImplemented)")
}

func uploadImages() {
	logrus.Warnln("Uploading missing images (NotImplemented)")
}
