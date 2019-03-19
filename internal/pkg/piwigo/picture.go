package piwigo

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type uploadChunkResponse struct {
	Status string      `json:"stat"`
	Result interface{} `json:"result"`
}

type fileAddResponse struct {
	Status string `json:"stat"`
	Result struct {
		ImageID int    `json:"image_id"`
		URL     string `json:"url"`
	} `json:"result"`
}

type imageExistResponse struct {
	Status string            `json:"stat"`
	Result map[string]string `json:"result"`
}

type checkFilesResponse struct {
	Status string `json:"stat"`
	Result struct {
		File string `json:"file"`
	} `json:"result"`
}

const (
	ImageStateInvalid   = -1
	ImageStateUptodate  = 0
	ImageStateDifferent = 1
)

func ImageCheckFile(context PiwigoFormPoster, piwigoId int, md5sum string) (int, error) {
	formData := url.Values{}
	formData.Set("method", "pwg.images.exist")
	formData.Set("image_id", strconv.Itoa(piwigoId))
	formData.Set("file_sum", md5sum)

	logrus.Tracef("Checking if file %s - %d needs to be uploaded", md5sum, piwigoId)

	response, err := context.postForm(formData)
	if err != nil {
		return ImageStateInvalid, err
	}
	defer response.Body.Close()

	var checkFilesResponse checkFilesResponse
	if err := json.NewDecoder(response.Body).Decode(&checkFilesResponse); err != nil {
		logrus.Errorln(err)
		return ImageStateInvalid, err
	}

	if checkFilesResponse.Result.File == "equals" {
		return ImageStateUptodate, nil
	}
	return ImageStateDifferent, nil
}

func ImagesExistOnPiwigo(context PiwigoFormPoster, md5sums []string) (map[string]int, error) {
	//TODO: make sure to split to multiple queries -> to honor max upload queries
	md5sumList := strings.Join(md5sums, ",")

	formData := url.Values{}
	formData.Set("method", "pwg.images.exist")
	formData.Set("md5sum_list", md5sumList)

	logrus.Tracef("Looking up if files exist: %s", md5sumList)

	response, err := context.postForm(formData)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var imageExistResponse imageExistResponse
	if err := json.NewDecoder(response.Body).Decode(&imageExistResponse); err != nil {
		logrus.Errorln(err)
		return nil, err
	}

	existResults := make(map[string]int, len(imageExistResponse.Result))

	for key, value := range imageExistResponse.Result {
		if value == "" {
			logrus.Tracef("Missing file with md5sum: %s", key)
			existResults[key] = 0
		} else {
			piwigoId, err := strconv.Atoi(value)
			if err != nil {
				logrus.Warnf("could not parse piwigoid of file %s", key)
				continue
			}
			logrus.Tracef("Found piwigo id %d for md5sum %s", piwigoId, key)
			existResults[key] = piwigoId
		}
	}

	return existResults, nil
}

func UploadImage(context PiwigoFormPoster, filePath string, md5sum string, category int) (int, error) {
	if context.getChunkSizeInKB() <= 0 {
		return 0, errors.New("Uploadchunk size is less or equal to zero. 512 is a recommendet value to begin with.")
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}

	fileSizeInKB := fileInfo.Size() / 1024
	logrus.Infof("Uploading %s using chunksize of %d KB and total size of %d KB", filePath, context.getChunkSizeInKB(), fileSizeInKB)

	err = uploadImageChunks(filePath, context, fileSizeInKB, md5sum)
	if err != nil {
		return 0, err
	}

	imageId, err := uploadImageFinal(context, fileInfo.Name(), md5sum, category)
	if err != nil {
		return 0, err
	}

	return imageId, nil
}

func uploadImageChunks(filePath string, context PiwigoFormPoster, fileSizeInKB int64, md5sum string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	bufferSize := 1024 * context.getChunkSizeInKB()
	buffer := make([]byte, bufferSize)
	numberOfChunks := (fileSizeInKB / int64(context.getChunkSizeInKB())) + 1
	currentChunk := int64(0)

	for {
		logrus.Tracef("Processing chunk %d of %d of %s", currentChunk, numberOfChunks, filePath)

		readBytes, readError := reader.Read(buffer)
		if readError == io.EOF && readBytes == 0 {
			break
		}
		if readError != io.EOF && readError != nil {
			return readError
		}

		encodedChunk := base64.StdEncoding.EncodeToString(buffer[:readBytes])

		uploadError := uploadImageChunk(context, encodedChunk, md5sum, currentChunk)
		if uploadError != nil {
			return uploadError
		}

		currentChunk++
	}

	return nil
}

func uploadImageChunk(context PiwigoFormPoster, base64chunk string, md5sum string, position int64) error {
	formData := url.Values{}
	formData.Set("method", "pwg.images.addChunk")
	formData.Set("data", base64chunk)
	formData.Set("original_sum", md5sum)
	// required by the API for compatibility
	formData.Set("type", "file")
	formData.Set("position", strconv.FormatInt(position, 10))

	logrus.Tracef("Uploading chunk %d of file with sum %s", position, md5sum)

	response, err := context.postForm(formData)
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

func uploadImageFinal(context PiwigoFormPoster, originalFilename string, md5sum string, categoryId int) (int, error) {
	formData := url.Values{}
	formData.Set("method", "pwg.images.add")
	formData.Set("original_sum", md5sum)
	formData.Set("original_filename", originalFilename)
	formData.Set("name", originalFilename)
	formData.Set("categories", strconv.Itoa(categoryId))

	logrus.Debugf("Finalizing upload of file %s with sum %s to category %d", originalFilename, md5sum, categoryId)

	response, err := context.postForm(formData)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	var fileAddResponse fileAddResponse
	if err := json.NewDecoder(response.Body).Decode(&fileAddResponse); err != nil {
		logrus.Errorln(err)
		return 0, err
	}

	if fileAddResponse.Status != "ok" {
		logrus.Errorf("Got state %s while adding image %s", fileAddResponse.Status, originalFilename)
		return 0, errors.New(fmt.Sprintf("Got state %s while adding image %s", fileAddResponse.Status, originalFilename))
	}

	return fileAddResponse.Result.ImageID, nil
}
