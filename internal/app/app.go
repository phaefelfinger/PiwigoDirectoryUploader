package app

import (
	"errors"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/localFileStructure"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/matcher"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/piwigo"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/piwigo/authentication"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/piwigo/category"
	"os"
)

var (
	imagesRootPath = flag.String("imagesRootPath", "", "This is the images root path that should be mirrored to piwigo.")
	piwigoUrl      = flag.String("piwigoUrl", "", "The root url without tailing slash to your piwigo installation.")
	piwigoUser     = flag.String("piwigoUser", "", "The username to use during sync.")
	piwigoPassword = flag.String("piwigoPassword", "", "This is password to the given username.")
)

func Run() {
	context, err := configureContext()
	if err != nil {
		os.Exit(1)
	}

	err = loginToPiwigoAndConfigureContext(context)
	if err != nil {
		os.Exit(2)
	}

	filesystemNodes := scanLocalDirectories(context)
	categories := getAllCategoriesFromServer(context)

	synchronizeCategories(filesystemNodes, categories)


	findMissingImages()
	uploadImages()

	_ = authentication.Logout(context.Piwigo)
}

func synchronizeCategories(filesystemNodes map[string]*localFileStructure.FilesystemNode, categories map[string]*category.PiwigoCategory) {
	missingCategories := findMissingCategories(filesystemNodes, categories)
	createMissingCategories(missingCategories)
}

func scanLocalDirectories(context *AppContext) map[string]*localFileStructure.FilesystemNode {
	fileNodes, err := localFileStructure.ScanLocalFileStructure(context.LocalRootPath)
	if err != nil {
		os.Exit(3)
	}
	return fileNodes
}

func getAllCategoriesFromServer(context *AppContext) map[string]*category.PiwigoCategory {

	categories, err := category.GetAllCategories(context.Piwigo)
	if err != nil {
		os.Exit(4)
	}

	return categories
}

func findMissingCategories(fileSystem map[string]*localFileStructure.FilesystemNode, categories map[string]*category.PiwigoCategory) []string {
	return matcher.FindMissingCategories(fileSystem, categories)
}

func createMissingCategories(categories []string) {
	logrus.Warnln("Creating missing albums (NotImplemented)")
	for _, c := range categories  {
		logrus.Debug(c)
	}
}

func findMissingImages() {
	logrus.Warnln("Finding missing images (NotImplemented)")
}

func uploadImages() {
	logrus.Warnln("Uploading missing images (NotImplemented)")
}

func configureContext() (*AppContext, error) {
	logrus.Infoln("Preparing application context and configuration")

	if *piwigoUrl == "" {
		return nil, errors.New("missing piwigo url!")
	}

	if *piwigoUser == "" {
		return nil, errors.New("missing piwigo user!")
	}

	if *piwigoPassword == "" {
		return nil, errors.New("missing piwigo password!")
	}

	context := new(AppContext)
	context.LocalRootPath = *imagesRootPath
	context.Piwigo = new(piwigo.PiwigoContext)
	context.Piwigo.Url = fmt.Sprintf("%s/ws.php?format=json", *piwigoUrl)
	context.Piwigo.Username = *piwigoUser
	context.Piwigo.Password = *piwigoPassword

	return context, nil
}

func loginToPiwigoAndConfigureContext(context *AppContext) error {
	logrus.Infoln("Logging in to piwigo and getting chunk size configuration for uploads")
	err := authentication.Login(context.Piwigo)
	if err != nil {
		return err
	}
	return initializeUploadChunkSize(context)
}

func initializeUploadChunkSize(context *AppContext) error {
	userStatus, err := authentication.GetStatus(context.Piwigo)
	if err != nil {
		return err
	}
	context.ChunkSizeBytes = userStatus.Result.UploadFormChunkSize * 1024
	logrus.Debugln(context.ChunkSizeBytes)
	return nil
}
