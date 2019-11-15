package integration_test

import (
	"github.com/pivotal/deplab/test/test_utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("deplab", func() {
	Context("when called with --metadata-file", func() {
		Describe("and metadata can be written", func() {
			It("succeeds", func() {
				metadataDestinationPath := test_utils.ExistingFileName()
				defer test_utils.CleanupFile(metadataDestinationPath)

				_, _ = runDepLab([]string{
					"--image-tar", getTestAssetPath("tiny.tgz"),
					"--git", pathToGitRepo,
					"--metadata-file", metadataDestinationPath,
				}, 0)
			})
		})

		Describe("and metadata can't be written", func() {
			It("exits with 1 and throws an error about the file missing", func() {
				_, stdErr := runDepLab([]string{
					"--image-tar", getTestAssetPath("tiny.tgz"),
					"--git", pathToGitRepo,
					"--metadata-file", "a-path-that-does-not-exist/foo.json",
				}, 1)

				Expect(string(getContentsOfReader(stdErr))).To(
					ContainSubstring("a-path-that-does-not-exist/foo.json"))
			})
		})
	})
})
