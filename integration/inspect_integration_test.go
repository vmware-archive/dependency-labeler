package integration_test

import (
	"encoding/json"
	"path/filepath"

	"github.com/pivotal/deplab/metadata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("deplab inspect", func() {
	Context("with a deplab-d image tar path", func() {
		BeforeEach(func() {

		})

		It("prints the label", func() {
			inputTarPath := filepath.Join("assets", "tiny-deplabd.tgz")
			inputTarPath, err := filepath.Abs(inputTarPath)
			Expect(err).ToNot(HaveOccurred())
			stdOut, _ := runDepLab([]string{
				"inspect",
				"--image-tar", inputTarPath,
			}, 0)

			md := metadata.Metadata{}
			err = json.NewDecoder(stdOut).Decode(&md)

			Expect(err).ToNot(HaveOccurred())

			Expect(md.Provenance[0].Name).To(Equal("deplab"))
		})
	})
	Context("with a undeplabbed image tar path", func() {
		It("provides an error", func() {
			inputTarPath := filepath.Join("assets", "tiny.tgz")
			inputTarPath, err := filepath.Abs(inputTarPath)
			Expect(err).ToNot(HaveOccurred())
			_, stderr := runDepLab([]string{
				"inspect",
				"--image-tar", inputTarPath,
			}, 1)

			Expect(getContentsOfReader(stderr)).To(
				SatisfyAll(
					ContainSubstring("deplab cannot find the 'io.pivotal.metadata' label on the provided image"),
					ContainSubstring("tiny.tgz")))
		})
	})

	Context("with a invalid image", func() {
		It("provides an error", func() {
			inputTarPath := filepath.Join("assets", "invalid-image-archive.tgz")
			inputTarPath, err := filepath.Abs(inputTarPath)
			Expect(err).ToNot(HaveOccurred())
			_, stderr := runDepLab([]string{
				"inspect",
				"--image-tar", inputTarPath,
			}, 1)

			Expect(getContentsOfReader(stderr)).To(
				SatisfyAll(
					ContainSubstring("deplab cannot open the provided image"),
					ContainSubstring("invalid-image-archive.tgz")))
		})
	})

	Context("with a valid image but invalid json label", func() {
		It("provides an error", func() {
			inputTarPath := filepath.Join("assets", "tiny-with-invalid-label.tgz")
			inputTarPath, err := filepath.Abs(inputTarPath)
			Expect(err).ToNot(HaveOccurred())
			_, stderr := runDepLab([]string{
				"inspect",
				"--image-tar", inputTarPath,
			}, 1)

			Expect(getContentsOfReader(stderr)).To(
				SatisfyAll(
					ContainSubstring("deplab cannot parse the label"),
					ContainSubstring("tiny-with-invalid-label.tgz")))
		})
	})
})
