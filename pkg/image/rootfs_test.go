package image_test

import (
	"path/filepath"

	v1 "github.com/google/go-containerregistry/pkg/v1"

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
		})

		Context("when rootFS is cleaned", func() {
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
	Context("when the tar contains a directory with no permissions", func() {
		var image v1.Image
		BeforeEach(func() {
			inputTarPath, err := filepath.Abs("../../test/integration/assets/image-archives/broken-files.tgz")
			Expect(err).ToNot(HaveOccurred())

			image, err = crane.Load(inputTarPath)
			Expect(err).ToNot(HaveOccurred())
		})
		It("can no longer retrieve the content", func() {
			_, err := NewRootFS(image, nil)
			Expect(err).To(HaveOccurred())
		})
		Context("when the offending file is excluded", func() {
			It("succeeds creating the rootfs", func() {
				_, err := NewRootFS(image, []string{"all-files/broken-folder/"})
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
