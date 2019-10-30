package test_utils

import (
	"io/ioutil"
	"os"
	"path"

	. "github.com/onsi/gomega"
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
	err := os.Remove(filePath)
	Expect(err).ToNot(HaveOccurred())
}
