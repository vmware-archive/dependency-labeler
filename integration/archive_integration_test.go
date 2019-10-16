package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/deplab/metadata"
	"github.com/pivotal/deplab/providers"
)

var _ = Describe("deplab additional-source-url", func() {

	Context("when I supply additional source url(s) as argument(s)", func() {
		var (
			metadataLabel       metadata.Metadata
			additionalArguments []string
		)

		JustBeforeEach(func() {
			inputImage := "ubuntu:bionic"
			metadataLabel = runDeplabAgainstImage(inputImage, additionalArguments...)
		})

		Context("when I supply only one --additional-source-url argument", func() {
			BeforeEach(func() {
				additionalArguments = []string{"--additional-source-url", "https://example.com"}
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
				Expect(archiveSourceMetadata["url"]).To(Equal("https://example.com"))
			})
		})

		Context("when I supply multiple additional-source-url as separate arguments", func() {
			BeforeEach(func() {
				additionalArguments = []string{"--additional-source-url", "https://example.com", "--additional-source-url", "https://example.com/foobar"}
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
