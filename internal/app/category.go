package app

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/localFileStructure"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/piwigo/category"
	"path/filepath"
	"sort"
)

func getAllCategoriesFromServer(context *AppContext) (map[string]*category.PiwigoCategory, error) {
	logrus.Debugln("Starting GetAllCategories")
	categories, err := category.GetAllCategories(context.Piwigo)
	return categories, err
}

func synchronizeCategories(context *AppContext, filesystemNodes map[string]*localFileStructure.FilesystemNode, existingCategories map[string]*category.PiwigoCategory) error {
	logrus.Infoln("Synchronizing categories...")

	missingCategories := findMissingCategories(filesystemNodes, existingCategories)

	return createMissingCategories(context, missingCategories, existingCategories)
}

func findMissingCategories(fileSystem map[string]*localFileStructure.FilesystemNode, existingCategories map[string]*category.PiwigoCategory) []string {
	missingCategories := []string{}

	for _, file := range fileSystem {
		if !file.IsDir {
			continue
		}

		_, exists := existingCategories[file.Key]

		if exists {
			logrus.Debugf("Found existing category %s", file.Key)
		} else {
			logrus.Infof("Missing category detected %s", file.Key)
			missingCategories = append(missingCategories, file.Key)
		}
	}

	return missingCategories
}

func createMissingCategories(context *AppContext, missingCategories []string, existingCategories map[string]*category.PiwigoCategory) error {
	// we sort them to make sure the categories gets created
	// in the right order and we have the parent available while creating the children
	sort.Strings(missingCategories)

	for _, categoryKey := range missingCategories {
		logrus.Infof("Creating category %s", categoryKey)

		name, parentId, err := getNameAndParentId(categoryKey, existingCategories)
		if err != nil {
			return err
		}

		// create category on piwigo
		id, err := category.CreateCategory(context.Piwigo, parentId, name)
		if err != nil {
			return errors.New(fmt.Sprintf("Could not create category on piwigo: %s", err))
		}

		newCategory := category.PiwigoCategory{Id: id, Name: name, ParentId: parentId, Key: categoryKey}
		logrus.Println(newCategory)
		existingCategories[newCategory.Key] = &newCategory
	}

	return errors.New("createMissingCategories: NOT IMPLEMENTED")
}

func getNameAndParentId(category string, categories map[string]*category.PiwigoCategory) (string, int, error) {
	parentKey := filepath.Dir(category)
	_, name := filepath.Split(category)
	if name == category {
		logrus.Debugf("The category %s is a root category, there is no parent", name)
		return name, 0, nil
	}

	logrus.Debugf("Looking up parent with key %s", parentKey)
	parent, exists := categories[parentKey]
	if !exists {
		return "", 0, errors.New(fmt.Sprintf("could not find parent with key %s", parentKey))
	}

	parentId := parent.Id
	logrus.Debugf("Found parent %s with id %d", parentKey, parentId)

	return name, parentId, nil
}
