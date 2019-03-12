package piwigo

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type PiwigoFormPoster interface {
	getChunkSizeInKB() int
	postForm(formData url.Values) (resp *http.Response, err error)
}

type PiwigoContext struct {
	url           string
	username      string
	password      string
	chunkSizeInKB int
	Cookies       *cookiejar.Jar
}

func (context *PiwigoContext) Initialize(baseUrl string, username string, password string, chunkSizeInKB int) error {
	if baseUrl == "" {
		return errors.New("Please provide a valid piwigo server base URL")
	}
	_, err := url.Parse(baseUrl)
	if err != nil {
		return err
	}

	if username == "" {
		return errors.New("Please provide a valid username for the given piwigo server.")
	}

	if chunkSizeInKB < 256 {
		return errors.New("The minimum chunksize is 256KB. Please provide a value above. Default is 512KB")
	}

	context.url = fmt.Sprintf("%s/ws.php?format=json", baseUrl)
	context.username = username
	context.password = password
	context.chunkSizeInKB = chunkSizeInKB

	return nil
}

func (context *PiwigoContext) LoginToPiwigoAndConfigureContext() error {
	logrus.Infoln("Logging in to piwigo and getting chunk size configuration for uploads")
	err := Login(context)
	if err != nil {
		return err
	}
	return initializeUploadChunkSize(context)
}

func (context *PiwigoContext) getChunkSizeInKB() int {
	return context.chunkSizeInKB
}

func (context *PiwigoContext) postForm(formData url.Values) (resp *http.Response, err error) {
	context.initializeCookieJarIfRequired()

	client := http.Client{Jar: context.Cookies}
	response, err := client.PostForm(context.url, formData)
	if err != nil {
		logrus.Errorln("The HTTP request failed with error %s", err)
		return nil, err
	}
	return response, nil
}

func (context *PiwigoContext) initializeCookieJarIfRequired() {
	if context.Cookies != nil {
		return
	}

	options := cookiejar.Options{}
	jar, _ := cookiejar.New(&options)
	context.Cookies = jar
}

func initializeUploadChunkSize(context *PiwigoContext) error {
	userStatus, err := GetStatus(context)
	if err != nil {
		return err
	}
	context.chunkSizeInKB = userStatus.Result.UploadFormChunkSize * 1024
	logrus.Debugf("Got chunksize of %d KB from server.", context.chunkSizeInKB)
	return nil
}
