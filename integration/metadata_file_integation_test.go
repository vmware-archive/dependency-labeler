package integration_test

import (
	. "github.com/onsi/ginkgo/extensions/table"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("deplab", func() {
	Context("when called with --metadata-file", func() {
		DescribeTable("and metadata can be written", func(metadataDestinationPath string) {
			defer cleanupFile(metadataDestinationPath)

			inputImage := "pivotalnavcon/ubuntu-additional-sources"

			_, _ = runDepLab([]string{
				"--image", inputImage,
				"--git", pathToGitRepo,
				"--metadata-file", metadataDestinationPath,
			}, 0)
		},
			Entry("when the file exists", existingFileName()),
			Entry("when the file does not exists", nonExistingFileName()),
		)

		Describe("and metadata can't be written", func() {
			It("writes the image metadata, return the sha and throws an error about the file missing", func() {
				inputImage := "pivotalnavcon/ubuntu-additional-sources"
				_, stdErr := runDepLab([]string{
					"--image", inputImage,
					"--git", pathToGitRepo,
					"--metadata-file", "a-path-that-does-not-exist/foo.json",
				}, 1)

				Expect(string(getContentsOfReader(stdErr))).To(
					ContainSubstring("a-path-that-does-not-exist/foo.json"))
			})
		})
	})
})
