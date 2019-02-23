package app

import (
	"github.com/sirupsen/logrus"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/localFileStructure"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/internal/pkg/piwigo/authentication"
)

func Run(rootPath string) {
	context := ConfigureContext(rootPath)

	loginToPiwigoAndConfigureContext(context)

	//ScanLocalDirectories(context)
	//GetAllCategoriesFromServer()

	//FindMissingAlbums()
	//CreateMissingAlbums()
	//FindMissingImages()
	//UploadImages()

	authentication.Logout(context.Piwigo)
}

func ConfigureContext(rootPath string) *AppContext {
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

func ScanLocalDirectories(context *AppContext) {
	var fileNodes map[string]localFileStructure.FilesystemNode = localFileStructure.ScanLocalFileStructure(context.LocalRootPath)
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

func loginToPiwigoAndConfigureContext(context *AppContext) {
	logrus.Infoln("Logging in to piwigo and getting chunk size configuration for uploads")
	authentication.Login(context.Piwigo)
	initializeUploadChunkSize(context)
}

func initializeUploadChunkSize(context *AppContext) {
	userStatus := authentication.GetStatus(context.Piwigo)
	context.ChunkSizeBytes = userStatus.Result.UploadFormChunkSize * 1024
	logrus.Debugln(context.ChunkSizeBytes)
}
