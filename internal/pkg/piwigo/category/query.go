package category

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/piwigo"
	"net/url"
	"os"
)

func GetAllCategories(context *piwigo.PiwigoContext) (map[string]*PiwigoCategory, error) {
	formData := url.Values{}
	formData.Set("method", "pwg.categories.getList")
	formData.Set("recursive", "true")

	response, err := context.PostForm(formData)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var statusResponse getCategoryListResponse
	if err := json.NewDecoder(response.Body).Decode(&statusResponse); err != nil {
		logrus.Errorln(err)
		return nil, err
	}

	if statusResponse.Status != "ok" {
		logrus.Errorf("Got state %s while loading categories", statusResponse.Status)
		return nil, errors.New("Could not load categories")
	}

	logrus.Infof("Successfully got all categories")

	categories := buildCategoryMap(&statusResponse)
	buildCategoryKeys(categories)
	categoryLookups := buildLookupMap(categories)

	return categoryLookups, nil
}

func buildLookupMap(categories map[int]*PiwigoCategory) map[string]*PiwigoCategory {
	categoryLookups := map[string]*PiwigoCategory{}
	for _, category := range categories {
		logrus.Debugf("Loaded existing category %s", category.Key)
		categoryLookups[category.Key] = category
	}
	return categoryLookups
}

func buildCategoryMap(statusResponse *getCategoryListResponse) map[int]*PiwigoCategory {
	categories := map[int]*PiwigoCategory{}
	for _, category := range statusResponse.Result.Categories {
		categories[category.ID] = &PiwigoCategory{Id: category.ID, ParentId: category.IDUppercat, Name: category.Name, Key: category.Name}
	}
	return categories
}

func buildCategoryKeys(categories map[int]*PiwigoCategory) {
	for _, category := range categories {
		if category.ParentId == 0 {
			category.Key = category.Name
			continue
		}

		key := category.Name
		parentId := category.ParentId
		for parentId != 0 {
			parent := categories[parentId]
			// as we build the category as a directory hierarchy,
			// we have to use the path separator to construct the path key
			key = fmt.Sprintf("%s%c%s", parent.Name, os.PathSeparator, key)
			parentId = parent.ParentId
		}

		category.Key = key
	}
}
