package integration_test

import (
	"archive/tar"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/pivotal/deplab/pkg/metadata"

	"github.com/google/go-containerregistry/pkg/crane"

	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/pivotal/deplab/test/test_utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("deplab", func() {
	var (
		tarDestinationPath     string
		outputFilesDestination string
	)

	Context("when called with --output-tar", func() {
		Describe("and tar can be written", func() {
			BeforeEach(func() {
				var err error
				outputFilesDestination, err = ioutil.TempDir("", "output-files-")
				Expect(err).ToNot(HaveOccurred())
			})

			DescribeTable("without a tag", func(inputImageTar, tarDestinationPath string) {
				defer test_utils.CleanupFile(tarDestinationPath)
				metadataFile, err := ioutil.TempFile("", "")
				Expect(err).ToNot(HaveOccurred())
				defer test_utils.CleanupFile(metadataFile.Name())

				_, _ = runDepLab([]string{
					"--image-tar", inputImageTar,
					"--git", pathToGitRepo,
					"--metadata-file", metadataFile.Name(),
					"--output-tar", tarDestinationPath,
				}, 0)

				metadataFileContent := metadata.Metadata{}
				err = json.NewDecoder(metadataFile).Decode(&metadataFileContent)
				Expect(err).ToNot(HaveOccurred())

				md := getMetadataFromImageTarball(tarDestinationPath)

				Expect(metadataFileContent).To(Equal(md))
			},
				Entry("non existing file", getTestAssetPath("image-archives/tiny.tgz"), test_utils.NonExistingFileName()),
				Entry("existing file", getTestAssetPath("image-archives/tiny.tgz"), test_utils.ExistingFileName()),
			)

			Context("when there is a tag", func() {
				It("writes the image as a tar", func() {
					tempDir, err := ioutil.TempDir(outputFilesDestination, "deplab-integration-output-tar-file-")
					Expect(err).ToNot(HaveOccurred())
					tarDestinationPath = path.Join(tempDir, "image.tar")

					Expect(err).ToNot(HaveOccurred())
					_ = runDeplabAgainstTar(getTestAssetPath("image-archives/tiny.tgz"), "--output-tar", tarDestinationPath, "--tag", "foo:bar")

					manifest := getManifestFromImageTarball(tarDestinationPath)
					Expect(manifest["RepoTags"]).To(ConsistOf("foo:bar"))
				})

				It("exits with an error if the tag passed is not valid", func() {
					_, stdErr := runDepLab([]string{"--image-tar", getTestAssetPath("image-archives/tiny.tgz"),
						"--git", pathToGitRepo,
						"--tag", "foo:testtag/bar",
						"--output-tar", test_utils.ExistingFileName(),
					}, 1)

					errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
					Expect(errorOutput).To(SatisfyAll(
						ContainSubstring("foo:testtag/bar"),
						ContainSubstring("error exporting tar"),
					))
				})
			})

			AfterEach(func() {
				err := os.RemoveAll(outputFilesDestination)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Describe("and file can't be written", func() {
			It("writes the image metadata, returns the sha and throws an error about the file location", func() {
				_, stdErr := runDepLab([]string{
					"--image-tar", getTestAssetPath("image-archives/tiny.tgz"),
					"--git", pathToGitRepo,
					"--output-tar", "a-path-that-does-not-exist/image.tar",
				}, 1)

				Expect(string(getContentsOfReader(stdErr))).To(
					SatisfyAll(
						ContainSubstring("a-path-that-does-not-exist"),
						ContainSubstring("could not export to"),
					))
			})
		})
	})
})

func getMetadataFromImageTarball(tarDestinationPath string) metadata.Metadata {
	image, err := crane.Load(tarDestinationPath)
	Expect(err).To(Not(HaveOccurred()))
	rawConfig, err := image.RawConfigFile()
	Expect(err).ToNot(HaveOccurred())

	config := make(map[string]interface{}, 0)
	err = json.Unmarshal(rawConfig, &config)
	Expect(err).ToNot(HaveOccurred())

	mdString := config["config"].(map[string]interface{})["Labels"].(map[string]interface{})["io.pivotal.metadata"].(string)

	md := metadata.Metadata{}

	err = json.Unmarshal([]byte(mdString), &md)
	Expect(err).ToNot(HaveOccurred())

	return md
}

func getManifestFromImageTarball(tarDestinationPath string) map[string]interface{} {
	tarDestinationFile, err := os.Open(tarDestinationPath)
	Expect(err).ToNot(HaveOccurred())
	defer tarDestinationFile.Close()

	tr := tar.NewReader(tarDestinationFile)
	manifest := make([]map[string]interface{}, 1)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			Expect(err).ToNot(HaveOccurred())
		}

		if strings.Contains(hdr.Name, ".json") {
			if hdr.Name == "manifest.json" {
				err = json.NewDecoder(tr).Decode(&manifest)
				Expect(err).ToNot(HaveOccurred())
				break
			}
		}
	}

	return manifest[0]
}
