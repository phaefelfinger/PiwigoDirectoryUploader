package authentication

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"go/types"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

func Login(context *PiwigoContext) error {
	logrus.Debugf("Logging in to %s using user %s", context.Url, context.Username)

	if !strings.HasPrefix(context.Url, "https") {
		logrus.Warnf("The server url %s does not use https! Credentials are not encrypted!", context.Url)
	}

	initializeCookieJarIfRequired(context)

	formData := url.Values{}
	formData.Set("method", "pwg.session.login")
	formData.Set("username", context.Username)
	formData.Set("password", context.Password)

	client := http.Client{Jar: context.Cookies}

	response, err := client.PostForm(context.Url, formData)

	if err != nil {
		logrus.Errorf("The HTTP request failed with error %s", err)
		return err
	}

	var loginResponse LoginResponse
	if err := json.NewDecoder(response.Body).Decode(&loginResponse); err != nil {
		logrus.Errorln(err)
		return err
	}

	if loginResponse.Status != "ok" {
		errorMessage := fmt.Sprintf("Login failed: %d - %s", loginResponse.ErrorNumber, loginResponse.Message)
		logrus.Errorln(errorMessage)
		return types.Error{Msg: errorMessage}
	}

	logrus.Infof("Login succeeded: %s", loginResponse.Status)
	return nil
}

func Logout(context *PiwigoContext) error {
	logrus.Debugf("Logging out from %s", context.Url)

	initializeCookieJarIfRequired(context)

	formData := url.Values{}
	formData.Set("method", "pwg.session.logout")

	client := http.Client{Jar: context.Cookies}
	response, err := client.PostForm(context.Url, formData)

	if err != nil {
		logrus.Errorln("The HTTP request failed with error %s", err)
		return err
	}

	var statusResponse LogoutResponse
	if err := json.NewDecoder(response.Body).Decode(&statusResponse); err != nil {
		logrus.Errorln(err)
	}

	if statusResponse.Status != "ok" {
		logrus.Errorf("Logout from %s failed", context.Url)
	} else {
		logrus.Infof("Successfully logged out from %s", context.Url)
	}

	return nil
}

func GetStatus(context *PiwigoContext) (*GetStatusResponse, error) {

	logrus.Debugln("Getting current login state...")

	initializeCookieJarIfRequired(context)

	formData := url.Values{}
	formData.Set("method", "pwg.session.getStatus")

	client := http.Client{Jar: context.Cookies}
	response, err := client.PostForm(context.Url, formData)

	if err != nil {
		logrus.Errorln("The HTTP request failed with error %s\n", err)
		return nil, err
	}

	var statusResponse GetStatusResponse
	if err := json.NewDecoder(response.Body).Decode(&statusResponse); err != nil {
		logrus.Errorln(err)
		return nil, err
	}

	if statusResponse.Status != "ok" {
		errorMessage := fmt.Sprintf("Could not get session state from %s", context.Url)
		logrus.Errorln(errorMessage)
		return nil, types.Error{Msg: errorMessage}
	}

	return &statusResponse, nil
}

func initializeCookieJarIfRequired(context *PiwigoContext) {
	if context.Cookies != nil {
		return
	}

	options := cookiejar.Options{}
	jar, _ := cookiejar.New(&options)
	context.Cookies = jar
}
