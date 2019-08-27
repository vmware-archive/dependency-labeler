package integration_test

import (
	"context"
	"io/ioutil"
	"os"
	"path"

	"github.com/docker/docker/api/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("deplab", func() {
	var (
		inputImage          string
		outputImage         string
		dpkgDestinationPath string
	)

	FContext("when called with --dpkg-file", func() {
		Describe("and dpkg can be written", func() {
			JustBeforeEach(func() {
				inputImage = "pivotalnavcon/ubuntu-additional-sources"
				outputImage, _, _, _ = runDeplabAgainstImage(inputImage, "--dpkg-file", dpkgDestinationPath)
			})

			Context("when the file exists", func() {
				BeforeEach(func() {
					dpkgDestination, err := ioutil.TempFile("", "dpkg-file.dpkg")
					Expect(err).ToNot(HaveOccurred())

					dpkgDestinationPath = dpkgDestination.Name()
				})

				It("overwrites the dpkg-file with the dpkg metadata content in dpkg -l format", func() {
					dpkgFileBytes, err := ioutil.ReadFile(dpkgDestinationPath)
					Expect(err).ToNot(HaveOccurred())
					Expect(string(dpkgFileBytes)).To(ContainSubstring(
						"Desired=Unknown/Install/Remove/Purge/Hold",
					))
					Expect(string(dpkgFileBytes)).To(ContainSubstring(
						"ii  zlib1g              1:1.2.11.dfsg-0ubuntu2   amd64",
					))
				})
			})

			Context("when the file does not exist", func() {
				BeforeEach(func() {
					tempDir, err := ioutil.TempDir("", "deplab-integration-dpkg-file")
					Expect(err).ToNot(HaveOccurred())
					dpkgDestinationPath = path.Join(tempDir, "dpkg-list.dpkg")
				})

				It("writes the dpkg-file with the dpkg metadata content in dpkg -l format", func() {
					dpkgFileBytes, err := ioutil.ReadFile(dpkgDestinationPath)
					Expect(err).ToNot(HaveOccurred())
					Expect(string(dpkgFileBytes)).To(ContainSubstring(
						"Desired=Unknown/Install/Remove/Purge/Hold",
					))
					Expect(string(dpkgFileBytes)).To(ContainSubstring(
						"ii  zlib1g              1:1.2.11.dfsg-0ubuntu2   amd64",
					))
				})
			})

			AfterEach(func() {
				err := os.Remove(dpkgDestinationPath)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Describe("and metadata can't be written", func() {
			It("writes the image metadata, returns the sha and throws an error about the file missing", func() {
				inputImage = "pivotalnavcon/ubuntu-additional-sources"
				stdOut, stdErr := runDepLab([]string{"--image", inputImage, "--git", pathToGitRepo, "--dpkg-file", "a-path-that-does-not-exist/foo.dpkg"}, 1)
				outputImage, _, _, _ = parseOutputAndValidate(stdOut)
				Expect(string(getContentsOfReader(stdErr))).To(ContainSubstring("a-path-that-does-not-exist/foo.dpkg"))
			})
		})

		AfterEach(func() {
			_, err := dockerCli.ImageRemove(context.TODO(), outputImage, types.ImageRemoveOptions{})
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
