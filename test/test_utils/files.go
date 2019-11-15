package test_utils

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
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

func AppendContent(filePath string) string {
	ginkgo.By("appending content")
	originalContent, err := ioutil.ReadFile(filePath)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	ginkgo.By("opening a file in append mode")
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0644)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	ginkgo.By("appending some content")
	_, err = f.WriteString("\n some additional content \n")
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	ginkgo.By("checking that both the original content and the appended are there")
	bytes, err := ioutil.ReadFile(filePath)
	originalContentString := string(originalContent)
	gomega.Expect(string(bytes), err).To(gomega.SatisfyAll(
		gomega.ContainSubstring(originalContentString),
		gomega.ContainSubstring("\n some additional content \n"),
	))

	return originalContentString
}
