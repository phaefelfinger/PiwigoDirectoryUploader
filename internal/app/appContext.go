package app

import (
	"errors"
	"git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/piwigo"
	"github.com/sirupsen/logrus"
)

type appContext struct {
	// think again if this is a good idea to have such a context!
	piwigo        *piwigo.PiwigoContext
	dataStore     *localDataStore
	sessionId     string
	localRootPath string
}

func (c *appContext) UseMetadataStore(connectionString string) error {
	if connectionString == "" {
		return errors.New("missing connectionString to use metadata store!")
	}

	logrus.Infof("Using SQL Lite data store with '%s'", connectionString)
	c.dataStore = &localDataStore{}
	err := c.dataStore.Initialize(connectionString)

	return err
}

func (c *appContext) UsePiwigo(url string, user string, password string) error {
	if url == "" {
		return errors.New("missing piwigo url!")
	}

	if user == "" {
		return errors.New("missing piwigo user!")
	}

	if password == "" {
		return errors.New("missing piwigo password!")
	}

	c.piwigo = new(piwigo.PiwigoContext)
	return c.piwigo.Initialize(*piwigoUrl, *piwigoUser, *piwigoPassword)
}

func newAppContext() (*appContext, error) {
	logrus.Infoln("Preparing application context and configuration")

	context := new(appContext)
	context.localRootPath = *imagesRootPath

	if *sqliteDb != "" {
		err := context.UseMetadataStore(*sqliteDb)
		if err != nil {
			return nil, err
		}
	} else {
		logrus.Warnln("No persistence configured. Skipping metadata storage. This might affect performance on large collections!")
	}

	err := context.UsePiwigo(*piwigoUrl, *piwigoUser, *piwigoPassword)

	return context, err
}
