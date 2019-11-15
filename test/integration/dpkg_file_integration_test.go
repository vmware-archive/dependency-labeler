package integration_test

import (
	"io/ioutil"

	"github.com/pivotal/deplab/test/test_utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("deplab", func() {
	Describe("when called with --dpkg-file", func() {
		Describe("and dpkg can be written", func() {
			It("succeeds", func() {
				dpkgDestinationPath := test_utils.ExistingFileName()
				defer test_utils.CleanupFile(dpkgDestinationPath)

				inputImage := "pivotalnavcon/test-asset-additional-sources"
				_ = runDeplabAgainstImage(inputImage,
					"--dpkg-file", dpkgDestinationPath)

				dpkgFileBytes, err := ioutil.ReadFile(dpkgDestinationPath)

				Expect(err).ToNot(HaveOccurred())
				Expect(string(dpkgFileBytes)).To(
					SatisfyAll(
						ContainSubstring("deplab SHASUM"),
						ContainSubstring("deplab version: 0.0.0-dev"),
						ContainSubstring("Desired=Unknown/Install/Remove/Purge/Hold"),
						ContainSubstring("ii  zlib1g              1:1.2.11.dfsg-0ubuntu2   amd64"),
					))
			})
		})

		Describe("and metadata can't be written", func() {
			It("writes the image metadata, returns the sha and throws an error about the file missing", func() {
				inputImage := "pivotalnavcon/test-asset-additional-sources"
				_, stdErr := runDepLab([]string{"--image", inputImage, "--git", pathToGitRepo, "--dpkg-file", "a-path-that-does-not-exist/foo.dpkg"}, 1)
				Expect(string(getContentsOfReader(stdErr))).To(
					ContainSubstring("a-path-that-does-not-exist/foo.dpkg"))
			})
		})
	})
})
