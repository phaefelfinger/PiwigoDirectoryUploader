/*
 * Copyright (C) 2019 Philipp Haefelfinger (http://www.haefelfinger.ch/). All Rights Reserved.
 * This application is licensed under GPLv2. See the LICENSE file in the root directory of the project.
 */

package localFileStructure

import (
	"crypto/md5"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

func CalculateFileCheckSums(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		logrus.Errorf("Could not open file %s", filePath)
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		logrus.Errorf("Could calculate md5 sum of file %s", filePath)
		return "", err
	}

	md5sum := fmt.Sprintf("%x", hash.Sum(nil))

	logrus.Tracef("Calculated md5 sum of %s - %s", filePath, md5sum)

	return md5sum, nil
}
