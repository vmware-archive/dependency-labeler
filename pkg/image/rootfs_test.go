package image_test

import (
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/crane"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/pivotal/deplab/pkg/image"
)

var rfs RootFS

var _ = Describe("rootFS", func() {
	Context("when there is a valid image archive", func() {
		BeforeEach(func() {
			inputTarPath, err := filepath.Abs("../../test/integration/assets/image-archives/all-file-types.tgz")
			Expect(err).ToNot(HaveOccurred())

			image, err := crane.Load(inputTarPath)
			Expect(err).ToNot(HaveOccurred())

			rfs, err = NewRootFS(image, nil)
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

			It("retrieves all the file names inside a directory", func() {
				statusFiles, err := rfs.GetDirFileNames("/all-files")
				Expect(err).ToNot(HaveOccurred())

				Expect(statusFiles).To(
					SatisfyAll(
						HaveLen(3),
						ConsistOf(
							ContainSubstring("hard-link-file"),
							ContainSubstring("start-file"),
							ContainSubstring("symbolic-link-file"),
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
				Expect(err).To(MatchError(ContainSubstring("could not find file in rootFS")))
			})

			It("returns an error when trying to access a directory", func() {
				_, err := rfs.GetDirContents("/this/directory/does/not/exist")
				Expect(err).To(MatchError(ContainSubstring("could not find directory in rootFS")))
			})

			It("returns an error when trying to access a directory", func() {
				_, err := rfs.GetDirFileNames("/this/directory/does/not/exist")
				Expect(err).To(MatchError(ContainSubstring("could not find directory in rootFS")))
			})
		})

		Context("when rootFS is cleaned", func() {
			It("can no longer retrieve the content", func() {
				var err error

				_, err = rfs.GetFileContent("/all-files/start-file")
				Expect(err).ToNot(HaveOccurred())

				_, err = rfs.GetDirContents("/all-files/folder")
				Expect(err).ToNot(HaveOccurred())

				_, err = rfs.GetDirFileNames("/all-files")
				Expect(err).ToNot(HaveOccurred())

				rfs.Cleanup()

				_, err = rfs.GetFileContent("/all-files/start-file")
				Expect(err).To(HaveOccurred())

				_, err = rfs.GetDirContents("/all-files/folder")
				Expect(err).To(HaveOccurred())

				_, err = rfs.GetDirFileNames("/all-files")
				Expect(err).To(HaveOccurred())
			})
		})

		AfterEach(func() {
			rfs.Cleanup()
		})
	})
	Context("when the tar contains a directory with no permissions", func() {
		Context("when the offending file is excluded", func() {
			It("succeeds creating the rootfs", func() {
				inputTarPath, err := filepath.Abs("../../test/integration/assets/image-archives/broken-files.tgz")
				Expect(err).ToNot(HaveOccurred())

				image, err := crane.Load(inputTarPath)
				Expect(err).ToNot(HaveOccurred())

				rootfs, err := NewRootFS(image, []string{"all-files/broken-folder/"})
				defer rootfs.Cleanup()
				Expect(err).ToNot(HaveOccurred())

				_, err = rootfs.GetDirContents("all-files/broken-folder/no-permissions")
				Expect(err).To(MatchError(ContainSubstring("all-files/broken-folder/no-permissions")))
			})
		})
	})
	Context("when the image contains char device file", func() {
		It("successfully does something", func() {
			inputTarPath, err := filepath.Abs("../../test/integration/assets/image-archives/char-device.tgz")
			Expect(err).ToNot(HaveOccurred())

			image, err := crane.Load(inputTarPath)
			Expect(err).ToNot(HaveOccurred())

			rootfs, err := NewRootFS(image, []string{})
			defer rootfs.Cleanup()
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
