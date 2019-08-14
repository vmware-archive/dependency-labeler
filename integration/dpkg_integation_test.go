package integration_test

import (
	"context"

	"github.com/pivotal/deplab/metadata"

	"github.com/docker/docker/api/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("deplab dpkg", func() {
	var (
		inputImage          string
		outputImage         string
		metadataLabelString string
		metadataLabel       metadata.Metadata
	)

	AfterEach(func() {
		_, err := dockerCli.ImageRemove(context.TODO(), outputImage, types.ImageRemoveOptions{})
		Expect(err).ToNot(HaveOccurred())
	})

	Context("with an ubuntu:bionic image", func() {
		BeforeEach(func() {
			inputImage = "ubuntu:bionic"
			outputImage, metadataLabelString, metadataLabel = runDeplabAgainstImage(inputImage)
		})

		It("applies a metadata label", func() {
			Expect(metadataLabelString).ToNot(BeEmpty())

			By("listing debian package dependencies in the image")
			Expect(metadataLabel.Dependencies[0].Type).To(Equal("debian_package_list"))

			dependencyMetadata := metadataLabel.Dependencies[0].Source.Metadata
			dpkgMetadata := dependencyMetadata.(map[string]interface{})
			Expect(len(dpkgMetadata["packages"].([]interface{}))).To(Equal(89))

			By("generating an image with the input as the parent")
			inspectOutput, _, err := dockerCli.ImageInspectWithRaw(context.TODO(), outputImage)
			Expect(err).ToNot(HaveOccurred())

			inspectInput, _, err := dockerCli.ImageInspectWithRaw(context.TODO(), inputImage)
			Expect(err).ToNot(HaveOccurred())

			Expect(inspectOutput.Parent).To(Equal(inspectInput.ID))
		})
	})

	Context("with an image without dpkg", func() {
		BeforeEach(func() {
			inputImage = "alpine:latest"
			outputImage, metadataLabelString, metadataLabel = runDeplabAgainstImage(inputImage)
		})

		It("does not return a dpkg list", func() {
			Expect(len(metadataLabel.Dependencies)).To(Equal(0))
		})
	})
})
