/*
 * Copyright (C) 2019 Philipp Haefelfinger (http://www.haefelfinger.ch/). All Rights Reserved.
 * This application is licensed under GPLv2. See the LICENSE file in the root directory of the project.
 */

package app

import (
	"flag"
	"github.com/vharitonsky/iniflags"
	"strings"
)

var (
	logLevel        = flag.String("logLevel", "info", "The minimum log level required to write out a log message. (panic,fatal,error,warn,info,debug,trace)")
	imagesRootPath  = flag.String("imagesRootPath", "", "This is the images root path that should be mirrored to piwigo.")
	sqliteDb        = flag.String("sqliteDb", "./localstate.db", "The connection string to the sql lite database file.")
	noUpload        = flag.Bool("noUpload", false, "If set to true, the metadata gets prepared but the upload is not called and the application is exited with code 90")
	piwigoUrl       = flag.String("piwigoUrl", "", "The root url without tailing slash to your piwigo installation.")
	piwigoUser      = flag.String("piwigoUser", "", "The username to use during sync.")
	piwigoPassword  = flag.String("piwigoPassword", "", "This is password to the given username.")
	removeImages    = flag.Bool("removeImages", false, "If set to true, images scheduled to delete will be removed from the piwigo server. Be sure you want to delete images before enabling this flag.")
	parallelUploads = flag.Int("parallelUploads", 4, "Set the number of images that get uploaded in parallel.")
	dirSuffixToSkip = flag.Int("dirSuffixToSkip", 0, "Set the number of directories at the end of the filepath to remove to build the category (e.g. value of 1: /foo/png/img.png results in foo/img.png).")
	extensions      arrayFlags
	ignoreDirs      arrayFlags
)

type arrayFlags []string

func (arr *arrayFlags) String() string {
	b := strings.Builder{}
	for _, v := range *arr {
		if b.Len() > 0 {
			b.WriteString(",")
		}
		b.WriteString(v)
	}
	return b.String()
}

func (arr *arrayFlags) Set(value string) error {
	*arr = append(*arr, strings.TrimSpace(value))
	return nil
}

func initializeFlags() {
	flag.Var(&extensions, "extension", "Supported file extensions. Flag can be specified multiple times. Uses jpg and png if omitted.")
	flag.Var(&ignoreDirs, "ignoreDir", "Directories that should be ignored. Flag can be specified multiple times for more than one directory.")
	iniflags.Parse()
}
