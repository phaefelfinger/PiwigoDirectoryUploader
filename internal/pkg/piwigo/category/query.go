package category

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/piwigo"
	"net/http"
	"net/url"
)

func GetAllCategories(context *piwigo.PiwigoContext) error {
	logrus.Debugln("Starting GetAllCategories")
	if context.Cookies == nil {
		return errors.New("Not logged in and no cookies found! Can not get the category list!")
	}

	formData := url.Values{}
	formData.Set("method", "pwg.categories.getList")
	formData.Set("recursive", "true")

	client := http.Client{Jar: context.Cookies}
	response, err := client.PostForm(context.Url, formData)

	if err != nil {
		logrus.Errorln("The HTTP request failed with error %s", err)
		return err
	}

	var statusResponse getCategoryListResponse
	if err := json.NewDecoder(response.Body).Decode(&statusResponse); err != nil {
		logrus.Errorln(err)
		return err
	}

	if statusResponse.Status != "ok" {
		logrus.Errorf("Got state %s while loading categories", statusResponse.Status)
		return errors.New("Could not load categories")
	}

	logrus.Infof("Successfully got all categories from %s", context.Url)

	for _, category := range statusResponse.Result.Categories {
		logrus.Debugf("Category %d - %s", category.ID, category.Name)
	}

	return nil
}
