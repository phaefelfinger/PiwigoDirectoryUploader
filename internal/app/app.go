package app

import (
	"github.com/sirupsen/logrus"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/localFileStructure"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/piwigo/authentication"
)

func Run(rootPath string) {
	context := configureContext(rootPath)

	loginToPiwigoAndConfigureContext(context)

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

func configureContext(rootPath string) *AppContext {
	logrus.Infoln("Preparing application context and configuration")

	context := new(AppContext)
	context.LocalRootPath = rootPath
	context.Piwigo = new(authentication.PiwigoContext)

	//TODO: Move this values to configuration files
	//No, these are not real credentials :-P
	context.Piwigo.Url = "http://pictures.haefelfinger.net/ws.php?format=json"
	context.Piwigo.Username = "admin"
	context.Piwigo.Password = "asdf"

	return context
}

func loginToPiwigoAndConfigureContext(context *AppContext) {
	logrus.Infoln("Logging in to piwigo and getting chunk size configuration for uploads")
	err := authentication.Login(context.Piwigo)
	if err != nil {
		panic(err)
	}
	initializeUploadChunkSize(context)
}

func initializeUploadChunkSize(context *AppContext) {
	userStatus, err := authentication.GetStatus(context.Piwigo)
	if err != nil {
		panic(err)
	}
	context.ChunkSizeBytes = userStatus.Result.UploadFormChunkSize * 1024
	logrus.Debugln(context.ChunkSizeBytes)
}
