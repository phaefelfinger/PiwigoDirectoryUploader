package app

import "github.com/sirupsen/logrus"

func synchronizeImages() error {
	findMissingImages()
	uploadImages()
	return nil
}

func findMissingImages() {
	logrus.Warnln("Finding missing images (NotImplemented)")
}

func uploadImages() {
	logrus.Warnln("Uploading missing images (NotImplemented)")
}
