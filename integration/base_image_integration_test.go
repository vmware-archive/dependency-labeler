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
			inputImage = "ubuntu:bionic-20190718"
			outputImage, _, metadataLabel, _ = runDeplabAgainstImage(inputImage)
		})

		It("adds the base image metadata to the label", func() {
			Expect(metadataLabel.Base).To(
				SatisfyAll(
					HaveKeyWithValue("name", "Ubuntu"),
					HaveKeyWithValue("version", "18.04.2 LTS (Bionic Beaver)"),
					HaveKeyWithValue("version_id", "18.04"),
					HaveKeyWithValue("id_like", "debian"),
					HaveKeyWithValue("version_codename", "bionic"),
					HaveKeyWithValue("pretty_name", "Ubuntu 18.04.2 LTS"),
				))
		})
	})

	Context("with a non-ubuntu:bionic image with /etc/os-release", func() {
		BeforeEach(func() {
			inputImage = "alpine:3.10.1"
			outputImage, _, metadataLabel, _ = runDeplabAgainstImage(inputImage)
		})

		It("adds the base image metadata to the label", func() {
			Expect(metadataLabel.Base).To(
				SatisfyAll(
					HaveKeyWithValue("name", "Alpine Linux"),
					HaveKeyWithValue("version_id", "3.10.1"),
				))
		})
	})

	Context("with an image that doesn't have an os-release", func() {
		BeforeEach(func() {
			inputImage = "pivotalnavcon/ubuntu-no-os-release"
			outputImage, _, metadataLabel, _ = runDeplabAgainstImage(inputImage)
		})

		It("set all fields of base to unknown", func() {
			Expect(metadataLabel.Base).To(
				SatisfyAll(
					HaveKeyWithValue("name", "unknown"),
					HaveKeyWithValue("version_codename", "unknown"),
					HaveKeyWithValue("version_id", "unknown"),
				))
		})
	})

	XContext("with an image that doesn't have cat but has an os-release", func() {
		BeforeEach(func() {
			inputImage = "cloudfoundry/run:tiny"
			outputImage, _, metadataLabel, _ = runDeplabAgainstImage(inputImage)
		})

		It("adds the base image metadata to the label", func() {
			Expect(metadataLabel.Base).To(
				SatisfyAll(
					HaveKeyWithValue("name", "Pivotal Tiny"),
					HaveKeyWithValue("version_codename", "dev"),
					HaveKeyWithValue("version_id", "dev"),
					HaveKeyWithValue("pretty_name", "Pivotal Tiny"),
				))
		})
	})
})
