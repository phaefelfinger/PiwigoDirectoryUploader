/*
 * Copyright (C) 2019 Philipp Haefelfinger (http://www.haefelfinger.ch/). All Rights Reserved.
 * This application is licensed under GPLv2. See the LICENSE file in the root directory of the project.
 */

package category

import (
	"fmt"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/datastore"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/localFileStructure"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/piwigo"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

//go:generate mockgen -destination=./piwigo_mock_test.go -package=category git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/piwigo PiwigoApi,PiwigoCategoryApi,PiwigoImageApi
//go:generate mockgen -destination=./datastore_mock_test.go -package=category git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/datastore CategoryProvider

func Test_updatePiwigoCategoriesFromServer_adds_new_categories(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	piwigoCategories := createTwoServerCategories()
	dbCategories := createDbCategoriesFrom(piwigoCategories)

	dbmock := NewMockCategoryProvider(mockCtrl)
	dbmock.EXPECT().GetCategoryByPiwigoId(gomock.Any()).Return(datastore.CategoryData{}, datastore.ErrorRecordNotFound).Times(len(piwigoCategories))
	for _, cat := range dbCategories {
		dbmock.EXPECT().SaveCategory(cat).Times(1)
	}

	piwigoMock := NewMockPiwigoCategoryApi(mockCtrl)
	piwigoMock.EXPECT().GetAllCategories().Return(piwigoCategories, nil).Times(1)

	err := updatePiwigoCategoriesFromServer(piwigoMock, dbmock)
	if err != nil {
		t.Error(err)
	}
}

func Test_updatePiwigoCategoriesFromServer_updates_a_category(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	piwigoCategories := createTwoServerCategories()

	oldCategory := datastore.CategoryData{
		PiwigoId:       1,
		PiwigoParentId: 0,
		CategoryId:     1,
		Key:            "oldKey",
		Name:           "oldName",
	}

	expectedCategory := createDbRootCategory()

	dbmock := NewMockCategoryProvider(mockCtrl)
	dbmock.EXPECT().GetCategoryByPiwigoId(1).Return(oldCategory, nil).Times(1)
	dbmock.EXPECT().GetCategoryByPiwigoId(gomock.Any()).Return(datastore.CategoryData{}, datastore.ErrorRecordNotFound).Times(1)
	dbmock.EXPECT().SaveCategory(expectedCategory).Times(1)
	dbmock.EXPECT().SaveCategory(gomock.Any()).Times(1)

	piwigoMock := NewMockPiwigoCategoryApi(mockCtrl)
	piwigoMock.EXPECT().GetAllCategories().Return(piwigoCategories, nil).Times(1)

	err := updatePiwigoCategoriesFromServer(piwigoMock, dbmock)
	if err != nil {
		t.Error(err)
	}
}

func Test_createMissingCategories_does_not_call_piwigo_if_there_is_no_category_missing(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	var categoriesToCreate []datastore.CategoryData

	dbmock := NewMockCategoryProvider(mockCtrl)
	dbmock.EXPECT().GetCategoriesToCreate().Return(categoriesToCreate, nil).Times(1)

	piwigoMock := NewMockPiwigoCategoryApi(mockCtrl)
	piwigoMock.EXPECT().CreateCategory(gomock.Any(), gomock.Any()).Times(0)

	err := createMissingCategories(piwigoMock, dbmock)
	if err != nil {
		t.Error(err)
	}
}

func Test_addMissingPiwigoCategoriesToLocalDb_creates_category_in_database(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectedCategory := createDbRootCategory()
	expectedCategory.PiwigoParentId = 0
	expectedCategory.PiwigoId = 0
	expectedCategory.CategoryId = 0

	fileNode := &localFileStructure.FilesystemNode{
		Name:    expectedCategory.Name,
		Key:     expectedCategory.Key,
		Path:    fmt.Sprintf("/home/nonexisting/%s", expectedCategory.Name),
		ModTime: time.Now(),
		IsDir:   true,
	}

	fileSystemNodes := make(map[string]*localFileStructure.FilesystemNode)
	fileSystemNodes[fileNode.Key] = fileNode

	dbmock := NewMockCategoryProvider(mockCtrl)
	dbmock.EXPECT().GetCategoryByKey(fileNode.Key).Return(datastore.CategoryData{}, datastore.ErrorRecordNotFound).Times(1)
	dbmock.EXPECT().SaveCategory(expectedCategory).Return(nil).Times(1)

	err := addMissingPiwigoCategoriesToLocalDb(dbmock, fileSystemNodes)
	if err != nil {
		t.Error(err)
	}
}

func Test_addMissingPiwigoCategoriesToLocalDb_does_nothing_already_in_db(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	fileNode := &localFileStructure.FilesystemNode{
		Name:    "dir",
		Key:     "dir",
		Path:    "/home/nonexisting/dir",
		ModTime: time.Now(),
		IsDir:   true,
	}

	fileSystemNodes := make(map[string]*localFileStructure.FilesystemNode)
	fileSystemNodes[fileNode.Key] = fileNode

	dbmock := NewMockCategoryProvider(mockCtrl)
	dbmock.EXPECT().GetCategoryByKey(fileNode.Key).Return(datastore.CategoryData{}, nil).Times(1)

	err := addMissingPiwigoCategoriesToLocalDb(dbmock, fileSystemNodes)
	if err != nil {
		t.Error(err)
	}
}

func Test_addMissingPiwigoCategoriesToLocalDb_does_nothing_if_list_contains_only_files(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	fileNode := &localFileStructure.FilesystemNode{
		Name:    "file.jpg",
		Key:     "file",
		Path:    "/home/nonexisting/file.jpg",
		ModTime: time.Now(),
		IsDir:   false,
	}

	fileSystemNodes := make(map[string]*localFileStructure.FilesystemNode)
	fileSystemNodes[fileNode.Key] = fileNode

	dbmock := NewMockCategoryProvider(mockCtrl)

	err := addMissingPiwigoCategoriesToLocalDb(dbmock, fileSystemNodes)
	if err != nil {
		t.Error(err)
	}
}

func Test_addMissingPiwigoCategoriesToLocalDb_does_nothing_if_list_is_empty(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	var fileSystemNodes map[string]*localFileStructure.FilesystemNode

	dbmock := NewMockCategoryProvider(mockCtrl)

	err := addMissingPiwigoCategoriesToLocalDb(dbmock, fileSystemNodes)
	if err != nil {
		t.Error(err)
	}
}

func Test_createMissingCategories_calls_piwigo_api_and_saves_returned_id(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectedCategory := createDbRootCategory()
	category := createDbRootCategory()
	category.PiwigoId = 0

	categoriesToCreate := []datastore.CategoryData{category}

	dbmock := NewMockCategoryProvider(mockCtrl)
	dbmock.EXPECT().GetCategoriesToCreate().Return(categoriesToCreate, nil).Times(1)
	dbmock.EXPECT().SaveCategory(expectedCategory).Return(nil).Times(1)

	piwigoMock := NewMockPiwigoCategoryApi(mockCtrl)
	piwigoMock.EXPECT().CreateCategory(0, category.Name).Return(1, nil).Times(1)

	err := createMissingCategories(piwigoMock, dbmock)
	if err != nil {
		t.Error(err)
	}
}

func Test_getParentId_returns_0_for_root_nodes(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	category := createDbRootCategory()

	dbmock := NewMockCategoryProvider(mockCtrl)
	dbmock.EXPECT().GetCategoryByKey(gomock.Any()).Times(0)

	parentId, err := getParentId(category, dbmock)
	if err != nil {
		t.Error(err)
	}

	if parentId != 0 {
		t.Errorf("Found parent id %d but expected 0", parentId)
	}
}

func Test_getParentId_returns_error_if_parentkey_is_not_found(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	category := createDbSubCategory()

	dbmock := NewMockCategoryProvider(mockCtrl)
	dbmock.EXPECT().GetCategoryByKey(gomock.Any()).Return(datastore.CategoryData{}, datastore.ErrorRecordNotFound).Times(1)

	parentId, err := getParentId(category, dbmock)
	if err == nil {
		t.Error("There should an error be returned if category key value is not valid!")
	}
	if parentId != 0 {
		t.Errorf("Found parent id %d but expected 0", parentId)
	}
}

func Test_getParentId_returns_error_if_key_invalid(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	category := createDbRootCategory()
	category.Key = "."

	dbmock := NewMockCategoryProvider(mockCtrl)
	dbmock.EXPECT().GetCategoryByKey(gomock.Any()).Times(0)

	parentId, err := getParentId(category, dbmock)
	if err == nil {
		t.Error("There should an error be returned if category key value is not valid!")
	}
	if parentId != 0 {
		t.Errorf("Found parent id %d but expected 0", parentId)
	}
}

func Test_getParentId_finds_the_exptected_parent_id(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	parentCategory := createDbRootCategory()
	category := createDbSubCategory()

	dbmock := NewMockCategoryProvider(mockCtrl)
	dbmock.EXPECT().GetCategoryByKey("2019").Return(parentCategory, nil).Times(1)

	parentId, err := getParentId(category, dbmock)
	if err != nil {
		t.Error(err)
	}

	if parentId != 1 {
		t.Errorf("Found parent id %d but expected 1", parentId)
	}
}

func createDbRootCategory() datastore.CategoryData {
	parentCategory := datastore.CategoryData{
		PiwigoId:       1,
		PiwigoParentId: 0,
		CategoryId:     1,
		Key:            "2019",
		Name:           "2019",
	}
	return parentCategory
}

func createDbSubCategory() datastore.CategoryData {
	category := datastore.CategoryData{
		PiwigoId:       2,
		PiwigoParentId: 0,
		CategoryId:     2,
		Key:            "2019/testalbumb",
		Name:           "testalbumb",
	}
	return category
}

func createDbCategoriesFrom(categories map[string]*piwigo.PiwigoCategory) []datastore.CategoryData {
	var dbCategories []datastore.CategoryData
	for _, cat := range categories {
		dbCat := datastore.CategoryData{
			PiwigoId:       cat.Id,
			PiwigoParentId: cat.ParentId,
			Key:            cat.Key,
			Name:           cat.Name,
			CategoryId:     0,
		}
		dbCategories = append(dbCategories, dbCat)
	}
	return dbCategories
}

func createTwoServerCategories() map[string]*piwigo.PiwigoCategory {
	piwigoCategory1 := piwigo.PiwigoCategory{
		Id:       1,
		Name:     "2019",
		Key:      "2019",
		ParentId: 0,
	}
	piwigoCategory2 := piwigo.PiwigoCategory{
		Id:       2,
		Name:     "SubCategory",
		Key:      "2019/SubCategory",
		ParentId: 1,
	}
	piwigoCategories := make(map[string]*piwigo.PiwigoCategory)
	piwigoCategories[piwigoCategory1.Key] = &piwigoCategory1
	piwigoCategories[piwigoCategory2.Key] = &piwigoCategory2
	return piwigoCategories
}
