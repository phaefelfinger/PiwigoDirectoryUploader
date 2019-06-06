#!/bin/bash
CC=musl-gcc go build -o dist/PiwigoDirectoryUploaderMusl --ldflags '-linkmode external -extldflags "-static"' ./cmd/PiwigoDirectoryUploader/PiwigoDirectoryUploader.go