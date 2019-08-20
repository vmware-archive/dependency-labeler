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
			inputImage = "pivotalnavcon/ubuntu-additional-sources"
			outputImage, metadataLabelString, metadataLabel, _ = runDeplabAgainstImage(inputImage)
		})

		It("applies a metadata label", func() {
			Expect(metadataLabelString).ToNot(BeEmpty())

			dependencyMetadata := metadataLabel.Dependencies[0].Source.Metadata
			dpkgMetadata := dependencyMetadata.(map[string]interface{})

			By("listing the dpkg sources")
			sources, ok := dpkgMetadata["apt_sources"].([]interface{})
			Expect(ok).To(BeTrue())
			Expect(len(sources)).To(BeNumerically(">", 0))
			Expect(sources).To(ConsistOf(
				"deb http://archive.ubuntu.com/ubuntu/ bionic main restricted",
				"deb http://archive.ubuntu.com/ubuntu/ bionic universe",
				"deb http://archive.ubuntu.com/ubuntu/ bionic-updates main restricted",
				"deb http://archive.ubuntu.com/ubuntu/ bionic-updates universe",
				"deb http://archive.ubuntu.com/ubuntu/ bionic multiverse",
				"deb http://archive.ubuntu.com/ubuntu/ bionic-updates multiverse",
				"deb http://archive.ubuntu.com/ubuntu/ bionic-backports main restricted universe multiverse",
				"deb http://security.ubuntu.com/ubuntu/ bionic-security main restricted",
				"deb http://security.ubuntu.com/ubuntu/ bionic-security universe",
				"deb http://security.ubuntu.com/ubuntu/ bionic-security multiverse",
				"deb http://example.com/ubuntu getdeb example",
			))

			By("listing debian package dependencies in the image")
			Expect(metadataLabel.Dependencies[0].Type).To(Equal("debian_package_list"))

			Expect(len(dpkgMetadata["packages"].([]interface{}))).To(Equal(89))

			By("generating an image with the input as the parent")
			inspectOutput, _, err := dockerCli.ImageInspectWithRaw(context.TODO(), outputImage)
			Expect(err).ToNot(HaveOccurred())

			By("generating a sha256 digest of the metadata content as version")
			Expect(metadataLabel.Dependencies[0].Source.Version["sha256"]).To(MatchRegexp(`^[0-9a-f]{64}$`))

			inspectInput, _, err := dockerCli.ImageInspectWithRaw(context.TODO(), inputImage)
			Expect(err).ToNot(HaveOccurred())

			Expect(inspectOutput.Parent).To(Equal(inspectInput.ID))
		})
	})

	Context("with an image without dpkg", func() {
		BeforeEach(func() {
			inputImage = "alpine:latest"
			outputImage, metadataLabelString, metadataLabel, _ = runDeplabAgainstImage(inputImage)
		})

		It("does not return a dpkg list", func() {
			_, ok := filterDpkgDependency(metadataLabel.Dependencies)
			Expect(ok).To(BeFalse())
		})
	})

	Context("with an image with dpkg, but no apt sources", func() {
		BeforeEach(func() {
			inputImage = "pivotalnavcon/ubuntu-no-sources"
			outputImage, metadataLabelString, metadataLabel, _ = runDeplabAgainstImage(inputImage)
		})

		It("does not return a dpkg list", func() {
			Expect(metadataLabelString).ToNot(BeEmpty())

			dependencyMetadata := metadataLabel.Dependencies[0].Source.Metadata
			dpkgMetadata := dependencyMetadata.(map[string]interface{})

			sources, ok := dpkgMetadata["apt_sources"].([]interface{})

			Expect(ok).To(BeTrue())
			Expect(sources).To(BeEmpty())
		})
	})

	Context("with an image with dpkg, but no grep", func() {
		BeforeEach(func() {
			inputImage = "pivotalnavcon/ubuntu-no-grep"
			outputImage, metadataLabelString, metadataLabel, _ = runDeplabAgainstImage(inputImage)
		})

		It("does not return a dpkg list", func() {
			Expect(metadataLabelString).ToNot(BeEmpty())

			dependencyMetadata := metadataLabel.Dependencies[0].Source.Metadata
			dpkgMetadata := dependencyMetadata.(map[string]interface{})

			sources, ok := dpkgMetadata["apt_sources"].([]interface{})

			Expect(ok).To(BeTrue())
			Expect(sources).To(BeEmpty())
		})
	})
})

func filterDpkgDependency(dependencies []metadata.Dependency) (metadata.Dependency, bool) {
	for _, dependency := range dependencies {
		if dependency.Source.Type == "debian_package_list" {
			return dependency, true
		}
	}
	return metadata.Dependency{}, false //should never be reached
}
