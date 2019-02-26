package category

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/piwigo"
	"net/http"
	"net/url"
)

func CreateCategory(context *piwigo.PiwigoContext, parentId int, name string) (int, error) {
	if context.Cookies == nil {
		return 0, errors.New("Not logged in and no cookies found! Can not get the category list!")
	}

	formData := url.Values{}
	formData.Set("method", "pwg.categories.add")
	formData.Set("name", name)

	// we only submit the parentid if there is one.
	if parentId > 0 {
		formData.Set("parent", fmt.Sprint(parentId))
	}

	client := http.Client{Jar: context.Cookies}
	response, err := client.PostForm(context.Url, formData)
	if err != nil {
		logrus.Errorln("The HTTP request failed with error %s", err)
		return 0, err
	}
	defer response.Body.Close()

	var createResponse createCategoryResponse
	if err := json.NewDecoder(response.Body).Decode(&createResponse); err != nil {
		logrus.Errorln(err)
		return 0, err
	}

	if createResponse.Status != "ok" {
		logrus.Errorf("Got state %s while loading categories", createResponse.Status)
		return 0, errors.New("Could not create category")
	}

	logrus.Infof("Successfully got all categories from %s", context.Url)

	return createResponse.Result.ID, nil
}
