package app

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

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

type ImageMetadataLoader interface {
	GetImageMetadata(relativePath string) (ImageMetaData, error)
}

type ImageMetadataSaver interface {
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
	db, err := d.openDatabase()
	if err != nil {
		return ImageMetaData{}, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return ImageMetaData{}, err
	}

	//TODO: select entry by path
	//stmt, err := tx.Prepare("select * from image WHERE relativePath = '?'")
	//if err != nil {
	//	log.Fatal(err)
	//}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return ImageMetaData{}, nil
}

func (d *localDataStore) SaveImageMetadata(m ImageMetaData) error {
	db, err := d.openDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if m.ImageId <= 0 {
		err = d.insertImageMetaData(tx, m)
		if err != nil {
			return err
		}
	} else {
		// TODO: update existing entry
	}

	err = tx.Commit()
	return err
}

func (d *localDataStore) insertImageMetaData(tx *sql.Tx, m ImageMetaData) error {
	stmt, err := tx.Prepare("INSERT INTO image (piwigoId, relativePath, fileName, md5sum, lastChanged, categoryPath, categoryId) VALUES (?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(m.PiwigoId, m.RelativeImagePath, m.Filename, m.Md5Sum, m.LastChange, m.CategoryPath, m.CategoryId)
	return err
}

func (d *localDataStore) openDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", d.connectionString)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)

	return db, err
}

func (d *localDataStore) createTablesIfNeeded(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS image (" +
		"imageId INTEGER PRIMARY KEY AUTOINCREMENT," +
		"piwigoId INTEGER NULL," +
		"relativePath NVARCHAR(1000) NOT NULL," +
		"fileName NVARCHAR(255) NOT NULL," +
		"md5sum NVARCHAR(50) NOT NULL," +
		"lastChanged DATETIME NOT NULL," +
		"categoryPath NVARCHAR(1000) NOT NULL," +
		"categoryId INTEGER NULL" +
		");")
	return err
}
