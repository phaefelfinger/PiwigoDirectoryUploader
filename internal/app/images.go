package app

import (
	"errors"
	"github.com/sirupsen/logrus"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/localFileStructure"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/piwigo/category"
)

func synchronizeImages(context *appContext, fileSystem map[string]*localFileStructure.FilesystemNode, existingCategories map[string]*category.PiwigoCategory) error {

	imageFiles, err := localFileStructure.GetImageList(fileSystem)
	if err != nil {
		return err
	}

	missingFiles := findMissingImages(imageFiles)
	uploadImages(missingFiles)

	return errors.New("synchronizeImages: NOT IMPLEMENTED")
}

func findMissingImages(imageFiles []*localFileStructure.ImageNode) []string {

	logrus.Warnln("Finding missing images (NotImplemented)")

	return nil
}

func uploadImages(missingFiles []string) {
	logrus.Warnln("Uploading missing images (NotImplemented)")
}
