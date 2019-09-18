package integration_test

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/pivotal/deplab/metadata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("deplab base image", func() {
	var (
		inputImage    string
		outputImage   string
		metadataLabel metadata.Metadata
	)

	AfterEach(func() {
		_, err := dockerCli.ImageRemove(context.TODO(), outputImage, types.ImageRemoveOptions{})
		Expect(err).ToNot(HaveOccurred())
	})

	Context("with an ubuntu:bionic image", func() {
		BeforeEach(func() {
			inputImage = "ubuntu:bionic"
			outputImage, _, metadataLabel, _ = runDeplabAgainstImage(inputImage)
		})

		It("adds the base image metadata to the label", func() {
			Expect(metadataLabel.Base.Name).To(Equal("Ubuntu"))
			Expect(metadataLabel.Base.VersionCodename).To(Equal("bionic"))
			Expect(metadataLabel.Base.VersionID).To(Equal("18.04"))
		})
	})

	Context("with a non-ubuntu:bionic image with /etc/os-release", func() {
		BeforeEach(func() {
			inputImage = "alpine:3.10.1"
			outputImage, _, metadataLabel, _ = runDeplabAgainstImage(inputImage)
		})

		It("adds the base image metadata to the label", func() {
			Expect(metadataLabel.Base.Name).To(Equal("Alpine Linux"))
			Expect(metadataLabel.Base.VersionCodename).To(Equal(""))
			Expect(metadataLabel.Base.VersionID).To(Equal("3.10.1"))
		})
	})

	Context("with an image that doesn't have an os-release", func() {
		BeforeEach(func() {
			inputImage = "pivotalnavcon/ubuntu-no-os-release"
			outputImage, _, metadataLabel, _ = runDeplabAgainstImage(inputImage)
		})

		It("set all fields of base to unknown", func() {
			Expect(metadataLabel.Base.Name).To(Equal("unknown"))
			Expect(metadataLabel.Base.VersionCodename).To(Equal("unknown"))
			Expect(metadataLabel.Base.VersionID).To(Equal("unknown"))
		})
	})

	Context("with an image that doesn't have cat but has an os-release", func() {
		BeforeEach(func() {
			inputImage = "pivotalnavcon/tiny"
			outputImage, _, metadataLabel, _ = runDeplabAgainstImage(inputImage)
		})

		It("adds the base image metadata to the label", func() {
			Expect(metadataLabel.Base.Name).To(Equal("Pivotal Tiny"))
			Expect(metadataLabel.Base.VersionCodename).To(Equal("dev"))
			Expect(metadataLabel.Base.VersionID).To(Equal("dev"))
		})
	})
})
