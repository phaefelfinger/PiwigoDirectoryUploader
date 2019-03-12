package app

import (
	"errors"
	"time"
)

type ImageMetaData struct {
	ImageId           int
	RelativeImagePath string
	Filename          string
	Md5Sum            string
	LastChange        time.Time
	CategoryPath      string
	CategoryId        int
}

type ImageMetadataLoader interface {
	GetImageMetadata(relativePath string) (ImageMetaData, error)
}

type ImageMetadataSaver interface {
	SaveImageMetadata(m ImageMetaData) error
}

type localDataStore struct {
	connectionString string
}

func (d *localDataStore) Open(connectionString string) error {
	if connectionString == "" {
		return errors.New("connection string could not be empty.")
	}

	d.connectionString = connectionString

	//TODO: open and test connection
	return nil
}

func (d *localDataStore) GetImageMetadata(relativePath string) (ImageMetaData, error) {
	return ImageMetaData{}, nil
}

func (d *localDataStore) SaveImageMetadata(m ImageMetaData) error {
	return nil
}
