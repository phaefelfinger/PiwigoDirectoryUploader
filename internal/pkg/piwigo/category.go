package piwigo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/url"
	"os"
)

type PiwigoCategory struct {
	Id       int
	ParentId int
	Name     string
	Key      string
}

type getCategoryListResponse struct {
	Status string `json:"stat"`
	Result struct {
		Categories []struct {
			ID                      int    `json:"id"`
			Name                    string `json:"name"`
			Comment                 string `json:"comment,omitempty"`
			Permalink               string `json:"permalink,omitempty"`
			Status                  string `json:"status,omitempty"`
			Uppercats               string `json:"uppercats,omitempty"`
			GlobalRank              string `json:"global_rank,omitempty"`
			IDUppercat              int    `json:"id_uppercat,string,omitempty"`
			NbImages                int    `json:"nb_images,omitempty"`
			TotalNbImages           int    `json:"total_nb_images,omitempty"`
			RepresentativePictureID string `json:"representative_picture_id,omitempty"`
			DateLast                string `json:"date_last,omitempty"`
			MaxDateLast             string `json:"max_date_last,omitempty"`
			NbCategories            int    `json:"nb_categories,omitempty"`
			URL                     string `json:"url,omitempty"`
			TnURL                   string `json:"tn_url,omitempty"`
		} `json:"categories"`
	} `json:"result"`
}

type createCategoryResponse struct {
	Status  string `json:"stat"`
	Err     int    `json:"err"`
	Message string `json:"message"`
	Result  struct {
		Info string `json:"info"`
		ID   int    `json:"id"`
	} `json:"result"`
}

func GetAllCategories(context *PiwigoContext) (map[string]*PiwigoCategory, error) {
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

func CreateCategory(context *PiwigoContext, parentId int, name string) (int, error) {
	formData := url.Values{}
	formData.Set("method", "pwg.categories.add")
	formData.Set("name", name)

	// we only submit the parentid if there is one.
	if parentId > 0 {
		formData.Set("parent", fmt.Sprint(parentId))
	}

	response, err := context.PostForm(formData)
	if err != nil {
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
