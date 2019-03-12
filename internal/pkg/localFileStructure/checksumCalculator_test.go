package localFileStructure

import (
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"testing"
)

func TestCalculateFileCheckSumsWithValidFile(t *testing.T) {
	expectedSum := "2e7c66bd6657b1a8659ba05af26a0f7e"

	sum, err := calculateFileCheckSums("../../../test/md5testfile.txt")
	if err != nil {
		t.Error(err)
	}

	if sum != expectedSum {
		t.Errorf("wrong md5 sum provided: expected %s - got %s", expectedSum, sum)
	}
}

func TestCalculateFileCheckSumsWithWrongPath(t *testing.T) {
	hook := test.NewGlobal()
	hook.Reset()

	sum, err := calculateFileCheckSums("unknownPath")
	if err == nil {
		t.Error("there was no error using an invalid and unknown path.")
	}

	if sum != "" {
		t.Error("found a checksum of an invalid file. This should not happen!")
	}

	if hook.LastEntry().Level != logrus.ErrorLevel {
		t.Errorf("the error was not logged")
	}
}
