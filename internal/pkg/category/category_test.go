/*
 * Copyright (C) 2019 Philipp Haefelfinger (http://www.haefelfinger.ch/). All Rights Reserved.
 * This application is licensed under GPLv2. See the LICENSE file in the root directory of the project.
 */

package category

import (
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/datastore"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/piwigo"
	"github.com/golang/mock/gomock"
	"testing"
)

//go:generate mockgen -destination=./piwigo_mock_test.go -package=category git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/piwigo PiwigoApi,PiwigoCategoryApi,PiwigoImageApi
//go:generate mockgen -destination=./datastore_mock_test.go -package=category git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/datastore CategoryProvider

func Test_SynchronizePiwigoCategories_adds_new_categories(t *testing.T) {
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

	err := SynchronizePiwigoCategories(piwigoMock, dbmock)
	if err != nil {
		t.Error(err)
	}
}

func Test_SynchronizePiwigoCategories_updates_a_category(t *testing.T) {
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

	expectedCategory := datastore.CategoryData{
		PiwigoId:       1,
		PiwigoParentId: 0,
		CategoryId:     1,
		Name:           "2019",
		Key:            "2019",
	}

	dbmock := NewMockCategoryProvider(mockCtrl)
	dbmock.EXPECT().GetCategoryByPiwigoId(1).Return(oldCategory, nil).Times(1)
	dbmock.EXPECT().GetCategoryByPiwigoId(gomock.Any()).Return(datastore.CategoryData{}, datastore.ErrorRecordNotFound).Times(1)
	dbmock.EXPECT().SaveCategory(expectedCategory).Times(1)
	dbmock.EXPECT().SaveCategory(gomock.Any()).Times(1)

	piwigoMock := NewMockPiwigoCategoryApi(mockCtrl)
	piwigoMock.EXPECT().GetAllCategories().Return(piwigoCategories, nil).Times(1)

	err := SynchronizePiwigoCategories(piwigoMock, dbmock)
	if err != nil {
		t.Error(err)
	}
}

func createDbCategoriesFrom(categories map[string]*piwigo.PiwigoCategory) []datastore.CategoryData {
	dbCategories := []datastore.CategoryData{}
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
