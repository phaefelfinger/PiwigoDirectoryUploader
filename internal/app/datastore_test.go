/*
 * Copyright (C) 2019 Philipp Haefelfinger (http://www.haefelfinger.ch/). All Rights Reserved.
 * This application is licensed under GPLv2. See the LICENSE file in the root directory of the project.
 */

package app

import (
	"os"
	"strings"
	"testing"
	"time"
)

var databaseFile = "./metadatatest.db"
var dbinitOk bool

func TestDataStoreInitialize(t *testing.T) {
	_ = setupDatabase(t)
	cleanupDatabase(t)
}

func Test_save_and_load_metadata(t *testing.T) {
	if !dbinitOk {
		t.Skip("Skipping test as TestDataStoreInitialize failed!")
	}
	dataStore := setupDatabase(t)
	defer cleanupDatabase(t)

	filePath := "blah/foo/bar.jpg"
	img := getExampleImageMetadata(filePath)

	saveImageShouldNotFail("insert", dataStore, img, t)
	img.ImageId = 1

	imgLoad := loadMetadataShouldNotFail("insert", dataStore, filePath, t)
	EnsureMetadataAreEqual("insert", img, imgLoad, t)

	// updated the image again
	img.Md5Sum = "123456"
	saveImageShouldNotFail("update", dataStore, img, t)

	imgLoad = loadMetadataShouldNotFail("update", dataStore, filePath, t)
	EnsureMetadataAreEqual("update", img, imgLoad, t)
}

func Test_save_and_query_for_upload_records(t *testing.T) {
	if !dbinitOk {
		t.Skip("Skipping test as TestDataStoreInitialize failed!")
	}
	dataStore := setupDatabase(t)
	defer cleanupDatabase(t)

	img := getExampleImageMetadata("blah/foo/bar.jpg")

	saveImageShouldNotFail("toupload", dataStore, img, t)
	img.ImageId = 1

	images, err := dataStore.ImageMetadataToUpload()
	if err != nil {
		t.Fatalf("Could not query images to upload! %s", err)
	}

	if len(images) != 1 {
		t.Fatal("Did not get any saved images to upload!")
	}

	imgLoad := images[0]
	EnsureMetadataAreEqual("toupload", img, imgLoad, t)

}

func Test_save_and_query_for_upload_records_do_not_contain_images_to_delete(t *testing.T) {
	if !dbinitOk {
		t.Skip("Skipping test as TestDataStoreInitialize failed!")
	}
	dataStore := setupDatabase(t)
	defer cleanupDatabase(t)

	img1 := getExampleImageMetadata("blah/foo/bar.jpg")

	img2 := getExampleImageMetadata("blah/foo/bar2.jpg")
	img2.DeleteRequired = true

	saveImageShouldNotFail("toupload1", dataStore, img1, t)
	img1.ImageId = 1

	saveImageShouldNotFail("toupload2", dataStore, img2, t)
	img2.ImageId = 2

	images, err := dataStore.ImageMetadataToUpload()
	if err != nil {
		t.Fatalf("Could not query images to upload! %s", err)
	}

	if len(images) > 1 {
		t.Fatal("Got more than one image to upload but only one is expected")
	}

	if len(images) != 1 {
		t.Fatal("Did not get the saved images to upload!")
	}

	imgLoad := images[0]
	EnsureMetadataAreEqual("toupload", img1, imgLoad, t)
}

func Test_load_metadata_not_found(t *testing.T) {
	if !dbinitOk {
		t.Skip("Skipping test as TestDataStoreInitialize failed!")
	}
	dataStore := setupDatabase(t)
	defer cleanupDatabase(t)

	filePath := "blah/foo/bar.jpg"
	imgLoad, err := dataStore.ImageMetadata(filePath)
	if err != ErrorRecordNotFound {
		t.Errorf("Unexpected error on loading non existing file %s: %s", filePath, err)
	}
	if imgLoad.ImageId > 0 {
		t.Error("Found an image metadata that should not exist on an emtpy database.")
	}
}

