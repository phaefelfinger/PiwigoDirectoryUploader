/*
 * Copyright (C) 2019 Philipp Haefelfinger (http://www.haefelfinger.ch/). All Rights Reserved.
 * This application is licensed under GPLv2. See the LICENSE file in the root directory of the project.
 */

package localFileStructure

import "testing"

func Test_ScanLocalFileStructure_should_find_testfile(t *testing.T) {

	images, err := ScanLocalFileStructure("../../../test/")
	if err != nil {
		t.Fatal(err)
	}

	if len(images) != 1 {
		t.Error("Did not find expected testfiles. Expected at least one!")
	}

	for _, img := range images {
		if img.Name != "testimage.jpg" {
			t.Errorf("Did not find the expected testimage.")
		}
	}

}
