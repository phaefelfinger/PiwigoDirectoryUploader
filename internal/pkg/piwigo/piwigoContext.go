/*
 * Copyright (C) 2019 Philipp Haefelfinger (http://www.haefelfinger.ch/). All Rights Reserved.
 * This application is licensed under GPLv2. See the LICENSE file in the root directory of the project.
 */

package piwigo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type PiwigoCategoryApi interface {
	GetAllCategories() (map[string]*PiwigoCategory, error)
	CreateCategory(parentId int, name string) (int, error)
}

type PiwigoImageApi interface {
	ImageCheckFile(piwigoId int, md5sum string) (int, error)
	ImagesExistOnPiwigo(md5sums []string) (map[string]int, error)
	UploadImage(piwigoId int, filePath string, md5sum string, category int) (int, error)
	DeleteImages(imageIds []int) error
}

type PiwigoContext struct {
	url           string
	username      string
	password      string
	chunkSizeInKB int
	cookies       *cookiejar.Jar
}

func (context *PiwigoContext) Initialize(baseUrl string, username string, password string) error {
	if baseUrl == "" {
		return errors.New("please provide a valid piwigo server base URL")
	}
	_, err := url.Parse(baseUrl)
	if err != nil {
		return err
	}

	if username == "" {
		return errors.New("please provide a valid username for the given piwigo server")
	}

	context.url = fmt.Sprintf("%s/ws.php?format=json", baseUrl)
	context.username = username
	context.password = password
	context.chunkSizeInKB = 512

	return nil
}

func (context *PiwigoContext) Login() error {
	logrus.Infoln("Logging in to piwigo and getting chunk size configuration for uploads")
	logrus.Debugf("Logging in to %s using user %s", context.url, context.username)

	if !strings.HasPrefix(context.url, "https") {
		logrus.Warnf("The server url %s does not use https! Credentials are not encrypted!", context.url)
	}

	formData := url.Values{}
	formData.Set("method", "pwg.session.login")
	formData.Set("username", context.username)
	formData.Set("password", context.password)

	var response loginResponse
	err := context.executePiwigoRequest(formData, &response)
	if err != nil {
		errorMessage := fmt.Sprintf("Login failed: %d - %s", response.ErrorNumber, response.Message)
		logrus.Errorln(errorMessage)
		return errors.New(errorMessage)
	}

	logrus.Infof("Login succeeded: %s", response.Status)
	return context.initializeUploadChunkSize()
}

func (context *PiwigoContext) Logout() error {
	logrus.Debugf("Logging out from %s", context.url)

	formData := url.Values{}
	formData.Set("method", "pwg.session.logout")

	var response logoutResponse
	err := context.executePiwigoRequest(formData, &response)
	if err != nil {
		logrus.Errorf("Logout from %s failed", context.url)
		return err
	}
	logrus.Infof("Successfully logged out from %s", context.url)

	return nil
}

func (context *PiwigoContext) getStatus() (*getStatusResponse, error) {
	logrus.Debugln("Getting current login state...")

	formData := url.Values{}
	formData.Set("method", "pwg.session.getStatus")

	var response getStatusResponse
	err := context.executePiwigoRequest(formData, &response)
	if err != nil {
		errorMessage := fmt.Sprintln("Could not get session state from server")
		logrus.Errorln(errorMessage)
		return nil, errors.New(errorMessage)
	}

	return &response, nil
}

func (context *PiwigoContext) GetAllCategories() (map[string]*PiwigoCategory, error) {
	formData := url.Values{}
	formData.Set("method", "pwg.categories.getList")
	formData.Set("recursive", "true")

	var response getCategoryListResponse
	err := context.executePiwigoRequest(formData, &response)
	if err != nil {
		logrus.Errorf("Got error while loading categories: %s", err)
		return nil, errors.New("could not load categories")
	}

	logrus.Infof("Successfully got all categories")
	categories := buildCategoryMap(&response)
	buildCategoryKeys(categories)
	categoryLookups := buildLookupMap(categories)

	return categoryLookups, nil
}

func (context *PiwigoContext) CreateCategory(parentId int, name string) (int, error) {
	formData := url.Values{}
	formData.Set("method", "pwg.categories.add")
	formData.Set("name", name)

	// we only submit the parentid if there is one.
	if parentId > 0 {
		formData.Set("parent", fmt.Sprint(parentId))
	}

	var response createCategoryResponse
	err := context.executePiwigoRequest(formData, &response)
	if err != nil {
		logrus.Errorln(err)
		return 0, err
	}

	logrus.Infof("Successfully created category %s with id %d", name, response.Result.ID)
	return response.Result.ID, nil
}

