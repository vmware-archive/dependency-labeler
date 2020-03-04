package dpkg_test

import (
	"github.com/pivotal/deplab/pkg/common"
	"github.com/pivotal/deplab/test/test_utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal/deplab/pkg/dpkg"
	"github.com/pivotal/deplab/pkg/metadata"
)

var _ = Describe("Dpkg", func() {
	Describe("ParseStatDBEntry", func() {
		It("parses a StatDB entry string", func() {
			Expect(ParseStatDBEntry(`Package: libgcc1
Status: install ok installed
Priority: optional
Section: libs
Installed-Size: 112
Maintainer: Ubuntu Core developers <ubuntu-devel-discuss@lists.ubuntu.com>
Architecture: amd64
Multi-Arch: same
Source: gcc-8 (8.3.0-6ubuntu1~18.04.1)
Version: 1:8.3.0-6ubuntu1~18.04.1
Depends: gcc-8-base (= 8.3.0-6ubuntu1~18.04.1), libc6 (>= 2.14)
Breaks: gcc-4.3 (<< 4.3.6-1), gcc-4.4 (<< 4.4.6-4), gcc-4.5 (<< 4.5.3-2)
Description: GCC support library
 Shared version of the support: library, a library of internal subroutines
 that GCC uses to overcome shortcomings of particular machines, or
 special needs for some languages.
Homepage: http://gcc.gnu.org/
Original-Maintainer: Debian GCC Maintainers <debian-gcc@lists.debian.org>`)).To(Equal(metadata.DpkgPackage{
				Package:      "libgcc1",
				Version:      "1:8.3.0-6ubuntu1~18.04.1",
				Architecture: "amd64",
				Source: metadata.PackageSource{
					Package:         "gcc-8",
					Version:         "8.3.0-6ubuntu1~18.04.1",
					UpstreamVersion: "8.3.0",
				},
			}))
		})

		It("returns error if entry does not contain DpkgPackage:", func() {
			_, err := ParseStatDBEntry("\n")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Provider", func() {
		Context("when the image has no database packages", func() {
			It("does not modify the metadata content", func() {
				md, err := Provider(test_utils.MockImage{}, common.RunParams{}, metadata.Metadata{})
				Expect(err).NotTo(HaveOccurred())

				Expect(md).To(Equal(metadata.Metadata{}))
			})
		})
	})
})
