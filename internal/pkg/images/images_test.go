/*
 * Copyright (C) 2019 Philipp Haefelfinger (http://www.haefelfinger.ch/). All Rights Reserved.
 * This application is licensed under GPLv2. See the LICENSE file in the root directory of the project.
 */

package images

//go:generate mockgen -destination=./piwigo_mock_test.go -package=app git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/piwigo PiwigoApi,PiwigoCategoryApi,PiwigoImageApi
//go:generate mockgen -destination=./datastore_mock_test.go -package=app git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/datastore ImageMetadataProvider

import (
	"errors"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/datastore"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/localFileStructure"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/piwigo"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func Test_synchronize_local_image_metadata_should_find_nothing_if_empty(t *testing.T) {
	categories := make(map[string]*piwigo.PiwigoCategory)
	categories["2019/shooting1"] = &piwigo.PiwigoCategory{Id: 1}

	db := NewtestStore()
	fileSystemNodes := map[string]*localFileStructure.FilesystemNode{}

	err := SynchronizeLocalImageMetadata(db, fileSystemNodes, categories, testChecksumCalculator)
	if err != nil {
		t.Error(err)
	}

	if len(db.savedMetadata) > 0 {
		t.Error("There were metadata records saved but non expected!")
	}
}

func Test_synchronize_local_image_metadata_should_add_new_metadata(t *testing.T) {

	categories := make(map[string]*piwigo.PiwigoCategory)
	categories["2019/shooting1"] = &piwigo.PiwigoCategory{Id: 1}

	db := NewtestStore()

	testFileSystemNode := &localFileStructure.FilesystemNode{
		Key:     "2019/shooting1/abc.jpg",
		ModTime: time.Date(2019, 01, 01, 01, 0, 0, 0, time.UTC),
		Name:    "abc.jpg",
		Path:    "2019/shooting1/abc.jpg",
		IsDir:   false}

	fileSystemNodes := map[string]*localFileStructure.FilesystemNode{}
	fileSystemNodes[testFileSystemNode.Key] = testFileSystemNode

	// execute the sync metadata based on the file system results
	err := SynchronizeLocalImageMetadata(db, fileSystemNodes, categories, testChecksumCalculator)
	if err != nil {
		t.Error(err)
	}

	// check if data are saved
	savedData, exist := db.savedMetadata[testFileSystemNode.Key]
	if !exist {
		t.Fatal("Could not find correct metadata!")
	}
	if savedData.FullImagePath != testFileSystemNode.Key {
		t.Errorf("fullImagePath %s on db image metadata is not set to %s!", savedData.FullImagePath, testFileSystemNode.Key)
	}
	if savedData.LastChange != testFileSystemNode.ModTime {
		t.Error("lastChange on db image metadata is not set to the right date!")
	}
	if savedData.Filename != "abc.jpg" {
		t.Error("filename on db image metadata is not set to abc.jpg!")
	}
	if savedData.Md5Sum != testFileSystemNode.Key {
		t.Errorf("md5sum %s on db image metadata is not set to %s!", savedData.Md5Sum, testFileSystemNode.Key)
	}
	if savedData.UploadRequired != true {
		t.Errorf("uploadRequired on db image metadata is not set to true!")
	}
}

func Test_synchronize_local_image_metadata_should_mark_unchanged_entries_without_piwigoid_as_uploads_and_reset_deleted(t *testing.T) {

	categories := make(map[string]*piwigo.PiwigoCategory)
	categories["2019/shooting1"] = &piwigo.PiwigoCategory{Id: 1}

	db := NewtestStore()
	db.savedMetadata["2019/shooting1/abc.jpg"] = datastore.ImageMetaData{
		Md5Sum:         "2019/shooting1/abc.jpg",
		FullImagePath:  "2019/shooting1/abc.jpg",
		PiwigoId:       0,
		UploadRequired: false,
		LastChange:     time.Date(2019, 01, 01, 00, 0, 0, 0, time.UTC),
		Filename:       "abc.jpg",
		DeleteRequired: true,
	}

	testFileSystemNode := &localFileStructure.FilesystemNode{
		Key:     "2019/shooting1/abc.jpg",
		ModTime: time.Date(2019, 01, 01, 00, 0, 0, 0, time.UTC),
		Name:    "abc.jpg",
		Path:    "2019/shooting1/abc.jpg",
		IsDir:   false}

	fileSystemNodes := map[string]*localFileStructure.FilesystemNode{}
	fileSystemNodes[testFileSystemNode.Key] = testFileSystemNode

	// execute the sync metadata based on the file system results
	err := SynchronizeLocalImageMetadata(db, fileSystemNodes, categories, testChecksumCalculator)
	if err != nil {
		t.Error(err)
	}

	// check if data are saved
	savedData, exist := db.savedMetadata[testFileSystemNode.Key]
	if !exist {
		t.Fatal("Could not find correct metadata!")
	}
	if savedData.LastChange != testFileSystemNode.ModTime {
		t.Error("lastChange on db image metadata is not set to the right date!")
	}
	if savedData.UploadRequired != true {
		t.Errorf("uploadRequired on db image metadata is not set to true!")
	}
	if savedData.DeleteRequired != false {
		t.Errorf("deleteRequired on db image metadata is not set to false!")
	}
}

