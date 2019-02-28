package localFileStructure

import (
	"github.com/sirupsen/logrus"
	"path/filepath"
	"time"
)

func GetImageList(fileSystem map[string]*FilesystemNode) ([]*ImageNode, error) {
	imageFiles := make([]*ImageNode, 0, len(fileSystem))

	for _, file := range fileSystem {
		if file.IsDir {
			continue
		}

		md5sum, err := calculateFileCheckSums(file.Path)
		if err != nil {
			return nil, err
		}

		logrus.Debugf("Local Image %s - %s - %s", md5sum, file.ModTime.Format(time.RFC3339), file.Path)

		imageFiles = append(imageFiles, &ImageNode{
			Path:         file.Path,
			CategoryName: filepath.Dir(file.Key),
			ModTime:      file.ModTime,
			Md5Sum:       md5sum,
		})
	}

	logrus.Infof("Found %d local images to process", len(imageFiles))

	return imageFiles, nil
}
