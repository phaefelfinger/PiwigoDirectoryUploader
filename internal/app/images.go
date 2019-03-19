package app

import (
	"errors"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/localFileStructure"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/piwigo"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

type fileChecksumCalculator func(filePath string) (string, error)

// to make use of the new local data store, we have to rethink and refactor the whole local detection process
// extend the storage of the images to keep track of upload state
// TBD: How to deal with updates -> delete / upload all based on md5 sums
func synchronizeLocalImageMetadata(metadataStorage ImageMetadataProvider, fileSystemNodes map[string]*localFileStructure.FilesystemNode, checksumCalculator fileChecksumCalculator) error {
	// STEP 1 - update and sync local datastore with filesystem
	// - walk through all files of the fileSystem map
	// - get file metadata from filesystem (date, filename, dir, modtime etc.)
	// - recalculate md5 sum if file changed referring to the stored record (reduces load after first calculation a lot)
	// - mark metadata as upload required if changed or new
	logrus.Debugf("Starting synchronizeLocalImageMetadata")

	logrus.Info("Synchronizing local image metadata database with local available images")

	for _, file := range fileSystemNodes {
		if file.IsDir {
			// we are only interested in files not directories
			continue
		}

		metadata, err := metadataStorage.ImageMetadata(file.Key)
		if err == ErrorRecordNotFound {
			logrus.Debugf("No metadata for %s found. Creating new entry.", file.Key)
			metadata = ImageMetaData{}
			metadata.Filename = file.Name
			metadata.RelativeImagePath = file.Key
			metadata.CategoryPath = filepath.Dir(file.Key)
		} else if err != nil {
			logrus.Errorf("Could not get metadata due to trouble. Cancelling - %s", err)
			return err
		}

		if metadata.LastChange.Equal(file.ModTime) {
			logrus.Infof("No changed detected on file %s -> keeping current state", file.Key)
			continue
		}

		metadata.LastChange = file.ModTime
		metadata.UploadRequired = true
		metadata.Md5Sum, err = checksumCalculator(file.Path)
		if err != nil {
			logrus.Warnf("Could not calculate checksum for file %s. Skipping...", file.Path)
			continue
		}

		err = metadataStorage.SaveImageMetadata(metadata)
		if err != nil {
			return err
		}
	}

	logrus.Debugf("Finished synchronizeLocalImageMetadata")
	return nil
}

func synchronizePiwigoMetadata(piwigo *piwigo.PiwigoContext, metadataStorage ImageMetadataProvider) error  {
	// STEP 2 - get file states from piwigo (pwg.images.checkFiles)
	// - make bulk query possible
	// - get upload status of md5 sum from piwigo for all marked to upload
	// - build method to update database with new piwigoId by md5 sum
	// - set UploadRequired to false and the returned piwigoid
	// - check if category has to be assigned (image possibly added to two albums -> only uploaded once but assigned multiple times) -> implement later

	logrus.Debugf("Starting synchronizePiwigoMetadata")
	images, err := metadataStorage.ImageMetadataToUpload()
	if err != nil  {
		return err
	}

	logrus.Infof("Synchronizing metadata with existing Piwigo data...")
	for img := range images {
		logrus.Debug(img)
	}

	logrus.Debugf("Finished synchronizePiwigoMetadata")
	return errors.New("N/A")
}

// STEP 3: Upload missing images
// - upload file in chunks
// - assign image to category

func synchronizeImages(piwigo *piwigo.PiwigoContext, metadataStorage ImageMetadataProvider, existingCategories map[string]*piwigo.PiwigoCategory) error {
	//imageFiles, err := localFileStructure.GetImageList(fileSystem)
	//if err != nil {
	//	return err
	//}
	//
	//missingFiles, err := findMissingImages(context, imageFiles)
	//if err != nil {
	//	return err
	//}
	//
	//err = uploadImages(context, missingFiles, existingCategories)
	//if err != nil {
	//	return err
	//}
	//
	//logrus.Infof("Synchronized %d files.", len(missingFiles))

	return nil
}

//func findMissingImages(context *appContext, imageFiles []*localFileStructure.ImageNode) ([]*localFileStructure.ImageNode, error) {
//
//	logrus.Debugln("Preparing lookuplist for missing files...")
//
//	files := make([]string, 0, len(imageFiles))
//	md5map := make(map[string]*localFileStructure.ImageNode, len(imageFiles))
//	for _, file := range imageFiles {
//		md5map[file.Md5Sum] = file
//		files = append(files, file.Md5Sum)
//	}
//
//	missingSums, err := piwigo.ImageUploadRequired(context.piwigo, files)
//	if err != nil {
//		return nil, err
//	}
//
//	missingFiles := make([]*localFileStructure.ImageNode, 0, len(missingSums))
//	for _, sum := range missingSums {
//		file := md5map[sum]
//		logrus.Infof("Found missing file %s", file.Path)
//		missingFiles = append(missingFiles, file)
//	}
//
//	logrus.Infof("Found %d missing files", len(missingFiles))
//
//	return missingFiles, nil
//}
//
//func uploadImages(context *appContext, missingFiles []*localFileStructure.ImageNode, existingCategories map[string]*piwigo.PiwigoCategory) error {
//
//	// We sort the files by path to populate per category and not random by file
//	sort.Slice(missingFiles, func(i, j int) bool {
//		return missingFiles[i].Path < missingFiles[j].Path
//	})
//
//	for _, file := range missingFiles {
//		categoryId := existingCategories[file.CategoryName].Id
//
//		imageId, err := piwigo.UploadImage(context.piwigo, file.Path, file.Md5Sum, categoryId)
//		if err != nil {
//			return err
//		}
//		file.ImageId = imageId
//	}
//
//	return nil
//}
