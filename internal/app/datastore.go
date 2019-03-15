package app

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
	ImageId           int
	PiwigoId          int
	RelativeImagePath string
	Filename          string
	Md5Sum            string
	LastChange        time.Time
	CategoryPath      string
	CategoryId        int
}

func (img *ImageMetaData) String() string {
	return fmt.Sprintf("ImageMetaData{ImageId:%d, PiwigoId:%d, CategoryId:%d, RelPath:%s, File:%s, Md5:%s, Change:%sS, catpath:%s}", img.ImageId, img.PiwigoId, img.CategoryId, img.RelativeImagePath, img.Filename, img.Md5Sum, img.LastChange.String(), img.CategoryPath)
}

type ImageMetadataProvider interface {
	GetImageMetadata(relativePath string) (ImageMetaData, error)
	SaveImageMetadata(m ImageMetaData) error
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

func (d *localDataStore) GetImageMetadata(relativePath string) (ImageMetaData, error) {
	logrus.Debugf("Query image metadata for file %s", relativePath)
	img := ImageMetaData{}

	db, err := d.openDatabase()
	if err != nil {
		return img, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT imageId, piwigoId, relativePath, fileName, md5sum, lastChanged, categoryPath, categoryId FROM image WHERE relativePath = ?")
	if err != nil {
		return img, err
	}

	rows, err := stmt.Query(relativePath)
	if err != nil {
		return img, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&img.ImageId, &img.PiwigoId, &img.RelativeImagePath, &img.Filename, &img.Md5Sum, &img.LastChange, &img.CategoryPath, &img.CategoryId)
		if err != nil {
			return img, err
		}
	} else {
		return img, ErrorRecordNotFound
	}
	err = rows.Err()

	return img, err
}

func (d *localDataStore) SaveImageMetadata(img ImageMetaData) error {
	logrus.Debugf("Saving imagemetadata: %s", img.String())
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
		logrus.Errorf("Rolling back transaction for metadata of %s", img.RelativeImagePath)
		errTx := tx.Rollback()
		if errTx != nil {
			logrus.Errorf("Rollback of transaction for metadata of %s failed!", img.RelativeImagePath)
		}
		return err
	}

	logrus.Debugf("Commiting metadata for image %s", img.String())
	return tx.Commit()
}

func (d *localDataStore) insertImageMetaData(tx *sql.Tx, data ImageMetaData) error {
	stmt, err := tx.Prepare("INSERT INTO image (piwigoId, relativePath, fileName, md5sum, lastChanged, categoryPath, categoryId) VALUES (?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(data.PiwigoId, data.RelativeImagePath, data.Filename, data.Md5Sum, data.LastChange, data.CategoryPath, data.CategoryId)
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
		"relativePath NVARCHAR(1000) NOT NULL," +
		"fileName NVARCHAR(255) NOT NULL," +
		"md5sum NVARCHAR(50) NOT NULL," +
		"lastChanged DATETIME NOT NULL," +
		"categoryPath NVARCHAR(1000) NOT NULL," +
		"categoryId INTEGER NULL" +
		");")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS UX_ImageRelativePath ON image (relativePath);")
	return err
}

func (d *localDataStore) updateImageMetaData(tx *sql.Tx, data ImageMetaData) error {
	stmt, err := tx.Prepare("UPDATE image SET piwigoId = ?, relativePath = ?, fileName = ?, md5sum = ?, lastChanged = ?, categoryPath = ?, categoryId = ? WHERE imageId = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(data.PiwigoId, data.RelativeImagePath, data.Filename, data.Md5Sum, data.LastChange, data.CategoryPath, data.CategoryId, data.ImageId)
	return err
}
