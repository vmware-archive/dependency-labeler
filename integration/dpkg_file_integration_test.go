package integration_test

import (
	"context"
	"io/ioutil"

	"github.com/docker/docker/api/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("deplab", func() {
	var (
		outputImage string
	)

	Describe("when called with --dpkg-file", func() {

		DescribeTable("and dpkg can be written", func(dpkgDestinationPath string) {
			defer cleanupFile(dpkgDestinationPath)

			inputImage := "pivotalnavcon/ubuntu-additional-sources"
			outputImage, _, _, _ = runDeplabAgainstImage(inputImage, "--dpkg-file", dpkgDestinationPath)

			dpkgFileBytes, err := ioutil.ReadFile(dpkgDestinationPath)

			Expect(err).ToNot(HaveOccurred())
			Expect(string(dpkgFileBytes)).To(ContainSubstring(
				"deplab SHASUM",
			))
			Expect(string(dpkgFileBytes)).To(ContainSubstring(
				"Desired=Unknown/Install/Remove/Purge/Hold",
			))
			Expect(string(dpkgFileBytes)).To(ContainSubstring(
				"ii  zlib1g              1:1.2.11.dfsg-0ubuntu2   amd64",
			))
		},
			Entry("when the file exists", existingFileName()),
			Entry("when the file does not exists", nonExistingFileName()),
		)

		Describe("and metadata can't be written", func() {
			It("writes the image metadata, returns the sha and throws an error about the file missing", func() {
				inputImage := "pivotalnavcon/ubuntu-additional-sources"
				stdOut, stdErr := runDepLab([]string{"--image", inputImage, "--git", pathToGitRepo, "--dpkg-file", "a-path-that-does-not-exist/foo.dpkg"}, 1)
				outputImage, _, _, _ = parseOutputAndValidate(stdOut)
				Expect(string(getContentsOfReader(stdErr))).To(ContainSubstring("a-path-that-does-not-exist/foo.dpkg"))
			})
		})
	})

	AfterEach(func() {
		_, err := dockerCli.ImageRemove(context.TODO(), outputImage, types.ImageRemoveOptions{})
		Expect(err).ToNot(HaveOccurred())
	})
})
