package app

import (
	"errors"
	"github.com/sirupsen/logrus"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/localFileStructure"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/piwigo/category"
)

func synchronizeImages(fileSystem map[string]*localFileStructure.FilesystemNode, existingCategories map[string]*category.PiwigoCategory) error {
	findMissingImages()
	uploadImages()
	return errors.New("synchronizeImages: NOT IMPLEMENTED")
}

func findMissingImages() {
	logrus.Warnln("Finding missing images (NotImplemented)")
}

func uploadImages() {
	logrus.Warnln("Uploading missing images (NotImplemented)")
}
