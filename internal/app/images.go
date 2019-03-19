package app

import (
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/localFileStructure"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/piwigo"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

// to make use of the new local data store, we have to rethink and refactor the whole local detection process
// extend the storage of the images to keep track of upload state
// TBD: How to deal with updates -> delete / upload all based on md5 sums

type fileChecksumCalculator func(filePath string) (string, error)

// Update the local image metadata by walking through all found files and check if the modification date has changed
// or if they are new to the local database. If the files is new or changed, the md5sum will be rebuilt as well.
func synchronizeLocalImageMetadata(metadataStorage ImageMetadataProvider, fileSystemNodes map[string]*localFileStructure.FilesystemNode, checksumCalculator fileChecksumCalculator) error {
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

// This method agregates the check for files with missing piwigoids and if changed files need to be uploaded again.
func synchronizePiwigoMetadata(piwigoCtx *piwigo.PiwigoContext, metadataStorage ImageMetadataProvider) error {
	// TODO: check if category has to be assigned (image possibly added to two albums -> only uploaded once but assigned multiple times) -> implement later
	logrus.Debugf("Starting synchronizePiwigoMetadata")
	err := updatePiwigoIdIfAlreadyUploaded(metadataStorage, piwigoCtx)
	if err != nil {
		return err
	}

	err = checkPiwigoForChangedImages(metadataStorage, piwigoCtx)
	if err != nil {
		return err
	}

	return nil
}

// Check all images with upload required if they are really changed and need to be uploaded to the server.
func checkPiwigoForChangedImages(provider ImageMetadataProvider, piwigoCtx *piwigo.PiwigoContext) error {
	logrus.Infof("checking for pending files that are already on piwigo and updating piwigoids...")

	images, err := provider.ImageMetadataToUpload()
	if err != nil {
		return err
	}

	for _, img := range images {
		if img.PiwigoId == 0 {
			continue
		}
		state, err := piwigo.ImageCheckFile(piwigoCtx, img.PiwigoId, img.Md5Sum)
		if err != nil {
			logrus.Warnf("Error during file change check of file %s", img.RelativeImagePath)
			continue
		}

		if state == piwigo.ImageStateUptodate {
			logrus.Debugf("File %s - %d has not changed", img.RelativeImagePath, img.PiwigoId)
			img.UploadRequired = false
			err = provider.SaveImageMetadata(*img)
			if err != nil {
				logrus.Warnf("Could not save image data of image %s", img.RelativeImagePath)
			}
		}
	}

	return nil
}

// This function calls piwigo and checks if the given md5sum is already present.
// Only files without a piwigo id are used to query the server.
func updatePiwigoIdIfAlreadyUploaded(provider ImageMetadataProvider, piwiCtx *piwigo.PiwigoContext) error {
	logrus.Infof("checking for pending files that are already on piwigo and updating piwigoids...")
	images, err := provider.ImageMetadataToUpload()
	if err != nil {
		return err
	}

	logrus.Debugln("Preparing lookuplist for missing piwigo ids...")
	files := make([]string, 0, len(images))
	for _, img := range images {
		if img.PiwigoId == 0 {
			files = append(files, img.Md5Sum)
		}
	}
	missingResults, err := piwigo.ImagesExistOnPiwigo(piwiCtx, files)
	if err != nil {
		return err
	}
	for md5sum, piwigoId := range missingResults {
		logrus.Debugf("Setting piwigo id of %s to %d", md5sum, piwigoId)
		err = provider.SavePiwigoIdAndUpdateUploadFlag(md5sum, piwigoId)
		if err != nil {
			logrus.Warnf("Could not save piwigo id %d for file %s", piwigoId, md5sum)
		}
	}
	return nil
}

// STEP 3: Upload missing images
// - upload file in chunks
// - assign image to category

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
