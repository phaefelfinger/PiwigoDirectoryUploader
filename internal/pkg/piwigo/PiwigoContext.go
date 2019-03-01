package piwigo

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type PiwigoContext struct {
	Url           string
	Username      string
	Password      string
	ChunkSizeInKB int
	Cookies       *cookiejar.Jar
}

func (context *PiwigoContext) PostForm(formData url.Values) (resp *http.Response, err error) {
	context.initializeCookieJarIfRequired()

	client := http.Client{Jar: context.Cookies}
	response, err := client.PostForm(context.Url, formData)
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
