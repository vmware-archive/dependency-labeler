package providers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/deplab/metadata"
	"github.com/pivotal/deplab/providers"
)

var _ = Describe("Providers/Dpkg", func() {
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
