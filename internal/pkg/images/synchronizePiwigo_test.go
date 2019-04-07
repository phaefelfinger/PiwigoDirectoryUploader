/*
 * Copyright (C) 2019 Philipp Haefelfinger (http://www.haefelfinger.ch/). All Rights Reserved.
 * This application is licensed under GPLv2. See the LICENSE file in the root directory of the project.
 */

package images

import (
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/datastore"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/piwigo"
	"github.com/golang/mock/gomock"
	"testing"
)

func Test_checkPiwigoForChangedImages_none_with_piwigoId(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	img := datastore.ImageMetaData{ImageId: 1, UploadRequired: true}
	images := []datastore.ImageMetaData{img}

	dbmock := NewMockImageMetadataProvider(mockCtrl)
	dbmock.EXPECT().ImageMetadataToUpload().Return(images, nil)

	piwigomock := NewMockPiwigoImageApi(mockCtrl)
	piwigomock.EXPECT().ImagesExistOnPiwigo(gomock.Any()).Times(0)
	piwigomock.EXPECT().ImageCheckFile(gomock.Any(), gomock.Any()).Times(0)

	err := checkPiwigoForChangedImages(dbmock, piwigomock)
	if err != nil {
		t.Error(err)
	}
}

func Test_checkPiwigoForChangedImages_with_empty_list(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	images := []datastore.ImageMetaData{}

	dbmock := NewMockImageMetadataProvider(mockCtrl)
	dbmock.EXPECT().ImageMetadataToUpload().Return(images, nil)

	piwigomock := NewMockPiwigoImageApi(mockCtrl)
	piwigomock.EXPECT().ImagesExistOnPiwigo(gomock.Any()).Times(0)
	piwigomock.EXPECT().ImageCheckFile(gomock.Any(), gomock.Any()).Times(0)

	err := checkPiwigoForChangedImages(dbmock, piwigomock)
	if err != nil {
		t.Error(err)
	}
}

func Test_checkPiwigoForChangedImages_should_call_piwigo_set_uploadRequired_to_false(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	img := datastore.ImageMetaData{
		ImageId:        1,
		PiwigoId:       1,
		UploadRequired: true,
		Md5Sum:         "1234",
	}
	images := []datastore.ImageMetaData{img}
	dbmock := NewMockImageMetadataProvider(mockCtrl)
	dbmock.EXPECT().ImageMetadataToUpload().Return(images, nil)

	imgExpected := img
	imgExpected.UploadRequired = false
	dbmock.EXPECT().SaveImageMetadata(imgExpected).Times(1)

	piwigomock := NewMockPiwigoImageApi(mockCtrl)
	piwigomock.EXPECT().ImageCheckFile(1, "1234").Return(piwigo.ImageStateUptodate, nil)

	err := checkPiwigoForChangedImages(dbmock, piwigomock)
	if err != nil {
		t.Error(err)
	}
}

func Test_checkPiwigoForChangedImages_return_image_differs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	img := datastore.ImageMetaData{
		ImageId:        1,
		PiwigoId:       1,
		UploadRequired: true,
		Md5Sum:         "1234",
	}
	images := []datastore.ImageMetaData{img}
	dbmock := NewMockImageMetadataProvider(mockCtrl)
	dbmock.EXPECT().ImageMetadataToUpload().Return(images, nil)
	dbmock.EXPECT().SaveImageMetadata(gomock.Any()).Times(0)

	piwigomock := NewMockPiwigoImageApi(mockCtrl)
	piwigomock.EXPECT().ImageCheckFile(1, "1234").Return(piwigo.ImageStateDifferent, nil)

	err := checkPiwigoForChangedImages(dbmock, piwigomock)
	if err != nil {
		t.Error(err)
	}
}

func Test_updatePiwigoIdIfAlreadyUploaded_without_images_to_upload(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	images := []datastore.ImageMetaData{}

	dbmock := NewMockImageMetadataProvider(mockCtrl)
	dbmock.EXPECT().ImageMetadataToUpload().Return(images, nil)
	dbmock.EXPECT().SavePiwigoIdAndUpdateUploadFlag(gomock.Any(), gomock.Any()).Times(0)

	piwigomock := NewMockPiwigoImageApi(mockCtrl)
	piwigomock.EXPECT().ImagesExistOnPiwigo(gomock.Any()).Times(0)

	err := updatePiwigoIdIfAlreadyUploaded(dbmock, piwigomock)
	if err != nil {
		t.Error(err)
	}
}

func Test_updatePiwigoIdIfAlreadyUploaded_without_image_to_check(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	img := datastore.ImageMetaData{
		ImageId:        1,
		PiwigoId:       1,
		UploadRequired: true,
		Md5Sum:         "1234",
	}
	images := []datastore.ImageMetaData{img}

	dbmock := NewMockImageMetadataProvider(mockCtrl)
	dbmock.EXPECT().ImageMetadataToUpload().Return(images, nil)
	dbmock.EXPECT().SavePiwigoIdAndUpdateUploadFlag(gomock.Any(), gomock.Any()).Times(0)

	piwigomock := NewMockPiwigoImageApi(mockCtrl)
	piwigomock.EXPECT().ImagesExistOnPiwigo(gomock.Any()).Times(0)

	err := updatePiwigoIdIfAlreadyUploaded(dbmock, piwigomock)
	if err != nil {
		t.Error(err)
	}
}

func Test_updatePiwigoIdIfAlreadyUploaded_with_image_to_check(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	img := datastore.ImageMetaData{
		ImageId:        1,
		PiwigoId:       0,
		UploadRequired: true,
		Md5Sum:         "1234",
	}
	images := []datastore.ImageMetaData{img}

	dbmock := NewMockImageMetadataProvider(mockCtrl)
	dbmock.EXPECT().ImageMetadataToUpload().Return(images, nil)
	dbmock.EXPECT().SavePiwigoIdAndUpdateUploadFlag("1234", 1).Times(1)

	piwigoResponose := make(map[string]int)
	piwigoResponose["1234"] = 1

	piwigomock := NewMockPiwigoImageApi(mockCtrl)
	piwigomock.EXPECT().ImagesExistOnPiwigo(gomock.Any()).Times(1).Return(piwigoResponose, nil)

	err := updatePiwigoIdIfAlreadyUploaded(dbmock, piwigomock)
	if err != nil {
		t.Error(err)
	}
}

func Test_updatePiwigoIdIfAlreadyUploaded_with_image_to_check_missing_on_server(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	img := datastore.ImageMetaData{
		ImageId:        1,
		PiwigoId:       0,
		UploadRequired: true,
		Md5Sum:         "1234",
	}
	images := []datastore.ImageMetaData{img}

	dbmock := NewMockImageMetadataProvider(mockCtrl)
	dbmock.EXPECT().ImageMetadataToUpload().Return(images, nil)
	dbmock.EXPECT().SavePiwigoIdAndUpdateUploadFlag("1234", 1).Times(0)

	piwigoResponose := make(map[string]int)

	piwigomock := NewMockPiwigoImageApi(mockCtrl)
	piwigomock.EXPECT().ImagesExistOnPiwigo(gomock.Any()).Times(1).Return(piwigoResponose, nil)

	err := updatePiwigoIdIfAlreadyUploaded(dbmock, piwigomock)
	if err != nil {
		t.Error(err)
	}
}
