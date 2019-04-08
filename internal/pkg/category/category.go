/*
 * Copyright (C) 2019 Philipp Haefelfinger (http://www.haefelfinger.ch/). All Rights Reserved.
 * This application is licensed under GPLv2. See the LICENSE file in the root directory of the project.
 */

package category

import (
	"errors"
	"fmt"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/datastore"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/localFileStructure"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/piwigo"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

func SynchronizeCategories(filesystemNodes map[string]*localFileStructure.FilesystemNode, piwigoApi piwigo.CategoryApi, db datastore.CategoryProvider) error {
	logrus.Debug("Entering SynchronizeCategories...")
	defer logrus.Debug("Leaving SynchronizeCategories...")

	err := updatePiwigoCategoriesFromServer(piwigoApi, db)
	if err != nil {
		return err
	}

	logrus.Infoln("Adding missing categories to local db...")
	err = addMissingPiwigoCategoriesToLocalDb(db, filesystemNodes)
	if err != nil {
		return err
	}

	return createMissingCategories(piwigoApi, db)
}

func addMissingPiwigoCategoriesToLocalDb(db datastore.CategoryProvider, fileSystemNodes map[string]*localFileStructure.FilesystemNode) error {
	logrus.Debug("Entering addMissingPiwigoCategoriesToLocalDb...")
	defer logrus.Debug("Leave addMissingPiwigoCategoriesToLocalDb...")

	for _, file := range fileSystemNodes {
		if !file.IsDir {
			logrus.Tracef("%s: Skipping as no directory", file.Key)
			continue
		}

		_, err := db.GetCategoryByKey(file.Key)
		if err == nil {
			logrus.Debugf("%s already exists.", file.Key)
			continue
		}
		if err != datastore.ErrorRecordNotFound {
			return err
		}

		logrus.Debugf("Creating missing category %s", file.Key)
		category := datastore.CategoryData{
			Key:            file.Key,
			Name:           file.Name,
			PiwigoParentId: 0,
			PiwigoId:       0,
		}

		err = db.SaveCategory(category)
		if err != nil {
			return err
		}
	}
	return nil
}

func updatePiwigoCategoriesFromServer(piwigoApi piwigo.CategoryApi, db datastore.CategoryProvider) error {
	logrus.Debug("Entering updatePiwigoCategoriesFromServer")
	defer logrus.Debug("Leaving updatePiwigoCategoriesFromServer")

	categories, err := piwigoApi.GetAllCategories()
	if err != nil {
		return err
	}

	for _, pwgCat := range categories {
		dbCat, err := db.GetCategoryByPiwigoId(pwgCat.Id)
		if err == datastore.ErrorRecordNotFound {
			logrus.Debugf("Adding category %s", pwgCat.Key)
			dbCat = datastore.CategoryData{
				PiwigoId: pwgCat.Id,
			}
		} else if err != nil {
			return err
		}

		if dbCat.Name == pwgCat.Name && dbCat.Key == pwgCat.Key && dbCat.PiwigoParentId == pwgCat.ParentId {
			logrus.Debugf("No changes for category %s", dbCat.Key)
			continue
		}

		dbCat.Name = pwgCat.Name
		dbCat.Key = pwgCat.Key
		dbCat.PiwigoParentId = pwgCat.ParentId

		err = db.SaveCategory(dbCat)
		if err != nil {
			return err
		}
	}

	return nil
}

func createMissingCategories(piwigoApi piwigo.CategoryApi, db datastore.CategoryProvider) error {
	logrus.Debug("Entering createMissingCategories...")
	defer logrus.Debug("Leaving createMissingCategories...")

	missingCategories, err := db.GetCategoriesToCreate()
	if err != nil {
		return err
	}

	if len(missingCategories) == 0 {
		logrus.Info("No categories missing on piwigo.")
		return nil
	}

	logrus.Infof("Creating %d categories", len(missingCategories))

	for _, category := range missingCategories {
		logrus.Infof("Creating category %s", category.Key)

		parentId, err := getParentId(category, db)
		if err != nil {
			return err
		}

		// create category on piwigo
		id, err := piwigoApi.CreateCategory(parentId, category.Name)
		if err != nil {
			return errors.New(fmt.Sprintf("Could not create category on piwigo: %s", err))
		}

		// update local category information
		category.PiwigoId = id
		category.PiwigoParentId = parentId

		err = db.SaveCategory(category)
		if err != nil {
			return err
		}
	}

	return nil
}

func getParentId(category datastore.CategoryData, db datastore.CategoryProvider) (int, error) {
	if category.Key == "" || category.Key == "." {
		msg := fmt.Sprintf("Category with id %d has a invalid value in the keyfield!", category.CategoryId)
		logrus.Warnf(msg)
		return 0, errors.New(msg)
	}

	parentKey := filepath.Dir(category.Key)
	if category.Name == parentKey || parentKey == "." || parentKey == "" {
		logrus.Debugf("The category %s is a root category, there is no parent", category.Name)
		return 0, nil
	}

	logrus.Debugf("Looking up parent with key %s", parentKey)
	parentCategory, err := db.GetCategoryByKey(parentKey)
	if err != nil {
		return 0, err
	}

	return parentCategory.PiwigoId, nil
}
