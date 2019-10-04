package integration_test

import (
	"context"
	"github.com/docker/docker/api/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/deplab/metadata"
	"path/filepath"
)

var _ = Describe("deplab blob", func(){

	Context("when I supply an artefacts file as an argument", func() {
		var (
			metadataLabel       metadata.Metadata
			additionalArguments []string
			outputImage         string
			blobDependency metadata.Dependency
			blobSourceMetadata   map[string]interface{}
		)

		JustBeforeEach(func() {
			inputImage := "ubuntu:bionic"
			outputImage, _, metadataLabel, _ = runDeplabAgainstImage(inputImage, additionalArguments...)
		})

		AfterEach(func() {
			_, err := dockerCli.ImageRemove(context.TODO(), outputImage, types.ImageRemoveOptions{})
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when I supply an artefacts file with only one blob", func() {
			BeforeEach(func() {
				inputArtefactsPath, err := filepath.Abs(filepath.Join("assets", "artefacts-single-blob.yml"))
				Expect(err).ToNot(HaveOccurred())
				additionalArguments = []string{"--artefacts-file", inputArtefactsPath}
			})

			It("adds a blob dependency", func() {
				blobDependency = filterBlobDependency(metadataLabel.Dependencies)
				Expect(blobDependency).NotTo(BeNil())
				Expect(blobDependency.Source.Metadata).NotTo(BeNil())
				blobSourceMetadata = blobDependency.Source.Metadata.(map[string]interface{})
				Expect(blobDependency.Type).ToNot(BeEmpty())

				By("adding the blob url to the blob dependency")
				Expect(blobDependency.Type).To(Equal("package"))
				Expect(blobDependency.Source.Type).To(Equal("blob"))
				Expect(blobSourceMetadata["url"]).To(Equal("http://archive.ubuntu.com/ubuntu/pool/main/c/ca-certificates/ca-certificates_20180409.tar.xz"))
			})
		})

		Context("when I supply an artefacts file with multiple blobs", func() {
			BeforeEach(func() {
				inputArtefactsPath, err := filepath.Abs(filepath.Join("assets", "artefacts-multiple-blobs.yml"))
				Expect(err).ToNot(HaveOccurred())
				additionalArguments = []string{"--artefacts-file", inputArtefactsPath}
			})

			It("adds multiple blobDependency entries", func() {
				i := 0
				for _, dep := range metadataLabel.Dependencies {
					if dep.Source.Type == "blob" {
						i++
					}
				}

				Expect(i).To(Equal(2))
			})
		})

		Context("when I supply an artefacts file with no blobs", func() {
			BeforeEach(func() {
				inputArtefactsPath, err := filepath.Abs(filepath.Join("assets", "artefacts-empty.yml"))
				Expect(err).ToNot(HaveOccurred())
				additionalArguments = []string{"--artefacts-file", inputArtefactsPath}
			})

			It("adds zero blobDependency entries", func() {
				i := 0
				for _, dep := range metadataLabel.Dependencies {
					if dep.Source.Type == "blob" {
						i++
					}
				}

				Expect(i).To(Equal(0))
			})
		})

		Context("when I supply multiple artefacts files", func() {
			BeforeEach(func() {
				inputArtefactsPath1, err := filepath.Abs(filepath.Join("assets", "artefacts-multiple-blobs.yml"))
				Expect(err).ToNot(HaveOccurred())
				inputArtefactsPath2, err := filepath.Abs(filepath.Join("assets", "artefacts-single-blob.yml"))
				Expect(err).ToNot(HaveOccurred())
				additionalArguments = []string{"--artefacts-file", inputArtefactsPath1, "--artefacts-file", inputArtefactsPath2}
			})

			It("adds multiple blobDependency entries", func() {
				i := 0
				for _, dep := range metadataLabel.Dependencies {
					if dep.Source.Type == "blob" {
						i++
					}
				}

				Expect(i).To(Equal(3))
			})
		})
	})
})
