# PiwigoDirectoryUploader

This tools mirrors the directory structure of the given root directory as albums and subalbums in piwigo
and uploads all images to the albums.

## Features

Currently the following features are supported

- Creating directory structure as album hierarchy in Piwigo
- Check if an image needs to be uploaded (only md5sum variant currently supported)
- Upload image and assign it to the album based on the directory structure
- Upload updated images that changed locally
- Local metadata storage using sqlite to make change detection easier
- Rebuild the local metadata database without uploading any pictures. Though, The categories get created!

There are some features planned but not ready yet:

- Optimize performance on initial matadata build up.
- Upload more than one file at a time
- Fully support files within multiple albums
- Specify more than one root path to gather images on the local system
- Remove images no longer present (configurable)


## Build and run the application

### checkout

To get the latest version, you should check out https://git.haefelfinger.net/piwigo/PiwigoDirectoryUploader.git to
your local go source directory.

### Build

Get all dependencies first.

```
go get ./...
```

Build the main executable by using the following command. By default it gets the name PiwigoDirectoryUploader.go but
can be renamed to your favorite application name.

```
go build cmd/PiwigoDirectoryUploader/PiwigoDirectoryUploader.go
```

### Configuration

#### Command line

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
  -piwigoPassword string
        This is password to the given username.
  -piwigoUrl string
        The root url without tailing slash to your piwigo installation.
  -piwigoUser string
        The username to use during sync.
  -sqliteDb string
        The connection string to the sql lite database file. (default "./localstate.db")
```

#### Configuration file

It is also possible to use a configuration file to save the settings to be used with multiple piwigo instances.
To use configuration files, just copy the default one and edit the parameters to your wish.

```
cp ./configs/defaultConfig.ini ./localConfig.ini
nano ./localConfig.ini
```

### Run the uploader

Finally you may run the application using the following example command.

```
./PiwigoDirectoryUploader -config=./localConfig.ini
```