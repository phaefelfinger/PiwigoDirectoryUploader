/*
 * Copyright (C) 2019 Philipp Haefelfinger (http://www.haefelfinger.ch/). All Rights Reserved.
 * This application is licensed under GPLv2. See the LICENSE file in the root directory of the project.
 */

package datastore

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"time"
)

var ErrorRecordNotFound = errors.New("Record not found")

type CategoryData struct {
	CategoryId     int
	PiwigoId       int
	PiwigoParentId int
	Name           string
	Key            string
}

func (cat *CategoryData) String() string {
	return fmt.Sprintf("CategoryData{CategoryId:%d, PiwigoId:%d, PiwigoParentId:%d, Name:%s, Key:%s}", cat.CategoryId, cat.PiwigoId, cat.PiwigoParentId, cat.Name, cat.Key)
}

type ImageMetaData struct {
	ImageId          int
	PiwigoId         int
	FullImagePath    string
	Filename         string
	Md5Sum           string
	LastChange       time.Time
	CategoryPath     string
	CategoryPiwigoId int
	UploadRequired   bool
	DeleteRequired   bool
}

func (img *ImageMetaData) String() string {
	return fmt.Sprintf("ImageMetaData{ImageId:%d, PiwigoId:%d, CategoryPiwigoId:%d, RelPath:%s, File:%s, Md5:%s, Change:%sS, catpath:%s, UploadRequired: %t, DeleteRequired: %t}", img.ImageId, img.PiwigoId, img.CategoryPiwigoId, img.FullImagePath, img.Filename, img.Md5Sum, img.LastChange.String(), img.CategoryPath, img.UploadRequired, img.DeleteRequired)
}

type CategoryProvider interface {
	SaveCategory(category CategoryData) error
	GetCategoryByPiwigoId(piwigoId int) (CategoryData, error)
	GetCategoryByKey(key string) (CategoryData, error)
	GetCategoriesToCreate() ([]CategoryData, error)
}

type ImageMetadataProvider interface {
	ImageMetadata(fullImagePath string) (ImageMetaData, error)
	ImageMetadataToUpload() ([]ImageMetaData, error)
	ImageMetadataToDelete() ([]ImageMetaData, error)
	ImageMetadataAll() ([]ImageMetaData, error)
	SaveImageMetadata(m ImageMetaData) error
	SavePiwigoIdAndUpdateUploadFlag(md5Sum string, piwigoId int) error
	DeleteMarkedImages() error
}

type LocalDataStore struct {
	connectionString string
}

func NewLocalDataStore() *LocalDataStore {
	return &LocalDataStore{}
}

