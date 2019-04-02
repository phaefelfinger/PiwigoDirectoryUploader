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

	piwigoCategory := createTestPiwigoCategory(1)
	piwigoCategories := []piwigo.PiwigoCategory{piwigoCategory}

	category := createTestCategoryData(1)

	dbmock := NewMockCategoryProvider(mockCtrl)
	dbmock.EXPECT().SaveCategory(category).Times(1)

	piwigoMock := NewMockPiwigoCategoryApi(mockCtrl)
	piwigoMock.EXPECT().GetAllCategories().Return(piwigoCategories).Times(1)

	err := SynchronizePiwigoCategories(piwigoMock, dbmock)
	if err != nil {
		t.Error(err)
	}
}

func createTestPiwigoCategory(piwigoId int) piwigo.PiwigoCategory {
	cat := piwigo.PiwigoCategory{
		Id:       piwigoId,
		Name:     "2019",
		Key:      "2019",
		ParentId: 0,
	}
	return cat
}

func createTestCategoryData(piwigoId int) datastore.CategoryData {
	cat := datastore.CategoryData{
		CategoryId:     1,
		PiwigoId:       piwigoId,
		PiwigoParentId: 0,
		Name:           "2019",
		Key:            "2019",
	}
	return cat
}
