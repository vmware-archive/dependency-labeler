package kpack_test

import (
	v1 "github.com/google/go-containerregistry/pkg/v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/deplab/pkg/common"
	. "github.com/pivotal/deplab/pkg/kpack"
	"github.com/pivotal/deplab/pkg/metadata"
	"github.com/pivotal/deplab/test/test_utils"
)

type MockImage struct {
	labels map[string]string
}

func (m MockImage) GetConfig() (*v1.ConfigFile, error) {
	config :=  &v1.ConfigFile{}

	config.Config.Labels = m.labels

	return config, nil
}

func (m MockImage) Cleanup() {
	panic("implement me")
}

func (m MockImage) GetFileContent(string) (string, error) {
	panic("implement me")
}

func (m MockImage) GetDirFileNames(string) ([]string, error) {
	panic("implement me")
}

func (m MockImage) GetDirContents(string) ([]string, error) {
	panic("implement me")
}

func (m MockImage) AbsolutePath(string) (string, error) {
	panic("implement me")
}

func (m MockImage) ExportWithMetadata(metadata.Metadata, string, string) error {
	panic("implement me")
}

var _ = Describe("Kpack", func() {
	Describe("Provider", func() {
		Context("when the image has no kpack label", func() {
			It("does not modify the metadata content", func() {
				Expect(Provider(test_utils.NewMockImageWithEmptyConfig(), common.RunParams{}, metadata.Metadata{})).To(Equal(metadata.Metadata{}))
			})
		})
		Context("when the image has kpack label", func() {
			It("returns properly formatted metadata", func() {
				expectedMd := metadata.KpackRepoSourceMetadata {
					Url: "https://github.com/zmackie/github-actions-automate-projects.git",
					Refs: []interface{}{},
				}
				expectedResult := metadata.Dependency{
					Type:   "package",
					Source: metadata.Source{
						Type: "git",
						Version: map[string]interface{}{
							"commit": "1736b5e3b43a8cf40b3640821ee0e26049e1a58c",
						},
						Metadata: expectedMd,
					},
				}
				labels := map[string]string{}
				labels["io.buildpacks.project.metadata"] = "{\"source\":{\"type\":\"git\",\"version\":{\"commit\":\"1736b5e3b43a8cf40b3640821ee0e26049e1a58c\"},\"metadata\":{\"repository\":\"https://github.com/zmackie/github-actions-automate-projects.git\",\"revision\":\"1736b5e3b43a8cf40b3640821ee0e26049e1a58c\"}}}"

				md, err := Provider(MockImage{labels}, common.RunParams{}, metadata.Metadata{})

				Expect(err).To(Not(HaveOccurred()))
				Expect(md.Dependencies).To(HaveLen(1))
				Expect(md.Dependencies[0]).To(Equal(expectedResult))
			})
		})
		Context("when the image has kpack label with malformed data", func() {
			It("returns generic error message that supported metadata is not present", func() {
				labels := map[string]string{}
				labels["io.buildpacks.project.metadata"] = "broken-source\"\"type\":\"git\",\"version\":{\"commit\":\"1736b5e3b43a8cf40b3640821ee0e26049e1a58c\"},\"metadata\":{\"repository\":\"https://github.com/zmackie/github-actions-automate-projects.git\",\"revision\":\"1736b5e3b43a8cf40b3640821ee0e26049e1a58c\"}}}"

				_, err := Provider(MockImage{labels}, common.RunParams{}, metadata.Metadata{})

				Expect(err).To(MatchError(SatisfyAll(
					ContainSubstring("could not parse kpack metadata"),
					ContainSubstring("could not decode json"))))
				})
		})
	})
})