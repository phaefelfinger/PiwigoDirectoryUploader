package app

import (
	"errors"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/localFileStructure"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/piwigo"
	"testing"
	"time"
)

func TestSynchronizeLocalImageMetadataShouldDoNothingIfEmpty(t *testing.T) {
	categories := make(map[string]*piwigo.PiwigoCategory)
	categories["2019/shooting1"] = &piwigo.PiwigoCategory{Id: 1}

	db := NewtestStore()
	fileSystemNodes := map[string]*localFileStructure.FilesystemNode{}

	err := synchronizeLocalImageMetadata(db, fileSystemNodes, categories, testChecksumCalculator)
	if err != nil {
		t.Error(err)
	}

	if len(db.savedMetadata) > 0 {
		t.Error("There were metadata records saved but non expected!")
	}
}

func TestSynchronizeLocalImageMetadataShouldAddNewMetadata(t *testing.T) {

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
	err := synchronizeLocalImageMetadata(db, fileSystemNodes, categories, testChecksumCalculator)
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

func TestSynchronizeLocalImageMetadataShouldMarkChangedEntriesAsUploads(t *testing.T) {

	categories := make(map[string]*piwigo.PiwigoCategory)
	categories["2019/shooting1"] = &piwigo.PiwigoCategory{Id: 1}

	db := NewtestStore()
	db.savedMetadata["2019/shooting1/abc.jpg"] = ImageMetaData{
		Md5Sum:         "2019/shooting1/abc.jpg",
		FullImagePath:  "2019/shooting1/abc.jpg",
		UploadRequired: false,
		LastChange:     time.Date(2019, 01, 01, 00, 0, 0, 0, time.UTC),
		Filename:       "abc.jpg",
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
	err := synchronizeLocalImageMetadata(db, fileSystemNodes, categories, testChecksumCalculator)
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

	categories := make(map[string]*piwigo.PiwigoCategory)
	categories["2019/shooting1"] = &piwigo.PiwigoCategory{Id: 1}

	db.savedMetadata["2019/shooting1/abc.jpg"] = ImageMetaData{
		Md5Sum:         "2019/shooting1/abc.jpg",
		FullImagePath:  "2019/shooting1/abc.jpg",
		UploadRequired: false,
		LastChange:     time.Date(2019, 01, 01, 01, 0, 0, 0, time.UTC),
		Filename:       "abc.jpg",
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
	err := synchronizeLocalImageMetadata(db, fileSystemNodes, categories, testChecksumCalculator)
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
	err := synchronizeLocalImageMetadata(db, fileSystemNodes, categories, testChecksumCalculator)
	if err != nil {
		t.Error(err)
	}

	if len(db.savedMetadata) > 0 {
		t.Error("There were metadata records saved but non expected!")
	}
}

func TestSynchronizePiwigoMetadata(t *testing.T) {
	db := NewtestStore()
	db.savedMetadata["2019/shooting1/abc.jpg"] = ImageMetaData{
		Md5Sum:         "2019/shooting1/abc.jpg",
		FullImagePath:  "2019/shooting1/abc.jpg",
		UploadRequired: false,
		LastChange:     time.Date(2019, 01, 01, 00, 0, 0, 0, time.UTC),
		Filename:       "abc.jpg",
	}

	// execute the sync metadata based on the file system results
	//err := synchronizeLocalImageMetadata( db)
	//if err != nil {
	//	t.Error(err)
	//}
	t.Skip("Not yet implemented!")
}

// test metadata store to store save the metadat and simulate the database
type testStore struct {
	savedMetadata map[string]ImageMetaData
}

func NewtestStore() *testStore {
	return &testStore{savedMetadata: make(map[string]ImageMetaData)}
}

func (s *testStore) ImageMetadata(fullImagePath string) (ImageMetaData, error) {
	metadata, exist := s.savedMetadata[fullImagePath]
	if !exist {
		return ImageMetaData{}, ErrorRecordNotFound
	}
	return metadata, nil
}

func (s *testStore) SaveImageMetadata(m ImageMetaData) error {
	s.savedMetadata[m.FullImagePath] = m
	return nil
}

func (d *testStore) ImageMetadataToUpload() ([]ImageMetaData, error) {
	return nil, errors.New("N/A")
}

func (d *testStore) SavePiwigoIdAndUpdateUploadFlag(md5Sum string, piwigoId int) error {
	return errors.New("N/A")
}

// to make the sync testable, we pass in a simple mock that returns the filepath as checksum
func testChecksumCalculator(file string) (string, error) {
	return file, nil
}
