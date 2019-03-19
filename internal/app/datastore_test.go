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

func TestSaveAndLoadMetadata(t *testing.T) {
	if !dbinitOk {
		t.Skip("Skipping test as TestDataStoreInitialize failed!")
	}
	dataStore := setupDatabase(t)

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

	cleanupDatabase(t)
}

func TestSaveAndQueryForUploadRecords(t *testing.T) {
	if !dbinitOk {
		t.Skip("Skipping test as TestDataStoreInitialize failed!")
	}
	dataStore := setupDatabase(t)

	filePath := "blah/foo/bar.jpg"
	img := getExampleImageMetadata(filePath)

	saveImageShouldNotFail("toupload", dataStore, img, t)
	img.ImageId = 1

	images, err := dataStore.ImageMetadataToUpload()
	if err != nil {
		t.Fatalf("Could not query images to upload! %s", err)
	}

	if len(images)<1 {
		t.Fatal("Did not get any saved images to upload!")
	}

	imgLoad := images[0]
	EnsureMetadataAreEqual("toupload", img, *imgLoad, t)

	cleanupDatabase(t)
}

func TestLoadMetadataNotFound(t *testing.T) {
	if !dbinitOk {
		t.Skip("Skipping test as TestDataStoreInitialize failed!")
	}
	dataStore := setupDatabase(t)

	filePath := "blah/foo/bar.jpg"
	imgLoad, err := dataStore.ImageMetadata(filePath)
	if err != ErrorRecordNotFound {
		t.Errorf("Unexpected error on loading non existing file %s: %s", filePath, err)
	}
	if imgLoad.ImageId > 0 {
		t.Error("Found an image metadata that should not exist on an emtpy database.")
	}

	cleanupDatabase(t)
}

func TestUniqueIndexOnRelativeFilePath(t *testing.T) {
	if !dbinitOk {
		t.Skip("Skipping test as TestDataStoreInitialize failed!")
	}
	dataStore := setupDatabase(t)

	filePath := "blah/foo/bar.jpg"
	img := getExampleImageMetadata(filePath)

	saveImageShouldNotFail("insert", dataStore, img, t)

	err := dataStore.SaveImageMetadata(img)
	if err == nil {
		t.Errorf("Could save duplicated image metadata. Expected error but got none!")
	}

	// check if the error contains the expected column as name. If not, this indicates another problem than
	// the expected duplicated insert error.
	if !strings.Contains(err.Error(), "relativePath") {
		t.Errorf("Got a unexpected error on saving duplicate records: %s", err)
	}

	cleanupDatabase(t)
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
		RelativeImagePath: filePath,
		PiwigoId:          1,
		Md5Sum:            "aabbccddeeff",
		LastChange:        time.Now().UTC(),
		Filename:          "bar.jpg",
		CategoryPath:      "blah/foo",
		CategoryId:        100,
		UploadRequired:    true,
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
