// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package common_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/vmware-tanzu/dependency-labeler/pkg/common"
	"github.com/vmware-tanzu/dependency-labeler/pkg/metadata"
)

var _ = Describe("Digest", func() {
	It("generate a sha256 of the input", func() {
		out, err := Digest(metadata.DebianPackageListSourceMetadata{
			AptSources: []string{"deb example.com bionic main universe"},
			Packages: []metadata.DpkgPackage{
				{
					Package: "foo",
					Source: metadata.PackageSource{
						Version: "4.2",
					},
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())

		Expect(out).To(Equal("7ecf7aa2c71ee01ec2b90f37a3b8e944158e9aea6b8cee0290a7cb187884cf4c"))
	})

	It("generates the same digest for 2 different input instances with the same content", func() {
		input1 := metadata.DebianPackageListSourceMetadata{
			AptSources: []string{"deb example.com bionic main universe"},
			Packages: []metadata.DpkgPackage{
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
			Packages: []metadata.DpkgPackage{
				{
					Package: "foo",
					Source: metadata.PackageSource{
						Version: "4.2",
					},
				},
			},
		}

		out1, err := Digest(input1)
		Expect(err).ToNot(HaveOccurred())
		out2, err := Digest(input2)
		Expect(err).ToNot(HaveOccurred())

		Expect(&input1 != &input2).To(BeTrue())
		Expect(out2).To(Equal(out1))
	})

	It("generates a different digest for 2 different input instances with the same content in different order", func() {
		input1 := metadata.DebianPackageListSourceMetadata{
			AptSources: []string{"deb example.com bionic main universe"},
			Packages: []metadata.DpkgPackage{
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
			Packages: []metadata.DpkgPackage{
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

		out1, err := Digest(input1)
		Expect(err).ToNot(HaveOccurred())
		out2, err := Digest(input2)
		Expect(err).ToNot(HaveOccurred())

		Expect(&input1 != &input2).To(BeTrue())
		Expect(out2).ToNot(Equal(out1))
	})

	It("generates different digest for 2 different inputs", func() {
		input1 := metadata.DebianPackageListSourceMetadata{
			AptSources: []string{"deb example.com bionic main universe"},
			Packages: []metadata.DpkgPackage{
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
			Packages: []metadata.DpkgPackage{
				{
					Package: "foo",
					Source: metadata.PackageSource{
						Version: "4.2",
					},
				},
			},
		}

		out1, err := Digest(input1)
		Expect(err).ToNot(HaveOccurred())
		out2, err := Digest(input2)
		Expect(err).ToNot(HaveOccurred())

		Expect(out2).ToNot(Equal(out1))
	})
})
