package localFileStructure

import (
	"crypto/md5"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

func calculateFileCheckSums(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		logrus.Errorf("Could not open file %s", filePath)
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		logrus.Errorf("Could calculate md5 sum of file %s", filePath)
		return "", err
	}

	md5sum := fmt.Sprintf("%x", md5.Sum(nil))

	logrus.Tracef("Calculated md5 sum of %s - %s", filePath, md5sum)

	return md5sum, nil
}
