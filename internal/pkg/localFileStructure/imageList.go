package localFileStructure

import (
	"github.com/sirupsen/logrus"
	"time"
)

func GetImageList(fileSystem map[string]*FilesystemNode) ([]*ImageNode, error) {
	imageFiles := []*ImageNode{}

	for _, file := range fileSystem {
		if file.IsDir {
			continue
		}

		md5sum, err := calculateFileCheckSums(file.Path)
		if err != nil {
			return nil, err
		}

		logrus.Debugf("Image %s - %s - %s", md5sum, file.ModTime.Format(time.RFC3339), file.Path)

		imageFiles = append(imageFiles, &ImageNode{
			Path:    file.Path,
			ModTime: file.ModTime,
			Md5Sum:  md5sum,
		})
	}

	return imageFiles, nil
}