func (context *PiwigoContext) ImageCheckFile(piwigoId int, md5sum string) (int, error) {
	formData := url.Values{}
	formData.Set("method", "pwg.images.checkFiles")
	formData.Set("image_id", strconv.Itoa(piwigoId))
	formData.Set("file_sum", md5sum)

	logrus.Tracef("Checking if file %s - %d needs to be uploaded", md5sum, piwigoId)

	var response checkFilesResponse
	err := context.executePiwigoRequest(formData, &response)
	if err != nil {
		return imageStateInvalid, err
	}

	if response.Result["file"] == "equals" {
		return ImageStateUptodate, nil
	}
	return ImageStateDifferent, nil
}

func (context *PiwigoContext) ImagesExistOnPiwigo(md5sums []string) (map[string]int, error) {
	existResults := make(map[string]int, len(md5sums))

	batchSize := 2000
	for i := 0; i < len(md5sums); i += batchSize {
		j := i + batchSize
		if j > len(md5sums) {
			j = len(md5sums)
		}

		err := context.imagesExistOnPiwigoBatch(md5sums[i:j], existResults)
		if err != nil {
			return nil, err
		}
	}
	return existResults, nil
}

func (context *PiwigoContext) imagesExistOnPiwigoBatch(md5sums []string, existResults map[string]int) error {
	md5sumList := strings.Join(md5sums, "|")

	formData := url.Values{}
	formData.Set("method", "pwg.images.exist")
	formData.Set("md5sum_list", md5sumList)

	logrus.Tracef("Looking up if files exist: %s", md5sumList)

	var response imageExistResponse
	err := context.executePiwigoRequest(formData, &response)
	if err != nil {
		return err
	}

	for key, value := range response.Result {
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

	return nil
}

func (context *PiwigoContext) UploadImage(piwigoId int, filePath string, md5sum string, category int) (int, error) {
	if context.chunkSizeInKB <= 0 {
		return 0, errors.New("uploadchunk size is less or equal to zero. 512 is a recommendet value to begin with")
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}

	fileSizeInKB := fileInfo.Size() / 1024
	logrus.Infof("Uploading %s using chunksize of %d KB and total size of %d KB", filePath, context.chunkSizeInKB, fileSizeInKB)

	err = uploadImageChunks(filePath, context, fileSizeInKB, md5sum)
	if err != nil {
		return 0, err
	}

	imageId, err := uploadImageFinal(context, piwigoId, fileInfo.Name(), md5sum, category)
	if err != nil {
		return 0, err
	}

	return imageId, nil
}

func (context *PiwigoContext) DeleteImages(imageIds []int) error {
	logrus.Debug("Entering DeleteImages")
	defer logrus.Debug("Leaving DeleteImages")

	pwgToken, err := context.getPiwigoToken()
	if err != nil {
		return err
	}

	parts := make([]string, len(imageIds))
	for id := range imageIds {
		parts = append(parts, strconv.Itoa(id))
	}
	joinedIds := strings.Join(parts, "|")

	logrus.Infof("Deleting images: %s", joinedIds)

	formData := url.Values{}
	formData.Set("method", "pwg.images.delete")
	formData.Set("image_id", joinedIds)
	formData.Set("pwg_token", pwgToken)

	var response deleteResponse
	return context.executePiwigoRequest(formData, &response)
}

func (context *PiwigoContext) getPiwigoToken() (string, error) {
	logrus.Debug("Entering getPiwigoToken")
	defer logrus.Debug("Leaving getPiwigoToken")

	status, err := context.getStatus()
	if err != nil {
		logrus.Error("Could not get piwigo status.")
		return "", err
	}
	pwgToken := status.Result.PwgToken
	if pwgToken == "" {
		return "", errors.New("did not get a valid piwigo token. Could not delete the images")
	}
	return pwgToken, nil
}

func (context *PiwigoContext) initializeCookieJarIfRequired() {
	if context.cookies != nil {
		return
	}

	options := cookiejar.Options{}
	jar, _ := cookiejar.New(&options)
	context.cookies = jar
}

func (context *PiwigoContext) initializeUploadChunkSize() error {
	userStatus, err := context.getStatus()
	if err != nil {
		return err
	}
	context.chunkSizeInKB = userStatus.Result.UploadFormChunkSize
	logrus.Debugf("Got chunksize of %d KB from server.", context.chunkSizeInKB)
	return nil
}

func (context *PiwigoContext) executePiwigoRequest(formData url.Values, decodedResponse responseStatuser) error {
	context.initializeCookieJarIfRequired()

	client := http.Client{Jar: context.cookies}
	response, err := client.PostForm(context.url, formData)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(decodedResponse); err != nil {
		logrus.Errorln(err)
		return err
	}

	if decodedResponse.responseStatus() != "ok" {
		errorMessage := fmt.Sprintf("Error on handling piwigo response: %s", decodedResponse)
		logrus.Error(errorMessage)
		return errors.New(errorMessage)
	}
	return nil
}
