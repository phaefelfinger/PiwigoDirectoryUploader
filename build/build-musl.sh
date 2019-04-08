#!/bin/bash
CC=/usr/bin/musl-gcc go build --ldflags '-linkmode external -extldflags "-static"' ./cmd/PiwigoDirectoryUploader/PiwigoDirectoryUploader.go