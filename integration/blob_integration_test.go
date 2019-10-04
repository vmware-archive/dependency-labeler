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
			blobDependency metadata.Dependency
			blobSourceMetadata   map[string]interface{}
		)

		JustBeforeEach(func() {
			inputImage := "ubuntu:bionic"
			outputImage, _, metadataLabel, _ = runDeplabAgainstImage(inputImage, additionalArguments...)
			blobDependency = filterBlobDependency(metadataLabel.Dependencies)
			blobSourceMetadata = blobDependency.Source.Metadata.(map[string]interface{})
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
				Expect(blobDependency.Type).ToNot(BeEmpty())

				By("adding the blob url to the blob dependency")
				Expect(blobDependency.Type).To(Equal("package"))
				Expect(blobDependency.Source.Type).To(Equal("blob"))
				Expect(blobSourceMetadata["url"]).To(Equal("url_to_a_blob.com"))
			})
		})

		Context("when I supply multiple blobs as separate arguments", func() {
			BeforeEach(func() {
				additionalArguments = []string{"--blob", "url_to_a_blob.com", "--blob", "second_url_to_blob.com"}
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
	})
})


func filterBlobDependency(dependencies []metadata.Dependency) metadata.Dependency {
	for _, dependency := range dependencies {
		if dependency.Source.Type == "blob" {
			return dependency
		}
	}
	return metadata.Dependency{}
}
