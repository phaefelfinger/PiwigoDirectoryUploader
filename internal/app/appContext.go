/*
 * Copyright (C) 2019 Philipp Haefelfinger (http://www.haefelfinger.ch/). All Rights Reserved.
 * This application is licensed under GPLv2. See the LICENSE file in the root directory of the project.
 */

package app

import (
	"errors"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/datastore"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/piwigo"
	"github.com/sirupsen/logrus"
)

type appContext struct {
	// think again if this is a good idea to have such a context!
	piwigo        *piwigo.ServerContext
	dataStore     *datastore.LocalDataStore
	sessionId     string
	localRootPath string
}

func (c *appContext) useMetadataStore(connectionString string) error {
	if connectionString == "" {
		return errors.New("missing connectionString to use metadata store")
	}

	logrus.Infof("Using SQL Lite data store with '%s'", connectionString)
	c.dataStore = datastore.NewLocalDataStore()
	err := c.dataStore.Initialize(connectionString)

	return err
}

func (c *appContext) usePiwigo(url string, user string, password string) error {
	if url == "" {
		return errors.New("missing piwigo url")
	}

	if user == "" {
		return errors.New("missing piwigo user")
	}

	if password == "" {
		return errors.New("missing piwigo password")
	}

	c.piwigo = new(piwigo.ServerContext)
	return c.piwigo.Initialize(url, user, password)
}

func newAppContext() (*appContext, error) {
	logrus.Infoln("Preparing application context and configuration")

	context := new(appContext)
	context.localRootPath = *imagesRootPath

	if *sqliteDb != "" {
		err := context.useMetadataStore(*sqliteDb)
		if err != nil {
			return nil, err
		}
	} else {
		logrus.Warnln("No persistence configured. Skipping metadata storage. This might affect performance on large collections!")
	}

	err := context.usePiwigo(*piwigoUrl, *piwigoUser, *piwigoPassword)

	return context, err
}
