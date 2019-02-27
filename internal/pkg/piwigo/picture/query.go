package picture

import (
	"github.com/sirupsen/logrus"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/piwigo"
)

func ImageUploadRequired(context *piwigo.PiwigoContext, files []string) (bool, error) {

	for file := range files {
		logrus.Debug(file)
	}

	/*
			http://pictures.haefelfinger.net/ws.php?format=json
		{
		    "md5sum_list": "d327416a83452b91764ed2888a5630a3,6d5f122e2b98bc1a192850e89fc2ae8c,40bfe8dd8349ccdedd4a939f9191cafa"
		}
	*/

	return false, nil
}
