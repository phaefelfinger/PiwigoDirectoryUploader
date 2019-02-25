package app

import (
	"errors"
	"github.com/sirupsen/logrus"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/localFileStructure"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/piwigo/category"
	"sort"
)

func getAllCategoriesFromServer(context *AppContext) (map[string]*category.PiwigoCategory, error) {
	logrus.Debugln("Starting GetAllCategories")
	categories, err := category.GetAllCategories(context.Piwigo)
	return categories, err
}

func synchronizeCategories(filesystemNodes map[string]*localFileStructure.FilesystemNode, existingCategories map[string]*category.PiwigoCategory) error {
	logrus.Infoln("Synchronizing categories...")

	missingCategories := findMissingCategories(filesystemNodes, existingCategories)

	return createMissingCategories(missingCategories, existingCategories)
}

func findMissingCategories(fileSystem map[string]*localFileStructure.FilesystemNode, existingCategories map[string]*category.PiwigoCategory) []string {
	missingCategories := []string{}

	for _, file := range fileSystem {
		if !file.IsDir {
			continue
		}

		_, exists := existingCategories[file.Key]

		if !exists {
			logrus.Infof("Missing category detected %s", file.Key)
			missingCategories = append(missingCategories, file.Key)
		} else {
			logrus.Debugf("Found existing category %s", file.Key)
		}
	}

	return missingCategories
}

func createMissingCategories(missingCategories []string, existingCategories map[string]*category.PiwigoCategory) error {

	// we sort them to make sure the categories gets created
	// in the right order and we have the parent available while creating the children
	sort.Strings(missingCategories)

	for _, c := range missingCategories {
		logrus.Infof("Creating category %s",c)

		// create category on piwigo

		// build new map entry

		// get parent entry by path
		// set parent entry id

		// calculate new map key
		// add to existing map

	}

	return errors.New("NOT IMPLEMENTED")
}
