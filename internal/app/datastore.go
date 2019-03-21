package app

//go:generate mockgen -destination=./datastore_mock_test.go -package=app git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/app ImageMetadataProvider

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"time"
)

var ErrorRecordNotFound = errors.New("Record not found")

type ImageMetaData struct {
	ImageId        int
	PiwigoId       int
	FullImagePath  string
	Filename       string
	Md5Sum         string
	LastChange     time.Time
	CategoryPath   string
	CategoryId     int
	UploadRequired bool
}

func (img *ImageMetaData) String() string {
	return fmt.Sprintf("ImageMetaData{ImageId:%d, PiwigoId:%d, CategoryId:%d, RelPath:%s, File:%s, Md5:%s, Change:%sS, catpath:%s, UploadRequired: %t}", img.ImageId, img.PiwigoId, img.CategoryId, img.FullImagePath, img.Filename, img.Md5Sum, img.LastChange.String(), img.CategoryPath, img.UploadRequired)
}

type ImageMetadataProvider interface {
	ImageMetadata(fullImagePath string) (ImageMetaData, error)
	ImageMetadataToUpload() ([]ImageMetaData, error)
	SaveImageMetadata(m ImageMetaData) error
	SavePiwigoIdAndUpdateUploadFlag(md5Sum string, piwigoId int) error
}

type localDataStore struct {
	connectionString string
}

func (d *localDataStore) Initialize(connectionString string) error {
	if connectionString == "" {
		return errors.New("connection string could not be empty.")
	}

	d.connectionString = connectionString

	db, err := d.openDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	err = d.createTablesIfNeeded(db)

	return err
}

func (d *localDataStore) ImageMetadata(fullImagePath string) (ImageMetaData, error) {
	logrus.Tracef("Query image metadata for file %s", fullImagePath)
	img := ImageMetaData{}

	db, err := d.openDatabase()
	if err != nil {
		return img, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT imageId, piwigoId, fullImagePath, fileName, md5sum, lastChanged, categoryPath, categoryId, uploadRequired FROM image WHERE fullImagePath = ?")
	if err != nil {
		return img, err
	}

	rows, err := stmt.Query(fullImagePath)
	if err != nil {
		return img, err
	}
	defer rows.Close()

	if rows.Next() {
		err = ReadImageMetadataFromRow(rows, &img)
		if err != nil {
			return img, err
		}
	} else {
		return img, ErrorRecordNotFound
	}
	err = rows.Err()

	return img, err
}

func (d *localDataStore) ImageMetadataToUpload() ([]ImageMetaData, error) {
	logrus.Tracef("Query all image metadata that represent files queued to upload")

	db, err := d.openDatabase()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT imageId, piwigoId, fullImagePath, fileName, md5sum, lastChanged, categoryPath, categoryId, uploadRequired FROM image WHERE uploadRequired = 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	images := []ImageMetaData{}
	for rows.Next() {
		img := &ImageMetaData{}
		err = ReadImageMetadataFromRow(rows, img)
		if err != nil {
			return nil, err
		}
		images = append(images, *img)
	}
	err = rows.Err()

	return images, err
}

func ReadImageMetadataFromRow(rows *sql.Rows, img *ImageMetaData) error {
	err := rows.Scan(&img.ImageId, &img.PiwigoId, &img.FullImagePath, &img.Filename, &img.Md5Sum, &img.LastChange, &img.CategoryPath, &img.CategoryId, &img.UploadRequired)
	return err
}

func (d *localDataStore) SaveImageMetadata(img ImageMetaData) error {
	logrus.Tracef("Saving imagemetadata: %s", img.String())
	db, err := d.openDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if img.ImageId <= 0 {
		err = d.insertImageMetaData(tx, img)
	} else {
		err = d.updateImageMetaData(tx, img)
	}

	if err != nil {
		logrus.Errorf("Rolling back transaction for metadata of %s", img.FullImagePath)
		errTx := tx.Rollback()
		if errTx != nil {
			logrus.Errorf("Rollback of transaction for metadata of %s failed!", img.FullImagePath)
		}
		return err
	}

	logrus.Tracef("Commiting metadata for image %s", img.String())
	return tx.Commit()
}

func (d *localDataStore) SavePiwigoIdAndUpdateUploadFlag(md5Sum string, piwigoId int) error {
	logrus.Tracef("Saving piwigo id %d for file with md5sum %s", piwigoId, md5Sum)
	db, err := d.openDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	uploadRequired := 1
	if piwigoId > 0 {
		uploadRequired = 0
	}

	stmt, err := tx.Prepare("UPDATE image SET piwigoId = ?, uploadRequired = ? WHERE md5sum = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(piwigoId, uploadRequired, md5Sum)
	if err != nil {
		return err
	}

	if err != nil {
		logrus.Errorf("Rolling back transaction for piwigo id update of file %s", md5Sum)
		errTx := tx.Rollback()
		if errTx != nil {
			logrus.Errorf("Rollback of transaction for piwigo id update of file %s failed!", md5Sum)
		}
		return err
	}

	logrus.Tracef("Commiting piwigo id %d for file with md5sum %s", piwigoId, md5Sum)
	return tx.Commit()
}

func (d *localDataStore) insertImageMetaData(tx *sql.Tx, data ImageMetaData) error {
	stmt, err := tx.Prepare("INSERT INTO image (piwigoId, fullImagePath, fileName, md5sum, lastChanged, categoryPath, categoryId, uploadRequired) VALUES (?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(data.PiwigoId, data.FullImagePath, data.Filename, data.Md5Sum, data.LastChange, data.CategoryPath, data.CategoryId, data.UploadRequired)
	return err
}

func (d *localDataStore) openDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", d.connectionString)
	if err != nil {
		logrus.Warnf("Could not open database %s", d.connectionString)
		return nil, err
	}
	db.SetMaxOpenConns(1)

	return db, err
}

func (d *localDataStore) createTablesIfNeeded(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS image (" +
		"imageId INTEGER PRIMARY KEY," +
		"piwigoId INTEGER NULL," +
		"fullImagePath NVARCHAR(1000) NOT NULL," +
		"fileName NVARCHAR(255) NOT NULL," +
		"md5sum NVARCHAR(50) NOT NULL," +
		"lastChanged DATETIME NOT NULL," +
		"categoryPath NVARCHAR(1000) NOT NULL," +
		"categoryId INTEGER NULL," +
		"uploadRequired BIT NOT NULL" +
		");")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS UX_ImageFullImagePath ON image (fullImagePath);")
	return err
}

func (d *localDataStore) updateImageMetaData(tx *sql.Tx, data ImageMetaData) error {
	stmt, err := tx.Prepare("UPDATE image SET piwigoId = ?, fullImagePath = ?, fileName = ?, md5sum = ?, lastChanged = ?, categoryPath = ?, categoryId = ?, uploadRequired = ? WHERE imageId = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(data.PiwigoId, data.FullImagePath, data.Filename, data.Md5Sum, data.LastChange, data.CategoryPath, data.CategoryId, data.UploadRequired, data.ImageId)
	return err
}
