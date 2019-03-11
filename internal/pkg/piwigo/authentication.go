package piwigo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/url"
	"strings"
)

type LoginResponse struct {
	Status      string `json:"stat"`
	Result      bool   `json:"result"`
	ErrorNumber int    `json:"err"`
	Message     string `json:"message"`
}

type GetStatusResponse struct {
	Status string `json:"stat"`
	Result struct {
		Username            string   `json:"username"`
		Status              string   `json:"status"`
		Theme               string   `json:"theme"`
		Language            string   `json:"language"`
		PwgToken            string   `json:"pwg_token"`
		Charset             string   `json:"charset"`
		CurrentDatetime     string   `json:"current_datetime"`
		Version             string   `json:"version"`
		AvailableSizes      []string `json:"available_sizes"`
		UploadFileTypes     string   `json:"upload_file_types"`
		UploadFormChunkSize int      `json:"upload_form_chunk_size"`
	} `json:"result"`
}

type LogoutResponse struct {
	Status string `json:"stat"`
	Result bool   `json:"result"`
}

func Login(context *PiwigoContext) error {
	logrus.Debugf("Logging in to %s using user %s", context.Url, context.Username)

	if !strings.HasPrefix(context.Url, "https") {
		logrus.Warnf("The server url %s does not use https! Credentials are not encrypted!", context.Url)
	}

	formData := url.Values{}
	formData.Set("method", "pwg.session.login")
	formData.Set("username", context.Username)
	formData.Set("password", context.Password)

	response, err := context.PostForm(formData)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var loginResponse LoginResponse
	if err := json.NewDecoder(response.Body).Decode(&loginResponse); err != nil {
		logrus.Errorln(err)
		return err
	}

	if loginResponse.Status != "ok" {
		errorMessage := fmt.Sprintf("Login failed: %d - %s", loginResponse.ErrorNumber, loginResponse.Message)
		logrus.Errorln(errorMessage)
		return errors.New(errorMessage)
	}

	logrus.Infof("Login succeeded: %s", loginResponse.Status)
	return nil
}

func Logout(context *PiwigoContext) error {
	logrus.Debugf("Logging out from %s", context.Url)

	formData := url.Values{}
	formData.Set("method", "pwg.session.logout")

	response, err := context.PostForm(formData)
	if err != nil {
		return err
	}
	defer response.Body.Close()

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

	formData := url.Values{}
	formData.Set("method", "pwg.session.getStatus")

	response, err := context.PostForm(formData)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var statusResponse GetStatusResponse
	if err := json.NewDecoder(response.Body).Decode(&statusResponse); err != nil {
		logrus.Errorln(err)
		return nil, err
	}

	if statusResponse.Status != "ok" {
		errorMessage := fmt.Sprintf("Could not get session state from %s", context.Url)
		logrus.Errorln(errorMessage)
		return nil, errors.New(errorMessage)
	}

	return &statusResponse, nil
}
