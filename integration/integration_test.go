package integration_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/pivotal/deplab/metadata"
	"github.com/pivotal/deplab/providers"

	"github.com/onsi/gomega/ghttp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("deplab", func() {
	Context("with an image reference", func() {
		It("throws an error if scratch image is provided", func() {
			By("executing it")
			inputImage := "scratch"
			_, stdErr := runDepLab([]string{
				"--image", inputImage,
				"--git", pathToGitRepo,
				"--metadata-file", "doesnotmatter8",
			}, 1)
			errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
			Expect(errorOutput).To(
				SatisfyAll(
					ContainSubstring("could not load image"),
					ContainSubstring("scratch"),
				))
		})

		It("throws an error if trying to pull an invalid image", func() {
			By("executing it")
			inputImage := "swkichtlsmhasd" // random string unlikely for an image ever to exist
			_, stdErr := runDepLab([]string{
				"--image", inputImage,
				"--git", pathToGitRepo,
				"--metadata-file", "doesnotmatter9",
			}, 1)

			errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
			Expect(errorOutput).To(
				SatisfyAll(
					ContainSubstring("could not load image"),
					ContainSubstring("swkichtlsmhasd"),
				))
		})

		It("exits with an error if neither image or image-tar flags are set", func() {
			_, stdErr := runDepLab([]string{"--git", "does-not-matter"}, 1)
			errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
			Expect(errorOutput).To(ContainSubstring("ERROR: requires one of --image or --image-tar"))
		})

		It("exits with an error if neither metadata-file, dpkg-list, output-tar flags are set", func() {
			_, stdErr := runDepLab([]string{"--git", pathToGitRepo,
				"--image", "ubuntu:bionic"}, 1)
			errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
			Expect(errorOutput).To(ContainSubstring("ERROR: requires one of --metadata-file, --dpkg-file, or --output-tar"))
		})

		It("exits with an error if both image and image-tar flags are set", func() {
			_, stdErr := runDepLab([]string{"--image", "foo",
				"--image-tar", "path/to/image.tar",
				"--git", "does-not-matter",
				"--metadata-file", "doesnotmatter10",
			}, 1)
			errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
			Expect(errorOutput).To(ContainSubstring("ERROR: cannot accept both --image and --image-tar"))
		})

		It("throws an error if invalid characters are in image name", func() {
			By("executing it")
			inputImage := "£$Invalid_image_name$£"
			_, stdErr := runDepLab([]string{
				"--image", inputImage,
				"--git", pathToGitRepo,
				"--metadata-file", "doesnotmatter11",
			}, 1)

			errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
			Expect(errorOutput).To(ContainSubstring("could not parse reference"))
		})

		It("exits with an error if additional-source-url is not valid", func() {
			_, stdErr := runDepLab([]string{
				"--image", "ubuntu:bionic",
				"--git", pathToGitRepo,
				"--metadata-file", "doesnotmatter12",
				"--additional-source-url", "/foo/bar",
			}, 1)

			errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
			Expect(errorOutput).To(SatisfyAll(
				ContainSubstring("/foo/bar"),
				ContainSubstring("error"),
				ContainSubstring("failed to validate additional source url")))
		})

		It("exits with an error if additional-source-url is not reachable ", func() {
			_, stdErr := runDepLab([]string{
				"--image", "ubuntu:bionic",
				"--git", pathToGitRepo,
				"--metadata-file", "doesnotmatter13",
				"--additional-source-url", "https://package.some.invalid/cool-package",
			}, 1)

			errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
			Expect(errorOutput).To(SatisfyAll(
				ContainSubstring("https://package.some.invalid/cool-package"),
				ContainSubstring("error"),
				ContainSubstring("failed to validate additional source url")))
		})

		It("exits with an error if additional-source-url is not returning a success status code ", func() {
			server := startServer(ghttp.RespondWith(http.StatusNotFound, []byte("HTTP status not found code returned")))
			defer server.Close()

			address := server.URL() + "/cool-package"

			_, stdErr := runDepLab([]string{
				"--image", "ubuntu:bionic",
				"--git", pathToGitRepo,
				"--additional-source-url", address,
				"--metadata-file", "doesnotmatter14",
			}, 1)

			errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
			Expect(errorOutput).To(SatisfyAll(
				ContainSubstring(address),
				ContainSubstring("error"),
				ContainSubstring("failed to validate additional source url")))
		})

		Context("when ignore-validation-errors flag is set", func() {
			It("succeeds with a warning if additional-source-url is not valid", func() {
				d, err := ioutil.TempDir("", "deplab-integration-test-")
				metadataFileName := d + "/metadata-file.yml"
				Expect(err).To(Not(HaveOccurred()))
				defer os.Remove(d)
				_, stdErr := runDepLab([]string{
					"--image", "ubuntu:bionic",
					"--git", pathToGitRepo,
					"--additional-source-url", "/foo/bar",
					"--metadata-file", metadataFileName,
					"--ignore-validation-errors",
				}, 0)

				errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
				Expect(errorOutput).To(SatisfyAll(
					ContainSubstring("/foo/bar"),
					ContainSubstring("warning"),
					ContainSubstring("failed to validate additional source url")))
			})

			It("succeeds with a warning for multiple invalid additional-source-urls", func() {
				f, err := ioutil.TempFile("", "")
				Expect(err).ToNot(HaveOccurred())

				defer os.Remove(f.Name())

				server := startServer(
					ghttp.RespondWith(http.StatusNotFound, []byte("HTTP status not found code returned")),
					ghttp.RespondWith(http.StatusOK, ""),
					ghttp.RespondWith(http.StatusOK, ""))
				defer server.Close()

				addresses := []string{
					"/foo/bar",
					"https://package.some.invalid/unreachable-package",
					server.URL() + "/404-package",
					server.URL() + "/invalid-extension",
					server.URL() + "/should-pass.tgz",
				}

				_, stdErr := runDepLab([]string{
					"--image", "ubuntu:bionic",
					"--git", pathToGitRepo,
					"--additional-source-url", addresses[0],
					"--additional-source-url", addresses[1],
					"--additional-source-url", addresses[2],
					"--additional-source-url", addresses[3],
					"--additional-source-url", addresses[4],
					"--metadata-file", f.Name(),
					"--ignore-validation-errors",
				}, 0)

				By("reporting a validation warning for all invalid urls")
				errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
				Expect(errorOutput).To(SatisfyAll(
					ContainSubstring(addresses[0]),
					ContainSubstring(addresses[1]),
					ContainSubstring(addresses[2]),
					ContainSubstring(addresses[3]),
					Not(ContainSubstring(addresses[4])),
					ContainSubstring("warning"),
					ContainSubstring("failed to validate additional source url"),
				))

				By("by including an archive dependency")
				metadataLabel := metadata.Metadata{}
				err = json.NewDecoder(f).Decode(&metadataLabel)
				Expect(err).ToNot(HaveOccurred())

				archiveDependencies := selectArchiveDependencies(metadataLabel.Dependencies)
				Expect(archiveDependencies).To(HaveLen(5))
				for i, archiveDependency := range archiveDependencies {

					By("adding the invalid additional-source-url to the archive dependency")
					Expect(archiveDependency.Type).ToNot(BeEmpty())
					Expect(archiveDependency.Type).To(Equal("package"))
					Expect(archiveDependency.Source.Type).To(Equal(providers.ArchiveType))

					archiveSourceMetadata := archiveDependency.Source.Metadata.(map[string]interface{})
					Expect(archiveSourceMetadata["url"]).To(Equal(addresses[i]))
				}
			})
		})
	})
})
