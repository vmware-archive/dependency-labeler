package integration_test

import (
 "context"
	"github.com/docker/docker/api/types"
	. "github.com/onsi/ginkgo"
. "github.com/onsi/gomega"
	"github.com/pivotal/deplab/metadata"
	"github.com/pivotal/deplab/providers"
)

var _ = Describe("deplab additional-source-url", func(){

	Context("when I supply additional source url(s) as argument(s)", func() {
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

		Context("when I supply only one --additional-source-url argument", func() {
			BeforeEach(func() {
				additionalArguments = []string{"--additional-source-url", "url_to_an_archive.com"}
			})

			It("adds a additional-source-url dependency", func() {
				archiveDependencies := selectArchiveDependencies(metadataLabel.Dependencies)
				Expect(archiveDependencies).To(HaveLen(1))
				archiveDependency := archiveDependencies[0]

				By("adding the additional-source-url to the archive dependency")
				Expect(archiveDependency.Type).ToNot(BeEmpty())
				Expect(archiveDependency.Type).To(Equal("package"))
				Expect(archiveDependency.Source.Type).To(Equal(providers.ArchiveType))

				archiveSourceMetadata := archiveDependency.Source.Metadata.(map[string]interface{})
				Expect(archiveSourceMetadata["url"]).To(Equal("url_to_an_archive.com"))
			})
		})

		Context("when I supply multiple additional-source-url as separate arguments", func() {
			BeforeEach(func() {
				additionalArguments = []string{"--additional-source-url", "url_to_an_archive.com", "--additional-source-url", "second_url_to_an_archive.com"}
			})

			It("adds multiple archive entries", func() {
				archiveDependencies := selectArchiveDependencies(metadataLabel.Dependencies)
				Expect(archiveDependencies).To(HaveLen(2))
			})
		})
	})
})


func selectArchiveDependencies(dependencies []metadata.Dependency) []metadata.Dependency {
	var archives []metadata.Dependency
	for _, dependency := range dependencies {
		if dependency.Source.Type == providers.ArchiveType {
			archives = append(archives, dependency)
		}
	}
	return archives
}
