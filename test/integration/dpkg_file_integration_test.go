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

				_ = runDeplabAgainstTar(
					getTestAssetPath("image-archives/tiny.tgz"),
					"--dpkg-file",
					dpkgDestinationPath)

				dpkgFileBytes, err := ioutil.ReadFile(dpkgDestinationPath)

				Expect(err).ToNot(HaveOccurred())
				Expect(string(dpkgFileBytes)).To(
					SatisfyAll(
						ContainSubstring("deplab SHASUM"),
						ContainSubstring("deplab version: 0.0.0-dev"),
						ContainSubstring("Desired=Unknown/Install/Remove/Purge/Hold"),
						ContainSubstring("ii  tzdata          2019c-0ubuntu0.18.04     all"),
					))
			})
		})

		Describe("and metadata can't be written", func() {
			It("writes the image metadata, returns the sha and throws an error about the file missing", func() {
				_, stdErr := runDepLab([]string{
					"--image-tar", getTestAssetPath("image-archives/tiny.tgz"),
					"--git", pathToGitRepo,
					"--dpkg-file", "a-path-that-does-not-exist/foo.dpkg",
				}, 1)

				Expect(string(getContentsOfReader(stdErr))).To(
					ContainSubstring("a-path-that-does-not-exist/foo.dpkg"))
			})
		})
	})
})
