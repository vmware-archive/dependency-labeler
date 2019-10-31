package rootfs_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/pivotal/deplab/metadata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/deplab/rootfs"
)

var _ = Describe("Image", func() {

	Describe("NewDeplabImage", func() {
		Context("with valid inputs", func() {
			var (
				image rootfs.Image
				err   error
			)

			It("instantiates an image starting from a remote source", func() {
				image, err = rootfs.NewDeplabImage("pivotalnavcon/test-asset-all-file-types", "")

				Expect(err).ToNot(HaveOccurred())
				Expect(image).ToNot(BeNil())
			})

			It("instantiates an image starting from a tarball", func() {
				inputTarPath, err := filepath.Abs("../integration/assets/all-file-types.tgz")
				Expect(err).ToNot(HaveOccurred())

				image, err = rootfs.NewDeplabImage("", inputTarPath)

				Expect(err).ToNot(HaveOccurred())
				Expect(image).ToNot(BeNil())
			})

			AfterEach(func() {
				image.Cleanup()
			})
		})

		Context("when cannot be instantiated", func() {
			It("returns an error if no image at the remote source", func() {
				_, err := rootfs.NewDeplabImage("pivotgarbledile-types", "")

				Expect(err).To(HaveOccurred())
			})

			It("returns an error if an invalid image at the tarball path", func() {
				inputTarPath, err := filepath.Abs("integration/assets/invalid-image-archive.tgz")
				Expect(err).ToNot(HaveOccurred())

				_, err = rootfs.NewDeplabImage("", inputTarPath)
				Expect(err).To(HaveOccurred())
			})

			It("returns an error if no image at the tarball path", func() {
				_, err := rootfs.NewDeplabImage("", "non-existing-tar-ball")

				Expect(err).To(HaveOccurred())
			})
		})

	})

	Describe("ExportWithMetadata", func() {
		Context("when saving the image to a tar", func() {
			var (
				image rootfs.Image
				dir   string
			)

			AfterEach(func() {
				image.Cleanup()
				err := os.RemoveAll(dir)
				Expect(err).ToNot(HaveOccurred())
			})

			It("includes metadata in the label", func() {
				inputTarPath, err := filepath.Abs("../integration/assets/all-file-types.tgz")
				Expect(err).ToNot(HaveOccurred())

				image, err = rootfs.NewDeplabImage("", inputTarPath)
				Expect(err).ToNot(HaveOccurred())

				dir, err = ioutil.TempDir("", "deplab-")
				Expect(err).ToNot(HaveOccurred())

				destinationImage := filepath.Join(dir, "output-image.tar")
				err = image.ExportWithMetadata(metadata.Metadata{}, destinationImage, "")
				Expect(err).ToNot(HaveOccurred())

				labelledImage, err := crane.Load(destinationImage)
				Expect(err).ToNot(HaveOccurred())

				cf, err := labelledImage.ConfigFile()
				Expect(err).ToNot(HaveOccurred())

				Expect(cf.Config.Labels["io.pivotal.metadata"]).To(MatchJSON(`{
					"base": null,
					"provenance": null,
					"dependencies": null
				}`))

				By("keeping the existing labels")
				Expect(cf.Config.Labels["foo"]).To(Equal("bar"))
			})

			It("returns an error if the destination path is invalid", func() {
				inputTarPath, err := filepath.Abs("../integration/assets/all-file-types.tgz")
				Expect(err).ToNot(HaveOccurred())

				image, err = rootfs.NewDeplabImage("", inputTarPath)
				Expect(err).ToNot(HaveOccurred())

				err = image.ExportWithMetadata(metadata.Metadata{}, "/tmp/this-path-does-not-exist/this-file-does-not-matter", "")
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
