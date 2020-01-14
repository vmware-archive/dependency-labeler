package dpkg_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	"github.com/onsi/gomega"
	"github.com/pivotal/deplab/pkg/dpkg"
	"github.com/pivotal/deplab/pkg/metadata"
	"github.com/pivotal/deplab/test/test_utils"
)

var _ = Describe("dpkg", func() {
	Describe("WriteDpkgFile", func() {
		table.DescribeTable("when the file can be written", func(path string) {
			defer test_utils.CleanupFile(path)

			err := dpkg.WriteDpkgFile(test_utils.MetadataSample, path, "0.1.0-dev")
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			dpkgFileBytes, err := ioutil.ReadFile(path)

			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(string(dpkgFileBytes)).To(
				gomega.SatisfyAll(
					gomega.ContainSubstring("deplab SHASUM: some-sha"),
					gomega.ContainSubstring("deplab version: 0.1.0-dev"),
					gomega.ContainSubstring("Desired=Unknown/Install/Remove/Purge/Hold"),
					gomega.ContainSubstring("ii  foobar 0.42.0-version amd46"),
				))
		},
			table.Entry("the file exists", test_utils.ExistingFileName()),
			table.Entry("the file does not exists", test_utils.NonExistingFileName()),
		)

		Describe("when metadata does not have a debian_package_list", func() {
			It("returns an error", func() {
				path := test_utils.ExistingFileName()
				defer test_utils.CleanupFile(path)

				err := dpkg.WriteDpkgFile(metadata.Metadata{}, path, "0.1.0-dev")
				gomega.Expect(err).To(gomega.MatchError(
					gomega.ContainSubstring(metadata.DebianPackageListSourceType)))
			})
		})

		Describe("and dpkg file can't be written", func() {
			It("returns an error", func() {
				err := dpkg.WriteDpkgFile(test_utils.MetadataSample, "a-path-that-does-not-exist/foo.dpkg", "0.1.0-dev")

				gomega.Expect(err).To(gomega.MatchError(gomega.ContainSubstring("a-path-that-does-not-exist/foo.dpkg")))
			})
		})

		Describe("when the file already has content", func() {
			It("wipes the original file", func() {
				By("creating a dpkg file")
				filePath := test_utils.ExistingFileName()
				err := dpkg.WriteDpkgFile(test_utils.MetadataSample, filePath, "0.0.0-dev")
				gomega.Expect(err).ToNot(gomega.HaveOccurred())

				originalContent := test_utils.AppendContent(filePath)

				By("running it again against the same file")
				err = dpkg.WriteDpkgFile(test_utils.MetadataSample, filePath, "0.0.0-dev")
				gomega.Expect(err).ToNot(gomega.HaveOccurred())

				newBytes, err := ioutil.ReadFile(filePath)
				gomega.Expect(string(newBytes), err).To(gomega.Equal(string(originalContent)))
			})
		})
	})
})
