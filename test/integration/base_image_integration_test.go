package integration_test

import (
	. "github.com/onsi/ginkgo/extensions/table"
	types2 "github.com/onsi/gomega/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("deplab", func() {
	DescribeTable("generates base property", func(imageAssetName string, matchers ...types2.GomegaMatcher) {
		metadataLabel := runDeplabAgainstTar(getTestAssetPath(imageAssetName))

		Expect(metadataLabel.Base).To(SatisfyAll(matchers...))
	},
		Entry("ubuntu:bionic image", "os-release-on-scratch.tgz",
			HaveKeyWithValue("name", "Ubuntu"),
			HaveKeyWithValue("version", "18.04.3 LTS (Bionic Beaver)"),
			HaveKeyWithValue("version_id", "18.04"),
			HaveKeyWithValue("id_like", "debian"),
			HaveKeyWithValue("version_codename", "bionic"),
			HaveKeyWithValue("pretty_name", "Ubuntu 18.04.3 LTS"),
		),
		Entry("an image that doesn't have an os-release", "all-file-types.tgz",
			HaveKeyWithValue("name", "unknown"),
			HaveKeyWithValue("version_codename", "unknown"),
			HaveKeyWithValue("version_id", "unknown"),
		),
	)
})
