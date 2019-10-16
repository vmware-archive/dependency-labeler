package rootfs_test

import (
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/crane"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/deplab/rootfs"
)

var rfs rootfs.RootFS

var _ = Describe("rootfs", func() {
	Context("when there is a valid image archive", func() {
		BeforeEach(func() {
			inputTarPath, err := filepath.Abs("../integration/assets/all-file-types.tgz")
			Expect(err).ToNot(HaveOccurred())

			image, err := crane.Load(inputTarPath)
			Expect(err).ToNot(HaveOccurred())

			rfs, err = rootfs.New(image)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when the path exists", func() {
			It("retrieves the content of a file", func() {
				aFile, err := rfs.GetFileContent("/all-files/start-file")
				Expect(err).ToNot(HaveOccurred())

				Expect(aFile).To(ContainSubstring("hello world"))

			})

			It("retrieves the content of all files of a directory", func() {
				statusFiles, err := rfs.GetDirContents("/all-files/folder")
				Expect(err).ToNot(HaveOccurred())

				Expect(statusFiles).To(
					SatisfyAll(
						HaveLen(2),
						ConsistOf(
							ContainSubstring("foo"),
							ContainSubstring("bar"),
						)))
			})

			It("ignores subdirectory in the given directory", func() {
				statusFiles, err := rfs.GetDirContents("/all-files")
				Expect(err).ToNot(HaveOccurred())

				Expect(statusFiles).To(
					SatisfyAll(
						HaveLen(3),
						ConsistOf(
							ContainSubstring("hello world"),
							ContainSubstring("hello world"),
							ContainSubstring("hello world"),
						)))
			})
		})

		Context("when the path does not exists", func() {
			It("returns an error when trying to access a file", func() {
				_, err := rfs.GetFileContent("/this/file/does/not/exist.txt")
				Expect(err).To(MatchError(ContainSubstring("could not find file in rootfs")))
			})

			It("returns an error when trying to access a directory", func() {
				_, err := rfs.GetDirContents("/this/directory/does/not/exist")
				Expect(err).To(MatchError(ContainSubstring("could not find directory in rootfs")))
			})
		})

		Context("when rootfs is cleaned", func() {
			It("can no longer retrieve the content", func() {
				var err error

				_, err = rfs.GetFileContent("/all-files/start-file")
				Expect(err).ToNot(HaveOccurred())

				_, err = rfs.GetDirContents("/all-files/folder")
				Expect(err).ToNot(HaveOccurred())

				rfs.Cleanup()

				_, err = rfs.GetFileContent("/all-files/start-file")
				Expect(err).To(HaveOccurred())

				_, err = rfs.GetDirContents("/all-files/folder")
				Expect(err).To(HaveOccurred())
			})
		})

		AfterEach(func() {
			rfs.Cleanup()
		})
	})

	// test to convert for DeplabImage instantation

	//Context("when there is not a valid image archive", func() {
	//	It("returns an error", func() {
	//		inputTarPath, err := filepath.Abs(filepath.Join("..", "integration", "assets", "invalid-image-archive.tgz"))
	//		Expect(err).ToNot(HaveOccurred())
	//
	//		image, err := crane.Load(inputTarPath)
	//		Expect(err).ToNot(HaveOccurred())
	//
	//		_, err = rootfs.New(image)
	//		Expect(err).To(MatchError(
	//			SatisfyAll(
	//				ContainSubstring("Could not load image from tar"),
	//				ContainSubstring("invalid-image-archive.tgz"),
	//			)))
	//	})
	//})
})
