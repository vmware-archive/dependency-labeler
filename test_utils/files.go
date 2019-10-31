package test_utils

import (
	"io/ioutil"
	"os"
	"path"
)

func ExistingFileName() string {
	dpkgDestination, _ := ioutil.TempFile("", "deplab-test-file")
	return dpkgDestination.Name()
}

func NonExistingFileName() string {
	tempDir, _ := ioutil.TempDir("", "deplab-test-dir")
	return path.Join(tempDir, "deplab-test-file")
}

func CleanupFile(filePath string) {
	os.Remove(filePath)
}
