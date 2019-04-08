/*
 * Copyright (C) 2019 Philipp Haefelfinger (http://www.haefelfinger.ch/). All Rights Reserved.
 * This application is licensed under GPLv2. See the LICENSE file in the root directory of the project.
 */

package piwigo

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

type Category struct {
	Id       int
	ParentId int
	Name     string
	Key      string
}

func buildLookupMap(categories map[int]*Category) map[string]*Category {
	categoryLookups := map[string]*Category{}
	for _, category := range categories {
		logrus.Debugf("Loaded existing category %s", category.Key)
		categoryLookups[category.Key] = category
	}
	return categoryLookups
}

func buildCategoryMap(statusResponse *getCategoryListResponse) map[int]*Category {
	categories := map[int]*Category{}
	for _, category := range statusResponse.Result.Categories {
		categories[category.ID] = &Category{Id: category.ID, ParentId: category.IDUppercat, Name: category.Name, Key: category.Name}
	}
	return categories
}

func buildCategoryKeys(categories map[int]*Category) {
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
