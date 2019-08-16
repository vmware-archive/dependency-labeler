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
		inputImage          string
		outputImage         string
		metadataLabelString string
		metadataDestination *os.File
	)

	Context("when called with --metadata-file", func() {
		Describe("and the metadata file exists", func() {
			BeforeEach(func() {
				var err error
				metadataDestination, err = ioutil.TempFile("", "metadata-file")
				Expect(err).ToNot(HaveOccurred())

				inputImage = "pivotalnavcon/ubuntu-additional-sources"
				outputImage, metadataLabelString, _ = runDeplabAgainstImage(inputImage, "--metadata-file", metadataDestination.Name())
			})

			AfterEach(func() {
				_, err := dockerCli.ImageRemove(context.TODO(), outputImage, types.ImageRemoveOptions{})
				Expect(err).ToNot(HaveOccurred())

				err = os.Remove(metadataDestination.Name())
				Expect(err).ToNot(HaveOccurred())
			})

			It("write the metadata content in json format into the metadata-file value", func() {
				metadataFileBytes, err := ioutil.ReadFile(
					metadataDestination.Name(),
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(metadataFileBytes)).To(Equal(fmt.Sprintf("%s\n", metadataLabelString)))
			})
		})

		Describe("when the metadata file does not exist", func() {
			It("throws an error", func() {
				inputImage = "pivotalnavcon/ubuntu-additional-sources"
				_, stdErrBuffer := runDepLab([]string{"--image", inputImage, "--git", pathToGitRepo, "--metadata-file", "a-path-that-does-not-exist/foo.json"}, 1)

				Expect(stdErrBuffer.String()).To(ContainSubstring("a-path-that-does-not-exist/foo.json"))
			})
		})
	})
})
