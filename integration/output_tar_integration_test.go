package integration_test

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/pivotal/deplab/metadata"

	"github.com/docker/docker/api/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("deplab", func() {
	var (
		inputImage             string
		outputImage            string
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

			DescribeTable("without a tag", func(tarDestinationPath string) {
				defer cleanupFile(tarDestinationPath)

				inputImage = "pivotalnavcon/ubuntu-additional-sources"
				outputImage, _, _, _ = runDeplabAgainstImage(inputImage, "--output-tar", tarDestinationPath)

				md := getMetadataFromImageTarball(tarDestinationPath)

				Expect(md.Base).To(
					SatisfyAll(
						HaveKeyWithValue("name", "Ubuntu"),
						HaveKeyWithValue("version_codename", "bionic"),
					))
			},
				Entry("when the file exists", existingFileName()),
				Entry("when the file does not exists", nonExistingFileName()),
			)

			Context("when there is a tag", func() {
				BeforeEach(func() {
					tempDir, err := ioutil.TempDir("", "deplab-integration-output-tar-file-")
					Expect(err).ToNot(HaveOccurred())
					tarDestinationPath = path.Join(tempDir, "image.tar")

					Expect(err).ToNot(HaveOccurred())
					inputImage = "pivotalnavcon/ubuntu-additional-sources"
					outputImage, _, _, _ = runDeplabAgainstImage(inputImage, "--output-tar", tarDestinationPath, "--tag", "foo:bar")
				})

				It("writes the image as a tar", func() {
					manifest := getManifestFromImageTarball(tarDestinationPath)
					Expect(manifest[0]["RepoTags"].([]interface{})[0].(string)).To(Equal("foo:bar"))
				})
			})

			AfterEach(func() {
				err := os.RemoveAll(outputFilesDestination)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Describe("and file can't be written", func() {
			It("writes the image metadata, returns the sha and throws an error about the file location", func() {
				inputImage = "pivotalnavcon/ubuntu-additional-sources"
				stdOut, stdErr := runDepLab([]string{"--image", inputImage, "--git", pathToGitRepo, "--output-tar", "a-path-that-does-not-exist/image.tar"}, 1)
				outputImage, _, _, _ = parseOutputAndValidate(stdOut)
				Expect(string(getContentsOfReader(stdErr))).To(ContainSubstring("directory \"a-path-that-does-not-exist\" does not exist"))
			})
		})

		AfterEach(func() {
			_, err := dockerCli.ImageRemove(context.TODO(), outputImage, types.ImageRemoveOptions{})
			Expect(err).ToNot(HaveOccurred())
		})
	})
})

func getMetadataFromImageTarball(tarDestinationPath string) metadata.Metadata {
	_, configBuffer := getManifestAndConfigFromImageTarball(tarDestinationPath)

	config := make(map[string]interface{})
	err := json.NewDecoder(configBuffer).Decode(&config)
	Expect(err).ToNot(HaveOccurred())
	mdString := config["config"].(map[string]interface{})["Labels"].(map[string]interface{})["io.pivotal.metadata"].(string)

	md := metadata.Metadata{}

	err = json.Unmarshal([]byte(mdString), &md)
	Expect(err).ToNot(HaveOccurred())

	return md
}

func getManifestFromImageTarball(tarDestinationPath string) []map[string]interface{} {
	manifest, _ := getManifestAndConfigFromImageTarball(tarDestinationPath)
	return manifest
}

func getManifestAndConfigFromImageTarball(tarDestinationPath string) ([]map[string]interface{}, *bytes.Buffer) {
	tarDestinationFile, err := os.Open(tarDestinationPath)
	Expect(err).ToNot(HaveOccurred())
	defer tarDestinationFile.Close()

	tr := tar.NewReader(tarDestinationFile)

	manifestBuffer := bytes.Buffer{}
	configBuffer := bytes.Buffer{}

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
				_, err := io.Copy(&manifestBuffer, tr)
				Expect(err).ToNot(HaveOccurred())
			} else {
				_, err := io.Copy(&configBuffer, tr)
				Expect(err).ToNot(HaveOccurred())
			}
		}
	}

	manifest := make([]map[string]interface{}, 0)
	err = json.NewDecoder(&manifestBuffer).Decode(&manifest)
	Expect(err).ToNot(HaveOccurred())

	return manifest, &configBuffer
}
