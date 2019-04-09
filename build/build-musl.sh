#!/bin/bash
CC=musl-gcc go build --ldflags '-linkmode external -extldflags "-static"' ./cmd/PiwigoDirectoryUploader/PiwigoDirectoryUploader.go