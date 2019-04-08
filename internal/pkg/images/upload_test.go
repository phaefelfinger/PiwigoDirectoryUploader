/*
 * Copyright (C) 2019 Philipp Haefelfinger (http://www.haefelfinger.ch/). All Rights Reserved.
 * This application is licensed under GPLv2. See the LICENSE file in the root directory of the project.
 */

package images

import (
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/datastore"
	"github.com/golang/mock/gomock"
	"testing"
)

func Test_uploadImages_saves_new_id_to_db(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	img := createTestImageMetaData(0)
	images := []datastore.ImageMetaData{img}

	imgToSave := img
	imgToSave.PiwigoId = 5
	imgToSave.UploadRequired = false

	dbmock := NewMockImageMetadataProvider(mockCtrl)
	dbmock.EXPECT().ImageMetadataToUpload().Times(1).Return(images, nil)
	dbmock.EXPECT().SaveImageMetadata(imgToSave).Times(1)

	piwigomock := NewMockImageApi(mockCtrl)
	piwigomock.EXPECT().UploadImage(0, "/nonexisting/file.jpg", "1234", 2).Times(1).Return(5, nil)

	err := UploadImages(piwigomock, dbmock, 1)
	if err != nil {
		t.Error(err)
	}
}

func Test_uploadImages_saves_same_id_to_db(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	img := createTestImageMetaData(5)
	images := []datastore.ImageMetaData{img}

	imgToSave := img
	imgToSave.UploadRequired = false

	dbmock := NewMockImageMetadataProvider(mockCtrl)
	dbmock.EXPECT().ImageMetadataToUpload().Times(1).Return(images, nil)
	dbmock.EXPECT().SaveImageMetadata(imgToSave).Times(1)

	piwigomock := NewMockImageApi(mockCtrl)
	piwigomock.EXPECT().UploadImage(5, "/nonexisting/file.jpg", "1234", 2).Times(1).Return(5, nil)

	err := UploadImages(piwigomock, dbmock, 1)
	if err != nil {
		t.Error(err)
	}
}
