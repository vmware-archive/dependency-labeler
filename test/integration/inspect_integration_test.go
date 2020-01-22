package integration_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	. "github.com/onsi/gomega/gstruct"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/pivotal/deplab/pkg/metadata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("deplab inspect", func() {
	It("exits with an error if neither image or image-tar flags are set", func() {
		_, stdErr := runDepLab([]string{"inspect"}, 1)
		errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
		Expect(errorOutput).To(ContainSubstring("ERROR: requires one of --image or --image-tar"))
	})
	It("exits with an error if both image and image-tar flags are set", func() {
		_, stdErr := runDepLab([]string{"inspect",
			"--image", "foo",
			"--image-tar", "path/to/image.tar",
		}, 1)
		errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
		Expect(errorOutput).To(ContainSubstring("ERROR: cannot accept both --image and --image-tar"))
	})

	It("throws an error if invalid characters are in image name", func() {
		By("executing it")
		inputImage := "£$Invalid_image_name$£"
		_, stdErr := runDepLab([]string{
			"inspect",
			"--image", inputImage,
		}, 1)

		errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
		Expect(errorOutput).To(ContainSubstring("could not parse reference"))
	})

	Context("image is already deplab'd", func() {
		DescribeTable("prints all the metadata", func(flag, path string) {
			stdOut, _ := runDepLab([]string{
				"inspect",
				flag, path,
			}, 0)

			md := metadata.Metadata{}
			err := json.NewDecoder(stdOut).Decode(&md)

			Expect(err).ToNot(HaveOccurred())
			Expect(md.Provenance[0].Name).To(Equal("deplab"))
			gitDependencies := selectGitDependencies(md.Dependencies)
			Expect(gitDependencies).ToNot(BeEmpty())
		},
			Entry("with a deplab'd image tarball", "--image-tar", getTestAssetPath("image-archives/tiny-deplabd.tgz")),
			Entry("[remote-image][private-registry] with a deplab'd image from a registry", "--image", "dev.registry.pivotal.io/navcon/deplab-test-asset:tiny-deplabd"),
		)

		It("merge metadata according to the rules and returns a warning", func() {
			provenance := metadata.Provenance{
				Name:    "not-deplab",
				Version: "0.42.42",
				URL:     "",
			}

			base := metadata.Base{
				"some-key": "some-value",
			}

			gitDependency := metadata.Dependency{
				Type: metadata.PackageType,
				Source: metadata.Source{
					Type: metadata.GitSourceType,
					Version: map[string]interface{}{
						"commit": "git-commit",
					},
				},
			}

			archiveDependency := metadata.Dependency{
				Type: metadata.PackageType,
				Source: metadata.Source{
					Type: metadata.ArchiveType,
					Version: map[string]interface{}{
						"sha256": "somesha",
					},
				},
			}

			dpkgDependency := metadata.Dependency{
				Type: metadata.DebianPackageListSourceType,
				Source: metadata.Source{
					Version: map[string]interface{}{
						"sha256": "some-sha",
					},
				},
			}

			imagePath := CreateTinyImageWithDeplabLabel(metadata.Metadata{
				Base:       base,
				Provenance: []metadata.Provenance{provenance},
				Dependencies: []metadata.Dependency{
					gitDependency,
					archiveDependency,
					dpkgDependency,
				},
			})

			defer os.Remove(imagePath)

			stdOut, stderr := runDepLab([]string{
				"inspect",
				"--image-tar", imagePath,
			}, 0)

			md := metadata.Metadata{}
			err := json.NewDecoder(stdOut).Decode(&md)
			Expect(err).ToNot(HaveOccurred())

			Expect(md.Provenance).To(
				SatisfyAll(
					HaveLen(2),
					ContainElement(provenance),
				))
			Expect(md.Base).To(
				SatisfyAll(
					Not(Equal(base)),
					HaveKeyWithValue("pretty_name", "Cloud Foundry Tiny"),
				))
			Expect(md.Dependencies).To(
				SatisfyAll(
					ContainElement(gitDependency),
					ContainElement(archiveDependency),
					Not(ContainElement(dpkgDependency)),
					ContainElement(MatchFields(IgnoreExtras, Fields{
						"Type": Equal(metadata.DebianPackageListSourceType),
					})),
				))

			By("emitting warning for duplicated items")
			Expect(ioutil.ReadAll(stderr)).To(
				SatisfyAll(
					ContainSubstring("base"),
					ContainSubstring("Metadata elements already present on image"),
				))
		})
	})

	Context("image is not previously deplab'd", func() {
		Context("image does not have an /etc/os-release file", func() {
			It("prints the label", func() {
				stdOut, stdErr := runDepLab([]string{
					"inspect",
					"--image-tar", getTestAssetPath("image-archives/scratch.tgz"),
				}, 0)

				md := metadata.Metadata{}
				err := json.NewDecoder(stdOut).Decode(&md)

				Expect(err).ToNot(HaveOccurred())
				Expect(md.Provenance[0].Name).To(Equal("deplab"))

				errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
				Expect(errorOutput).ToNot(ContainSubstring("base"))
			})
		})
		It("prints the label", func() {
			stdOut, _ := runDepLab([]string{
				"inspect",
				"--image-tar", getTestAssetPath("image-archives/tiny.tgz"),
			}, 0)

			md := metadata.Metadata{}
			err := json.NewDecoder(stdOut).Decode(&md)

			Expect(err).ToNot(HaveOccurred())
			Expect(md.Provenance[0].Name).To(Equal("deplab"))
		})
	})

	DescribeTable("provides an error", func(flag, path, errorMsg string) {
		_, stderr := runDepLab([]string{
			"inspect",
			flag, path,
		}, 1)

		Expect(getContentsOfReader(stderr)).To(
			SatisfyAll(
				ContainSubstring(errorMsg),
				ContainSubstring(path)))
	},
		Entry("with a invalid image tarball", "--image-tar", getTestAssetPath("image-archives/invalid-image-archive.tgz"), "cannot open the provided image"),
		Entry("with a non-existent image from registry", "--image", "pivotalnavcon/does-not-exist", "cannot open the provided image"),
		Entry("with a valid image tar ball with invalid json label", "--image-tar", getTestAssetPath("image-archives/tiny-with-invalid-label.tgz"), "cannot parse the label"),
		Entry("[remote-image][private-registry] with a valid image from a registry with invalid json label", "--image", "dev.registry.pivotal.io/navcon/deplab-test-asset:tiny-with-invalid-label", "cannot parse the label"),
	)
})

func CreateTinyImageWithDeplabLabel(m metadata.Metadata) string {
	i, _ := crane.Load(getTestAssetPath("image-archives/tiny.tgz"))

	config, err := i.ConfigFile()
	Expect(err).ToNot(HaveOccurred())

	md, err := json.Marshal(m)
	Expect(err).ToNot(HaveOccurred())

	if config.Config.Labels == nil {
		config.Config.Labels = map[string]string{}
	}

	config.Config.Labels["io.pivotal.metadata"] = string(md)

	i, err = mutate.Config(i, config.Config)
	Expect(err).ToNot(HaveOccurred())

	imagePath, err := ioutil.TempFile("", "")
	Expect(err).ToNot(HaveOccurred())

	err = crane.Save(i, "latest", imagePath.Name())
	Expect(err).ToNot(HaveOccurred())

	return imagePath.Name()
}
