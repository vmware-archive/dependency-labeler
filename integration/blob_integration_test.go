package integration_test

import (
 "context"
	"github.com/docker/docker/api/types"
	. "github.com/onsi/ginkgo"
. "github.com/onsi/gomega"
	"github.com/pivotal/deplab/metadata"
)

var _ = Describe("deplab blob", func(){

	Context("when I supply blob url(s) as argument(s)", func() {
		var (
			metadataLabel       metadata.Metadata
			additionalArguments []string
			outputImage         string
		)

		JustBeforeEach(func() {
			inputImage := "ubuntu:bionic"
			outputImage, _, metadataLabel, _ = runDeplabAgainstImage(inputImage, additionalArguments...)
		})

		AfterEach(func() {
			_, err := dockerCli.ImageRemove(context.TODO(), outputImage, types.ImageRemoveOptions{})
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when I supply only one --blob argument", func() {
			BeforeEach(func() {
				additionalArguments = []string{"--blob", "url_to_a_blob.com"}
			})

			It("adds a blob dependency", func() {
				blobDependencies := selectBlobDependencies(metadataLabel.Dependencies)
				Expect(len(blobDependencies)).To(Equal(1))
				blobDependency := blobDependencies[0]

				By("adding the blob url to the blob dependency")
				Expect(blobDependency.Type).ToNot(BeEmpty())
				Expect(blobDependency.Type).To(Equal("package"))
				Expect(blobDependency.Source.Type).To(Equal("blob"))

				blobSourceMetadata := blobDependency.Source.Metadata.(map[string]interface{})
				Expect(blobSourceMetadata["url"]).To(Equal("url_to_a_blob.com"))
			})
		})

		Context("when I supply multiple blobs as separate arguments", func() {
			BeforeEach(func() {
				additionalArguments = []string{"--blob", "url_to_a_blob.com", "--blob", "second_url_to_blob.com"}
			})

			It("adds multiple blobDependency entries", func() {
				blobDependencies := selectBlobDependencies(metadataLabel.Dependencies)
				Expect(len(blobDependencies)).To(Equal(2))
			})
		})
	})
})


func selectBlobDependencies(dependencies []metadata.Dependency) []metadata.Dependency {
	var blobs []metadata.Dependency
	for _, dependency := range dependencies {
		if dependency.Source.Type == "blob" {
			blobs = append(blobs, dependency)
		}
	}
	return blobs
}
