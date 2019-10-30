package outputs_test

import (
	"io/ioutil"
	"os"

	"github.com/pivotal/deplab/test_utils"

	"github.com/pivotal/deplab/providers"

	"github.com/pivotal/deplab/outputs"

	"github.com/pivotal/deplab/metadata"

	. "github.com/onsi/ginkgo/extensions/table"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var md = metadata.Metadata{
	Dependencies: []metadata.Dependency{
		{
			Type: "debian_package_list",
			Source: metadata.Source{
				Version: map[string]interface{}{
					"sha256": "some-sha",
				},
				Metadata: metadata.DebianPackageListSourceMetadata{
					Packages: []metadata.Package{
						{
							Package:      "foobar",
							Version:      "0.42.0-version",
							Architecture: "amd46",
							Source: metadata.PackageSource{
								Package:         "foobar",
								Version:         "0.42.0-source",
								UpstreamVersion: "0.42.0-upstream",
							},
						},
					},
					AptSources: nil,
				},
			},
		},
	},
}

var _ = Describe("outputs", func() {
	Describe("WriteMetadataFile", func() {
		DescribeTable("when the file can be written", func(path string) {
			defer test_utils.CleanupFile(path)

			err := outputs.WriteMetadataFile(md, path)
			Expect(err).ToNot(HaveOccurred())

			_, err = ioutil.ReadFile(path)

			Expect(err).ToNot(HaveOccurred())
		},
			Entry("the file exists", test_utils.ExistingFileName()),
			Entry("the file does not exists", test_utils.NonExistingFileName()),
		)

		Describe("and metadata can't be written", func() {
			It("returns an error", func() {
				err := outputs.WriteMetadataFile(md, "a-path-that-does-not-exist/foo.dpkg")

				Expect(err).To(MatchError(ContainSubstring("a-path-that-does-not-exist/foo.dpkg")))
			})
		})

		Describe("when the file already has content", func() {
			It("wipes the original file", func() {
				By("creating a metadata file")
				filePath := test_utils.ExistingFileName()
				err := outputs.WriteMetadataFile(md, filePath)
				Expect(err).ToNot(HaveOccurred())

				originalContent := appendContent(filePath)

				By("running it again against the same file")
				err = outputs.WriteMetadataFile(md, filePath)
				Expect(err).ToNot(HaveOccurred())

				newBytes, err := ioutil.ReadFile(filePath)
				Expect(string(newBytes), err).To(Equal(string(originalContent)))
			})
		})
	})

	Describe("WriteDpkgFile", func() {
		DescribeTable("when the file can be written", func(path string) {
			defer test_utils.CleanupFile(path)

			err := outputs.WriteDpkgFile(md, path, "0.1.0-dev")
			Expect(err).ToNot(HaveOccurred())

			dpkgFileBytes, err := ioutil.ReadFile(path)

			Expect(err).ToNot(HaveOccurred())
			Expect(string(dpkgFileBytes)).To(
				SatisfyAll(
					ContainSubstring("deplab SHASUM: some-sha"),
					ContainSubstring("deplab version: 0.1.0-dev"),
					ContainSubstring("Desired=Unknown/Install/Remove/Purge/Hold"),
					ContainSubstring("ii  foobar 0.42.0-version amd46"),
				))
		},
			Entry("the file exists", test_utils.ExistingFileName()),
			Entry("the file does not exists", test_utils.NonExistingFileName()),
		)

		Describe("when metadata does not have a debian_package_list", func() {
			It("returns an error", func() {
				path := test_utils.ExistingFileName()
				defer test_utils.CleanupFile(path)

				err := outputs.WriteDpkgFile(metadata.Metadata{}, path, "0.1.0-dev")
				Expect(err).To(MatchError(
					ContainSubstring(providers.DebianPackageListSourceType)))
			})
		})

		Describe("and dpkg file can't be written", func() {
			It("returns an error", func() {
				err := outputs.WriteDpkgFile(md, "a-path-that-does-not-exist/foo.dpkg", "0.1.0-dev")

				Expect(err).To(MatchError(ContainSubstring("a-path-that-does-not-exist/foo.dpkg")))
			})
		})

		Describe("when the file already has content", func() {
			It("wipes the original file", func() {
				By("creating a dpkg file")
				filePath := test_utils.ExistingFileName()
				err := outputs.WriteDpkgFile(md, filePath, "0.0.0-dev")
				Expect(err).ToNot(HaveOccurred())

				originalContent := appendContent(filePath)

				By("running it again against the same file")
				err = outputs.WriteDpkgFile(md, filePath, "0.0.0-dev")
				Expect(err).ToNot(HaveOccurred())

				newBytes, err := ioutil.ReadFile(filePath)
				Expect(string(newBytes), err).To(Equal(string(originalContent)))
			})
		})
	})
})

func appendContent(filePath string) string {
	By("appending content")
	originalContent, err := ioutil.ReadFile(filePath)
	Expect(err).ToNot(HaveOccurred())

	By("opening a file in append mode")
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0644)
	Expect(err).ToNot(HaveOccurred())

	By("appending some content")
	_, err = f.WriteString("\n some additional content \n")
	Expect(err).ToNot(HaveOccurred())

	By("checking that both the original content and the appended are there")
	bytes, err := ioutil.ReadFile(filePath)
	originalContentString := string(originalContent)
	Expect(string(bytes), err).To(SatisfyAll(
		ContainSubstring(originalContentString),
		ContainSubstring("\n some additional content \n"),
	))

	return originalContentString
}
