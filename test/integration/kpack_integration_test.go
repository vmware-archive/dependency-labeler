// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package integration_test

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-tanzu/dependency-labeler/pkg/metadata"
)

var _ = Describe("deplab", func() {
	Describe("inspect", func() {
		Context("with an image with kpack metadata", func() {
			It("can find a kpack dependency", func() {
				stdOut, _ := runDepLab([]string{
					"inspect",
					"--image-tar", getTestAssetPath("image-archives/kpack-image.tgz"),
				}, 0)

				md := metadata.Metadata{}
				err := json.NewDecoder(stdOut).Decode(&md)

				Expect(err).ToNot(HaveOccurred())

				kpackDependency, ok := filterKpackDependencies(md.Dependencies)
				Expect(ok).To(BeTrue())

				Expect(kpackDependency.Source.Type).To(Equal("git"))

				kpackSourceVersion := kpackDependency.Source.Version
				kpackSourceMetadata := kpackDependency.Source.Metadata.(map[string]interface {})

				Expect(kpackSourceVersion["commit"]).To(Equal("1736b5e3b43a8cf40b3640821ee0e26049e1a58c"))
				Expect(kpackSourceMetadata["url"]).To(Equal("https://github.com/zmackie/github-actions-automate-projects.git"))
				Expect(kpackSourceMetadata["refs"]).To(Equal([]interface{}{}))
			})
		})
	})
})

func filterKpackDependencies(dependencies []metadata.Dependency) (metadata.Dependency, bool) {
	for _, dependency := range dependencies {
		if dependency.Type == metadata.PackageType {
			return dependency, true
		}
	}
	return metadata.Dependency{}, false
}
