package integration_test

import (
	"encoding/json"
	"log"
	"path/filepath"

	"github.com/pivotal/deplab/metadata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("deplab inspect", func() {
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
	)

	DescribeTable("provides an error", func(flag, path string) {
		_, stderr := runDepLab([]string{
			"inspect",
			flag, path,
		}, 1)

		Expect(getContentsOfReader(stderr)).To(
			SatisfyAll(
				ContainSubstring("deplab cannot find the 'io.pivotal.metadata' label on the provided image"),
				ContainSubstring("tiny.tgz")))
	},
		Entry("with a undeplab'd image tar path", "--image-tar", getTestAssetPath("tiny.tgz")),
	)

	DescribeTable("provides an error", func(flag, path string) {
		_, stderr := runDepLab([]string{
			"inspect",
			flag, path,
		}, 1)

		Expect(getContentsOfReader(stderr)).To(
			SatisfyAll(
				ContainSubstring("deplab cannot open the provided image"),
				ContainSubstring("invalid-image-archive.tgz")))
	},
		Entry("with a invalid image tarball", "--image-tar", getTestAssetPath("invalid-image-archive.tgz")),
	)

	DescribeTable("provides an error", func(flag, path string) {
		_, stderr := runDepLab([]string{
			"inspect",
			flag, path,
		}, 1)

		Expect(getContentsOfReader(stderr)).To(
			SatisfyAll(
				ContainSubstring("deplab cannot parse the label"),
				ContainSubstring("tiny-with-invalid-label.tgz")))
	},
		Entry("with a valid image but invalid json label", "--image-tar", getTestAssetPath("tiny-with-invalid-label.tgz")),
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