func (d *LocalDataStore) Initialize(connectionString string) error {
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

func (d *LocalDataStore) ImageMetadata(fullImagePath string) (ImageMetaData, error) {
	logrus.Tracef("Query image metadata for file %s", fullImagePath)
	img := ImageMetaData{}

	db, err := d.openDatabase()
	if err != nil {
		return img, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT imageId, piwigoId, fullImagePath, fileName, md5sum, lastChanged, categoryPath, categoryPiwigoId, uploadRequired, deleteRequired FROM image WHERE fullImagePath = ?")
	if err != nil {
		return img, err
	}

	rows, err := stmt.Query(fullImagePath)
	if err != nil {
		return img, err
	}
	defer rows.Close()

	if rows.Next() {
		err = readImageMetadataFromRow(rows, &img)
		if err != nil {
			return img, err
		}
	} else {
		return img, ErrorRecordNotFound
	}
	err = rows.Err()

	return img, err
}

func (d *LocalDataStore) ImageMetadataAll() ([]ImageMetaData, error) {
	logrus.Tracef("Query all image metadata that represent files on the disk")

	db, err := d.openDatabase()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT imageId, piwigoId, fullImagePath, fileName, md5sum, lastChanged, categoryPath, categoryPiwigoId, uploadRequired, deleteRequired FROM image")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	images := []ImageMetaData{}
	for rows.Next() {
		img := &ImageMetaData{}
		err = readImageMetadataFromRow(rows, img)
		if err != nil {
			return nil, err
		}
		images = append(images, *img)
	}
	err = rows.Err()

	return images, err
}

func (d *LocalDataStore) ImageMetadataToDelete() ([]ImageMetaData, error) {
	logrus.Tracef("Query all image metadata that represent files queued to delete")

	db, err := d.openDatabase()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT imageId, piwigoId, fullImagePath, fileName, md5sum, lastChanged, categoryPath, categoryPiwigoId, uploadRequired, deleteRequired FROM image WHERE deleteRequired = 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	images := []ImageMetaData{}
	for rows.Next() {
		img := &ImageMetaData{}
		err = readImageMetadataFromRow(rows, img)
		if err != nil {
			return nil, err
		}
		images = append(images, *img)
	}
	err = rows.Err()

	return images, err
}

func (d *LocalDataStore) ImageMetadataToUpload() ([]ImageMetaData, error) {
	logrus.Tracef("Query all image metadata that represent files queued to upload")

	db, err := d.openDatabase()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT imageId, piwigoId, fullImagePath, fileName, md5sum, lastChanged, categoryPath, categoryPiwigoId, uploadRequired, deleteRequired FROM image WHERE uploadRequired = 1 and deleteRequired = 0 order by fullImagePath asc")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	images := []ImageMetaData{}
	for rows.Next() {
		img := &ImageMetaData{}
		err = readImageMetadataFromRow(rows, img)
		if err != nil {
			return nil, err
		}
		images = append(images, *img)
	}
	err = rows.Err()

	return images, err
}

func (d *LocalDataStore) SaveImageMetadata(img ImageMetaData) error {
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

	logrus.Tracef("Committing metadata for image %s", img.String())
	return tx.Commit()
}

func (d *LocalDataStore) SavePiwigoIdAndUpdateUploadFlag(md5Sum string, piwigoId int) error {
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
		logrus.Errorf("Rolling back transaction for piwigo id update of file %s", md5Sum)
		errTx := tx.Rollback()
		if errTx != nil {
			logrus.Errorf("Rollback of transaction for piwigo id update of file %s failed!", md5Sum)
		}
		return err
	}

	logrus.Tracef("Committing piwigo id %d for file with md5sum %s", piwigoId, md5Sum)
	return tx.Commit()
}

func (d *LocalDataStore) DeleteMarkedImages() error {
	logrus.Trace("Deleting marked records from database...")
	db, err := d.openDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM image WHERE deleteRequired = 1")
	if err != nil {
		logrus.Errorf("Rolling back transaction of deleting marked images")
		errTx := tx.Rollback()
		if errTx != nil {
			logrus.Errorf("Rollback of transaction for piwigo delete failed!")
		}
		return err
	}

	logrus.Tracef("Committing deleted images from database")
	return tx.Commit()
}

func (d *LocalDataStore) SaveCategory(category CategoryData) error {
	logrus.Tracef("Saving category: %s", category.String())
	db, err := d.openDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if category.CategoryId <= 0 {
		err = d.insertCategoryData(tx, category)
	} else {
		err = d.updateCategoryData(tx, category)
	}

	if err != nil {
		logrus.Errorf("Rolling back transaction for category of %s", category.Key)
		errTx := tx.Rollback()
		if errTx != nil {
			logrus.Errorf("Rollback of transaction for category of %s failed!", category.Key)
		}
		return err
	}

	logrus.Tracef("Committing category for image %s", category.String())
	return tx.Commit()
}

func (d *LocalDataStore) GetCategoryByPiwigoId(piwigoId int) (CategoryData, error) {
	logrus.Tracef("Query category by piwigoid %d", piwigoId)
	cat := CategoryData{}

	db, err := d.openDatabase()
	if err != nil {
		return cat, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT categoryId, piwigoId, piwigoParentId, name, key FROM category WHERE piwigoId = ?")
	if err != nil {
		return cat, err
	}

	rows, err := stmt.Query(piwigoId)
	if err != nil {
		return cat, err
	}
	defer rows.Close()

	if rows.Next() {
		err = readCategoryFromRow(rows, &cat)
		if err != nil {
			return cat, err
		}
	} else {
		return cat, ErrorRecordNotFound
	}
	err = rows.Err()

	return cat, err
}

func (d *LocalDataStore) GetCategoryByKey(key string) (CategoryData, error) {
	logrus.Tracef("Query category %s", key)
	cat := CategoryData{}

	db, err := d.openDatabase()
	if err != nil {
		return cat, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT categoryId, piwigoId, piwigoParentId, name, key FROM category WHERE key = ?")
	if err != nil {
		return cat, err
	}

	rows, err := stmt.Query(key)
	if err != nil {
		return cat, err
	}
	defer rows.Close()

	if rows.Next() {
		err = readCategoryFromRow(rows, &cat)
		if err != nil {
			return cat, err
		}
	} else {
		return cat, ErrorRecordNotFound
	}
	err = rows.Err()

	return cat, err
}

func (d *LocalDataStore) GetCategoriesToCreate() ([]CategoryData, error) {
	logrus.Trace("Query categories to create on piwigo")

	db, err := d.openDatabase()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT categoryId, piwigoId, piwigoParentId, name, key FROM category WHERE piwigoId = 0 ORDER BY key")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []CategoryData
	for rows.Next() {
		cat := CategoryData{}
		err = readCategoryFromRow(rows, &cat)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	err = rows.Err()

	return categories, err
}

func (d *LocalDataStore) openDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", d.connectionString)
	if err != nil {
		logrus.Warnf("Could not open database %s", d.connectionString)
		return nil, err
	}
	db.SetMaxOpenConns(1)

	return db, err
}

func (d *LocalDataStore) createTablesIfNeeded(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS image (" +
		"imageId INTEGER PRIMARY KEY," +
		"piwigoId INTEGER NULL," +
		"fullImagePath NVARCHAR(1000) NOT NULL," +
		"fileName NVARCHAR(255) NOT NULL," +
		"md5sum NVARCHAR(50) NOT NULL," +
		"lastChanged DATETIME NOT NULL," +
		"categoryPath NVARCHAR(1000) NOT NULL," +
		"categoryPiwigoId INTEGER NULL," +
		"uploadRequired BIT NOT NULL," +
		"deleteRequired BIT NOT NULL" +
		");")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS UX_ImageFullImagePath ON image (fullImagePath);")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS category (" +
		"categoryId INTEGER PRIMARY KEY," +
		"piwigoId INTEGER NULL," +
		"piwigoParentId INTEGER NULL," +
		"name NVARCHAR(255) NOT NULL," +
		"key NVARCHAR(1000) NOT NULL" +
		");")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS UX_Category_Key ON category (key);")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS UX_Category_PiwigoId ON category (piwigoId);")
	if err != nil {
		return err
	}

	logrus.Debug("Database successfully initialized")
	return nil
}

func readImageMetadataFromRow(rows *sql.Rows, img *ImageMetaData) error {
	err := rows.Scan(&img.ImageId, &img.PiwigoId, &img.FullImagePath, &img.Filename, &img.Md5Sum, &img.LastChange, &img.CategoryPath, &img.CategoryPiwigoId, &img.UploadRequired, &img.DeleteRequired)
	return err
}

func (d *LocalDataStore) insertImageMetaData(tx *sql.Tx, data ImageMetaData) error {
	stmt, err := tx.Prepare("INSERT INTO image (piwigoId, fullImagePath, fileName, md5sum, lastChanged, categoryPath, categoryPiwigoId, uploadRequired, deleteRequired) VALUES (?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(data.PiwigoId, data.FullImagePath, data.Filename, data.Md5Sum, data.LastChange, data.CategoryPath, data.CategoryPiwigoId, data.UploadRequired, data.DeleteRequired)
	return err
}

func (d *LocalDataStore) updateImageMetaData(tx *sql.Tx, data ImageMetaData) error {
	stmt, err := tx.Prepare("UPDATE image SET piwigoId = ?, fullImagePath = ?, fileName = ?, md5sum = ?, lastChanged = ?, categoryPath = ?, categoryPiwigoId = ?, uploadRequired = ?, deleteRequired = ? WHERE imageId = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(data.PiwigoId, data.FullImagePath, data.Filename, data.Md5Sum, data.LastChange, data.CategoryPath, data.CategoryPiwigoId, data.UploadRequired, data.DeleteRequired, data.ImageId)
	return err
}

func readCategoryFromRow(rows *sql.Rows, cat *CategoryData) error {
	err := rows.Scan(&cat.CategoryId, &cat.PiwigoId, &cat.PiwigoParentId, &cat.Name, &cat.Key)
	return err
}

func (d *LocalDataStore) updateCategoryData(tx *sql.Tx, data CategoryData) error {
	stmt, err := tx.Prepare("UPDATE category SET piwigoId = ?, piwigoParentId = ?, name = ?, key = ? WHERE categoryId = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(data.PiwigoId, data.PiwigoParentId, data.Name, data.Key, data.CategoryId)
	return err
}

func (d *LocalDataStore) insertCategoryData(tx *sql.Tx, data CategoryData) error {
	stmt, err := tx.Prepare("INSERT INTO category (piwigoId, piwigoParentId, name, key) VALUES (?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(data.PiwigoId, data.PiwigoParentId, data.Name, data.Key)
	return err
}
