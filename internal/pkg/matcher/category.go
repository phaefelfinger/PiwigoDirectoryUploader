package matcher

import (
	"github.com/sirupsen/logrus"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/localFileStructure"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/piwigo/category"
)

func FindMissingCategories(fileSystem map[string]*localFileStructure.FilesystemNode, categories map[string]*category.PiwigoCategory) []string {
	missingCategories := []string{}
	for _, file := range fileSystem {
		if !file.IsDir {
			continue
		}

		_, exists := categories[file.Key]

		if !exists {
			logrus.Infof("Found missing category %s", file.Key)
			missingCategories = append(missingCategories, file.Key)
		} else {
			logrus.Debugf("Found existing category %s", file.Key)
		}
	}

	return missingCategories
}
