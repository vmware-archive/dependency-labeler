package integration_test

import (
	"encoding/json"
	"strings"

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
		Entry("with a deplab'd image tarball", "--image-tar", getTestAssetPath("image-archives/tiny-deplabd.tgz")),
		Entry("[remote-image][private-registry] with a deplab'd image from a registry", "--image", "dev.registry.pivotal.io/navcon/deplab-test-asset:tiny-deplabd"),
	)

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
		Entry("with a undeplab'd image tar path", "--image-tar", getTestAssetPath("image-archives/tiny.tgz"), "deplab cannot find the 'io.pivotal.metadata' label on the provided image"),
		Entry("[remote-image] with a undeplab'd image from a registry", "--image", "cloudfoundry/run:tiny", "deplab cannot find the 'io.pivotal.metadata' label on the provided image"),
		Entry("with a invalid image tarball", "--image-tar", getTestAssetPath("image-archives/invalid-image-archive.tgz"), "deplab cannot open the provided image"),
		Entry("with a non-existent image from registry", "--image", "pivotalnavcon/does-not-exist", "deplab cannot retrieve the Config file"),
		Entry("with a valid image tar ball with invalid json label", "--image-tar", getTestAssetPath("image-archives/tiny-with-invalid-label.tgz"), "deplab cannot parse the label"),
		Entry("[remote-image][private-registry] with a valid image from a registry with invalid json label", "--image", "dev.registry.pivotal.io/navcon/deplab-test-asset:tiny-with-invalid-label", "deplab cannot parse the label"),
	)
})
