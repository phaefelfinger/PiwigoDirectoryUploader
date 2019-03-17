package app

import (
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/localFileStructure"
	"testing"
	"time"
)

func TestSynchronizeLocalImageMetadataShouldDoNothingIfEmpty(t *testing.T) {
	db := NewtestStore()
	fileSystemNodes := map[string]*localFileStructure.FilesystemNode{}

	err := synchronizeLocalImageMetadata(db, fileSystemNodes, testChecksumCalculator)
	if err != nil {
		t.Error(err)
	}

	if len(db.savedMetadata) > 0 {
		t.Error("There were metadata records saved but non expected!")
	}
}

func TestSynchronizeLocalImageMetadataShouldAddNewMetadata(t *testing.T) {
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
	err := synchronizeLocalImageMetadata(db, fileSystemNodes, testChecksumCalculator)
	if err != nil {
		t.Error(err)
	}

	// check if data are saved
	savedData, exist := db.savedMetadata[testFileSystemNode.Key]
	if !exist {
		t.Fatal("Could not find correct metadata!")
	}
	if savedData.RelativeImagePath != testFileSystemNode.Key {
		t.Errorf("relativeImagePath %s on db image metadata is not set to %s!", savedData.RelativeImagePath, testFileSystemNode.Key)
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

func TestSynchronizeLocalImageMetadataShouldMarkChangedEntriesAsUploads(t *testing.T) {
	db := NewtestStore()
	db.savedMetadata["2019/shooting1/abc.jpg"] = ImageMetaData{
		Md5Sum:            "2019/shooting1/abc.jpg",
		RelativeImagePath: "2019/shooting1/abc.jpg",
		UploadRequired:    false,
		LastChange:        time.Date(2019, 01, 01, 00, 0, 0, 0, time.UTC),
		Filename:          "abc.jpg",
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
	err := synchronizeLocalImageMetadata(db, fileSystemNodes, testChecksumCalculator)
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
}

func TestSynchronizeLocalImageMetadataShouldNotMarkUnchangedFilesToUpload(t *testing.T) {
	db := NewtestStore()
	db.savedMetadata["2019/shooting1/abc.jpg"] = ImageMetaData{
		Md5Sum:            "2019/shooting1/abc.jpg",
		RelativeImagePath: "2019/shooting1/abc.jpg",
		UploadRequired:    false,
		LastChange:        time.Date(2019, 01, 01, 01, 0, 0, 0, time.UTC),
		Filename:          "abc.jpg",
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
	err := synchronizeLocalImageMetadata(db, fileSystemNodes, testChecksumCalculator)
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
}

func TestSynchronizeLocalImageMetadataShouldNotProcessDirectories(t *testing.T) {
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
	err := synchronizeLocalImageMetadata(db, fileSystemNodes, testChecksumCalculator)
	if err != nil {
		t.Error(err)
	}

	if len(db.savedMetadata) > 0 {
		t.Error("There were metadata records saved but non expected!")
	}
}

// test metadata store to store save the metadat and simulate the database
type testStore struct {
	savedMetadata map[string]ImageMetaData
}

func NewtestStore() *testStore {
	return &testStore{savedMetadata: make(map[string]ImageMetaData)}
}

func (s *testStore) GetImageMetadata(relativePath string) (ImageMetaData, error) {
	metadata, exist := s.savedMetadata[relativePath]
	if !exist {
		return ImageMetaData{}, ErrorRecordNotFound
	}
	return metadata, nil
}

func (s *testStore) SaveImageMetadata(m ImageMetaData) error {
	s.savedMetadata[m.RelativeImagePath] = m
	return nil
}

// to make the sync testable, we pass in a simple mock that returns the filepath as checksum
func testChecksumCalculator(file string) (string, error) {
	return file, nil
}
