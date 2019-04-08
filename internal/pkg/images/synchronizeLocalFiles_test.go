/*
 * Copyright (C) 2019 Philipp Haefelfinger (http://www.haefelfinger.ch/). All Rights Reserved.
 * This application is licensed under GPLv2. See the LICENSE file in the root directory of the project.
 */

package images

//go:generate mockgen -destination=./piwigo_mock_test.go -package=images git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/piwigo CategoryApi,ImageApi
//go:generate mockgen -destination=./datastore_mock_test.go -package=images git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/datastore ImageMetadataProvider,CategoryProvider

import (
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/datastore"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/localFileStructure"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func Test_synchronize_local_image_metadata_should_find_nothing_if_empty(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	category := datastore.CategoryData{CategoryId: 1, Name: "shooting1", PiwigoId: 1, Key: "2019/shooting1"}
	categoryMock := NewMockCategoryProvider(mockCtrl)
	categoryMock.EXPECT().GetCategoryByKey(category.Key).Return(category, nil).Times(0)

	db := NewMockImageMetadataProvider(mockCtrl)
	db.EXPECT().ImageMetadataAll().Times(1)
	db.EXPECT().SaveImageMetadata(gomock.Any()).Times(0)

	fileSystemNodes := map[string]*localFileStructure.FilesystemNode{}

	err := SynchronizeLocalImageMetadata(db, categoryMock, fileSystemNodes, testChecksumCalculator)
	if err != nil {
		t.Error(err)
	}
}

func Test_synchronize_local_image_metadata_should_add_new_metadata(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	category := datastore.CategoryData{CategoryId: 1, Name: "shooting1", PiwigoId: 1, Key: "2019/shooting1"}
	categoryMock := NewMockCategoryProvider(mockCtrl)
	categoryMock.EXPECT().GetCategoryByKey(category.Key).Return(category, nil).Times(1)

	testFileSystemNode := &localFileStructure.FilesystemNode{
		Key:     "2019/shooting1/abc.jpg",
		ModTime: time.Date(2019, 01, 01, 01, 0, 0, 0, time.UTC),
		Name:    "abc.jpg",
		Path:    "2019/shooting1/abc.jpg",
		IsDir:   false}

	fileSystemNodes := map[string]*localFileStructure.FilesystemNode{}
	fileSystemNodes[testFileSystemNode.Key] = testFileSystemNode

	image := createImageMetaDataFromFilesystem(testFileSystemNode, 0, true, false)
	image.CategoryPiwigoId = category.PiwigoId
	image.CategoryPath = category.Key

	db := NewMockImageMetadataProvider(mockCtrl)
	db.EXPECT().ImageMetadataAll().Times(1)
	db.EXPECT().ImageMetadata(testFileSystemNode.Key).Return(datastore.ImageMetaData{}, datastore.ErrorRecordNotFound).Times(1)
	db.EXPECT().SaveImageMetadata(image).Times(1)

	// execute the sync metadata based on the file system results
	err := SynchronizeLocalImageMetadata(db, categoryMock, fileSystemNodes, testChecksumCalculator)
	if err != nil {
		t.Error(err)
	}
}

func Test_synchronize_local_image_metadata_should_mark_unchanged_entries_without_piwigoid_as_uploads_and_reset_deleted(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	category := datastore.CategoryData{CategoryId: 1, Name: "shooting1", PiwigoId: 1, Key: "2019/shooting1"}
	categoryMock := NewMockCategoryProvider(mockCtrl)
	categoryMock.EXPECT().GetCategoryByKey(category.Key).Return(category, nil).Times(0)

	testFileSystemNode := &localFileStructure.FilesystemNode{
		Key:     "2019/shooting1/abc.jpg",
		ModTime: time.Date(2019, 01, 01, 00, 0, 0, 0, time.UTC),
		Name:    "abc.jpg",
		Path:    "2019/shooting1/abc.jpg",
		IsDir:   false}

	fileSystemNodes := map[string]*localFileStructure.FilesystemNode{}
	fileSystemNodes[testFileSystemNode.Key] = testFileSystemNode

	imageExptected := createImageMetaDataFromFilesystem(testFileSystemNode, 0, true, false)

	imageStored := imageExptected
	imageStored.DeleteRequired = true

	db := NewMockImageMetadataProvider(mockCtrl)
	db.EXPECT().ImageMetadataAll().Times(1)
	db.EXPECT().ImageMetadata(testFileSystemNode.Key).Return(imageStored, nil).Times(1)
	db.EXPECT().SaveImageMetadata(imageExptected).Times(1)

	// execute the sync metadata based on the file system results
	err := SynchronizeLocalImageMetadata(db, categoryMock, fileSystemNodes, testChecksumCalculator)
	if err != nil {
		t.Error(err)
	}
}

func Test_synchronize_local_image_metadata_should_mark_changed_entries_as_uploads_and_reset_deleted(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	category := datastore.CategoryData{CategoryId: 1, Name: "shooting1", PiwigoId: 1, Key: "2019/shooting1"}
	categoryMock := NewMockCategoryProvider(mockCtrl)
	categoryMock.EXPECT().GetCategoryByKey(category.Key).Return(category, nil).Times(0)

	testFileSystemNode := &localFileStructure.FilesystemNode{
		Key:     "2019/shooting1/abc.jpg",
		ModTime: time.Date(2019, 01, 01, 01, 0, 0, 0, time.UTC),
		Name:    "abc.jpg",
		Path:    "2019/shooting1/abc.jpg",
		IsDir:   false}

	fileSystemNodes := map[string]*localFileStructure.FilesystemNode{}
	fileSystemNodes[testFileSystemNode.Key] = testFileSystemNode

	imageExptected := createImageMetaDataFromFilesystem(testFileSystemNode, 0, true, false)

	imageStored := imageExptected
	imageStored.DeleteRequired = true
	imageStored.UploadRequired = false
	imageStored.LastChange = time.Date(2019, 01, 01, 00, 0, 0, 0, time.UTC)

	db := NewMockImageMetadataProvider(mockCtrl)
	db.EXPECT().ImageMetadataAll().Times(1)
	db.EXPECT().ImageMetadata(testFileSystemNode.Key).Return(imageStored, nil).Times(1)
	db.EXPECT().SaveImageMetadata(imageExptected).Times(1)

	// execute the sync metadata based on the file system results
	err := SynchronizeLocalImageMetadata(db, categoryMock, fileSystemNodes, testChecksumCalculator)
	if err != nil {
		t.Error(err)
	}
}

func Test_synchronize_local_image_metadata_should_not_mark_unchanged_files_to_upload_and_reset_deleted(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	category := datastore.CategoryData{CategoryId: 1, Name: "shooting1", PiwigoId: 1, Key: "2019/shooting1"}
	categoryMock := NewMockCategoryProvider(mockCtrl)
	categoryMock.EXPECT().GetCategoryByKey(category.Key).Return(category, nil).Times(0)

	testFileSystemNode := &localFileStructure.FilesystemNode{
		Key:     "2019/shooting1/abc.jpg",
		ModTime: time.Date(2019, 01, 01, 01, 0, 0, 0, time.UTC),
		Name:    "abc.jpg",
		Path:    "2019/shooting1/abc.jpg",
		IsDir:   false}

	fileSystemNodes := map[string]*localFileStructure.FilesystemNode{}
	fileSystemNodes[testFileSystemNode.Key] = testFileSystemNode

	imageExptected := createImageMetaDataFromFilesystem(testFileSystemNode, 5, false, false)

	imageStored := imageExptected
	imageStored.DeleteRequired = true
	imageStored.UploadRequired = false
	imageStored.LastChange = time.Date(2019, 01, 01, 01, 0, 0, 0, time.UTC)

	db := NewMockImageMetadataProvider(mockCtrl)
	db.EXPECT().ImageMetadataAll().Times(1)
	db.EXPECT().ImageMetadata(testFileSystemNode.Key).Return(imageStored, nil).Times(1)
	db.EXPECT().SaveImageMetadata(imageExptected).Times(1)

	// execute the sync metadata based on the file system results
	err := SynchronizeLocalImageMetadata(db, categoryMock, fileSystemNodes, testChecksumCalculator)
	if err != nil {
		t.Error(err)
	}
}

func Test_synchronize_local_image_metadata_should_not_process_directories(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	category := datastore.CategoryData{CategoryId: 1, Name: "shooting1", PiwigoId: 1, Key: "2019/shooting1"}
	categoryMock := NewMockCategoryProvider(mockCtrl)
	categoryMock.EXPECT().GetCategoryByKey(category.Key).Return(category, nil).Times(0)

	testFileSystemNode := &localFileStructure.FilesystemNode{
		Key:     "2019/shooting1",
		ModTime: time.Date(2019, 01, 01, 01, 0, 0, 0, time.UTC),
		Name:    "shooting1",
		Path:    "2019/shooting1/",
		IsDir:   true}

	fileSystemNodes := map[string]*localFileStructure.FilesystemNode{}
	fileSystemNodes[testFileSystemNode.Key] = testFileSystemNode

	db := NewMockImageMetadataProvider(mockCtrl)
	db.EXPECT().ImageMetadataAll().Times(1)
	db.EXPECT().SaveImageMetadata(gomock.Any()).Times(0)

	// execute the sync metadata based on the file system results
	err := SynchronizeLocalImageMetadata(db, categoryMock, fileSystemNodes, testChecksumCalculator)
	if err != nil {
		t.Error(err)
	}
}

func Test_synchronizeLocalImageMetadataFindFilesToDelete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	img := createTestImageMetaData(5)
	images := []datastore.ImageMetaData{img}

	imgToSave := img
	imgToSave.UploadRequired = false
	imgToSave.DeleteRequired = true

	dbmock := NewMockImageMetadataProvider(mockCtrl)
	dbmock.EXPECT().ImageMetadataAll().Times(1).Return(images, nil)
	dbmock.EXPECT().SaveImageMetadata(imgToSave).Times(1)

	err := synchronizeLocalImageMetadataFindFilesToDelete(dbmock)
	if err != nil {
		t.Error(err)
	}
}

// to make the sync testable, we pass in a simple mock that returns the filepath as checksum
func testChecksumCalculator(file string) (string, error) {
	return file, nil
}

func createTestImageMetaData(piwigoId int) datastore.ImageMetaData {
	img := datastore.ImageMetaData{
		ImageId:          1,
		PiwigoId:         piwigoId,
		FullImagePath:    "/nonexisting/file.jpg",
		UploadRequired:   true,
		Md5Sum:           "1234",
		CategoryPiwigoId: 2,
	}
	return img
}

func createImageMetaDataFromFilesystem(testFileSystemNode *localFileStructure.FilesystemNode, piwigoId int, uploadRequired bool, deleteRequired bool) datastore.ImageMetaData {
	imageExptected := datastore.ImageMetaData{
		Md5Sum:         testFileSystemNode.Key,
		FullImagePath:  testFileSystemNode.Key,
		PiwigoId:       piwigoId,
		UploadRequired: uploadRequired,
		LastChange:     testFileSystemNode.ModTime,
		Filename:       testFileSystemNode.Name,
		DeleteRequired: deleteRequired,
	}
	return imageExptected
}
