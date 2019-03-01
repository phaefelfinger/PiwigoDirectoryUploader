package picture

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"git.haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/piwigo"
	"net/url"
	"strconv"
)

func UploadImage(context *piwigo.PiwigoContext, filePath string, md5sum string, category int) (int, error) {
	logrus.Infof("Uploading %s", filePath)

	// split into chunks
	// upload chunks
	// finalize upload

	return 0, nil
}

func uploadImageChunk(context *piwigo.PiwigoContext, base64chunk string, md5sum string, position int) error {
	formData := url.Values{}
	formData.Set("method", "pwg.images.addChunk")
	formData.Set("data", base64chunk)
	formData.Set("original_sum", md5sum)
	// required by the API for compatibility
	formData.Set("type", "file")
	formData.Set("position", strconv.Itoa(position))

	logrus.Tracef("Uploading chunk %d of file with sum %s", position, md5sum)

	response, err := context.PostForm(formData)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var uploadChunkResponse uploadChunkResponse
	if err := json.NewDecoder(response.Body).Decode(&uploadChunkResponse); err != nil {
		logrus.Errorln(err)
		return err
	}

	if uploadChunkResponse.Status != "ok" {
		logrus.Errorf("Got state %s while uploading chunk %d of %s", uploadChunkResponse.Status, position, md5sum)
		return errors.New(fmt.Sprintf("Got state %s while uploading chunk %d of %s", uploadChunkResponse.Status, position, md5sum))
	}

	return nil
}

func uploadImageFinal(context *piwigo.PiwigoContext, originalFilename string, md5sum string, categoryId int) error {
	formData := url.Values{}
	formData.Set("method", "pwg.images.add")
	formData.Set("original_sum", md5sum)
	formData.Set("original_filename", originalFilename)
	formData.Set("categoriesi", strconv.Itoa(categoryId))

	logrus.Debugf("Finalizing upload of file %s with sum %s to category %d", originalFilename, md5sum, categoryId)

	response, err := context.PostForm(formData)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var fileAddResponse fileAddResponse
	if err := json.NewDecoder(response.Body).Decode(&fileAddResponse); err != nil {
		logrus.Errorln(err)
		return err
	}

	if fileAddResponse.Status != "ok" {
		logrus.Errorf("Got state %s while adding image %s", fileAddResponse.Status, originalFilename)
		return errors.New(fmt.Sprintf("Got state %s while adding image %s", fileAddResponse.Status, originalFilename))
	}

	return nil
}
