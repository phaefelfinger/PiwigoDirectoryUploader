package picture

import (
	"encoding/json"
	"git.haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/piwigo"
	"github.com/sirupsen/logrus"
	"net/url"
	"strings"
)

func ImageUploadRequired(context *piwigo.PiwigoContext, md5sums []string) ([]string, error) {
	//TODO: make sure to split to multiple queries -> to honor max upload queries
	//TODO: Make sure to return the found imageIds of the found sums to update the local image nodes

	md5sumList := strings.Join(md5sums, ",")

	formData := url.Values{}
	formData.Set("method", "pwg.images.exist")
	formData.Set("md5sum_list", md5sumList)

	logrus.Tracef("Looking up missing files: %s", md5sumList)

	response, err := context.PostForm(formData)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var imageExistResponse imageExistResponse
	if err := json.NewDecoder(response.Body).Decode(&imageExistResponse); err != nil {
		logrus.Errorln(err)
		return nil, err
	}

	missingFiles := make([]string, 0, len(imageExistResponse.Result))

	for key, value := range imageExistResponse.Result {
		if value == "" {
			logrus.Tracef("Missing file with md5sum: %s", key)
			missingFiles = append(missingFiles, key)
		}
	}

	return missingFiles, nil
}
