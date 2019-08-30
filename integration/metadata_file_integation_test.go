package integration_test

import (
	"context"
	"fmt"
	"io/ioutil"

	. "github.com/onsi/ginkgo/extensions/table"

	"github.com/docker/docker/api/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("deplab", func() {
	var (
		outputImage         string
		metadataLabelString string
	)

	Context("when called with --metadata-file", func() {
		DescribeTable("and metadata can be written", func(metadataDestinationPath string) {
			defer cleanupFile(metadataDestinationPath)

			inputImage := "pivotalnavcon/ubuntu-additional-sources"
			outputImage, metadataLabelString, _, _ = runDeplabAgainstImage(inputImage, "--metadata-file", metadataDestinationPath)
			metadataFileBytes, err := ioutil.ReadFile(metadataDestinationPath)

			Expect(err).ToNot(HaveOccurred())
			Expect(string(metadataFileBytes)).To(Equal(fmt.Sprintf("%s\n", metadataLabelString)))
		},
			Entry("when the file exists", existingFileName()),
			Entry("when the file does not exists", nonExistingFileName()),
		)

		Describe("and metadata can't be written", func() {
			It("writes the image metadata, return the sha and throws an error about the file missing", func() {
				inputImage := "pivotalnavcon/ubuntu-additional-sources"
				stdOut, stdErr := runDepLab([]string{"--image", inputImage, "--git", pathToGitRepo, "--metadata-file", "a-path-that-does-not-exist/foo.json"}, 1)

				outputImage, _, _, _ = parseOutputAndValidate(stdOut)

				Expect(string(getContentsOfReader(stdErr))).To(ContainSubstring("a-path-that-does-not-exist/foo.json"))
			})
		})
	})

	AfterEach(func() {
		_, err := dockerCli.ImageRemove(context.TODO(), outputImage, types.ImageRemoveOptions{})
		Expect(err).ToNot(HaveOccurred())
	})
})