func Test_synchronize_local_image_metadata_should_mark_changed_entries_as_uploads_and_reset_deleted(t *testing.T) {

	categories := make(map[string]*piwigo.PiwigoCategory)
	categories["2019/shooting1"] = &piwigo.PiwigoCategory{Id: 1}

	db := NewtestStore()
	db.savedMetadata["2019/shooting1/abc.jpg"] = datastore.ImageMetaData{
		Md5Sum:         "2019/shooting1/abc.jpg",
		FullImagePath:  "2019/shooting1/abc.jpg",
		UploadRequired: false,
		LastChange:     time.Date(2019, 01, 01, 00, 0, 0, 0, time.UTC),
		Filename:       "abc.jpg",
		DeleteRequired: true,
	}

	testFileSystemNode := &localFileStructure.FilesystemNode{
		Key:     "2019/shooting1/abc.jpg",
		ModTime: time.Date(2019, 01, 01, 01, 0, 0, 0, time.UTC),
		Name:    "abc.jpg",
		Path:    "2019/shooting1/abc.jpg",
		IsDir:   false}

	fileSystemNodes := map[string]*localFileStructure.FilesystemNode{}
	fileSystemNodes[testFileSystemNode.Key] = testFileSystemNode

	// execute the sync metadata based on the file system results
	err := SynchronizeLocalImageMetadata(db, fileSystemNodes, categories, testChecksumCalculator)
	if err != nil {
		t.Error(err)
	}

	// check if data are saved
	savedData, exist := db.savedMetadata[testFileSystemNode.Key]
	if !exist {
		t.Fatal("Could not find correct metadata!")
	}
	if savedData.LastChange != testFileSystemNode.ModTime {
		t.Error("lastChange on db image metadata is not set to the right date!")
	}
	if savedData.UploadRequired != true {
		t.Errorf("uploadRequired on db image metadata is not set to true!")
	}
	if savedData.DeleteRequired != false {
		t.Errorf("deleteRequired on db image metadata is not set to false!")
	}
}

func Test_synchronize_local_image_metadata_should_not_mark_unchanged_files_to_upload_and_reset_deleted(t *testing.T) {
	db := NewtestStore()

	categories := make(map[string]*piwigo.PiwigoCategory)
	categories["2019/shooting1"] = &piwigo.PiwigoCategory{Id: 1}

	db.savedMetadata["2019/shooting1/abc.jpg"] = datastore.ImageMetaData{
		Md5Sum:         "2019/shooting1/abc.jpg",
		FullImagePath:  "2019/shooting1/abc.jpg",
		PiwigoId:       5,
		UploadRequired: false,
		LastChange:     time.Date(2019, 01, 01, 01, 0, 0, 0, time.UTC),
		Filename:       "abc.jpg",
		DeleteRequired: true,
	}

	testFileSystemNode := &localFileStructure.FilesystemNode{
		Key:     "2019/shooting1/abc.jpg",
		ModTime: time.Date(2019, 01, 01, 01, 0, 0, 0, time.UTC),
		Name:    "abc.jpg",
		Path:    "2019/shooting1/abc.jpg",
		IsDir:   false}

	fileSystemNodes := map[string]*localFileStructure.FilesystemNode{}
	fileSystemNodes[testFileSystemNode.Key] = testFileSystemNode

	// execute the sync metadata based on the file system results
	err := SynchronizeLocalImageMetadata(db, fileSystemNodes, categories, testChecksumCalculator)
	if err != nil {
		t.Error(err)
	}

	// check if data are saved
	savedData, exist := db.savedMetadata[testFileSystemNode.Key]
	if !exist {
		t.Fatal("Could not find correct metadata!")
	}
	if savedData.UploadRequired {
		t.Errorf("uploadRequired on db image metadata is set to true, but should not be on unchanged items!")
	}
	if savedData.DeleteRequired != false {
		t.Errorf("deleteRequired on db image metadata is not set to false!")
	}
}

