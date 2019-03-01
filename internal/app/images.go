package app

import (
	"git.haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/localFileStructure"
	"git.haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/piwigo/category"
	"git.haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/piwigo/picture"
	"github.com/sirupsen/logrus"
	"sort"
)

func synchronizeImages(context *appContext, fileSystem map[string]*localFileStructure.FilesystemNode, existingCategories map[string]*category.PiwigoCategory) error {

	imageFiles, err := localFileStructure.GetImageList(fileSystem)
	if err != nil {
		return err
	}

	missingFiles, err := findMissingImages(context, imageFiles)
	if err != nil {
		return err
	}

	err = uploadImages(context, missingFiles, existingCategories)
	if err != nil {
		return err
	}

	logrus.Infof("Synchronized %d files.", len(missingFiles))

	return nil
}

func findMissingImages(context *appContext, imageFiles []*localFileStructure.ImageNode) ([]*localFileStructure.ImageNode, error) {

	logrus.Debugln("Preparing lookuplist for missing files...")

	files := make([]string, 0, len(imageFiles))
	md5map := make(map[string]*localFileStructure.ImageNode, len(imageFiles))
	for _, file := range imageFiles {
		md5map[file.Md5Sum] = file
		files = append(files, file.Md5Sum)
	}

	misingSums, err := picture.ImageUploadRequired(context.Piwigo, files)
	if err != nil {
		return nil, err
	}

	missingFiles := make([]*localFileStructure.ImageNode, 0, len(misingSums))
	for _, sum := range misingSums {
		file := md5map[sum]
		logrus.Infof("Found missing file %s", file.Path)
		missingFiles = append(missingFiles, file)
	}

	logrus.Infof("Found %d missing files", len(missingFiles))

	return missingFiles, nil
}

func uploadImages(context *appContext, missingFiles []*localFileStructure.ImageNode, existingCategories map[string]*category.PiwigoCategory) error {

	// We sort the files by path to populate per category and not random by file
	sort.Slice(missingFiles, func(i, j int) bool {
		return missingFiles[i].Path < missingFiles[j].Path
	})

	for _, file := range missingFiles {
		logrus.Infof("Uploading %s %s", file.CategoryName, file.Path)
		categoryId := existingCategories[file.CategoryName].Id

		imageId, err := picture.UploadImage(context.Piwigo, file.Path, file.Md5Sum, categoryId)
		if err != nil {
			return err
		}
		file.ImageId = imageId
	}

	return nil
}
