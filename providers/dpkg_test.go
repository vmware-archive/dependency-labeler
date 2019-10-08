package providers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/deplab/metadata"
	"github.com/pivotal/deplab/providers"
)

var _ = Describe("Providers/Dpkg", func() {
	Describe("Digest", func() {
		It("generate a sha256 of the input", func() {
			out := providers.Digest(metadata.DebianPackageListSourceMetadata{
				AptSources: []string{"deb example.com bionic main universe"},
				Packages: []metadata.Package{
					{
						Package: "foo",
						Source: metadata.PackageSource{
							Version: "4.2",
						},
					},
				},
			})

			Expect(out).To(Equal("7ecf7aa2c71ee01ec2b90f37a3b8e944158e9aea6b8cee0290a7cb187884cf4c"))
		})

		It("generates the same digest for 2 different input instances with the same content", func() {
			input1 := metadata.DebianPackageListSourceMetadata{
				AptSources: []string{"deb example.com bionic main universe"},
				Packages: []metadata.Package{
					{
						Package: "foo",
						Source: metadata.PackageSource{
							Version: "4.2",
						},
					},
				},
			}
			input2 := metadata.DebianPackageListSourceMetadata{
				AptSources: []string{"deb example.com bionic main universe"},
				Packages: []metadata.Package{
					{
						Package: "foo",
						Source: metadata.PackageSource{
							Version: "4.2",
						},
					},
				},
			}

			out1 := providers.Digest(input1)
			out2 := providers.Digest(input2)

			Expect(&input1 != &input2).To(BeTrue())
			Expect(out2).To(Equal(out1))
		})

		It("generates a different digest for 2 different input instances with the same content in different order", func() {
			input1 := metadata.DebianPackageListSourceMetadata{
				AptSources: []string{"deb example.com bionic main universe"},
				Packages: []metadata.Package{
					{
						Package: "foo",
						Source: metadata.PackageSource{
							Version: "4.2",
						},
					},
					{
						Package: "bar",
						Source: metadata.PackageSource{
							Version: "1.0",
						},
					},
				},
			}
			input2 := metadata.DebianPackageListSourceMetadata{
				AptSources: []string{"deb example.com bionic main universe"},
				Packages: []metadata.Package{
					{
						Package: "bar",
						Source: metadata.PackageSource{
							Version: "1.0",
						},
					},
					{
						Package: "foo",
						Source: metadata.PackageSource{
							Version: "4.2",
						},
					},
				},
			}

			out1 := providers.Digest(input1)
			out2 := providers.Digest(input2)

			Expect(&input1 != &input2).To(BeTrue())
			Expect(out2).ToNot(Equal(out1))
		})

		It("generates different digest for 2 different inputs", func() {
			input1 := metadata.DebianPackageListSourceMetadata{
				AptSources: []string{"deb example.com bionic main universe"},
				Packages: []metadata.Package{
					{
						Package: "bar",
						Source: metadata.PackageSource{
							Version: "4.2",
						},
					},
				},
			}
			input2 := metadata.DebianPackageListSourceMetadata{
				AptSources: []string{"deb example.com bionic main universe"},
				Packages: []metadata.Package{
					{
						Package: "foo",
						Source: metadata.PackageSource{
							Version: "4.2",
						},
					},
				},
			}

			out1 := providers.Digest(input1)
			out2 := providers.Digest(input2)

			Expect(out2).ToNot(Equal(out1))
		})
	})

	Describe("ParseStatDBEntry", func() {
		It("parses a StatDB entry string", func() {
			Expect(providers.ParseStatDBEntry(`Package: libgcc1
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
Original-Maintainer: Debian GCC Maintainers <debian-gcc@lists.debian.org>`)).To(Equal(metadata.Package{
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

		It("returns error if entry does not contain Package:", func() {
			_, err := providers.ParseStatDBEntry("\n")
			Expect(err).To(HaveOccurred())
		})
	})
})