func Test_synchronize_local_image_metadata_should_not_process_directories(t *testing.T) {
	categories := make(map[string]*piwigo.PiwigoCategory)
	categories["2019/shooting1"] = &piwigo.PiwigoCategory{Id: 1}

	db := NewtestStore()

	testFileSystemNode := &localFileStructure.FilesystemNode{
		Key:     "2019/shooting1",
		ModTime: time.Date(2019, 01, 01, 01, 0, 0, 0, time.UTC),
		Name:    "shooting1",
		Path:    "2019/shooting1/",
		IsDir:   true}

	fileSystemNodes := map[string]*localFileStructure.FilesystemNode{}
	fileSystemNodes[testFileSystemNode.Key] = testFileSystemNode

	// execute the sync metadata based on the file system results
	err := SynchronizeLocalImageMetadata(db, fileSystemNodes, categories, testChecksumCalculator)
	if err != nil {
		t.Error(err)
	}

	if len(db.savedMetadata) > 0 {
		t.Error("There were metadata records saved but non expected!")
	}
}

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

	piwigomock := NewMockPiwigoImageApi(mockCtrl)
	piwigomock.EXPECT().UploadImage(0, "/nonexisting/file.jpg", "1234", 2).Times(1).Return(5, nil)

	err := UploadImages(piwigomock, dbmock)
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

	piwigomock := NewMockPiwigoImageApi(mockCtrl)
	piwigomock.EXPECT().UploadImage(5, "/nonexisting/file.jpg", "1234", 2).Times(1).Return(5, nil)

	err := UploadImages(piwigomock, dbmock)
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

	piwigomock := NewMockPiwigoImageApi(mockCtrl)
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

	piwigomock := NewMockPiwigoImageApi(mockCtrl)
	piwigomock.EXPECT().DeleteImages(gomock.Any()).Times(0)

	err := DeleteImages(piwigomock, dbmock)
	if err != nil {
		t.Error(err)
	}
}

func Test_deleteImages_should_not_call_anything_if_no_images_are_marked_for_deletion(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	images := []datastore.ImageMetaData{}

	dbmock := NewMockImageMetadataProvider(mockCtrl)
	dbmock.EXPECT().ImageMetadataToDelete().Times(1).Return(images, nil)
	dbmock.EXPECT().DeleteMarkedImages().Times(0)

	piwigomock := NewMockPiwigoImageApi(mockCtrl)
	piwigomock.EXPECT().DeleteImages(gomock.Any()).Times(0)

	err := DeleteImages(piwigomock, dbmock)
	if err != nil {
		t.Error(err)
	}
}

// test metadata store to store save the metadat and simulate the database
//TODO: refactor to use generated test implementation
type testStore struct {
	savedMetadata map[string]datastore.ImageMetaData
}

func NewtestStore() *testStore {
	return &testStore{savedMetadata: make(map[string]datastore.ImageMetaData)}
}

func (s *testStore) ImageMetadata(fullImagePath string) (datastore.ImageMetaData, error) {
	metadata, exist := s.savedMetadata[fullImagePath]
	if !exist {
		return datastore.ImageMetaData{}, datastore.ErrorRecordNotFound
	}
	return metadata, nil
}

func (d *testStore) ImageMetadataAll() ([]datastore.ImageMetaData, error) {
	return []datastore.ImageMetaData{}, nil
}

func (s *testStore) SaveImageMetadata(m datastore.ImageMetaData) error {
	s.savedMetadata[m.FullImagePath] = m
	return nil
}

func (d *testStore) ImageMetadataToUpload() ([]datastore.ImageMetaData, error) {
	return nil, errors.New("N/A")
}

func (d *testStore) ImageMetadataToDelete() ([]datastore.ImageMetaData, error) {
	return nil, errors.New("N/A")
}

func (d *testStore) SavePiwigoIdAndUpdateUploadFlag(md5Sum string, piwigoId int) error {
	return errors.New("N/A")
}

func (d *testStore) DeleteMarkedImages() error {
	return errors.New("N/A")
}

// to make the sync testable, we pass in a simple mock that returns the filepath as checksum
func testChecksumCalculator(file string) (string, error) {
	return file, nil
}

func createTestImageMetaData(piwigoId int) datastore.ImageMetaData {
	img := datastore.ImageMetaData{
		ImageId:        1,
		PiwigoId:       piwigoId,
		FullImagePath:  "/nonexisting/file.jpg",
		UploadRequired: true,
		Md5Sum:         "1234",
		CategoryId:     2,
	}
	return img
}
