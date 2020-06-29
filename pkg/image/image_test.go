// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package image_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pivotal/deplab/pkg/metadata"

	"github.com/google/go-containerregistry/pkg/crane"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal/deplab/pkg/image"
)

var _ = Describe("Image", func() {

	Describe("NewDeplabImage", func() {
		Context("with valid inputs", func() {
			var (
				image RootFSImage
				err   error
			)

			It("[remote-image][private-registry] instantiates an image starting from a remote source", func() {
				image, err = NewDeplabImage("dev.registry.pivotal.io/navcon/deplab-test-asset:all-file-types", "")

				Expect(err).ToNot(HaveOccurred())
				Expect(image).ToNot(BeNil())
			})

			It("instantiates an image starting from a tarball", func() {
				inputTarPath, err := filepath.Abs("../../test/integration/assets/image-archives/all-file-types.tgz")
				Expect(err).ToNot(HaveOccurred())

				image, err = NewDeplabImage("", inputTarPath)

				Expect(err).ToNot(HaveOccurred())
				Expect(image).ToNot(BeNil())
			})

			AfterEach(func() {
				image.Cleanup()
			})
		})

		Context("when cannot be instantiated", func() {
			It("returns an error if no image at the remote source", func() {
				_, err := NewDeplabImage("pivotalnavcon/this-does-not-exists", "")

				Expect(err).To(HaveOccurred())
			})

			It("returns an error if an invalid image at the tarball path", func() {
				inputTarPath, err := filepath.Abs("../../test/integration/assets/image-archives/invalid-image-archive.tgz")
				Expect(err).ToNot(HaveOccurred())

				_, err = NewDeplabImage("", inputTarPath)
				Expect(err).To(HaveOccurred())
			})

			It("returns an error if no image at the tarball path", func() {
				_, err := NewDeplabImage("", "non-existing-tar-ball")

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("ExportWithMetadata", func() {
		Context("when saving the image to a tar", func() {
			var (
				image RootFSImage
				dir   string
			)

			AfterEach(func() {
				image.Cleanup()
				err := os.RemoveAll(dir)
				Expect(err).ToNot(HaveOccurred())
			})

			It("includes metadata in the label", func() {
				inputTarPath, err := filepath.Abs("../../test/integration/assets/image-archives/all-file-types.tgz")
				Expect(err).ToNot(HaveOccurred())

				image, err = NewDeplabImage("", inputTarPath)
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

				Expect(cf.Config.Labels["io.deplab.metadata"]).To(MatchJSON(`{
					"base": null,
					"provenance": null,
					"dependencies": null
				}`))

				By("keeping the existing labels")
				Expect(cf.Config.Labels["foo"]).To(Equal("bar"))
			})

			It("returns an error if the destination path is invalid", func() {
				inputTarPath, err := filepath.Abs("../../test/integration/assets/image-archives/all-file-types.tgz")
				Expect(err).ToNot(HaveOccurred())

				image, err = NewDeplabImage("", inputTarPath)
				Expect(err).ToNot(HaveOccurred())

				err = image.ExportWithMetadata(metadata.Metadata{}, "/tmp/this-path-does-not-exist/this-file-does-not-matter", "")
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("AbsolutePath", func() {
		Context("relative path to rootFS location", func() {
			var image RootFSImage
			It("provides the absolute path", func() {
				inputTarPath, err := filepath.Abs("../../test/integration/assets/image-archives/all-file-types.tgz")
				Expect(err).ToNot(HaveOccurred())

				image, err = NewDeplabImage("", inputTarPath)
				Expect(err).ToNot(HaveOccurred())

				path, err := image.AbsolutePath("/var/lib/rpm")
				Expect(err).ToNot(HaveOccurred())
				Expect(path).To(SatisfyAll(
					ContainSubstring("/var/lib/rpm"),
					ContainSubstring(RootfsPrefix)))
			})

			AfterEach(func() {
				image.Cleanup()
			})
		})
	})
})
