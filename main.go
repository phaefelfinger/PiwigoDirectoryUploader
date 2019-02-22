package main

import (
	"flag"
	"haefelfinger.net/piwigo/DirectoriesToAlbums/pkg/localFileStructure"
	"log"
	"os"
)

func main() {
	flag.Parse()
	root := flag.Arg(0)

	InitializeLog()

	AuthenticateToPiwigo()
	ScanLocalDirectories(root)
	GetAllCategoriesFromServer()
	FindMissingAlbums()
	CreateMissingAlbums()
	FindMissingImages()
	UploadImages()
}

func InitializeLog() {
	//TODO: make log configurable to file instead of console
	log.SetOutput(os.Stdout)
	log.Println("Starting Piwigo directories to albums...")
}

func AuthenticateToPiwigo() {
	log.Println("Authenticating to piwigo server (NotImplemented)")
}

func ScanLocalDirectories(root string) {
	fileNodes := localFileStructure.ScanLocalFileStructure(root)
	log.Printf("filepath.Walk() returned %v\n", fileNodes)
}

func GetAllCategoriesFromServer()  {
	// get all categories from server and flatten structure to match directory names
	// 2018/2018 album blah
	log.Println("Loading all categories from the server (NotImplemented)")
}

func FindMissingAlbums()  {
	log.Println("Looking up missing albums (NotImplemented)")
}

func CreateMissingAlbums()  {
	log.Println("Creating missing albums (NotImplemented)")
}

func FindMissingImages()  {
	log.Println("Finding missing images (NotImplemented)")
}

func UploadImages()  {
	log.Println("Uploading missing images (NotImplemented)")
}