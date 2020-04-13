/*
 * Copyright (C) 2020 Philipp Haefelfinger (http://www.haefelfinger.ch/). All Rights Reserved.
 * This application is licensed under GPLv2. See the LICENSE file in the root directory of the project.
 */

package localFileStructure

import "testing"

func Test_ScanLocalFileStructure_should_find_testfile(t *testing.T) {
	supportedExtensions := make([]string, 0)
	supportedExtensions = append(supportedExtensions, "jpg")

	images, err := ScanLocalFileStructure("../../../test/", supportedExtensions, make([]string, 0))
	if err != nil {
		t.Fatal(err)
	}

	if len(images) != 2 { // 1x folder, 1x image
		t.Error("Did not find expected testfiles. Expected at least one!")
	}

	containsTestImage := false
	for _, img := range images {
		if img.Name == "testimage.jpg" {
			containsTestImage = true
		}
	}

	if !containsTestImage {
		t.Errorf("Did not find the expected testimage.")
	}
}

func Test_ScanLocalFileStructure_should_ignore_test_directory(t *testing.T) {
	supportedExtensions := make([]string, 0)
	supportedExtensions = append(supportedExtensions, "jpg")

	ignores := make([]string, 0)
	ignores = append(ignores, "images")
	images, err := ScanLocalFileStructure("../../../test/", supportedExtensions, ignores)
	if err != nil {
		t.Fatal(err)
	}

	if len(images) != 0 {
		t.Error("Did find expected testfiles. Expected no files as test folder is excluded!")
	}
}

func Test_ScanLocalFileStructure_should_not_find_jpg_when_only_png_supported(t *testing.T) {
	supportedExtensions := make([]string, 0)
	supportedExtensions = append(supportedExtensions, "png")

	images, err := ScanLocalFileStructure("../../../test/", supportedExtensions, make([]string, 0))
	if err != nil {
		t.Fatal(err)
	}

	if len(images) != 1 {
		t.Error("Did find expected testfiles. Expected no files as extension is not supported!")
	}

	containsTestImage := false
	for _, img := range images {
		if img.Name == "testimage.jpg" {
			containsTestImage = true
		}
	}

	if containsTestImage {
		t.Errorf("Did find the testimage. This should not happen as png is searched but jpg found")
	}
}
