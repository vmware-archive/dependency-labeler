package rootfs_test

import (
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/deplab/rootfs"
)

var rfs rootfs.RootFS

var _ = Describe("rootfs", func() {
	Context("when there is a valid image archive", func() {
		BeforeEach(func() {
			inputTarPath, err := filepath.Abs(filepath.Join("..", "integration", "assets", "tiny.tgz"))
			Expect(err).ToNot(HaveOccurred())

			rfs, err = rootfs.New(inputTarPath)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when the path exists", func() {
			It("retrieves the content of a file", func() {
				osRelease, err := rfs.GetFileContent("/etc/os-release")
				Expect(err).ToNot(HaveOccurred())

				Expect(osRelease).To(
					SatisfyAll(
						ContainSubstring(`NAME="Pivotal Tiny`),
						ContainSubstring(`VERSION_ID="dev"`),
						ContainSubstring(`VERSION_CODENAME="dev"`),
					))

			})

			It("retrieves the content of all files of a directory", func() {
				statusFiles, err := rfs.GetDirContents("/var/lib/dpkg/status.d")
				Expect(err).ToNot(HaveOccurred())

				Expect(statusFiles).To(
					SatisfyAll(
						HaveLen(6),
						ConsistOf(
							ContainSubstring("Package: base-files"),
							ContainSubstring("Package: libc6"),
							ContainSubstring("Package: libssl1.1"),
							ContainSubstring("Package: netbase"),
							ContainSubstring("Package: openssl"),
							ContainSubstring("Package: tzdata"),
						)))
			})

			It("ignores subdirectory in the given directory", func() {
				statusFiles, err := rfs.GetDirContents("/var/lib/dpkg")
				Expect(err).ToNot(HaveOccurred())

				Expect(statusFiles).To(
					SatisfyAll(
						HaveLen(1),
						ConsistOf("\n")))
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

				_, err = rfs.GetFileContent("/etc/os-release")
				Expect(err).ToNot(HaveOccurred())

				_, err = rfs.GetDirContents("/var/lib/dpkg/status.d")
				Expect(err).ToNot(HaveOccurred())

				rfs.Cleanup()

				_, err = rfs.GetFileContent("/etc/os-release")
				Expect(err).To(HaveOccurred())

				_, err = rfs.GetDirContents("/var/lib/dpkg/status.d")
				Expect(err).To(HaveOccurred())
			})
		})

		AfterEach(func() {
			rfs.Cleanup()
		})
	})
	Context("when there is not a valid image archive", func() {
		It("returns an error", func() {
			inputTarPath, err := filepath.Abs(filepath.Join("..", "integration", "assets", "invalid-image-archive.tgz"))
			Expect(err).ToNot(HaveOccurred())

			_, err = rootfs.New(inputTarPath)
			Expect(err).To(MatchError(
				SatisfyAll(
					ContainSubstring("Could not load image from path"),
					ContainSubstring("invalid-image-archive.tgz"),
				)))
		})
	})
})
