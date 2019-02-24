package app

import (
	"errors"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/localFileStructure"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/piwigo/authentication"
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
	//ScanLocalDirectories(context)
	//GetAllCategoriesFromServer()

	//FindMissingAlbums()
	//CreateMissingAlbums()
	//FindMissingImages()
	//UploadImages()

	_ = authentication.Logout(context.Piwigo)
}

func ScanLocalDirectories(context *AppContext) {
	fileNodes, err := localFileStructure.ScanLocalFileStructure(context.LocalRootPath)
	if err != nil {
		panic(err)
	}
	for _, node := range fileNodes {
		logrus.Debugln("found path entry:", node.Key)
	}
}

func GetAllCategoriesFromServer() {
	// get all categories from server and flatten structure to match directory names
	// 2018/2018 album blah
	logrus.Warnln("Loading all categories from the server (NotImplemented)")
}

func FindMissingAlbums() {
	logrus.Warnln("Looking up missing albums (NotImplemented)")
}

func CreateMissingAlbums() {
	logrus.Warnln("Creating missing albums (NotImplemented)")
}

func FindMissingImages() {
	logrus.Warnln("Finding missing images (NotImplemented)")
}

func UploadImages() {
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
	context.Piwigo = new(authentication.PiwigoContext)
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
