# PiwigoDirectoryUploader

This tools mirrors the directory structure of the given root directory as albums and subalbums in piwigo
and uploads all images to the albums.

## Features

Currently the following features are supported

- Creating directory structure as album hierarchy in Piwigo
- Check if an image needs to be uploaded (only md5sum variant currently supported)
- Upload image and assign it to the album based on the directory structure
- Upload updated images that changed locally
- Local image metadata / category storage using sqlite to make change detection easier
- Rebuild the local metadata database without uploading any pictures. Though, The categories get created!
- Remove images no longer present (configurable)
- Uses all CPU Cores to calculate initial metadata
- Upload 4 files in parallel by default (configurable)
- Source uses go modules

There are some features planned but not ready yet:

- Fully support files within multiple albums
- Specify more than one root path to gather images on the local system
- Setup drone CI / CD and build a docker image

## Dependencies

There are some external dependencies to build the application.

- logrus: This is a little logging library that is quite handy
- iniflags: The iniflags makes handling configuration files and applications parameters quite easy.

## Get the source

To get the latest version, you should ``git clone https://git.haefelfinger.net/piwigo/PiwigoDirectoryUploader.git`` to
your local disk or into your local go source directory by using ``go get git.haefelfinger.net/piwigo/PiwigoDirectoryUploader`.

### GO modules

The repository supports gomodules so the modules should be resolved automatically during build or test run.
But just in case, here are the manual commands to install them. This is needed if you put this repo inside the GOPATH
as this sets the module config to ignore.

To get all dependencies at once:

```
go get ./...
go install github.com/golang/mock/mockgen
```

The installation of mockgen is required to use ``go generate ./...`` to build the mocks.

### Manual dependency installation

You may install the dependencies manually with the following commands
or just use the command under "GO modules" to get all dependencies.

```
go get github.com/sirupsen/logrus
go get github.com/vharitonsky/iniflags
```

To build the mocks there are two go:generate dependencies. The mockgen dependency must be installed to make it work:

```
go get github.com/golang/mock/gomock
go get github.com/golang/mock/mockgen
go install github.com/golang/mock/mockgen
```

To rebuild the mocks you can simply use the following command:

```
go generate ./...
```

## Build

### Dynamically linked using glibc

Build the main executable by using the following command. By default it gets the name PiwigoDirectoryUploader
but can be renamed to your favorite application name.

```
go build cmd/PiwigoDirectoryUploader/PiwigoDirectoryUploader.go
```

### Fully statically linked using musl

To get a fully static linked executable, you can use the build script build-musl.sh under the build folder.
You need to have musl installed and musl-gcc in your environment available to make this work. Under Arch Linux
the package is ``community/musl``.

```
./build/build-musl.sh
```

This static linked executable can be run in an absolute minimalistic linux image and without installing any
dependencies or additional packages.

## Configuration

### Command line

You get the following help information to the command line by using:

```
./PiwigoDirectoryUploader -help
```

The following options are supported to run the application from the command line.

```
Usage of ./PiwigoDirectoryUploader:
  -allowMissingConfig
        Don't terminate the app if the ini file cannot be read.
  -allowUnknownFlags
        Don't terminate the app if ini file contains unknown flags.
  -config string
        Path to ini config for using in go flags. May be relative to the current executable path.
  -configUpdateInterval duration
        Update interval for re-reading config file set via -config flag. Zero disables config file re-reading.
  -dumpflags
        Dumps values for all flags defined in the app into stdout in ini-compatible syntax and terminates the app.
  -imagesRootPath string
        This is the images root path that should be mirrored to piwigo.
  -logLevel string
        The minimum log level required to write out a log message. (panic,fatal,error,warn,info,debug,trace) (default "info")
  -noUpload
        If set to true, the metadata gets prepared but the upload is not called and the application is exited with code 90
  -parallelUploads int
        Set the number of images that get uploaded in parallel. (default 4)
  -piwigoPassword string
        This is password to the given username.
  -piwigoUrl string
        The root url without tailing slash to your piwigo installation.
  -piwigoUser string
        The username to use during sync.
  -removeImages
        If set to true, images scheduled to delete will be removed from the piwigo server. Be sure you want to delete images before enabling this flag.
  -sqliteDb string
        The connection string to the sql lite database file. (default "./localstate.db")
```

### Configuration file

It is also possible to use a configuration file to save the settings to be used with multiple piwigo instances.
To use configuration files, just copy the default one and edit the parameters to your wish.

```
cp ./configs/defaultConfig.ini ./localConfig.ini
nano ./localConfig.ini
```

## Run the uploader

Finally you may run the application using the following example command.

```
./PiwigoDirectoryUploader -config=./localConfig.ini
```

If you mess up the local database for some reason, you may just delete it and let the uploader regenerate the content.
The only thing you might loose in this situation is the track of the files that should be deleted during next sync as
this information is built upon existing records of the local database.