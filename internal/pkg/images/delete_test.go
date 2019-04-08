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

func Test_deleteImages_should_call_piwigo_and_remove_metadata(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	img := createTestImageMetaData(5)
	img.UploadRequired = false
	img.DeleteRequired = true
	images := []datastore.ImageMetaData{img}

	dbmock := NewMockImageMetadataProvider(mockCtrl)
	dbmock.EXPECT().ImageMetadataToDelete().Times(1).Return(images, nil)
	dbmock.EXPECT().DeleteMarkedImages().Times(1).Return(nil)

	piwigomock := NewMockImageApi(mockCtrl)
	piwigomock.EXPECT().DeleteImages([]int{5}).Times(1).Return(nil)

	err := DeleteImages(piwigomock, dbmock)
	if err != nil {
		t.Error(err)
	}
}

func Test_deleteImages_should_not_call_piwigo_for_not_uploaded_images_and_remove_metadata(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	img := createTestImageMetaData(0)
	img.UploadRequired = false
	img.DeleteRequired = true
	images := []datastore.ImageMetaData{img}

	dbmock := NewMockImageMetadataProvider(mockCtrl)
	dbmock.EXPECT().ImageMetadataToDelete().Times(1).Return(images, nil)
	dbmock.EXPECT().DeleteMarkedImages().Times(1).Return(nil)

	piwigomock := NewMockImageApi(mockCtrl)
	piwigomock.EXPECT().DeleteImages(gomock.Any()).Times(0)

	err := DeleteImages(piwigomock, dbmock)
	if err != nil {
		t.Error(err)
	}
}

func Test_deleteImages_should_not_call_anything_if_no_images_are_marked_for_deletion(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	var images []datastore.ImageMetaData

	dbmock := NewMockImageMetadataProvider(mockCtrl)
	dbmock.EXPECT().ImageMetadataToDelete().Times(1).Return(images, nil)
	dbmock.EXPECT().DeleteMarkedImages().Times(0)

	piwigomock := NewMockImageApi(mockCtrl)
	piwigomock.EXPECT().DeleteImages(gomock.Any()).Times(0)

	err := DeleteImages(piwigomock, dbmock)
	if err != nil {
		t.Error(err)
	}
}
