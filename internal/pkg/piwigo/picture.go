/*
 * Copyright (C) 2019 Philipp Haefelfinger (http://www.haefelfinger.ch/). All Rights Reserved.
 * This application is licensed under GPLv2. See the LICENSE file in the root directory of the project.
 */

package piwigo

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/url"
	"os"
	"strconv"
)

const (
	imageStateInvalid   = -1
	ImageStateUptodate  = 0
	ImageStateDifferent = 1
)

func uploadImageChunks(filePath string, context *PiwigoContext, fileSizeInKB int64, md5sum string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	bufferSize := 1024 * context.chunkSizeInKB
	buffer := make([]byte, bufferSize)
	numberOfChunks := (fileSizeInKB / int64(context.chunkSizeInKB)) + 1
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

func uploadImageChunk(context *PiwigoContext, base64chunk string, md5sum string, position int64) error {
	formData := url.Values{}
	formData.Set("method", "pwg.images.addChunk")
	formData.Set("data", base64chunk)
	formData.Set("original_sum", md5sum)
	// required by the API for compatibility
	formData.Set("type", "file")
	formData.Set("position", strconv.FormatInt(position, 10))

	logrus.Tracef("Uploading chunk %d of file with sum %s", position, md5sum)

	var response uploadChunkResponse
	err := context.executePiwigoRequest(formData, &response)
	if err != nil {
		logrus.Errorf("Got state %s while uploading chunk %d of %s", response.Status, position, md5sum)
		return errors.New(fmt.Sprintf("Got state %s while uploading chunk %d of %s", response.Status, position, md5sum))
	}

	return nil
}

func uploadImageFinal(context *PiwigoContext, piwigoId int, originalFilename string, md5sum string, categoryId int) (int, error) {
	formData := url.Values{}
	formData.Set("method", "pwg.images.add")
	formData.Set("original_sum", md5sum)
	formData.Set("original_filename", originalFilename)
	formData.Set("name", originalFilename)
	formData.Set("categories", strconv.Itoa(categoryId))

	// when there is a image id, we are updating an existing image and need to specify the piwigo image id.
	// if we skip the image id, a new id will be generated
	if piwigoId > 0 {
		formData.Set("image_id", strconv.Itoa(piwigoId))
	}

	logrus.Debugf("Finalizing upload of file %s with sum %s to category %d", originalFilename, md5sum, categoryId)

	var response fileAddResponse
	err := context.executePiwigoRequest(formData, &response)
	if err != nil {
		logrus.Errorf("Got state %s while adding image %s", response.Status, originalFilename)
		return 0, errors.New(fmt.Sprintf("Got state %s while adding image %s", response.Status, originalFilename))
	}

	return response.Result.ImageID, nil
}
