package integration_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/docker/docker/api/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("deplab", func() {
	var (
		inputImage              string
		outputImage             string
		metadataLabelString     string
		metadataDestinationPath string
	)

	Context("when called with --metadata-file", func() {
		Describe("and metadata can be written", func() {
			JustBeforeEach(func() {
				inputImage = "pivotalnavcon/ubuntu-additional-sources"
				outputImage, metadataLabelString, _ = runDeplabAgainstImage(inputImage, "--metadata-file", metadataDestinationPath)
			})

			Context("when the file exists", func() {
				BeforeEach(func() {
					metadataDestination, err := ioutil.TempFile("", "metadata-file.json")
					Expect(err).ToNot(HaveOccurred())

					metadataDestinationPath = metadataDestination.Name()
				})

				It("write the metadata content in json format into the metadata-file value", func() {
					metadataFileBytes, err := ioutil.ReadFile(
						metadataDestinationPath,
					)
					Expect(err).ToNot(HaveOccurred())
					Expect(string(metadataFileBytes)).To(Equal(fmt.Sprintf("%s\n", metadataLabelString)))
				})
			})
			Context("when the file does not exist", func() {
				BeforeEach(func() {
					metadataDestinationPath = "/tmp/metadata-file.json"
				})

				It("write the metadata content in json format into the metadata-file value", func() {
					metadataFileBytes, err := ioutil.ReadFile(
						metadataDestinationPath,
					)
					Expect(err).ToNot(HaveOccurred())
					Expect(string(metadataFileBytes)).To(Equal(fmt.Sprintf("%s\n", metadataLabelString)))
				})
			})

			AfterEach(func() {
				err := os.Remove(metadataDestinationPath)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Describe("and metadata can't be written", func() {
			It("writes the image metadata, return the sha and throws an error about the file missing", func() {
				inputImage = "pivotalnavcon/ubuntu-additional-sources"
				stdOutBuffer, stdErrBuffer := runDepLab([]string{"--image", inputImage, "--git", pathToGitRepo, "--metadata-file", "a-path-that-does-not-exist/foo.json"}, 1)

				outputImage, _, _ = parseOutputAndValidate(stdOutBuffer)

				Expect(stdErrBuffer.String()).To(ContainSubstring("a-path-that-does-not-exist/foo.json"))
			})
		})

		AfterEach(func() {
			_, err := dockerCli.ImageRemove(context.TODO(), outputImage, types.ImageRemoveOptions{})
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