func Test_unique_index_on_relativeFilePath(t *testing.T) {
	if !dbinitOk {
		t.Skip("Skipping test as TestDataStoreInitialize failed!")
	}
	dataStore := setupDatabase(t)
	defer cleanupDatabase(t)

	img := getExampleImageMetadata("blah/foo/bar.jpg")

	saveImageShouldNotFail("insert", dataStore, img, t)

	err := dataStore.SaveImageMetadata(img)
	if err == nil {
		t.Errorf("Could save duplicated image metadata. Expected error but got none!")
	}

	// check if the error contains the expected column as name. If not, this indicates another problem than
	// the expected duplicated insert error.
	if !strings.Contains(err.Error(), "fullImagePath") {
		t.Errorf("Got a unexpected error on saving duplicate records: %s", err)
	}
}

func Test_update_piwigoId_by_checksum(t *testing.T) {
	if !dbinitOk {
		t.Skip("Skipping test as TestDataStoreInitialize failed!")
	}
	dataStore := setupDatabase(t)
	defer cleanupDatabase(t)

	filePath := "blah/foo/bar.jpg"
	img := getExampleImageMetadata(filePath)

	saveImageShouldNotFail("SavePiwigoIdAndUpdateUploadFlag", dataStore, img, t)
	img.ImageId = 1
	img.PiwigoId = 1234
	img.UploadRequired = false

	err := dataStore.SavePiwigoIdAndUpdateUploadFlag(img.Md5Sum, img.PiwigoId)
	if err != nil {
		t.Errorf("SavePiwigoIdAndUpdateUploadFlag: Could not update piwigo id: %s", err)
	}

	imgLoad := loadMetadataShouldNotFail("update", dataStore, filePath, t)
	EnsureMetadataAreEqual("SavePiwigoIdAndUpdateUploadFlag", img, imgLoad, t)
}

func Test_update_piwigoId_by_checksum_found_no_image(t *testing.T) {
	if !dbinitOk {
		t.Skip("Skipping test as TestDataStoreInitialize failed!")
	}
	dataStore := setupDatabase(t)
	defer cleanupDatabase(t)

	filePath := "blah/foo/bar.jpg"
	img := getExampleImageMetadata(filePath)

	saveImageShouldNotFail("SavePiwigoIdAndUpdateUploadFlag", dataStore, img, t)
	img.ImageId = 1
	img.PiwigoId = 0
	img.UploadRequired = true

	err := dataStore.SavePiwigoIdAndUpdateUploadFlag(img.Md5Sum, img.PiwigoId)
	if err != nil {
		t.Errorf("SavePiwigoIdAndUpdateUploadFlag: Could not update piwigo id: %s", err)
	}

	imgLoad := loadMetadataShouldNotFail("update", dataStore, filePath, t)
	EnsureMetadataAreEqual("SavePiwigoIdAndUpdateUploadFlag", img, imgLoad, t)

}

func saveImageShouldNotFail(action string, dataStore *localDataStore, img ImageMetaData, t *testing.T) {
	err := dataStore.SaveImageMetadata(img)
	if err != nil {
		t.Errorf("%s: Could not save Metadata: %s", action, err)
	}
}

func loadMetadataShouldNotFail(action string, dataStore *localDataStore, filePath string, t *testing.T) ImageMetaData {
	imgLoad, err := dataStore.ImageMetadata(filePath)
	if err != nil {
		t.Errorf("%s: Could not load saved Metadata: %s - %s", action, filePath, err)
	}
	return imgLoad
}

func EnsureMetadataAreEqual(action string, img ImageMetaData, imgLoad ImageMetaData, t *testing.T) {
	// check if both instances serialize to the same string representation
	if img.String() != imgLoad.String() {
		t.Errorf("%s: Invalid image loaded! expected (ignore ImageId) %s but got %s", action, img.String(), imgLoad.String())
	}
}

func getExampleImageMetadata(filePath string) ImageMetaData {
	return ImageMetaData{
		FullImagePath:  filePath,
		PiwigoId:       1,
		Md5Sum:         "aabbccddeeff",
		LastChange:     time.Now().UTC(),
		Filename:       "bar.jpg",
		CategoryPath:   "blah/foo",
		CategoryId:     100,
		UploadRequired: true,
	}
}

func cleanupDatabase(t *testing.T) {
	err := os.Remove(databaseFile)
	if err != nil {
		t.Errorf("Failed remove test database %s: %s", databaseFile, err)
	}
}

func setupDatabase(t *testing.T) *localDataStore {
	dataStore := &localDataStore{}
	err := dataStore.Initialize(databaseFile)
	if err != nil {
		t.Errorf("Failed to init datastore: %s", err)
		return nil
	}
	dbinitOk = true
	return dataStore
}
