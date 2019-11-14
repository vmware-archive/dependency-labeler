package integration_test

import (
	"encoding/json"
	"log"
	"path/filepath"
	"strings"

	"github.com/pivotal/deplab/metadata"

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

	DescribeTable("prints the label", func(flag, path string) {
		stdOut, _ := runDepLab([]string{
			"inspect",
			flag, path,
		}, 0)

		md := metadata.Metadata{}
		err := json.NewDecoder(stdOut).Decode(&md)

		Expect(err).ToNot(HaveOccurred())
		Expect(md.Provenance[0].Name).To(Equal("deplab"))
	},
		Entry("with a deplab'd image tarball", "--image-tar", getTestAssetPath("tiny-deplabd.tgz")),
		Entry("with a deplab'd image from a registry", "--image", "pivotalnavcon/test-asset-tiny-deplabd"),
	)

	DescribeTable("provides an error", func(flag, path string) {
		_, stderr := runDepLab([]string{
			"inspect",
			flag, path,
		}, 1)

		Expect(getContentsOfReader(stderr)).To(
			SatisfyAll(
				ContainSubstring("deplab cannot find the 'io.pivotal.metadata' label on the provided image"),
				ContainSubstring(path)))
	},
		Entry("with a undeplab'd image tar path", "--image-tar", getTestAssetPath("tiny.tgz")),
		Entry("with a undeplab'd image from a registry", "--image", "cloudfoundry/run:tiny"),
	)

	DescribeTable("provides an error", func(flag, path string) {
		_, stderr := runDepLab([]string{
			"inspect",
			flag, path,
		}, 1)

		Expect(getContentsOfReader(stderr)).To(
			SatisfyAll(
				ContainSubstring("deplab cannot open the provided image"),
				ContainSubstring(path)))
	},
		Entry("with a invalid image tarball", "--image-tar", getTestAssetPath("invalid-image-archive.tgz")),
		Entry("with a non-existent image from registry", "--image", "pivotalnavcon/does-not-exist"),
	)

	DescribeTable("provides an error", func(flag, path string) {
		_, stderr := runDepLab([]string{
			"inspect",
			flag, path,
		}, 1)

		Expect(getContentsOfReader(stderr)).To(
			SatisfyAll(
				ContainSubstring("deplab cannot parse the label"),
				ContainSubstring(path)))
	},
		Entry("with a valid image tar ball with invalid json label", "--image-tar", getTestAssetPath("tiny-with-invalid-label.tgz")),
		Entry("with a valid image from a registry with invalid json label", "--image", "pivotalnavcon/test-asset-tiny-with-invalid-label"),
	)
})

func getTestAssetPath(path string) string {
	inputTarPath := filepath.Join("assets", path)
	inputTarPath, err := filepath.Abs(inputTarPath)
	if err != nil {
		log.Fatalf("Could not find test asset %s: %s", path, err)
	}
	return inputTarPath
}
