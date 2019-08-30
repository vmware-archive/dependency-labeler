package integration_test

import (
	"io/ioutil"
	"os"
	"path"

	. "github.com/onsi/gomega"
)

func existingFileName() string {
	dpkgDestination, _ := ioutil.TempFile("", "deplab-test-file")
	return dpkgDestination.Name()
}

func nonExistingFileName() string {
	tempDir, _ := ioutil.TempDir("", "deplab-test-dir")
	return path.Join(tempDir, "deplab-test-file")
}

func cleanupFile(filePath string) {
	err := os.Remove(filePath)
	Expect(err).ToNot(HaveOccurred())
}
