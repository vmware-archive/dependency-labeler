package integration_test

import (
	"net/http"

	"github.com/pivotal/deplab/pkg/metadata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("deplab additional-source-url", func() {

	Context("when I supply additional source url(s) as argument(s)", func() {
		var (
			metadataLabel       metadata.Metadata
			additionalArguments []string
			server              *ghttp.Server
		)

		JustBeforeEach(func() {
			metadataLabel = runDeplabAgainstTar(getTestAssetPath("image-archives/tiny.tgz"), additionalArguments...)
		})

		Context("when I supply only one --additional-source-url argument", func() {
			var address string
			BeforeEach(func() {
				server = startServer(
					ghttp.RespondWith(http.StatusOK, []byte("")))
				address = server.URL() + "/foo/bar/file.zip"
				additionalArguments = []string{"--additional-source-url", address}
			})

			AfterEach(func() {
				server.Close()
			})

			It("adds a additional-source-url dependency", func() {

				archiveDependencies := selectArchiveDependencies(metadataLabel.Dependencies)
				Expect(archiveDependencies).To(HaveLen(1))
				archiveDependency := archiveDependencies[0]

				By("adding the additional-source-url to the archive dependency")
				Expect(archiveDependency.Type).ToNot(BeEmpty())
				Expect(archiveDependency.Type).To(Equal("package"))
				Expect(archiveDependency.Source.Type).To(Equal(metadata.ArchiveType))

				archiveSourceMetadata := archiveDependency.Source.Metadata.(map[string]interface{})
				Expect(archiveSourceMetadata["url"]).To(Equal(address))
			})
		})

		Context("when I supply multiple additional-source-url as separate arguments", func() {
			BeforeEach(func() {
				server = startServer(
					ghttp.RespondWith(http.StatusOK, []byte("")),
					ghttp.RespondWith(http.StatusOK, []byte("")),
				)

				additionalArguments = []string{
					"--additional-source-url", server.URL() + "/deplab/file.zip",
					"--additional-source-url", server.URL() + "/foobar/file.zip"}
			})

			It("adds multiple archive entries", func() {
				archiveDependencies := selectArchiveDependencies(metadataLabel.Dependencies)
				Expect(archiveDependencies).To(HaveLen(2))
			})

			AfterEach(func() {
				server.Close()
			})
		})
	})
})

func selectArchiveDependencies(dependencies []metadata.Dependency) []metadata.Dependency {
	var archives []metadata.Dependency
	for _, dependency := range dependencies {
		if dependency.Source.Type == metadata.ArchiveType {
			archives = append(archives, dependency)
		}
	}
	return archives
}
