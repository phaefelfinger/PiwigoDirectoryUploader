# PiwigoDirectoryUploader

This tools mirrors the directory structure of the given root directory as albums and subalbums in piwigo
and uploads all images to the albums.

## Features

Currently the following features are supported

- Creating directory structure as album hierarchy in Piwigo
- Check if an image needs to be uploaded (only md5sum version currently supported)
- Upload image and assign it to the album based on the directory structure

Planned:

- Upload updated images that changed locally
- Remove images no longer present (configurable)
- Specify more than one root path to gather images on the local system
- Local metadata storage (sqlite or similar) to make change detection easier 


## Build and run the application

### checkout

To get the latest version, you should check out https://git.haefelfinger.net/piwigo/PiwigoDirectoryUploader.git to
your local go source directory.

### Build

Get all dependencies first.

```
go get ./...
```

Build your main executable by using the following command. By default it gets the name main.go but can be renamed to your
favorite application name.

```
go build cmd/PiwigoDirectoryUploader/main.go
```

### Configure

Next you need to prepare at least one configuration file.
You may create more than one configuration file if you have multiple Piwigo installations.

```
cp ./configs/defaultConfig.ini ./localConfig.ini
nano ./localConfig.ini
```

### Run

Finally you may run the application using the following example command.

```
./main -config=./localConfig.ini
```